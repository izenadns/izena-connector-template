package internal

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"testing"

	"github.com/izenadns/izena-connector-contract/abi"
	"github.com/izenadns/izena-connector-contract/connector"
)

func TestDispatchRoutesProviderAndZoneOperations(t *testing.T) {
	calls := 0
	p := Provider{HTTP: func(_ context.Context, request abi.HTTPRequest) (abi.HTTPResponse, error) {
		calls++
		if request.Method != http.MethodGet || !reflect.DeepEqual(request.Parts, []string{"health"}) ||
			request.Limit != 4096 || request.Headers.Get("Accept") != "application/json" {
			t.Fatal("unexpected example health request")
		}
		return abi.HTTPResponse{Status: http.StatusOK, Body: []byte(`{"version":"1.0.0"}`)}, nil
	}}
	settings := map[string]string{baseURLKey: "https://dns.example.test"}
	tests := []struct {
		op    abi.Operation
		input string
		want  any
		err   error
	}{
		{op: abi.Describe, want: p.Manifest()},
		{op: abi.PrepareSettings, want: connector.PreparedSettings{Settings: settings, TargetKey: settings[baseURLKey]}},
		{op: abi.Probe, want: connector.Profile{DaemonType: providerName, Version: "1.0.0"}},
		{op: abi.Capabilities, want: connector.ServerCapabilities{Provider: providerID, RecordModel: connector.RecordModelRRsets, RecordTypes: []string{"A"}, DNSSEC: unsupportedDNSSEC()}},
		{op: abi.ListZones, want: []connector.Zone{zone()}},
		{op: abi.GetZone, want: zone()},
		{op: abi.GetZoneMetadata, want: connector.ZoneMetadata{ID: exampleZoneID, Name: exampleZone, Kind: "example"}},
		{op: abi.GetRRset, want: rrset()},
		{op: abi.GetDNSSEC, err: connector.ErrUnsupported},
		{op: abi.SetDNSSEC, input: "true", err: connector.ErrUnsupported},
		{op: abi.CreateZone, input: `{"Name":"new.example.test."}`, err: connector.ErrUnsupported},
		{op: abi.UpdateZone, input: `{"Kind":"example"}`, want: connector.ZoneUpdateResult{}, err: connector.ErrUnsupported},
		{op: abi.DeleteZone, err: connector.ErrUnsupported},
		{op: abi.ReplaceRRset, input: `{"Name":"www.example.test.","Type":"A","TTL":300}`, err: connector.ErrUnsupported},
		{op: abi.DeleteRRset, err: connector.ErrUnsupported},
		{op: abi.RectifyZone, err: connector.ErrUnsupported},
		{op: abi.NotifyZone, err: connector.ErrUnsupported},
		{op: abi.AxfrZone, err: connector.ErrUnsupported},
		{op: "unknown_operation", err: connector.ErrUnsupported},
	}
	for _, test := range tests {
		t.Run(string(test.op), func(t *testing.T) {
			got, err := Dispatch(t.Context(), p, abi.Request{
				Operation: test.op, Input: []byte(test.input), Settings: settings,
				ZoneID: exampleZoneID, Name: "www." + exampleZone, Type: "A",
			})
			if !errors.Is(err, test.err) {
				t.Fatalf("error = %v, want %v", err, test.err)
			}
			if test.want != nil && !reflect.DeepEqual(got, test.want) {
				t.Fatalf("result = %#v, want %#v", got, test.want)
			}
		})
	}
	if calls != 1 {
		t.Fatalf("upstream calls = %d, want one probe request", calls)
	}
}

package internal

import (
	"context"
	"errors"
	"testing"

	"github.com/izenadns/izena-connector-contract/abi"
	"github.com/izenadns/izena-connector-contract/connector"
)

func TestIndividualRecordDispatchIsUnsupported(t *testing.T) {
	const recordInput = `{"Name":"www.example.test.","Type":"A","TTL":300,"Content":"192.0.2.1"}`
	for _, test := range []struct {
		operation abi.Operation
		input     string
		wantError error
	}{
		{operation: abi.ListRecords, wantError: connector.ErrUnsupported},
		{operation: abi.GetRecord, wantError: connector.ErrUnsupported},
		{operation: abi.CreateRecord, input: recordInput, wantError: connector.ErrUnsupported},
		{operation: abi.UpdateRecord, input: recordInput, wantError: connector.ErrUnsupported},
		{operation: abi.DeleteRecord, wantError: connector.ErrUnsupported},
		{operation: abi.CreateRecord, input: "{", wantError: connector.ErrValidation},
		{operation: abi.UpdateRecord, input: "{", wantError: connector.ErrValidation},
	} {
		t.Run(string(test.operation), func(t *testing.T) {
			req := abi.Request{
				Operation: test.operation, Input: []byte(test.input),
				ZoneID: "fixture-zone", RecordID: "fixture-record",
				Settings: map[string]string{baseURLKey: "https://dns.example.test"},
			}
			exchange := func(context.Context, abi.HTTPRequest) (abi.HTTPResponse, error) {
				t.Fatal("unsupported individual-record operation issued an HTTP request")
				return abi.HTTPResponse{}, nil
			}
			if _, err := Dispatch(t.Context(), Provider{HTTP: exchange}, req); !errors.Is(err, test.wantError) {
				t.Fatalf("Dispatch() = %v, want %v", err, test.wantError)
			}
		})
	}
}

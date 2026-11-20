package internal

import (
	"context"
	"encoding/json"

	"github.com/izenadns/izena-connector-contract/abi"
	"github.com/izenadns/izena-connector-contract/connector"
)

// Dispatch maps each v1 operation to the corresponding provider method.
func Dispatch(ctx context.Context, provider Provider, req abi.Request) (any, error) {
	switch req.Operation {
	case abi.Describe:
		return provider.Manifest(), nil
	case abi.PrepareSettings:
		return provider.PrepareSettings(ctx, req.Settings)
	}
	provider.HTTP = scopeHTTP(req.Operation, provider.HTTP)
	client, err := provider.connect(req.Settings)
	if err != nil {
		return nil, err
	}
	switch req.Operation {
	case abi.Probe:
		return client.Probe(ctx)
	case abi.Capabilities:
		return client.Capabilities(ctx)

	case abi.ListZones:
		return client.ListZones(ctx)
	case abi.ListZonesPage:
		return nil, connector.ErrUnsupported
	case abi.GetZone:
		return client.GetZone(ctx, req.ZoneID)
	case abi.GetZoneMetadata:
		return client.GetZoneMetadata(ctx, req.ZoneID)
	case abi.CreateZone:
		var input connector.ZoneCreateInput
		if err := json.Unmarshal(req.Input, &input); err != nil {
			return nil, connector.ErrValidation
		}
		return client.CreateZone(ctx, input)
	case abi.UpdateZone:
		var input connector.ZoneUpdateInput
		if err := json.Unmarshal(req.Input, &input); err != nil {
			return nil, connector.ErrValidation
		}
		return client.UpdateZone(ctx, req.ZoneID, input)
	case abi.DeleteZone:
		return nil, client.DeleteZone(ctx, req.ZoneID)

	case abi.GetZoneSettings:
		return client.GetZoneSettings(ctx, req.ZoneID)
	case abi.PatchZoneSettings:
		var patch connector.ZoneSettingsPatch
		if err := json.Unmarshal(req.Input, &patch); err != nil {
			return nil, connector.ErrValidation
		}
		return nil, client.PatchZoneSettings(ctx, req.ZoneID, patch)

	case abi.ListRRsetsPage:
		return nil, connector.ErrUnsupported
	case abi.GetRRset:
		return client.GetRRset(ctx, req.ZoneID, req.Name, req.Type)
	case abi.ReplaceRRset:
		var input connector.RRsetReplaceInput
		if err := json.Unmarshal(req.Input, &input); err != nil {
			return nil, connector.ErrValidation
		}
		return nil, client.ReplaceRRset(ctx, req.ZoneID, input)
	case abi.DeleteRRset:
		return nil, client.DeleteRRset(ctx, req.ZoneID, req.Name, req.Type)
	case abi.GetRRsetByID, abi.ReplaceRRsetByID, abi.DeleteRRsetByID:
		return nil, connector.ErrUnsupported

	case abi.ListRecords:
		return client.ListRecords(ctx, req.ZoneID)
	case abi.ListRecordsPage:
		return nil, connector.ErrUnsupported
	case abi.GetRecord:
		return client.GetRecord(ctx, req.ZoneID, req.RecordID)
	case abi.CreateRecord:
		var input connector.IndividualRecordInput
		if err := json.Unmarshal(req.Input, &input); err != nil {
			return nil, connector.ErrValidation
		}
		return client.CreateRecord(ctx, req.ZoneID, input)
	case abi.UpdateRecord:
		var input connector.IndividualRecordInput
		if err := json.Unmarshal(req.Input, &input); err != nil {
			return nil, connector.ErrValidation
		}
		return client.UpdateRecord(ctx, req.ZoneID, req.RecordID, input)
	case abi.DeleteRecord:
		return nil, client.DeleteRecord(ctx, req.ZoneID, req.RecordID)

	case abi.SubmitRecordBatch, abi.GetRecordBatchStatus:
		return nil, connector.ErrUnsupported

	case abi.GetDNSSEC:
		return client.GetDNSSEC(ctx, req.ZoneID)
	case abi.SetDNSSEC:
		var enabled bool
		if err := json.Unmarshal(req.Input, &enabled); err != nil {
			return nil, connector.ErrValidation
		}
		return nil, client.SetDNSSEC(ctx, req.ZoneID, enabled)

	case abi.RectifyZone:
		return nil, client.RectifyZone(ctx, req.ZoneID)
	case abi.NotifyZone:
		return nil, client.NotifyZone(ctx, req.ZoneID)
	case abi.AxfrZone:
		return nil, client.AxfrZone(ctx, req.ZoneID)
	default:
		return nil, connector.ErrUnsupported
	}
}

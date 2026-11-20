//go:build wasip1 && wasm

// Package main implements the WASI reactor entry point.
package main

import (
	"context"
	"encoding/json"
	"unsafe"

	"github.com/izenadns/izena-connector-contract/abi"
	"github.com/izenadns/izena-connector-contract/connector"
	"github.com/izenadns/izena-connector-template/internal"
)

//go:wasmimport izena_connector_v1 read_input
func readInput(dst unsafe.Pointer, size uint32) uint32

//go:wasmimport izena_connector_v1 write_output
func writeOutput(src unsafe.Pointer, size uint32) uint32

//go:wasmimport izena_connector_v1 http_request
func httpRequest(src unsafe.Pointer, size uint32) uint32

//go:wasmimport izena_connector_v1 http_response
func httpResponse(dst unsafe.Pointer, size uint32) uint32

//go:wasmexport abi_version
func abiVersion() uint32 { return abi.Version }

//go:wasmexport run
func run(size uint32) uint32 {
	if size == 0 || size > abi.MaxMessageBytes {
		return 0
	}
	raw := make([]byte, size)
	if readInput(unsafe.Pointer(unsafe.SliceData(raw)), size) != 1 {
		return 0
	}
	var request abi.Request
	if err := json.Unmarshal(raw, &request); err != nil {
		return 0
	}
	// Host imports enforce the calling request's context and deadline.
	value, err := internal.Dispatch(context.Background(), internal.Provider{HTTP: exchange}, request)
	result := abi.Result{Error: abi.EncodeFailure(err)}
	if value != nil {
		result.Value, err = abi.Marshal(value)
		if err != nil {
			return 0
		}
	}
	raw, err = abi.Marshal(result)
	if err != nil || len(raw) > abi.MaxMessageBytes {
		return 0
	}
	return writeOutput(unsafe.Pointer(unsafe.SliceData(raw)), uint32(len(raw)))
}

func exchange(_ context.Context, request abi.HTTPRequest) (abi.HTTPResponse, error) {
	raw, err := abi.Marshal(request)
	if err != nil || len(raw) > abi.MaxMessageBytes {
		return abi.HTTPResponse{}, connector.ErrValidation
	}
	size := httpRequest(unsafe.Pointer(unsafe.SliceData(raw)), uint32(len(raw)))
	if size == 0 || size > abi.MaxMessageBytes {
		return abi.HTTPResponse{}, connector.ErrUnavailable
	}
	raw = make([]byte, size)
	if httpResponse(unsafe.Pointer(unsafe.SliceData(raw)), size) != 1 {
		return abi.HTTPResponse{}, connector.ErrUnavailable
	}
	var response abi.HTTPResponse
	if err := json.Unmarshal(raw, &response); err != nil {
		return abi.HTTPResponse{}, connector.ErrMalformed
	}
	return response, response.Error.Decode()
}

func main() {}

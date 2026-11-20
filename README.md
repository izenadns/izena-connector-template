# Izena DNS connector

A WASI connector with a sample zone and an HTTP health probe. To adapt it,
update the module path in `go.mod` and the command import, rename the command
under `cmd/`, and implement your provider and tests under `internal/`.

## Development

```sh
make test
make test-integration-fixtures
make build
```

`make build` writes the Wasm artifact and its `.sha256` file. Run
`make artifact-path` for the artifact path and `make help` for other commands.

## License

[Apache License 2.0](LICENSE-TEMPLATE) applies to the template. You can choose
a license for your own changes while retaining the template's license and
notices when distributing its code.

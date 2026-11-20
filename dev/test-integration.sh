#!/usr/bin/env bash
set -euo pipefail

# The integration test starts and closes its own loopback HTTP server.
exec "$@"

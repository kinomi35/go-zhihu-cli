package configs

import _ "embed"

//go:embed endpoints.json
var DefaultEndpointsJSON []byte

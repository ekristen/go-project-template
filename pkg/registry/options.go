package registry

import "github.com/ekristen/go-telemetry/v2"

type RouteOptions struct {
	// TODO: you typically want to put some options here, like a database connection or something
	DB        interface{}
	Telemetry *telemetry.Telemetry
}

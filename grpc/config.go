package grpc

type ServiceConfig struct {
	Name string
	Port string
}

var Services = map[string]ServiceConfig{
	"trip-service": {
		Name: "trip-service",
		Port: "50003",
	},
	"identity-service": {
		Name: "identity-service",
		Port: "50002",
	},
	"driver-service": {
		Name: "driver-service",
		Port: "50004",
	},
	"rider-service": {
		Name: "rider-service",
		Port: "50001",
	},
}

// GetServiceConfig returns the ServiceConfig for a given service name.
// Returns the config and true if found, or an empty config and false if not found.
func GetServiceConfig(serviceName string) (ServiceConfig, bool) {
	cfg, ok := Services[serviceName]
	return cfg, ok
}

package grpc

type ServiceConfig struct {
	Name string
	Port string
}

var Services = map[string]ServiceConfig{
	"payment-service": {
		Name: "payment-service",
		Port: "50001", // Aligned with Dapr configuration
	},
	"trip-service": {
		Name: "trip-service",
		Port: "50002", // Aligned with Dapr configuration
	},
	"identity-service": {
		Name: "identity-service",
		Port: "50003", // Aligned with Dapr configuration
	},
	"driver-service": {
		Name: "driver-service",
		Port: "50004", // Aligned with Dapr configuration
	},
	"rider-service": {
		Name: "rider-service",
		Port: "50005", // Aligned with Dapr configuration
	},
}

// GetServiceConfig returns the ServiceConfig for a given service name.
// Returns the config and true if found, or an empty config and false if not found.
func GetServiceConfig(serviceName string) (ServiceConfig, bool) {
	cfg, ok := Services[serviceName]
	return cfg, ok
}

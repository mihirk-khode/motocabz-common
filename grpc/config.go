package grpc

type ServiceConfig struct {
	Name string
	Port string
}

var Services = map[string]ServiceConfig{
	"payment-service": {
		Name: "payment-service",
		Port: "50006", // Payment service gRPC port (from Kubernetes service)
	},
	"trip-service": {
		Name: "trip-service",
		Port: "50007", // Trip service gRPC port (from Kubernetes service)
	},
	"identity-service": {
		Name: "identity-service",
		Port: "50003", // Identity service gRPC port (from Kubernetes service)
	},
	"driver-service": {
		Name: "driver-service",
		Port: "50004", // Driver service gRPC port (from Kubernetes service)
	},
	"rider-service": {
		Name: "rider-service",
		Port: "50005", // Rider service gRPC port (from Kubernetes service)
	},
	"admin-service": {
		Name: "admin-service",
		Port: "50008", // Admin service gRPC port (from Kubernetes service)
	},
}

// GetServiceConfig returns the ServiceConfig for a given service name.
// Returns the config and true if found, or an empty config and false if not found.
func GetServiceConfig(serviceName string) (ServiceConfig, bool) {
	cfg, ok := Services[serviceName]
	return cfg, ok
}

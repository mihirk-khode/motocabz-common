package grpc

type ServiceConfig struct {
	Name string
	Host string
	Port string
}

var Services = map[string]ServiceConfig{
	"payment-service": {
		Name: "payment-service",
		Host: "",
		Port: "50055", // Payment service gRPC port (from Kubernetes service)
	},
	"trip-service": {
		Name: "trip-service",
		Host: "k8s-motocabz-tripserv-b73ea6daa4-2014578782.eu-north-1.elb.amazonaws.com",
		Port: "50051", // Trip service gRPC port (from Kubernetes service)
	},
	"identity-service": {
		Name: "identity-service",
		Host: "k8s-motocabz-identity-a1bc81e43d-816416611.eu-north-1.elb.amazonaws.com",
		Port: "50057", // Identity service gRPC port (from Kubernetes service)
	},
	"driver-service": {
		Name: "driver-service",
		Host: "k8s-motocabz-driverse-ffd6064118-284108539.eu-north-1.elb.amazonaws.com",
		Port: "50052", // Driver service Ingress port (ALB)
	},
	"rider-service": {
		Name: "rider-service",
		Host: "k8s-motocabz-riderser-0ca9d74497-328343501.eu-north-1.elb.amazonaws.com",
		Port: "50053", // Rider service gRPC port (from Kubernetes service)
	},
}

// GetServiceConfig returns the ServiceConfig for a given service name.
// Returns the config and true if found, or an empty config and false if not found.
func GetServiceConfig(serviceName string) (ServiceConfig, bool) {
	cfg, ok := Services[serviceName]
	return cfg, ok
}

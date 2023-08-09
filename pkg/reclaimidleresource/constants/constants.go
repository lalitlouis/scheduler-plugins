package constants

/*
Common constants
*/

// DevENVFlag for local builds
const DevENVFlag = "development"

// TestsENVFlag for unit testing
const TestsENVFlag = "testing"

// HTTPPrefix defines http prefix
const HTTPPrefix = `http://`

// HTTPSPrefix defines http prefix
const HTTPSPrefix = `https://`

/*
Prometheus constants
*/

// PrometheusServiceName Name of prometheus service
const PrometheusServiceName = "prometheus-kube-prometheus-prometheus"

// PrometheusNamespace is the namespace where prometheus is installed
const PrometheusNamespace = "prometheus"

// PrometheusPort is the port service
const PrometheusPort = "9090"

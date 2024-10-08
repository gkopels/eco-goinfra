package v1alpha1

const (

	// EnvVarFakeInClusterConfig is an environment variable to fake an in-cluster RESTConfig using
	// the current kubectl context (for development purposes)
	EnvVarFakeInClusterConfig = "ARGOCD_FAKE_IN_CLUSTER"

	// EnvK8sClientQPS is the QPS value used for the kubernetes client (default: 50)
	EnvK8sClientQPS = "ARGOCD_K8S_CLIENT_QPS"

	// EnvK8sClientBurst is the burst value used for the kubernetes client (default: twice the client QPS)
	EnvK8sClientBurst = "ARGOCD_K8S_CLIENT_BURST"

	// EnvK8sClientMaxIdleConnections is the number of max idle connections in K8s REST client HTTP transport (default: 500)
	EnvK8sClientMaxIdleConnections = "ARGOCD_K8S_CLIENT_MAX_IDLE_CONNECTIONS"

	// EnvK8sTCPTimeout is the duration for TCP timeouts when communicating with K8s API servers
	EnvK8sTCPTimeout = "ARGOCD_K8S_TCP_TIMEOUT"

	// EnvK8sTCPKeepalive is the interval for TCP keep alive probes to be sent when communicating with K8s API servers
	EnvK8sTCPKeepAlive = "ARGOCD_K8S_TCP_KEEPALIVE"

	// EnvK8sTLSHandshakeTimeout is the duration for TLS handshake timeouts when establishing connections to K8s API servers
	EnvK8sTLSHandshakeTimeout = "ARGOCD_K8S_TLS_HANDSHAKE_TIMEOUT"

	// EnvK8sTCPIdleConnTimeout is the duration when idle TCP connection to the K8s API servers should timeout
	EnvK8sTCPIdleConnTimeout = "ARGOCD_K8S_TCP_IDLE_TIMEOUT"
)

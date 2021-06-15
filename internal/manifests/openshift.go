package manifests

// OpenshiftOptions contains flags that activate Openshift features
type OpenshiftOptions struct {
	EnableCertificateSigningService bool
	EnableServiceMonitors           bool
	EnableTLSEnabledServiceMonitors bool
}

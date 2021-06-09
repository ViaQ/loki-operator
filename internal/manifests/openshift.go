package manifests

// UseCertificateSigningService is a flag to allow certificate signing for the http services.
var UseCertificateSigningService = false

// EnableOpenshiftFeatures activates all flags to add or configure the various
// objects to allow them to take advantage of Openshift features.
func EnableOpenshiftFeatures() {
	UseCertificateSigningService = true
}

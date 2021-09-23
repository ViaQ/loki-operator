package gateway

import (
	lokiv1beta1 "github.com/ViaQ/loki-operator/api/v1beta1"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuild(t *testing.T) {
	exp := `
tenants:
- name: test-a
  id: test
  oidc:
    clientID: test
    clientSecret: test123
    issuerCAPath: /tmp/ca/path

    issuerURL: https://127.0.0.1:5556/dex
    redirectURL: https://localhost:8443/oidc/test-a/callback
    
    usernameClaim: test
    groupClaim: test
  opa:
    url: http://127.0.0.1:8181/v1/data/observatorium/allow
`
	opts := Options{
		Stack: lokiv1beta1.LokiStackSpec{
			Tenants: &lokiv1beta1.TenantsSpec{
				Mode: lokiv1beta1.Dynamic,
				Authentication: []*lokiv1beta1.AuthenticationSpec{
					{
						Name: "test-a",
						ID:   "test",
						OIDC: &lokiv1beta1.OIDCSpec{
							Secret: &lokiv1beta1.TenantSecretSpec{
								Name: "test",
							},
							IssuerURL:     "https://127.0.0.1:5556/dex",
							RedirectURL:   "https://localhost:8443/oidc/test-a/callback",
							GroupClaim:    "test",
							UsernameClaim: "test",
						},
					},
				},
				Authorization: &lokiv1beta1.AuthorizationSpec{
					OPA: &lokiv1beta1.OPASpec{
						URL: "http://127.0.0.1:8181/v1/data/observatorium/allow",
					},
				},
			},
		},
		Namespace:     "test-ns",
		Name:          "test",
		GatewaySecret: []*Secret{
			{
				TenantName:   "test-a",
				ClientID:     "test",
				ClientSecret: "test123",
				IssuerCAPath: "/tmp/ca/path",
			},
		},
	}
	_, tenantsConfig, _, err := Build(opts)
	require.NoError(t, err)
	require.YAMLEq(t, exp, string(tenantsConfig))
}

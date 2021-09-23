package gateway

import (
	"testing"

	lokiv1beta1 "github.com/ViaQ/loki-operator/api/v1beta1"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestValidateModes_StaticMode(t *testing.T) {
	type test struct {
		name    string
		wantErr string
		stack   lokiv1beta1.LokiStack
	}
	table := []test{
		{
			name:    "missing authentication spec",
			wantErr: "mandatory configuration - missing tenants' authentication configuration",
			stack: lokiv1beta1.LokiStack{
				TypeMeta: metav1.TypeMeta{
					Kind: "LokiStack",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-stack",
					Namespace: "some-ns",
					UID:       "b23f9a38-9672-499f-8c29-15ede74d3ece",
				},
				Spec: lokiv1beta1.LokiStackSpec{
					Size: lokiv1beta1.SizeOneXExtraSmall,
					Tenants: &lokiv1beta1.TenantsSpec{
						Mode: "static",
					},
				},
			},
		},
		{
			name:    "missing roles spec",
			wantErr: "mandatory configuration - missing roles configuration",
			stack: lokiv1beta1.LokiStack{
				TypeMeta: metav1.TypeMeta{
					Kind: "LokiStack",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-stack",
					Namespace: "some-ns",
					UID:       "b23f9a38-9672-499f-8c29-15ede74d3ece",
				},
				Spec: lokiv1beta1.LokiStackSpec{
					Size: lokiv1beta1.SizeOneXExtraSmall,
					Tenants: &lokiv1beta1.TenantsSpec{
						Mode: "static",
						Authentication: []lokiv1beta1.AuthenticationSpec{
							{
								Name: "test",
								ID:   "1234",
								OIDC: &lokiv1beta1.OIDCSpec{
									IssuerURL:     "some-url",
									RedirectURL:   "some-other-url",
									GroupClaim:    "test",
									UsernameClaim: "test",
								},
							},
						},
						Authorization: &lokiv1beta1.AuthorizationSpec{
							Roles: nil,
						},
					},
				},
			},
		},
		{
			name:    "missing role bindings spec",
			wantErr: "mandatory configuration - missing role bindings configuration",
			stack: lokiv1beta1.LokiStack{
				TypeMeta: metav1.TypeMeta{
					Kind: "LokiStack",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-stack",
					Namespace: "some-ns",
					UID:       "b23f9a38-9672-499f-8c29-15ede74d3ece",
				},
				Spec: lokiv1beta1.LokiStackSpec{
					Size: lokiv1beta1.SizeOneXExtraSmall,
					Tenants: &lokiv1beta1.TenantsSpec{
						Mode: "static",
						Authentication: []lokiv1beta1.AuthenticationSpec{
							{
								Name: "test",
								ID:   "1234",
								OIDC: &lokiv1beta1.OIDCSpec{
									IssuerURL:     "some-url",
									RedirectURL:   "some-other-url",
									GroupClaim:    "test",
									UsernameClaim: "test",
								},
							},
						},
						Authorization: &lokiv1beta1.AuthorizationSpec{
							Roles: []lokiv1beta1.RoleSpec{
								{
									Name:        "some-name",
									Resources:   []string{"test"},
									Tenants:     []string{"test"},
									Permissions: []lokiv1beta1.PermissionType{"read"},
								},
							},
							RoleBindings: nil,
						},
					},
				},
			},
		},
		{
			name:    "extra OPA URL provided",
			wantErr: "extra configuration provided - OPA Url is not required",
			stack: lokiv1beta1.LokiStack{
				TypeMeta: metav1.TypeMeta{
					Kind: "LokiStack",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-stack",
					Namespace: "some-ns",
					UID:       "b23f9a38-9672-499f-8c29-15ede74d3ece",
				},
				Spec: lokiv1beta1.LokiStackSpec{
					Size: lokiv1beta1.SizeOneXExtraSmall,
					Tenants: &lokiv1beta1.TenantsSpec{
						Mode: "static",
						Authentication: []lokiv1beta1.AuthenticationSpec{
							{
								Name: "test",
								ID:   "1234",
								OIDC: &lokiv1beta1.OIDCSpec{
									IssuerURL:     "some-url",
									RedirectURL:   "some-other-url",
									GroupClaim:    "test",
									UsernameClaim: "test",
								},
							},
						},
						Authorization: &lokiv1beta1.AuthorizationSpec{
							OPA: &lokiv1beta1.OPASpec{
								URL: "some-url",
							},
							Roles: []lokiv1beta1.RoleSpec{
								{
									Name:        "some-name",
									Resources:   []string{"test"},
									Tenants:     []string{"test"},
									Permissions: []lokiv1beta1.PermissionType{"read"},
								},
							},
							RoleBindings: []lokiv1beta1.RoleBindingsSpec{
								{
									Name: "some-name",
									Subjects: []lokiv1beta1.Subject{
										{
											Name: "sub-1",
											Kind: "user",
										},
									},
									Roles: []string{"some-role"},
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "all set",
			wantErr: "",
			stack: lokiv1beta1.LokiStack{
				TypeMeta: metav1.TypeMeta{
					Kind: "LokiStack",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-stack",
					Namespace: "some-ns",
					UID:       "b23f9a38-9672-499f-8c29-15ede74d3ece",
				},
				Spec: lokiv1beta1.LokiStackSpec{
					Size: lokiv1beta1.SizeOneXExtraSmall,
					Tenants: &lokiv1beta1.TenantsSpec{
						Mode: "static",
						Authentication: []lokiv1beta1.AuthenticationSpec{
							{
								Name: "test",
								ID:   "1234",
								OIDC: &lokiv1beta1.OIDCSpec{
									IssuerURL:     "some-url",
									RedirectURL:   "some-other-url",
									GroupClaim:    "test",
									UsernameClaim: "test",
								},
							},
						},
						Authorization: &lokiv1beta1.AuthorizationSpec{
							Roles: []lokiv1beta1.RoleSpec{
								{
									Name:        "some-name",
									Resources:   []string{"test"},
									Tenants:     []string{"test"},
									Permissions: []lokiv1beta1.PermissionType{"read"},
								},
							},
							RoleBindings: []lokiv1beta1.RoleBindingsSpec{
								{
									Name: "some-name",
									Subjects: []lokiv1beta1.Subject{
										{
											Name: "sub-1",
											Kind: "user",
										},
									},
									Roles: []string{"some-role"},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tst := range table {
		tst := tst
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateModes(tst.stack)
			if tst.wantErr != "" {
				require.EqualError(t, err, tst.wantErr)
			}
		})
	}
}

func TestValidateModes_DynamicMode(t *testing.T) {
	type test struct {
		name    string
		wantErr string
		stack   lokiv1beta1.LokiStack
	}
	table := []test{
		{
			name:    "missing authentication spec",
			wantErr: "mandatory configuration - missing tenants configuration",
			stack: lokiv1beta1.LokiStack{
				TypeMeta: metav1.TypeMeta{
					Kind: "LokiStack",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-stack",
					Namespace: "some-ns",
					UID:       "b23f9a38-9672-499f-8c29-15ede74d3ece",
				},
				Spec: lokiv1beta1.LokiStackSpec{
					Size: lokiv1beta1.SizeOneXExtraSmall,
					Tenants: &lokiv1beta1.TenantsSpec{
						Mode: "dynamic",
					},
				},
			},
		},
		{
			name:    "missing OPA URL spec",
			wantErr: "mandatory configuration - missing OPA Url",
			stack: lokiv1beta1.LokiStack{
				TypeMeta: metav1.TypeMeta{
					Kind: "LokiStack",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-stack",
					Namespace: "some-ns",
					UID:       "b23f9a38-9672-499f-8c29-15ede74d3ece",
				},
				Spec: lokiv1beta1.LokiStackSpec{
					Size: lokiv1beta1.SizeOneXExtraSmall,
					Tenants: &lokiv1beta1.TenantsSpec{
						Mode: "dynamic",
						Authentication: []lokiv1beta1.AuthenticationSpec{
							{
								Name: "test",
								ID:   "1234",
								OIDC: &lokiv1beta1.OIDCSpec{
									IssuerURL:     "some-url",
									RedirectURL:   "some-other-url",
									GroupClaim:    "test",
									UsernameClaim: "test",
								},
							},
						},
						Authorization: &lokiv1beta1.AuthorizationSpec{
							OPA: nil,
						},
					},
				},
			},
		},
		{
			name:    "extra roles configuration provided",
			wantErr: "extra configuration provided - roles configuration is not required",
			stack: lokiv1beta1.LokiStack{
				TypeMeta: metav1.TypeMeta{
					Kind: "LokiStack",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-stack",
					Namespace: "some-ns",
					UID:       "b23f9a38-9672-499f-8c29-15ede74d3ece",
				},
				Spec: lokiv1beta1.LokiStackSpec{
					Size: lokiv1beta1.SizeOneXExtraSmall,
					Tenants: &lokiv1beta1.TenantsSpec{
						Mode: "dynamic",
						Authentication: []lokiv1beta1.AuthenticationSpec{
							{
								Name: "test",
								ID:   "1234",
								OIDC: &lokiv1beta1.OIDCSpec{
									IssuerURL:     "some-url",
									RedirectURL:   "some-other-url",
									GroupClaim:    "test",
									UsernameClaim: "test",
								},
							},
						},
						Authorization: &lokiv1beta1.AuthorizationSpec{
							OPA: &lokiv1beta1.OPASpec{
								URL: "some-url",
							},
							Roles: []lokiv1beta1.RoleSpec{
								{
									Name:        "some-name",
									Resources:   []string{"test"},
									Tenants:     []string{"test"},
									Permissions: []lokiv1beta1.PermissionType{"read"},
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "extra role bindings configuration provided",
			wantErr: "extra configuration provided - role bindings configuration is not required",
			stack: lokiv1beta1.LokiStack{
				TypeMeta: metav1.TypeMeta{
					Kind: "LokiStack",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-stack",
					Namespace: "some-ns",
					UID:       "b23f9a38-9672-499f-8c29-15ede74d3ece",
				},
				Spec: lokiv1beta1.LokiStackSpec{
					Size: lokiv1beta1.SizeOneXExtraSmall,
					Tenants: &lokiv1beta1.TenantsSpec{
						Mode: "dynamic",
						Authentication: []lokiv1beta1.AuthenticationSpec{
							{
								Name: "test",
								ID:   "1234",
								OIDC: &lokiv1beta1.OIDCSpec{
									IssuerURL:     "some-url",
									RedirectURL:   "some-other-url",
									GroupClaim:    "test",
									UsernameClaim: "test",
								},
							},
						},
						Authorization: &lokiv1beta1.AuthorizationSpec{
							OPA: &lokiv1beta1.OPASpec{
								URL: "some-url",
							},
							RoleBindings: []lokiv1beta1.RoleBindingsSpec{
								{
									Name: "some-name",
									Subjects: []lokiv1beta1.Subject{
										{
											Name: "sub-1",
											Kind: "user",
										},
									},
									Roles: []string{"some-role"},
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "all set",
			wantErr: "",
			stack: lokiv1beta1.LokiStack{
				TypeMeta: metav1.TypeMeta{
					Kind: "LokiStack",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-stack",
					Namespace: "some-ns",
					UID:       "b23f9a38-9672-499f-8c29-15ede74d3ece",
				},
				Spec: lokiv1beta1.LokiStackSpec{
					Size: lokiv1beta1.SizeOneXExtraSmall,
					Tenants: &lokiv1beta1.TenantsSpec{
						Mode: "dynamic",
						Authentication: []lokiv1beta1.AuthenticationSpec{
							{
								Name: "test",
								ID:   "1234",
								OIDC: &lokiv1beta1.OIDCSpec{
									IssuerURL:     "some-url",
									RedirectURL:   "some-other-url",
									GroupClaim:    "test",
									UsernameClaim: "test",
								},
							},
						},
						Authorization: &lokiv1beta1.AuthorizationSpec{
							OPA: &lokiv1beta1.OPASpec{
								URL: "some-url",
							},
						},
					},
				},
			},
		},
	}
	for _, tst := range table {
		tst := tst
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateModes(tst.stack)
			if tst.wantErr != "" {
				require.EqualError(t, err, tst.wantErr)
			}
		})
	}
}

func TestValidateModes_OpenshiftLoggingMode(t *testing.T) {
	type test struct {
		name    string
		wantErr string
		stack   lokiv1beta1.LokiStack
	}
	table := []test{
		{
			name:    "provided authentication spec",
			wantErr: "extra configuration provided - tenants configuration is not required.",
			stack: lokiv1beta1.LokiStack{
				TypeMeta: metav1.TypeMeta{
					Kind: "LokiStack",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-stack",
					Namespace: "some-ns",
					UID:       "b23f9a38-9672-499f-8c29-15ede74d3ece",
				},
				Spec: lokiv1beta1.LokiStackSpec{
					Size: lokiv1beta1.SizeOneXExtraSmall,
					Tenants: &lokiv1beta1.TenantsSpec{
						Mode: "openshift-logging",
						Authentication: []lokiv1beta1.AuthenticationSpec{
							{
								Name: "test",
								ID:   "1234",
								OIDC: &lokiv1beta1.OIDCSpec{
									IssuerURL:     "some-url",
									RedirectURL:   "some-other-url",
									GroupClaim:    "test",
									UsernameClaim: "test",
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "provided authorization spec",
			wantErr: "extra configuration provided - authorization configuration is not required.",
			stack: lokiv1beta1.LokiStack{
				TypeMeta: metav1.TypeMeta{
					Kind: "LokiStack",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-stack",
					Namespace: "some-ns",
					UID:       "b23f9a38-9672-499f-8c29-15ede74d3ece",
				},
				Spec: lokiv1beta1.LokiStackSpec{
					Size: lokiv1beta1.SizeOneXExtraSmall,
					Tenants: &lokiv1beta1.TenantsSpec{
						Mode:           "openshift-logging",
						Authentication: nil,
						Authorization: &lokiv1beta1.AuthorizationSpec{
							OPA: &lokiv1beta1.OPASpec{
								URL: "some-url",
							},
						},
					},
				},
			},
		},
		{
			name:    "all set",
			wantErr: "",
			stack: lokiv1beta1.LokiStack{
				TypeMeta: metav1.TypeMeta{
					Kind: "LokiStack",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-stack",
					Namespace: "some-ns",
					UID:       "b23f9a38-9672-499f-8c29-15ede74d3ece",
				},
				Spec: lokiv1beta1.LokiStackSpec{
					Size: lokiv1beta1.SizeOneXExtraSmall,
					Tenants: &lokiv1beta1.TenantsSpec{
						Mode: "openshift-logging",
					},
				},
			},
		},
	}
	for _, tst := range table {
		tst := tst
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()

			err := ValidateModes(tst.stack)
			if tst.wantErr != "" {
				require.EqualError(t, err, tst.wantErr)
			}
		})
	}
}

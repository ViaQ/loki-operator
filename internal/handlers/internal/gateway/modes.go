package gateway

import (
	"github.com/ViaQ/logerr/kverrors"
	lokiv1beta1 "github.com/ViaQ/loki-operator/api/v1beta1"
)

// ValidateModes validates the tenants mode specification.
func ValidateModes(stack lokiv1beta1.LokiStack) error {
	if stack.Spec.Tenants.Mode == lokiv1beta1.Static {
		if stack.Spec.Tenants.Authentication == nil {
			return kverrors.New("mandatory configuration - missing tenants configuration")
		}

		if stack.Spec.Tenants.Authorization == nil || stack.Spec.Tenants.Authorization.Roles == nil {
			return kverrors.New("mandatory configuration - missing roles configuration")
		}

		if stack.Spec.Tenants.Authorization == nil || stack.Spec.Tenants.Authorization.RoleBindings == nil {
			return kverrors.New("mandatory configuration - missing role bindings configuration")
		}
	}

	if stack.Spec.Tenants.Mode == lokiv1beta1.Dynamic {
		if stack.Spec.Tenants.Authentication == nil {
			return kverrors.New("mandatory configuration - missing tenants configuration")
		}

		if stack.Spec.Tenants.Authorization == nil || stack.Spec.Tenants.Authorization.OPA == nil {
			return kverrors.New("mandatory configuration - missing OPA Url")
		}
	}

	if stack.Spec.Tenants.Mode == lokiv1beta1.OpenshiftLogging {
		if stack.Spec.Tenants.Authentication != nil {
			return kverrors.New("extra configuration provided - tenants configuration is not required.")
		}

		if stack.Spec.Tenants.Authorization != nil {
			return kverrors.New("extra configuration provided - authorization configuration is not required.")
		}
	}

	return nil
}

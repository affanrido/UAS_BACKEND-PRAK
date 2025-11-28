package middleware

import "UAS_BACKEND/domain/service"

func RequirePermission(rbac *service.RBACService, userID uint, perm string) error {
	ok, err := rbac.UserHasPermission(userID, perm)
	if err != nil {
		return err
	}
	if !ok {
		return service.ErrForbidden
	}
	return nil
}

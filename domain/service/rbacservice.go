package service

import (
	model "UAS_BACKEND/domain/Model"
	"errors"
)

var ErrForbidden = errors.New("forbidden")

type RBACService struct {
	// FUNGSIONAL TANPA GORM
	GetUserByID      func(id uint) (*model.Users, error)
	GetRoleByID      func(id uint) (*model.RolePermission, error)
	GetPermissionByID func(id uint) (*model.Permissions, error)
}

func (r *RBACService) UserHasPermission(userID uint, permName string) (bool, error) {
	user, err := r.GetUserByID(userID)
	if err != nil {
		return false, err
	}

	for _, rid := range user.RoleID {
		role, err := r.GetRoleByID(uint(rid))
		if err != nil {
			continue
		}

		for _, pid := range role.PermissionID {
			p, err := r.GetPermissionByID(uint(pid))
			if err != nil {
				continue
			}

			if p.Name == permName {
				return true, nil
			}
		}
	}

	return false, nil
}

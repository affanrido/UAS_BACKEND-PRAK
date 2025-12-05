package repository

import (
	model "UAS_BACKEND/domain/Model"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

type RBACRepository struct {
	DB *sql.DB
}

func NewRBACRepository(db *sql.DB) *RBACRepository {
	return &RBACRepository{DB: db}
}

// GetRoleByID - Mendapatkan role berdasarkan ID
func (r *RBACRepository) GetRoleByID(roleID uuid.UUID) (*model.Roles, error) {
	query := `
		SELECT id, name, description, created_at
		FROM roles
		WHERE id = $1
	`

	var role model.Roles
	err := r.DB.QueryRow(query, roleID).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}

	return &role, nil
}

// GetUserPermissions - Mendapatkan semua permissions berdasarkan role_id
func (r *RBACRepository) GetUserPermissions(roleID uuid.UUID) ([]string, error) {
	query := `
		SELECT p.name
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
	`

	rows, err := r.DB.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var permission string
		if err := rows.Scan(&permission); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	return permissions, nil
}

// GetPermissionByName - Mendapatkan permission berdasarkan name
func (r *RBACRepository) GetPermissionByName(name string) (*model.Permission, error) {
	query := `
		SELECT id, name, resource, action, description
		FROM permissions
		WHERE name = $1
	`

	var perm model.Permission
	err := r.DB.QueryRow(query, name).Scan(
		&perm.ID,
		&perm.Name,
		&perm.Resource,
		&perm.Action,
		&perm.Description,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}

	return &perm, nil
}

// GetAllRoles - Mendapatkan semua roles
func (r *RBACRepository) GetAllRoles() ([]model.Roles, error) {
	query := `
		SELECT id, name, description, created_at
		FROM roles
		ORDER BY name
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []model.Roles
	for rows.Next() {
		var role model.Roles
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// GetAllPermissions - Mendapatkan semua permissions
func (r *RBACRepository) GetAllPermissions() ([]model.Permission, error) {
	query := `
		SELECT id, name, resource, action, description
		FROM permissions
		ORDER BY resource, action
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []model.Permission
	for rows.Next() {
		var perm model.Permission
		if err := rows.Scan(&perm.ID, &perm.Name, &perm.Resource, &perm.Action, &perm.Description); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

// GetRolePermissions - Mendapatkan semua permissions untuk role tertentu
func (r *RBACRepository) GetRolePermissions(roleID uuid.UUID) ([]model.Permission, error) {
	query := `
		SELECT p.id, p.name, p.resource, p.action, p.description
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.resource, p.action
	`

	rows, err := r.DB.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []model.Permission
	for rows.Next() {
		var perm model.Permission
		if err := rows.Scan(&perm.ID, &perm.Name, &perm.Resource, &perm.Action, &perm.Description); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

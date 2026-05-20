package seeders

import (
	"database/sql"
	"fmt"
)

var rolePermissions = []struct {
	RoleName       string
	PermissionName string
}{
	// SUPER ADMIN: Full access to everything
	{"ROLE_SUPER_ADMIN", "CREATE_ROLES"},
	{"ROLE_SUPER_ADMIN", "READ_ROLES"},
	{"ROLE_SUPER_ADMIN", "UPDATE_ROLES"},
	{"ROLE_SUPER_ADMIN", "DELETE_ROLES"},
	{"ROLE_SUPER_ADMIN", "MANAGE_ROLE_PERMISSIONS"},
	{"ROLE_SUPER_ADMIN", "CREATE_USERS"},
	{"ROLE_SUPER_ADMIN", "READ_USERS"},
	{"ROLE_SUPER_ADMIN", "UPDATE_USERS"},
	{"ROLE_SUPER_ADMIN", "DELETE_USERS"},
	{"ROLE_SUPER_ADMIN", "MANAGE_USER_ROLES"},
	{"ROLE_SUPER_ADMIN", "READ_PERMISSIONS"},
	{"ROLE_SUPER_ADMIN", "UPDATE_PERMISSIONS"},

	// ADMIN: CRUD for users
	{"ROLE_ADMIN", "CREATE_USERS"},
	{"ROLE_ADMIN", "READ_USERS"},
	{"ROLE_ADMIN", "UPDATE_USERS"},
	{"ROLE_ADMIN", "DELETE_USERS"},

	// USER: Read access
	{"ROLE_USER", "READ_USERS"},
}

func SeedRolePermissions(db *sql.DB) error {
	for _, rp := range rolePermissions {
		var roleId string
		if err := db.QueryRow(`SELECT id FROM roles WHERE name = $1`, rp.RoleName).Scan(&roleId); err != nil {
			return fmt.Errorf("failed to fetch role ID for %s: %v", rp.RoleName, err)
		}

		var permissionId string
		if err := db.QueryRow(`SELECT id FROM permissions WHERE name = $1`, rp.PermissionName).Scan(&permissionId); err != nil {
			return fmt.Errorf("failed to fetch permission ID for %s: %v", rp.PermissionName, err)
		}

		if _, err := db.Exec(`INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT (role_id, permission_id) DO NOTHING`, roleId, permissionId); err != nil {
			return fmt.Errorf("failed to seed role_permissions: %v", err)
		}
	}

	return nil
}

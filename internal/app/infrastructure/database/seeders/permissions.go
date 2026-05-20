package seeders

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/permission"
)

var permissions = []permission.Permission{
	{Id: ulid.MustParse("01FY7XRMMKB4FA7G0Q0D9S8CDN"), Name: "CREATE_ROLES", Description: "Allows the user to create roles"},
	{Id: ulid.MustParse("01FY7XP5V2EPJZFG361WRHJDVK"), Name: "READ_ROLES", Description: "Allows the user to have read access to roles"},
	{Id: ulid.MustParse("01FY7XTB323SXWWJ757AY5QJ7H"), Name: "UPDATE_ROLES", Description: "Allows the user to update roles"},
	{Id: ulid.MustParse("01FY7XVSJQHAC040RMMA37ZTNR"), Name: "DELETE_ROLES", Description: "Allows the user to delete roles"},
	{Id: ulid.MustParse("01FY7XXQMW888MCBXH67HADFY4"), Name: "MANAGE_ROLE_PERMISSIONS", Description: "Allows the user to manage role permissions"},

	{Id: ulid.MustParse("01FY7XRW3JSY2Y4Q8XRVDYSCZK"), Name: "CREATE_USERS", Description: "Allows the user to create users"},
	{Id: ulid.MustParse("01FY7XMMX83NKP6Y0BSDEJ1HQP"), Name: "READ_USERS", Description: "Allows the user to have read access to users"},
	{Id: ulid.MustParse("01FY7XTMCG9EWGWW0K2DBF4BJJ"), Name: "UPDATE_USERS", Description: "Allows the user to update users"},
	{Id: ulid.MustParse("01FY7XW2NRY5FKSKTQ748TAY0D"), Name: "DELETE_USERS", Description: "Allows the user to delete users"},
	{Id: ulid.MustParse("01FY7XYK0ZEDJ9Z4RBXZQD5FW4"), Name: "MANAGE_USER_ROLES", Description: "Allows the user to manage user roles"},

	{Id: ulid.MustParse("01FY7XQ0F9XK4B74QR4BGTXNFD"), Name: "READ_PERMISSIONS", Description: "Allows the user to have read access to permissions"},
	{Id: ulid.MustParse("01FY7XT3SRMCD2RNB9X8WR7VYJ"), Name: "UPDATE_PERMISSIONS", Description: "Allows the user to update permissions"},
}

func SeedPermissions(db *sql.DB) error {
	for _, p := range permissions {
		q := `INSERT INTO permissions (id, name, description, created_at) VALUES ($1, $2, $3, $4) ON CONFLICT (name) DO NOTHING`
		if _, err := db.Exec(q, p.Id.String(), p.Name, p.Description, time.Now()); err != nil {
			return fmt.Errorf("failed to seed permissions: %v", err)
		}
	}
	return nil
}

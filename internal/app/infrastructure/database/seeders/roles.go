package seeders

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/role"
)

var roles = []role.Role{
	{Id: ulid.MustParse("01FY87HWCDFCC3D525G552ZEN5"), Name: "ROLE_USER"},
	{Id: ulid.MustParse("01FY8BHR09JKSXBMYX325PM3SV"), Name: "ROLE_ADMIN"},
	{Id: ulid.MustParse("01FY87J3TQJKS6EN8J1CJ75BVR"), Name: "ROLE_SUPER_ADMIN"},
}

func SeedRoles(db *sql.DB) error {
	for _, r := range roles {
		q := `INSERT INTO roles (id, name, created_at) VALUES ($1, $2, $3) ON CONFLICT (name) DO NOTHING`
		if _, err := db.Exec(q, r.Id.String(), r.Name, time.Now()); err != nil {
			return fmt.Errorf("failed to seed roles: %v", err)
		}
	}
	return nil
}

package seeders

import (
	"database/sql"
	"log"
)

func SeedDatabase(db *sql.DB) {
	if err := SeedRoles(db); err != nil {
		log.Fatalf("Error seeding roles: %v", err)
	}
	log.Println("Successfully seeded roles.")

	if err := SeedPermissions(db); err != nil {
		log.Fatalf("Error seeding permissions: %v", err)
	}
	log.Println("Successfully seeded permissions.")

	if err := SeedRolePermissions(db); err != nil {
		log.Fatalf("Error seeding role_permissions: %v", err)
	}
	log.Println("Successfully seeded role_permissions.")
}

package main

import (
	"fmt"

	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/container"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/database/seeders"
)

func main() {
	c := container.InitializeContainer()

	fmt.Println("Starting database seeding...")
	seeders.SeedDatabase(c.DatabaseSession)
	fmt.Println("Database seeding completed successfully.")
}

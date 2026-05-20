package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/logger"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/permission"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/role"
)

type PostgresRoleRepository struct {
	session *sql.DB
	logger  logger.Logger
}

func NewPostgresRoleRepository(session *sql.DB, logger logger.Logger) role.Repository {
	return &PostgresRoleRepository{
		session: session,
		logger:  logger,
	}
}

func (r *PostgresRoleRepository) GetById(id ulid.ULID) (*role.Role, error) {
	r.logger.Debug("fetching role by id", "id", id.String())

	var ro role.Role
	var i string
	q := `SELECT 
			id,
			name,
			created_at,
			modified_at
		FROM roles 
		WHERE id = $1;`

	row := r.session.QueryRow(q, id.String())
	if err := row.Scan(&i, &ro.Name, &ro.CreatedAt, &ro.ModifiedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("role not found", "id", id.String())

			return nil, nil
		}

		r.logger.Error("failed to fetch role by id", "id", id.String(), "error", err)

		return nil, fmt.Errorf("error scanning Role by ID: %s", err)
	}

	ro.Id = ulid.MustParse(i)

	return &ro, nil
}

func (r *PostgresRoleRepository) GetByName(name string) (*role.Role, error) {
	r.logger.Debug("fetching role by name", "name", name)

	var ro role.Role
	var id string
	q := `SELECT 
			id,
			name,
			created_at,
			modified_at
		FROM roles 
		WHERE name = $1;`

	row := r.session.QueryRow(q, name)
	if err := row.Scan(&id, &ro.Name, &ro.CreatedAt, &ro.ModifiedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("role not found", "name", name)

			return nil, nil
		}

		r.logger.Error("failed to fetch role by name", "name", name, "error", err)

		return nil, fmt.Errorf("error scanning Role by name: %s", err)
	}

	ro.Id = ulid.MustParse(id)

	permissionsQuery := `SELECT p.id, p.name, p.description, p.created_at, p.modified_at FROM permissions p
						 INNER JOIN role_permissions rp ON p.id = rp.permission_id
						 WHERE rp.role_id = $1;`

	rows, err := r.session.Query(permissionsQuery, ro.Id.String())
	if err != nil {
		r.logger.Error("failed to fetch permissions for role", "id", ro.Id.String(), "name", ro.Name, "error", err)

		return nil, fmt.Errorf("error fetching permissions for Role: %s", err)
	}
	defer rows.Close()

	var permissions []permission.Permission
	for rows.Next() {
		var p permission.Permission
		var pid string
		if err := rows.Scan(&pid, &p.Name, &p.Description, &p.CreatedAt, &p.ModifiedAt); err != nil {
			r.logger.Error("failed to scan permission for role", "id", ro.Id.String(), "name", ro.Name, "error", err)

			return nil, fmt.Errorf("error scanning permission: %s", err)
		}
		p.Id = ulid.MustParse(pid)
		permissions = append(permissions, p)
	}

	ro.Permissions = permissions

	return &ro, nil
}

func (r *PostgresRoleRepository) GetAll() (*[]role.Role, error) {
	r.logger.Debug("fetching all roles")

	var roles []role.Role
	q := `SELECT 
			id,
			name,
			created_at,
			modified_at
		FROM roles;`

	rows, err := r.session.Query(q)
	if err != nil {
		r.logger.Error("failed to fetch all roles", "error", err)

		return nil, fmt.Errorf("error fetching all Roles: %s", err)
	}

	for rows.Next() {
		var id string
		var ro role.Role

		if err := rows.Scan(&id, &ro.Name, &ro.CreatedAt, &ro.ModifiedAt); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}

			r.logger.Error("failed to scan role during fetch all", "error", err)

			return nil, fmt.Errorf("error scanning all Roles: %s", err)
		}

		ro.Id = ulid.MustParse(id)
		roles = append(roles, ro)
	}

	return &roles, nil
}

func (r *PostgresRoleRepository) Add(role *role.Role) error {
	r.logger.Info("adding role", "id", role.Id.String(), "name", role.Name)

	q := `INSERT INTO roles(id, name, created_at) VALUES($1, $2, $3)`

	if _, err := r.session.Exec(q, role.Id.String(), role.Name, role.CreatedAt); err != nil {
		r.logger.Error("failed to add role", "id", role.Id.String(), "name", role.Name, "error", err)

		return fmt.Errorf("failed to persist Role to database: %v", err)
	}

	return nil
}

func (r *PostgresRoleRepository) Remove(role *role.Role) error {
	r.logger.Info("removing role", "id", role.Id.String(), "name", role.Name)

	q := `DELETE FROM roles WHERE id = $1`

	if _, err := r.session.Exec(q, role.Id.String()); err != nil {
		r.logger.Error("failed to remove role", "id", role.Id.String(), "name", role.Name, "error", err)

		return fmt.Errorf("failed to delete Role from database: %s", role)
	}

	return nil
}

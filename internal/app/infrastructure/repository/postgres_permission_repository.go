package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/logger"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/permission"
)

type PostgresPermissionRepository struct {
	session *sql.DB
	logger  logger.Logger
}

func NewPostgresPermissionRepository(session *sql.DB, logger logger.Logger) permission.Repository {
	return &PostgresPermissionRepository{
		session: session,
		logger:  logger,
	}
}

func (r *PostgresPermissionRepository) GetById(id ulid.ULID) (*permission.Permission, error) {
	r.logger.Debug("fetching permission by id", "id", id.String())

	var p permission.Permission
	var i string
	q := `SELECT 
			id,
			name,
			description,
			created_at,
			modified_at
		FROM permissions 
		WHERE id = $1;`

	row := r.session.QueryRow(q, id.String())
	if err := row.Scan(&i, &p.Name, &p.Description, &p.CreatedAt, &p.ModifiedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("permission not found", "id", id.String())

			return nil, nil
		}

		r.logger.Error("failed to fetch permission by id", "id", id.String(), "error", err)

		return nil, fmt.Errorf("error scanning Permission by ID: %s", err)
	}

	p.Id = ulid.MustParse(i)

	return &p, nil
}

func (r *PostgresPermissionRepository) GetByName(name string) (*permission.Permission, error) {
	r.logger.Debug("fetching permission by name", "name", name)

	var p permission.Permission
	var i string
	q := `SELECT 
			id,
			name,
			description,
			created_at,
			modified_at
		FROM permissions 
		WHERE name = $1;`

	row := r.session.QueryRow(q, name)
	if err := row.Scan(&i, &p.Name, &p.Description, &p.CreatedAt, &p.ModifiedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("permission not found", "name", name)

			return nil, nil
		}

		r.logger.Error("failed to fetch permission by name", "name", name, "error", err)

		return nil, fmt.Errorf("error scanning Permission by name: %s", err)
	}

	p.Id = ulid.MustParse(i)

	return &p, nil
}

func (r *PostgresPermissionRepository) GetAll() (*[]permission.Permission, error) {
	r.logger.Debug("fetching all permissions")

	var p permission.Permission
	var permissions []permission.Permission
	var id string
	q := `SELECT 
			id,
			name,
			description,
			created_at,
			modified_at
		FROM permissions;`

	rows, err := r.session.Query(q)
	if err != nil {
		r.logger.Error("failed to fetch all permissions", "error", err)

		return nil, fmt.Errorf("error fetching all Permissions: %s", err)
	}

	for rows.Next() {
		if err := rows.Scan(&id, &p.Name, &p.Description, &p.CreatedAt, &p.ModifiedAt); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, nil
			}

			r.logger.Error("failed to scan permission during fetch all", "error", err)

			return nil, fmt.Errorf("error scanning all Permissions: %s", err)
		}

		p.Id = ulid.MustParse(id)
		permissions = append(permissions, p)
	}

	return &permissions, nil
}

func (r *PostgresPermissionRepository) Add(permission *permission.Permission) error {
	r.logger.Info("adding permission", "id", permission.Id.String(), "name", permission.Name)

	q := `INSERT INTO permissions(id, name, description, created_at) VALUES($1, $2, $3, $4);`

	if _, err := r.session.Exec(q, permission.Id.String(), permission.Name, permission.Description, permission.CreatedAt); err != nil {
		r.logger.Error("failed to add permission", "id", permission.Id.String(), "name", permission.Name, "error", err)

		return fmt.Errorf("failed to persist Permission to database: %v", err)
	}

	return nil
}

func (r *PostgresPermissionRepository) Remove(permission *permission.Permission) error {
	r.logger.Info("removing permission", "id", permission.Id.String(), "name", permission.Name)

	q := `DELETE FROM permissions WHERE id = $1;`

	if _, err := r.session.Exec(q, permission.Id.String()); err != nil {
		r.logger.Error("failed to remove permission", "id", permission.Id.String(), "name", permission.Name, "error", err)

		return fmt.Errorf("failed to delete Permission from database: %s", permission)
	}

	return nil
}

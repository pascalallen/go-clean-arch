package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/logger"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/pagination"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/permission"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/role"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/user"
)

type PostgresUserRepository struct {
	session *sql.DB
	logger  logger.Logger
}

func NewPostgresUserRepository(session *sql.DB, logger logger.Logger) user.Repository {
	return &PostgresUserRepository{
		session: session,
		logger:  logger,
	}
}

func (r *PostgresUserRepository) GetById(id ulid.ULID) (*user.User, error) {
	r.logger.Debug("fetching user by id", "id", id.String())

	var u user.User
	var i string
	q := `SELECT id, first_name, last_name, email_address, password_hash, created_at, modified_at
		FROM users
		WHERE id = $1`

	row := r.session.QueryRow(q, id.String())
	if err := row.Scan(&i, &u.FirstName, &u.LastName, &u.EmailAddress, &u.PasswordHash, &u.CreatedAt, &u.ModifiedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("user not found", "id", id.String())
			return nil, nil
		}
		r.logger.Error("failed to fetch user by id", "id", id.String(), "error", err)
		return nil, fmt.Errorf("error scanning User by ID: %s", err)
	}

	u.Id = ulid.MustParse(i)
	if err := r.loadRoles(&u); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *PostgresUserRepository) GetByEmailAddress(emailAddress string) (*user.User, error) {
	r.logger.Debug("fetching user by email address", "emailAddress", emailAddress)

	var u user.User
	var id string
	q := `SELECT id, first_name, last_name, email_address, password_hash, created_at, modified_at
		FROM users
		WHERE email_address = $1`

	row := r.session.QueryRow(q, emailAddress)
	if err := row.Scan(&id, &u.FirstName, &u.LastName, &u.EmailAddress, &u.PasswordHash, &u.CreatedAt, &u.ModifiedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			r.logger.Debug("user not found", "emailAddress", emailAddress)
			return nil, nil
		}
		r.logger.Error("failed to fetch user by email address", "emailAddress", emailAddress, "error", err)
		return nil, fmt.Errorf("error scanning User by email address: %s", err)
	}

	u.Id = ulid.MustParse(id)
	if err := r.loadRoles(&u); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *PostgresUserRepository) GetAll(pageParams pagination.PageParams) (*pagination.Collection[user.User], error) {
	r.logger.Debug("fetching all users")

	var totalCount int
	if err := r.session.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&totalCount); err != nil {
		r.logger.Error("failed to count users", "error", err)
		return nil, fmt.Errorf("error counting Users: %v", err)
	}

	q := `SELECT id, first_name, last_name, email_address, password_hash, created_at, modified_at
		FROM users
		ORDER BY last_name ASC, first_name ASC
		LIMIT $1 OFFSET $2`

	rows, err := r.session.Query(q, pageParams.Limit, pageParams.Offset())
	if err != nil {
		r.logger.Error("failed to fetch users", "error", err)
		return nil, fmt.Errorf("error fetching Users: %s", err)
	}
	defer rows.Close()

	var users []user.User
	for rows.Next() {
		var u user.User
		var id string
		if err := rows.Scan(&id, &u.FirstName, &u.LastName, &u.EmailAddress, &u.PasswordHash, &u.CreatedAt, &u.ModifiedAt); err != nil {
			r.logger.Error("failed to scan user", "error", err)
			return nil, fmt.Errorf("error scanning User: %s", err)
		}
		u.Id = ulid.MustParse(id)
		if err := r.loadRoles(&u); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("error during rows iteration", "error", err)
		return nil, fmt.Errorf("error during rows iteration: %v", err)
	}

	return &pagination.Collection[user.User]{
		Items:      users,
		TotalCount: totalCount,
	}, nil
}

func (r *PostgresUserRepository) Add(u *user.User) error {
	r.logger.Info("adding user", "id", u.Id.String(), "emailAddress", u.EmailAddress)

	tx, err := r.session.Begin()
	if err != nil {
		r.logger.Error("failed to begin transaction", "error", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	q := `INSERT INTO users(id, first_name, last_name, email_address, password_hash, created_at) VALUES($1, $2, $3, $4, $5, $6)`
	if _, err := tx.Exec(q, u.Id.String(), u.FirstName, u.LastName, u.EmailAddress, u.PasswordHash, u.CreatedAt); err != nil {
		_ = tx.Rollback()
		r.logger.Error("failed to add user", "id", u.Id.String(), "error", err)
		return fmt.Errorf("failed to persist User to database: %v", err)
	}

	if len(u.Roles) > 0 {
		rq := `INSERT INTO user_roles(user_id, role_id) VALUES($1, $2)`
		for _, role := range u.Roles {
			if _, err := tx.Exec(rq, u.Id.String(), role.Id.String()); err != nil {
				_ = tx.Rollback()
				r.logger.Error("failed to add role for user", "userId", u.Id.String(), "roleId", role.Id.String(), "error", err)
				return fmt.Errorf("failed to persist User Role to database: %v", err)
			}
		}
	}

	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		r.logger.Error("failed to commit transaction", "error", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) Remove(u *user.User) error {
	r.logger.Info("removing user", "id", u.Id.String())

	tx, err := r.session.Begin()
	if err != nil {
		r.logger.Error("failed to begin transaction", "error", err)
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	if _, err := tx.Exec(`DELETE FROM user_roles WHERE user_id = $1`, u.Id.String()); err != nil {
		_ = tx.Rollback()
		r.logger.Error("failed to remove user roles", "id", u.Id.String(), "error", err)
		return fmt.Errorf("failed to remove user roles: %w", err)
	}

	if _, err := tx.Exec(`DELETE FROM users WHERE id = $1`, u.Id.String()); err != nil {
		_ = tx.Rollback()
		r.logger.Error("failed to remove user", "id", u.Id.String(), "error", err)
		return fmt.Errorf("failed to remove User from database: %w", err)
	}

	if err := tx.Commit(); err != nil {
		r.logger.Error("failed to commit transaction", "error", err)
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *PostgresUserRepository) Save(u *user.User) error {
	r.logger.Info("saving user", "id", u.Id.String())

	q := `UPDATE users SET first_name = $1, last_name = $2, email_address = $3, password_hash = $4, modified_at = $5 WHERE id = $6`
	res, err := r.session.Exec(q, u.FirstName, u.LastName, u.EmailAddress, u.PasswordHash, u.ModifiedAt, u.Id.String())
	if err != nil {
		r.logger.Error("failed to save user", "id", u.Id.String(), "error", err)
		return fmt.Errorf("failed to update User in database: %v", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		r.logger.Error("failed to verify save operation", "id", u.Id.String(), "error", err)
		return fmt.Errorf("failed to verify update operation: %v", err)
	}

	if rowsAffected == 0 {
		r.logger.Warn("save operation affected no rows", "id", u.Id.String())
		return fmt.Errorf("no User found with id: %s", u.Id.String())
	}

	return nil
}

func (r *PostgresUserRepository) loadRoles(u *user.User) error {
	q := `SELECT r.id, r.name, r.created_at, r.modified_at
		FROM roles r
		JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1`

	rows, err := r.session.Query(q, u.Id.String())
	if err != nil {
		r.logger.Error("failed to fetch roles for user", "userId", u.Id.String(), "error", err)
		return fmt.Errorf("error fetching roles for user: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var ro role.Role
		var rid string
		if err := rows.Scan(&rid, &ro.Name, &ro.CreatedAt, &ro.ModifiedAt); err != nil {
			r.logger.Error("failed to scan role for user", "userId", u.Id.String(), "error", err)
			return fmt.Errorf("error scanning role: %w", err)
		}
		ro.Id = ulid.MustParse(rid)

		pq := `SELECT p.id, p.name, p.description, p.created_at, p.modified_at
			FROM permissions p
			JOIN role_permissions rp ON rp.permission_id = p.id
			WHERE rp.role_id = $1`

		prows, err := r.session.Query(pq, rid)
		if err != nil {
			r.logger.Error("failed to fetch permissions for role", "roleId", rid, "error", err)
			return fmt.Errorf("error fetching permissions for role: %w", err)
		}

		var permissions []permission.Permission
		for prows.Next() {
			var p permission.Permission
			var pid string
			if err := prows.Scan(&pid, &p.Name, &p.Description, &p.CreatedAt, &p.ModifiedAt); err != nil {
				r.logger.Error("failed to scan permission for role", "roleId", rid, "error", err)
				prows.Close()
				return fmt.Errorf("error scanning permission: %w", err)
			}
			p.Id = ulid.MustParse(pid)
			permissions = append(permissions, p)
		}
		prows.Close()
		ro.Permissions = permissions

		u.Roles = append(u.Roles, ro)
	}

	return nil
}

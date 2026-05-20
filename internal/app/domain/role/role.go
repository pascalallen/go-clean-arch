package role

import (
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/permission"
)

type Role struct {
	Id          ulid.ULID               `json:"id"`
	Name        string                  `json:"name"`
	Permissions []permission.Permission `json:"permissions"`
	CreatedAt   time.Time               `json:"created_at"`
	ModifiedAt  *time.Time              `json:"modified_at,omitempty"`
}

func Define(id ulid.ULID, name string) *Role {
	createdAt := time.Now()

	return &Role{
		Id:        id,
		Name:      name,
		CreatedAt: createdAt,
	}
}

func (r *Role) UpdateName(name string) {
	r.Name = name
	now := time.Now()
	r.ModifiedAt = &now
}

func (r *Role) AddPermission(permission permission.Permission) {
	for _, p := range r.Permissions {
		if p.Id == permission.Id {
			return
		}
	}

	r.Permissions = append(r.Permissions, permission)
	now := time.Now()
	r.ModifiedAt = &now
}

func (r *Role) RemovePermission(permission permission.Permission) {
	for i, p := range r.Permissions {
		if p.Id == permission.Id {
			r.Permissions[i] = r.Permissions[len(r.Permissions)-1]
		}
	}

	r.Permissions = r.Permissions[:len(r.Permissions)-1]
}

func (r *Role) HasPermission(name string) bool {
	for _, p := range r.Permissions {
		if p.Name == name {
			return true
		}
	}

	return false
}

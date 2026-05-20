package permission

import "github.com/oklog/ulid/v2"

type Repository interface {
	GetById(id ulid.ULID) (*Permission, error)
	GetByName(name string) (*Permission, error)
	GetAll() (*[]Permission, error)
	Add(permission *Permission) error
	Remove(permission *Permission) error
}

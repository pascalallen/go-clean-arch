package role

import "github.com/oklog/ulid/v2"

type Repository interface {
	GetById(id ulid.ULID) (*Role, error)
	GetByName(name string) (*Role, error)
	GetAll() (*[]Role, error)
	Add(role *Role) error
	Remove(role *Role) error
}

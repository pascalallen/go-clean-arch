package user

import (
	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/pagination"
)

type Repository interface {
	GetById(id ulid.ULID) (*User, error)
	GetByEmailAddress(emailAddress string) (*User, error)
	GetAll(pageParams pagination.PageParams) (*pagination.Collection[User], error)
	Add(user *User) error
	Remove(user *User) error
	Save(user *User) error
}

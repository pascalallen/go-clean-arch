package permission

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Permission struct {
	Id          ulid.ULID  `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
	ModifiedAt  *time.Time `json:"modified_at,omitempty"`
}

func Define(id ulid.ULID, name string, description string) *Permission {
	createdAt := time.Now()

	return &Permission{
		Id:          id,
		Name:        name,
		Description: description,
		CreatedAt:   createdAt,
	}
}

func (p *Permission) UpdateName(name string) {
	p.Name = name
	now := time.Now()
	p.ModifiedAt = &now
}

func (p *Permission) UpdateDescription(description string) {
	p.Description = description
	now := time.Now()
	p.ModifiedAt = &now
}

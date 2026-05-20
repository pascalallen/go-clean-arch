package event

import "github.com/oklog/ulid/v2"

type UserRegistered struct {
	Id           ulid.ULID `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	EmailAddress string    `json:"email_address"`
}

func (e UserRegistered) EventName() string {
	return "UserRegistered"
}

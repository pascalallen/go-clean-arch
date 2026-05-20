package command

import "github.com/oklog/ulid/v2"

type RegisterUser struct {
	Id           ulid.ULID `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName     string    `json:"last_name"`
	EmailAddress string    `json:"email_address"`
}

func (c RegisterUser) CommandName() string {
	return "RegisterUser"
}

type DeleteUser struct {
	UserId ulid.ULID `json:"user_id"`
}

func (c DeleteUser) CommandName() string {
	return "DeleteUser"
}

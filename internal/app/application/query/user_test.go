package query

import (
	"testing"

	"github.com/oklog/ulid/v2"
)

func TestThatQueryNameReturnsExpectedValueGetUserById(t *testing.T) {
	qry := GetUserById{
		Id: ulid.Make(),
	}

	if qry.QueryName() != "GetUserById" {
		t.Fatal("test assertion failed for GetUserById.QueryName()")
	}
}

func TestThatQueryNameReturnsExpectedValueGetUserByEmailAddress(t *testing.T) {
	qry := GetUserByEmailAddress{
		EmailAddress: "foo@bar.com",
	}

	if qry.QueryName() != "GetUserByEmailAddress" {
		t.Fatal("test assertion failed for GetUserByEmailAddress.QueryName()")
	}
}

func TestThatQueryNameReturnsExpectedValueListUsers(t *testing.T) {
	qry := ListUsers{
		Page:  1,
		Limit: 10,
	}

	if qry.QueryName() != "ListUsers" {
		t.Fatal("test assertion failed for ListUsers.QueryName()")
	}
}

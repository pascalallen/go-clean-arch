package query_handler

import (
	"fmt"

	"github.com/pascalallen/go-clean-arch/internal/app/application/query"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/logger"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/pagination"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/user"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/messaging"
)

type ListUsersHandler struct {
	Logger         logger.Logger
	UserRepository user.Repository
}

func (h ListUsersHandler) Handle(qry messaging.Query) (any, error) {
	q, ok := qry.(query.ListUsers)
	if !ok {
		h.Logger.Error("invalid query type passed to ListUsersHandler", "query", qry)
		return nil, fmt.Errorf("invalid query type passed to ListUsersHandler: %v", qry)
	}

	h.Logger.Info("handling ListUsers query")

	pageParams := pagination.PageParams{
		Page:  q.Page,
		Limit: q.Limit,
	}

	u, err := h.UserRepository.GetAll(pageParams)
	if err != nil {
		h.Logger.Error("error attempting to list users from database", "error", err)
		return nil, fmt.Errorf("error attempting to list users from database: %s", err)
	}

	h.Logger.Info("successfully handled ListUsers query")

	return u, nil
}

type GetUserByIdHandler struct {
	Logger         logger.Logger
	UserRepository user.Repository
}

func (h GetUserByIdHandler) Handle(qry messaging.Query) (any, error) {
	q, ok := qry.(query.GetUserById)
	if !ok {
		h.Logger.Error("invalid query type passed to GetUserByIdHandler", "query", qry)
		return nil, fmt.Errorf("invalid query type passed to GetUserByIdHandler: %v", qry)
	}

	h.Logger.Info("handling GetUserById query", "user_id", q.Id.String())

	u, err := h.UserRepository.GetById(q.Id)
	if err != nil {
		h.Logger.Error("error attempting to retrieve user from database", "error", err, "user_id", q.Id.String())
		return nil, fmt.Errorf("error attempting to retrieve user from database: %s", err)
	}

	h.Logger.Info("successfully handled GetUserById query", "user_id", q.Id.String())

	return u, nil
}

type GetUserByEmailAddressHandler struct {
	Logger         logger.Logger
	UserRepository user.Repository
}

func (h GetUserByEmailAddressHandler) Handle(qry messaging.Query) (any, error) {
	q, ok := qry.(query.GetUserByEmailAddress)
	if !ok {
		h.Logger.Error("invalid query type passed to GetUserByEmailAddressHandler", "query", qry)
		return nil, fmt.Errorf("invalid query type passed to GetUserByEmailAddressHandler: %v", qry)
	}

	h.Logger.Info("handling GetUserByEmailAddress query", "email_address", q.EmailAddress)

	u, err := h.UserRepository.GetByEmailAddress(q.EmailAddress)
	if err != nil {
		h.Logger.Error("error attempting to retrieve user from database", "error", err, "email_address", q.EmailAddress)
		return nil, fmt.Errorf("error attempting to retrieve user from database: %s", err)
	}

	h.Logger.Info("successfully handled GetUserByEmailAddress query", "email_address", q.EmailAddress)

	return u, nil
}

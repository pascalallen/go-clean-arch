package auth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/go-clean-arch/internal/app/application/command"
	"github.com/pascalallen/go-clean-arch/internal/app/application/http/responder"
	"github.com/pascalallen/go-clean-arch/internal/app/application/query"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/user"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/messaging"
)

type RegisterRequestPayload struct {
	FirstName    string `form:"first_name" json:"first_name" binding:"required,max=100"`
	LastName     string `form:"last_name" json:"last_name" binding:"required,max=100"`
	EmailAddress string `form:"email_address" json:"email_address" binding:"required,max=100,email"`
}

func HandleRegisterUser(queryBus messaging.QueryBus, commandBus messaging.CommandBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request RegisterRequestPayload

		if err := c.ShouldBind(&request); err != nil {
			responder.BadRequestResponse(c, fmt.Errorf("request validation error: %s", err.Error()))
			return
		}

		request.FirstName = strings.TrimSpace(request.FirstName)
		request.LastName = strings.TrimSpace(request.LastName)
		request.EmailAddress = strings.ToLower(strings.TrimSpace(request.EmailAddress))

		result, err := queryBus.Fetch(query.GetUserByEmailAddress{EmailAddress: request.EmailAddress})
		if err != nil {
			responder.InternalServerErrorResponse(c, err)
			return
		}

		if u, _ := result.(*user.User); u != nil {
			responder.UnprocessableEntityResponse(c, errors.New("an account with that email address already exists"))
			return
		}

		cmd := &command.RegisterUser{
			Id:           ulid.Make(),
			FirstName:    request.FirstName,
			LastName:     request.LastName,
			EmailAddress: request.EmailAddress,
		}

		if err := commandBus.Execute(cmd); err != nil {
			responder.InternalServerErrorResponse(c, err)
			return
		}

		responder.AcceptedResponse(c)
	}
}

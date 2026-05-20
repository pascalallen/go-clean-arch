package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/pascalallen/go-clean-arch/internal/app/application/http/responder"
	"github.com/pascalallen/go-clean-arch/internal/app/application/query"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/user"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/messaging"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/service"
)

type LoginRequestPayload struct {
	EmailAddress string `form:"email_address" json:"email_address" binding:"required,max=100,email"`
	Password     string `form:"password" json:"password" binding:"required"`
}

func HandleLoginUser(queryBus messaging.QueryBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		var request LoginRequestPayload

		if err := c.ShouldBind(&request); err != nil {
			responder.BadRequestResponse(c, fmt.Errorf("request validation error: %s", err.Error()))
			return
		}

		result, err := queryBus.Fetch(query.GetUserByEmailAddress{EmailAddress: request.EmailAddress})
		u, ok := result.(*user.User)
		if u == nil || err != nil || !ok {
			responder.UnauthorizedResponse(c, errors.New("invalid credentials"))
			return
		}

		if !u.PasswordHash.Compare(request.Password) {
			responder.UnauthorizedResponse(c, errors.New("invalid credentials"))
			return
		}

		userClaims := service.UserClaims{
			Id:    u.Id.String(),
			First: u.FirstName,
			Last:  u.LastName,
			StandardClaims: jwt.StandardClaims{
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
			},
		}

		signedAccessToken, err := service.NewAccessToken(userClaims)
		if err != nil {
			responder.InternalServerErrorResponse(c, errors.New("error creating access token"))
			return
		}

		signedRefreshToken, err := service.NewRefreshToken(jwt.StandardClaims{
			Subject:   u.Id.String(),
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
		})
		if err != nil {
			responder.InternalServerErrorResponse(c, errors.New("error creating refresh token"))
			return
		}

		var roles []string
		for _, r := range u.Roles {
			roles = append(roles, r.Name)
		}

		var permissions []string
		for _, p := range u.Permissions() {
			permissions = append(permissions, p.Name)
		}

		userData := UserData{
			Id:           u.Id.String(),
			FirstName:    u.FirstName,
			LastName:     u.LastName,
			EmailAddress: u.EmailAddress,
			CreatedAt:    u.CreatedAt.String(),
		}

		if u.ModifiedAt != nil {
			userData.ModifiedAt = u.ModifiedAt.String()
		}

		responder.CreatedResponse(c, &TokenResponsePayload{
			AccessToken:  signedAccessToken,
			RefreshToken: signedRefreshToken,
			User:         userData,
			Roles:        roles,
			Permissions:  permissions,
		})
	}
}

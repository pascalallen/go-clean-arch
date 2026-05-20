package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
	"github.com/pascalallen/go-clean-arch/internal/app/application/http/responder"
	"github.com/pascalallen/go-clean-arch/internal/app/application/query"
	"github.com/pascalallen/go-clean-arch/internal/app/domain/user"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/messaging"
	"github.com/pascalallen/go-clean-arch/internal/app/infrastructure/service"
)

func AuthRequired(queryBus messaging.QueryBus) gin.HandlerFunc {
	return func(c *gin.Context) {
		var accessToken string
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				responder.BadRequestResponse(c, errors.New("invalid authorization header"))
				return
			}
			accessToken = parts[1]
		}

		if accessToken == "" {
			accessToken = c.Query("token")
		}

		if accessToken == "" {
			responder.BadRequestResponse(c, errors.New("missing authorization header or token query parameter"))
			return
		}

		userClaims := service.ParseAccessToken(accessToken)
		if userClaims == nil {
			responder.UnauthorizedResponse(c, errors.New("invalid or expired token"))
			return
		}

		q := query.GetUserById{Id: ulid.MustParse(userClaims.Id)}
		result, err := queryBus.Fetch(q)
		u, ok := result.(*user.User)
		if u == nil || err != nil || !ok {
			responder.UnauthorizedResponse(c, errors.New("invalid credentials"))
			return
		}

		c.Set("userId", u.Id)

		// If a request-scoped logger exists, enrich it with the authenticated user ID.
		if l, ok := Get(c); ok {
			Set(c, l.With("user_id", u.Id.String()))
		}

		c.Next()
	}
}

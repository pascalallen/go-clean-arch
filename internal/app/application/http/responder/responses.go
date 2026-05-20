package responder

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSendErrorResponse[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code,omitempty"`
	Data    T      `json:"data,omitempty"`
}

type JSendFailResponse[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data"`
}

type JSendSuccessResponse[T any] struct {
	Status string `json:"status"`
	Data   T      `json:"data,omitempty"`
}

func InternalServerErrorResponse(c *gin.Context, error error) {
	c.JSON(
		http.StatusInternalServerError,
		JSendErrorResponse[string]{
			Status:  "error",
			Message: error.Error(),
		},
	)

	return
}

func UnprocessableEntityResponse(c *gin.Context, error error) {
	c.JSON(
		http.StatusUnprocessableEntity,
		JSendErrorResponse[string]{
			Status:  "error",
			Message: error.Error(),
		},
	)

	return
}

func UnauthorizedResponse(c *gin.Context, error error) {
	c.JSON(
		http.StatusUnauthorized,
		JSendFailResponse[string]{
			Status: "fail",
			Data:   error.Error(),
		},
	)

	return
}

func BadRequestResponse(c *gin.Context, error error) {
	c.JSON(
		http.StatusBadRequest,
		JSendFailResponse[string]{
			Status: "fail",
			Data:   error.Error(),
		},
	)

	return
}

func CreatedResponse[T interface{}](c *gin.Context, i *T) {
	c.JSON(
		http.StatusCreated,
		JSendSuccessResponse[T]{
			Status: "success",
			Data:   *i,
		},
	)

	return
}

func OkResponse[T interface{}](c *gin.Context, i *T) {
	c.JSON(
		http.StatusOK,
		JSendSuccessResponse[T]{
			Status: "success",
			Data:   *i,
		},
	)

	return
}

func AcceptedResponse(c *gin.Context) {
	c.JSON(http.StatusAccepted, JSendSuccessResponse[string]{Status: "success"})
	return
}

func NotFoundResponse(c *gin.Context, error error) {
	c.JSON(
		http.StatusNotFound,
		JSendErrorResponse[string]{
			Status:  "error",
			Message: error.Error(),
		},
	)

	return
}

func ForbiddenResponse(c *gin.Context, error error) {
	c.JSON(
		http.StatusForbidden,
		JSendErrorResponse[string]{
			Status:  "error",
			Message: error.Error(),
		},
	)

	return
}

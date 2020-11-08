package delivery

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bekaza/go-clean/domain"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

// Handler - user handler http
type Handler struct {
	us domain.UserService
}

// NewHandler - new handler http user
func NewHandler(e *echo.Group, us domain.UserService) *Handler {
	h := &Handler{us}
	e.POST("/user/login", h.Login)
	e.POST("/user/register", h.Register)
	return h
}

func isRequestValid(m *domain.User) (bool, error) {
	validate := validator.New()
	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Register ...
func (h *Handler) Register(c echo.Context) (err error) {
	var user domain.User
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = c.Bind(&user)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	var ok bool
	if ok, err = isRequestValid(&user); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err = h.us.Register(ctx, user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, echo.Map{"message": err.Error()})
	}
	return c.JSON(http.StatusCreated, user)
}

// Login ...
func (h *Handler) Login(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	type userLoginRequest struct {
		Username string `json:"username" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	u := new(userLoginRequest)
	if err := c.Bind(u); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": err.Error()})
	}
	validate := validator.New().Struct(u)
	if validate != nil {
		if _, ok := validate.(*validator.InvalidValidationError); ok {
			fmt.Println(validate)
			return nil
		}

		for _, err := range validate.(validator.ValidationErrors) {

			fmt.Println(err.Namespace())
			fmt.Println(err.Field())
			fmt.Println(err.StructNamespace())
			fmt.Println(err.StructField())
			fmt.Println(err.Tag())
			fmt.Println(err.ActualTag())
			fmt.Println(err.Kind())
			fmt.Println(err.Type())
			fmt.Println(err.Value())
			fmt.Println(err.Param())
			fmt.Println()
		}
	}

	return c.JSON(http.StatusOK, echo.Map{"success": true})
}

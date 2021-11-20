package http

import (
	"homework-9/domain"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ResponseError struct {
	Message string `json:"message"`
}

type UserHandler struct {
	userUsecase domain.UserUsecase
}

func NewUserHandler(e *echo.Echo, uuc domain.UserUsecase) {
	handler := &UserHandler{uuc}

	userGroup := e.Group("/app/user")

	userGroup.Use(middleware.BasicAuth(func(username string, password string, c echo.Context) (bool, error) {
		user, err := handler.userUsecase.GetByUsername(username)
		if err != nil {
			log.Printf("BasicAuth: %s\n", err)
			return false, nil
		}

		if user.Username == username && user.Password == password {
			return true, nil
		}

		return false, nil
	}))

	userGroup.GET("/:id", handler.GetUser)
	userGroup.POST("/:id", handler.AddUser)

}

func (h *UserHandler) GetUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Printf("GetUser: %s\n", err)
		return c.JSON(http.StatusBadRequest, ResponseError{"invalid user id"})
	}

	user, err := h.userUsecase.GetById(id)
	if err != nil {
		log.Printf("GetUser: %s\n", err)
		return c.JSON(http.StatusNotFound, ResponseError{"user doesn't exist"})
	}

	return c.JSONPretty(http.StatusOK, user, "	")
}

func (h *UserHandler) AddUser(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		log.Printf("AddUser: %s\n", "invalid user id")
		return c.JSON(http.StatusBadRequest, ResponseError{"invalid user id"})
	}

	user := domain.User{}

	user.ID = id
	username := c.Request().URL.Query().Get("username")
	password := c.Request().URL.Query().Get("password")

	if len(username) == 0 || len(password) == 0 {
		log.Printf("AddUser: %s\n", "invalid query parameters")
		return c.JSON(http.StatusBadRequest, ResponseError{"invalid query parameters"})
	}

	user.Username = username
	user.Password = password

	err = h.userUsecase.Add(user)
	if err != nil {
		log.Printf("AddUser: %s\n", err)
		return c.JSON(http.StatusBadRequest, ResponseError{"user already exists"})
	}

	return c.JSONPretty(http.StatusOK, user, "	")
}

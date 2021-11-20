package http

import (
	"context"
	"fmt"
	"homework-9/domain"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ResponseError struct {
	Message string `json:"message"`
}

type WalletHandler struct {
	walletUsecase domain.WalletUsecase
	userUsecase   domain.UserUsecase
}

type key struct {
}

func NewWalletHandler(e *echo.Echo, wuc domain.WalletUsecase, uuc domain.UserUsecase) {
	handler := &WalletHandler{wuc, uuc}

	walletGroup := e.Group("/app/wallet")

	walletGroup.Use(middleware.BasicAuth(func(username string, password string, c echo.Context) (bool, error) {
		user, err := handler.userUsecase.GetByUsername(username)
		if err != nil {
			log.Printf("BasicAuth: %s\n", err)
			return false, nil
		}

		if user.Username == username && user.Password == password {
			ctx := context.WithValue(c.Request().Context(), key{}, username)
			c.SetRequest(c.Request().WithContext(ctx))

			return true, nil
		}

		return false, nil
	}))

	walletGroup.GET("/:name", handler.GetWallet)
	walletGroup.POST("/:name", handler.AddWallet)
	walletGroup.DELETE("/:name", handler.DeleteWallet)
	walletGroup.OPTIONS("/:name/start", handler.StartMiningWallet)
	walletGroup.OPTIONS("/:name/stop", handler.StopMiningWallet)
}

func (h *WalletHandler) GetWallet(c echo.Context) error {
	name := c.Param("name")

	username, ok := c.Request().Context().Value(key{}).(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ResponseError{"internal server error"})
	}

	user, err := h.userUsecase.GetByUsername(username)
	if err != nil {
		log.Printf("GetWallet: %s\n", err)
		return c.JSON(http.StatusInternalServerError, ResponseError{"internal server error"})
	}

	wallet, err := h.walletUsecase.GetByName(user.ID, name)
	if err != nil {
		log.Printf("GetWallet: %s\n", err)
		return c.JSON(http.StatusNotFound, ResponseError{"wallet doesn't exist"})
	}

	return c.JSON(http.StatusOK, &wallet)
}

func (h *WalletHandler) AddWallet(c echo.Context) error {
	name := c.Param("name")

	username, ok := c.Request().Context().Value(key{}).(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ResponseError{"internal server error"})
	}

	user, err := h.userUsecase.GetByUsername(username)
	if err != nil {
		log.Printf("AddWallet: %s\n", err)
		return c.JSON(http.StatusInternalServerError, ResponseError{"internal server error"})
	}

	err = h.walletUsecase.Add(user.ID, name)
	if err != nil {
		log.Printf("AddWallet: %s\n", err)
		return c.JSON(http.StatusBadRequest, ResponseError{"wallet already exists"})
	}

	return c.JSON(http.StatusCreated, ResponseError{"wallet created"})
}

func (h *WalletHandler) DeleteWallet(c echo.Context) error {
	name := c.Param("name")

	username, ok := c.Request().Context().Value(key{}).(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ResponseError{"internal server error"})
	}

	user, err := h.userUsecase.GetByUsername(username)
	if err != nil {
		log.Printf("DeleteWallet: %s\n", err)
		return c.JSON(http.StatusInternalServerError, ResponseError{"internal server error"})
	}

	err = h.walletUsecase.DeleteByName(user.ID, name)
	if err != nil {
		log.Printf("DeleteWallet: %s\n", err)
		return c.JSON(http.StatusNotFound, ResponseError{"wallet doesn't exist"})
	}

	return c.JSON(http.StatusOK, ResponseError{fmt.Sprintf("wallet '%s' successfully deleted", name)})
}

func (h *WalletHandler) StartMiningWallet(c echo.Context) error {
	name := c.Param("name")

	username, ok := c.Request().Context().Value(key{}).(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ResponseError{"internal server error"})
	}

	user, err := h.userUsecase.GetByUsername(username)
	if err != nil {
		log.Printf("StartMiningWallet: %s\n", err)
		return c.JSON(http.StatusInternalServerError, ResponseError{"internal server error"})
	}

	err = h.walletUsecase.Mine(user.ID, name)
	if err != nil {
		log.Printf("StartMiningWallet: %s\n", err)
		return c.JSON(http.StatusNotFound, ResponseError{"wallet is already mining"})
	}

	return c.JSON(http.StatusOK, ResponseError{fmt.Sprintf("wallet '%s' started mining", name)})
}

func (h *WalletHandler) StopMiningWallet(c echo.Context) error {
	name := c.Param("name")

	username, ok := c.Request().Context().Value(key{}).(string)
	if !ok {
		return c.JSON(http.StatusInternalServerError, ResponseError{"internal server error"})
	}

	user, err := h.userUsecase.GetByUsername(username)
	if err != nil {
		log.Printf("StopMiningWallet: %s\n", err)
		return c.JSON(http.StatusInternalServerError, ResponseError{"internal server error"})
	}

	err = h.walletUsecase.StopMining(user.ID, name)
	if err != nil {
		log.Printf("StopMiningWallet: %s\n", err)
		return c.JSON(http.StatusNotFound, ResponseError{"wallet is not mining"})
	}

	return c.JSON(http.StatusOK, ResponseError{fmt.Sprintf("wallet '%s' stopped mining", name)})
}

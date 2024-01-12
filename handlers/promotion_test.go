//go:build unit

package handlers_test

import (
	"errors"
	"fmt"
	"gotest/handlers"
	"gotest/services"
	"io"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestPromotionCalculateDiscount(t *testing.T) {
	t.Run("applied 100", func(t *testing.T) {
		//Arrange
		amount := 100
		expected := 80

		promoService := services.NewPromotionServiceMock()
		promoService.On("CalculateDiscount", amount).Return(expected, nil)

		promoHandler := handlers.NewPromotionHandler(promoService)

		app := fiber.New()
		app.Get("/calculate", promoHandler.CalculateDiscount)

		req := httptest.NewRequest("GET", fmt.Sprintf("/calculate?amount=%v", amount), nil)

		//Act
		res, _ := app.Test(req)

		//Assert
		if assert.Equal(t, fiber.StatusOK, res.StatusCode) {
			body, _ := io.ReadAll(res.Body)
			assert.Equal(t, strconv.Itoa(expected), string(body))
		}

	})
	t.Run("StatusBadRequest", func(t *testing.T) {
		// Arrange
		amount := "invalid"
		expected := 80

		promoService := services.NewPromotionServiceMock()
		promoService.On("CalculateDiscount", amount).Return(expected, nil)

		promoHandler := handlers.NewPromotionHandler(promoService)

		app := fiber.New()
		app.Get("/calculate", promoHandler.CalculateDiscount)

		req := httptest.NewRequest("GET", fmt.Sprintf("/calculate?amount=%v", amount), nil)

		// Act
		res, err := app.Test(req)

		// Assert
		assert.Nil(t, err)
		assert.Equal(t, fiber.StatusBadRequest, res.StatusCode)

	})

	t.Run("StatusNotFound", func(t *testing.T) {
		// Arrange
		amount := 0
		expectedStatusCode := fiber.StatusNotFound

		promoService := services.NewPromotionServiceMock()
		promoService.On("CalculateDiscount", amount).Return(0, errors.New(""))

		promoHandler := handlers.NewPromotionHandler(promoService)

		app := fiber.New()
		app.Get("/calculate", promoHandler.CalculateDiscount)

		req := httptest.NewRequest("GET", fmt.Sprintf("/calculate?amount=%v", amount), nil)

		// Act
		res, err := app.Test(req)
		// Assert
		assert.Nil(t, err)
		assert.Equal(t, expectedStatusCode, res.StatusCode)

	})
}

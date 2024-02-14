package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	ErrorResponse struct {
		Title  string `json:"title" example:"Error"`
		Status int    `json:"status" example:"500"`
		Detail string `json:"detail" example:"Something went wrong"`
	}
)

func NewHTTPErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		res := ErrorResponse{
			Title:  "ServerError",
			Status: http.StatusInternalServerError,
			Detail: err.Error(),
		}

		// TODO: improver errors
		if he, ok := err.(*echo.HTTPError); ok {
			res.Status = he.Code
			res.Title = he.Message.(string)
		}

		c.JSON(res.Status, res)
	}
}

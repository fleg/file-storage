package server

import (
	"file-storage/internal/errors"
	"fmt"
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

func NewHTTPErrorHandler(extended bool) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		res := ErrorResponse{
			Title:  "InternalServerError",
			Status: http.StatusInternalServerError,
		}

		if extended {
			res.Detail = err.Error()
		} else {
			res.Detail = "Internal server error"
		}

		if he, ok := err.(*echo.HTTPError); ok {
			res.Status = he.Code
			res.Title = he.Message.(string)
		}

		if be, ok := err.(*errors.BaseError); ok {
			res.Status = be.Status
			res.Title = be.Title

			if extended {
				res.Detail = fmt.Sprintf("%s: %s", be.Detail, be.Internal.Error())
			} else {
				res.Detail = be.Detail
			}
		}

		c.JSON(res.Status, res)
	}
}

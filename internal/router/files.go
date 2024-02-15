package router

import (
	"net/http"

	"file-storage/internal/services"

	"github.com/labstack/echo/v4"
)

type (
	FilesController struct {
		filesService *services.FilesService
	}

	UploadResponse struct {
		ID string `json:"id"`
		UploadedAt int64 `json:"uploadedAt"`
		Size uint `json:"size"`
		Mime string `json:"mime"`
		Name string `json:"name"`
	}
)

func (fc *FilesController) Upload(c echo.Context) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return err
	}

	file, err := fc.filesService.Upload(c.Request().Context(), fileHeader)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, UploadResponse{
		ID: file.ID,
		UploadedAt: file.UploadedAt.UTC().UnixMilli(),
		Size: file.Size,
		Mime: file.Mime,
		Name: file.Name,
	})
}

func NewFilesController(fs *services.FilesService) *FilesController {
	return &FilesController{filesService: fs}
}

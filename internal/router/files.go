package router

import (
	"net/http"

	"file-storage/internal/errors"
	"file-storage/internal/services"

	"github.com/labstack/echo/v4"
)

type (
	FilesController struct {
		filesService *services.FilesService
	}

	FilesUploadResponse struct {
		ID         string `json:"id"`
		UploadedAt int64  `json:"uploadedAt"`
		Size       uint   `json:"size"`
		Mime       string `json:"mime"`
		Name       string `json:"name"`
	}

	FilesGetOneResponse struct {
		ID         string `json:"id" example:"a6da224f-a0b6-4803-82a7-268fb98cd8d4" format:"uuid"`
		UploadedAt int64  `json:"uploadedAt" example:"1708339471468" minimum:"0"`
		Size       uint   `json:"size" example:"36354" minimum:"0"`
		Mime       string `json:"mime" example:"image/jpeg" maxLength:"255"`
		Name       string `json:"name" example:"image.jpg" maxLength:"255"`
	}
)

func (fc *FilesController) Upload(c echo.Context) error {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		return errors.NewUnsupportedMediaTypeError("can't parse multipart form", err)
	}

	file, err := fc.filesService.Upload(c.Request().Context(), fileHeader)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, FilesUploadResponse{
		ID:         file.ID,
		UploadedAt: file.UploadedAt.UTC().UnixMilli(),
		Size:       file.Size,
		Mime:       file.Mime,
		Name:       file.Name,
	})
}

func (fc *FilesController) Download(c echo.Context) error {
	id := c.Param("id")
	file, err := fc.filesService.FindOne(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.Attachment(fc.filesService.GetFilePath(file), file.Name)
}

// Files.GetOne godoc
//	@Summary	Get file information by ID
//	@Tags		files
//	@Produce	json
//	@Param		id	path		string	true	"File ID"	Format(uuid)
//	@Success	200	{object}	router.FilesGetOneResponse
//	@Failure	400	{object}	server.ErrorResponse
//	@Failure	404	{object}	server.ErrorResponse
//	@Failure	500	{object}	server.ErrorResponse
//	@Router		/files/{id} [get]
func (fc *FilesController) GetOne(c echo.Context) error {
	id := c.Param("id")
	file, err := fc.filesService.FindOne(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, FilesUploadResponse{
		ID:         file.ID,
		UploadedAt: file.UploadedAt.UTC().UnixMilli(),
		Size:       file.Size,
		Mime:       file.Mime,
		Name:       file.Name,
	})
}

// Files.Unlink godoc
//	@Summary	Remove file by ID
//	@Tags		files
//	@Produce	json
//	@Param		id	path	string	true	"File ID"	Format(uuid)
//	@Success	204
//	@Router		/files/{id} [delete]
func (fc *FilesController) Unlink(c echo.Context) error {
	id := c.Param("id")
	if err := fc.filesService.Unlink(c.Request().Context(), id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func NewFilesController(fs *services.FilesService) *FilesController {
	return &FilesController{filesService: fs}
}

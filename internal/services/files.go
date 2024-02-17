package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"file-storage/internal/config"
	"file-storage/internal/repositories"
)

const DefaultMimeType = "application/octet-stream"

type (
	FilesService struct {
		storagePath string
		maxFileSize uint

		filesRepository *repositories.FilesRepository
	}

	File struct {
		ID         string
		UploadedAt time.Time
		Size       uint
		Mime       string
		Name       string
	}
)

func (fs *FilesService) Upload(ctx context.Context, fileHeader *multipart.FileHeader) (*File, error) {
	if fileHeader.Size > int64(fs.maxFileSize) {
		return nil, fmt.Errorf("file is too large (%d bytes), maximum is %d", fileHeader.Size, fs.maxFileSize)
	}

	mimeType := fileHeader.Header.Get("content-type")
	if len(mimeType) == 0 {
		mimeType = DefaultMimeType
	}

	id, err := fs.filesRepository.Create(ctx, &repositories.CreateFileParams{
		Size: uint(fileHeader.Size),
		Mime: mimeType,
		Name: fileHeader.Filename,
	})
	if err != nil {
		return nil, err
	}

	// NOTE: multipart will dump file to the disk if it doesn't fit in the ram
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	fileName := fmt.Sprintf("%s%s", id, strings.ToLower(filepath.Ext(fileHeader.Filename)))
	dst, err := os.Create(filepath.Join(fs.storagePath, fileName))
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return nil, err
	}

	fileEntity, err := fs.filesRepository.FindOneById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &File{
		ID:         id,
		UploadedAt: fileEntity.UploadedAt,
		Size:       fileEntity.Size,
		Mime:       fileEntity.Mime,
		Name:       fileEntity.Name,
	}, nil
}

func NewFilesService(c *config.Config, fr *repositories.FilesRepository) *FilesService {
	return &FilesService{
		storagePath:     c.Storage.Path,
		maxFileSize:     c.Storage.MaxFileSize,
		filesRepository: fr,
	}
}

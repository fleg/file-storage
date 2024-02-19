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
	"file-storage/internal/errors"
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

func (f *File) getFileName() string {
	return fmt.Sprintf("%s%s", f.ID, strings.ToLower(filepath.Ext(f.Name)))
}

func (fs *FilesService) getFilePath(f *File) string {
	return filepath.Join(fs.storagePath, f.getFileName())
}

func (fs *FilesService) fileFromEntity(fe *repositories.FileEntity) *File {
	return &File{
		ID:         fe.ID,
		UploadedAt: fe.UploadedAt,
		Size:       fe.Size,
		Mime:       fe.Mime,
		Name:       fe.Name,
	}
}

func (fs *FilesService) Upload(ctx context.Context, fileHeader *multipart.FileHeader) (*File, error) {
	if fileHeader.Size > int64(fs.maxFileSize) {
		return nil, errors.NewBadRequestError(
			fmt.Sprintf("file is too large (%d bytes), maximum is %d", fileHeader.Size, fs.maxFileSize),
			nil,
		)
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

	f := &File{
		ID:   id,
		Size: uint(fileHeader.Size),
		Mime: mimeType,
		Name: fileHeader.Filename,
	}

	// NOTE: multipart will dump file to the disk if it doesn't fit in the ram
	src, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	dst, err := os.Create(fs.getFilePath(f))
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

	return fs.fileFromEntity(fileEntity), nil
}

func (fs *FilesService) FindOne(ctx context.Context, id string) (*File, error) {
	fileEntity, err := fs.filesRepository.FindOneById(ctx, id)
	if err != nil {
		if err == repositories.NotFoundError {
			return nil, errors.NewNotFoundError(fmt.Sprintf("file with id %s is not found", id), err)
		}
		return nil, err
	}

	return fs.fileFromEntity(fileEntity), nil
}

func (fs *FilesService) Unlink(ctx context.Context, id string) error {
	fileEntity, err := fs.filesRepository.FindOneById(ctx, id)
	if err != nil {
		if err == repositories.NotFoundError {
			// file already deleted or doesn't exists, do nothing
			return nil
		}
		return err
	}

	f := fs.fileFromEntity(fileEntity)

	if err := fs.filesRepository.RemoveOne(ctx, id); err != nil {
		return err
	}

	if err := os.Remove(fs.getFilePath(f)); err != nil {
		return err
	}

	return nil
}

func NewFilesService(c *config.Config, fr *repositories.FilesRepository) *FilesService {
	return &FilesService{
		storagePath:     c.Storage.Path,
		maxFileSize:     c.Storage.MaxFileSize,
		filesRepository: fr,
	}
}

package repositories

import (
	"context"
	"time"

	"file-storage/internal/postgresql"

	"github.com/jackc/pgx/v5"
)

type (
	FilesRepository struct {
		pg *postgresql.PostgreSQL
	}

	CreateFileParams struct {
		Size uint
		Mime string
		Name string
	}

	// 1 to 1 mapping to sql representation
	FileEntity struct {
		ID         string
		UploadedAt time.Time
		Size       uint
		Mime       string
		Name       string
	}
)

func (fr *FilesRepository) Create(ctx context.Context, p *CreateFileParams) (string, error) {
	var id string

	err := fr.pg.Pool.QueryRow(ctx, "insert into files(size, mime, name) values ($1, $2, $3) returning id", p.Size, p.Mime, p.Name).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (fr *FilesRepository) FindOneById(ctx context.Context, id string) (*FileEntity, error) {
	var f FileEntity

	// TODO: use github.com/georgysavva/scany?
	err := fr.pg.Pool.
		QueryRow(ctx, "select id, uploaded_at, size, mime, name from files where id = $1", id).
		Scan(
			&f.ID,
			&f.UploadedAt,
			&f.Size,
			&f.Mime,
			&f.Name,
		)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, NotFoundError
		}
		return nil, err
	}

	return &f, nil
}

func NewFilesRepository(pg *postgresql.PostgreSQL) *FilesRepository {
	return &FilesRepository{
		pg: pg,
	}
}

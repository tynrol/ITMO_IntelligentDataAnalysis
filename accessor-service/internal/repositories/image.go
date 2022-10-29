package repositories

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/model/domain"
	"log"
)

const tableName = "Images"

func imageColumns() []string {
	return []string{
		"id",
		"width",
		"height",
		"description",
		"url",
		"path",
		"created_at",
	}
}

type Repo struct {
	db *sql.DB

	log *log.Logger
}

func NewRepo(db *sql.DB, logger *log.Logger) *Repo {
	return &Repo{
		db:  db,
		log: logger,
	}
}

func (r *Repo) Create(ctx context.Context, img domain.Image) error {
	const op = "ImageRepository_Create"

	sql, args, err := squirrel.Insert(tableName).
		Columns(imageColumns()...).
		Values(img.Values()...).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return errors.Wrap(err, op)
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (r *Repo) UpdatePathById(ctx context.Context, img domain.Image) (err error) {
	const op = "ImageRepository_UpdatePathById"

	sql, args, err := squirrel.Update(tableName).
		Set("path", img.Path).
		Where(squirrel.Eq{"id": img.ID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return errors.Wrap(err, op)
	}

	_, err = r.db.ExecContext(ctx, sql, args...)
	if err != nil {
		return errors.Wrap(err, op)
	}

	return nil
}

func (r *Repo) GetById(ctx context.Context, imageId string) (image domain.Image, err error) {
	const op = "ImageRepository_GetById"

	sql, args, err := squirrel.Select(imageColumns()...).
		From(tableName).
		Where(squirrel.Eq{"id": imageId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return domain.Image{}, errors.Wrap(err, op)
	}

	rows := r.db.QueryRowContext(ctx, sql, args...)
	err = rows.Scan(image.ScanValues()...)
	if err != nil {
		return domain.Image{}, errors.Wrap(err, op)
	}

	return image, err
}

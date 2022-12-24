package repositories

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/model/domain"
	"log"
)

const imagesTableName = "images"

func imageColumns() []string {
	return []string{
		"id",
		"session_id",
		"type",
		"width",
		"height",
		"description",
		"url",
		"path",
		"created_at",
	}
}

type ImageRepo struct {
	db *sql.DB

	log *log.Logger
}

func NewImageRepo(db *sql.DB, logger *log.Logger) *ImageRepo {
	return &ImageRepo{
		db:  db,
		log: logger,
	}
}

func (r *ImageRepo) Create(ctx context.Context, img domain.Image) error {
	const op = "ImageRepository_Create"

	sql, args, err := squirrel.Insert(imagesTableName).
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

func (r *ImageRepo) UpdatePathById(ctx context.Context, imgId string, path string, imgType string, sessionId string) (err error) {
	const op = "ImageRepository_UpdatePathById"

	sql, args, err := squirrel.Update(imagesTableName).
		Set("path", path).
		Set("type", imgType).
		Set("session_id", sessionId).
		Where(squirrel.Eq{"id": imgId}).
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

func (r *ImageRepo) GetById(ctx context.Context, imageId string) (image domain.Image, err error) {
	const op = "ImageRepository_GetById"

	sql, args, err := squirrel.Select(imageColumns()...).
		From(imagesTableName).
		Where(squirrel.Eq{"id": imageId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return image, errors.Wrap(err, op)
	}

	rows := r.db.QueryRowContext(ctx, sql, args...)
	err = rows.Scan(image.ScanValues()...)
	if rows != nil && err != nil {
		return image, errors.Wrap(err, op)
	}

	return image, err
}

func (r *ImageRepo) GetHoney(ctx context.Context) (image domain.Image, err error) {
	const op = "ImageRepository_GetHoney"

	sql, args, err := squirrel.Select(imageColumns()...).
		From(imagesTableName).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return domain.Image{}, errors.Wrap(err, op)
	}

	rows := r.db.QueryRowContext(ctx, sql, args...)
	err = rows.Scan(image.ScanValues()...)
	if err != nil && rows.Err() != nil {
		return domain.Image{}, errors.Wrap(err, op)
	}

	return image, nil
}

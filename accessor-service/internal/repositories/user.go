package repositories

import (
	"context"
	"database/sql"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/tynrol/ITMO_IntelligentDataAnalysis/accessor-service/internal/model/domain"
	"log"
)

const usersTableName = "users"

func usersColumns() []string {
	return []string{
		"session_id",
		"is_lying",
	}
}

type UserRepo struct {
	db *sql.DB

	log *log.Logger
}

func NewUserRepo(db *sql.DB, logger *log.Logger) *ImageRepo {
	return &ImageRepo{
		db:  db,
		log: logger,
	}
}

func (r *UserRepo) Create(ctx context.Context, user domain.User) error {
	const op = "UsersRepository_Create"

	sql, args, err := squirrel.Insert(usersTableName).
		Columns(usersColumns()...).
		Values(user.Values()...).
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

func (r *UserRepo) GetBySessionId(ctx context.Context, sessionId string) (user domain.User, err error) {
	const op = "UsersRepository_GetById"

	sql, args, err := squirrel.Select(usersColumns()...).
		From(usersTableName).
		Where(squirrel.Eq{"session_id": sessionId}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return domain.User{}, errors.Wrap(err, op)
	}

	rows := r.db.QueryRowContext(ctx, sql, args...)
	err = rows.Scan(user.ScanValues()...)
	if err != nil {
		return domain.User{}, errors.Wrap(err, op)
	}

	return user, err
}

func (r *UserRepo) GetAllLying(ctx context.Context) (user domain.User, err error) {
	const op = "UsersRepository_GetById"

	sql, args, err := squirrel.Select(usersColumns()...).
		From(usersTableName).
		Where(squirrel.Eq{"is_lying": true}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()
	if err != nil {
		return domain.User{}, errors.Wrap(err, op)
	}

	rows := r.db.QueryRowContext(ctx, sql, args...)
	err = rows.Scan(user.ScanValues()...)
	if err != nil {
		return domain.User{}, errors.Wrap(err, op)
	}

	return user, err
}

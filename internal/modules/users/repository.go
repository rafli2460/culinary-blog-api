package users

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/rafli2460/culinary-blog-api/internal/domain"
	"github.com/rafli2460/culinary-blog-api/internal/platform/database"
	"github.com/rafli2460/culinary-blog-api/internal/server"
)

type UserRepository interface {
	Insert(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id int64) error
	UpdateRole(ctx context.Context, id int64, role string) error
	FindByUsername(ctx context.Context, username string) (domain.User, error)
}

type Statement struct {
	insert         *sqlx.NamedStmt
	delete         *sqlx.NamedStmt
	updateRole     *sqlx.NamedStmt
	findByUsername *sqlx.Stmt
}

type Repository struct {
	app  *server.App
	stmt Statement
}

func initRepository(ctx context.Context, app *server.App) UserRepository {
	stmts := Statement{
		insert:         database.PrepareNamed(ctx, app.Ds.WriterDB, insert),
		findByUsername: database.Prepare(ctx, app.Ds.ReaderDB, findByUsername),
		delete:         database.PrepareNamed(ctx, app.Ds.WriterDB, delete),
		updateRole:     database.PrepareNamed(ctx, app.Ds.WriterDB, updateRole),
	}

	r := Repository{
		app:  app,
		stmt: stmts,
	}

	return &r
}

func (r *Repository) Insert(ctx context.Context, user domain.User) error {
	_, err := r.stmt.insert.ExecContext(ctx, user)

	if err != nil {
		r.app.Logger.Error().Stack().
			Str("Context", "Insert new user").
			Err(err).Msg("")
	}

	return err
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	panic("not implemented") // TODO: Implement
}

func (r *Repository) UpdateRole(ctx context.Context, id int64, role string) error {
	panic("not implemented") // TODO: Implement
}

func (r *Repository) FindByUsername(ctx context.Context, username string) (domain.User, error) {
	panic("not implemented") // TODO: Implement
}

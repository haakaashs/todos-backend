package repository

import (
	"context"
	"database/sql"

	"log"

	"github.com/google/uuid"
	"github.com/haakaashs/todos-backend/internal/model"
)

type Repository struct {
	db *sql.DB

	createStmt *sql.Stmt
	getStmt    *sql.Stmt
	listStmt   *sql.Stmt
	updateStmt *sql.Stmt
	deleteStmt *sql.Stmt
}

func NewRepository(db *sql.DB) (*Repository, error) {
	r := &Repository{db: db}

	var err error

	r.createStmt, err = db.Prepare(`
		INSERT INTO todos (id, title, completed)
		VALUES ($1, $2, false)
	`)
	if err != nil {
		return nil, err
	}

	r.getStmt, err = db.Prepare(`
		SELECT id, title, completed
		FROM todos
		WHERE id = $1
	`)
	if err != nil {
		return nil, err
	}

	r.listStmt, err = db.Prepare(`
		SELECT id, title, completed
		FROM todos
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}

	r.updateStmt, err = db.Prepare(`
		UPDATE todos
		SET title = $1, completed = $2
		WHERE id = $3
	`)
	if err != nil {
		return nil, err
	}

	r.deleteStmt, err = db.Prepare(`
		DELETE FROM todos WHERE id = $1
	`)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Repository) Create(ctx context.Context, title string) (model.Todo, error) {
	t := model.Todo{
		Id:    uuid.NewString(),
		Title: title,
	}

	_, err := r.createStmt.ExecContext(ctx, t.Id, t.Title)
	if err != nil {
		log.Default().Println("repository: failed to create todo:", err)
		return model.Todo{}, err
	}

	log.Default().Println("repository: Created todo successfully:", t.Id)
	return t, nil
}

func (r *Repository) Get(ctx context.Context, id string) (model.Todo, error) {
	var t model.Todo

	err := r.getStmt.
		QueryRowContext(ctx, id).
		Scan(&t.Id, &t.Title, &t.Completed)

	if err == sql.ErrNoRows {
		log.Default().Println("repository: todo not found:", id)
		return model.Todo{}, nil
	}
	if err != nil {
		log.Default().Println("repository: failed to get todo:", err)
		return model.Todo{}, err
	}

	log.Default().Println("repository: Fetched todo successfully:", t.Id)
	return t, nil
}

func (r *Repository) Update(ctx context.Context, t *model.Todo) (model.Todo, error) {
	_, err := r.updateStmt.ExecContext(
		ctx,
		t.Title,
		t.Completed,
		t.Id,
	)
	if err != nil {
		log.Default().Println("repository: failed to update todo:", err)
		return model.Todo{}, err
	}

	log.Default().Println("repository: Updated todo successfully:", t.Id)
	return *t, nil
}

func (r *Repository) Delete(ctx context.Context, id string) error {
	_, err := r.deleteStmt.ExecContext(ctx, id)
	if err != nil {
		log.Default().Println("repository: failed to delete todo:", err)
		return err
	}
	return nil
}

func (r *Repository) List(ctx context.Context) ([]model.Todo, error) {
	rows, err := r.listStmt.QueryContext(ctx)
	if err != nil {
		log.Default().Println("repository: failed to list todos:", err)
		return nil, err
	}
	defer rows.Close()

	var result []model.Todo
	for rows.Next() {
		var t model.Todo
		if err := rows.Scan(&t.Id, &t.Title, &t.Completed); err != nil {
			log.Default().Println("repository: scan failed:", err)
			continue
		}
		result = append(result, t)
	}

	log.Default().Println("repository: Listed todos successfully, count:", len(result))
	return result, nil
}

func (r *Repository) Close() {
	stmts := []*sql.Stmt{
		r.createStmt,
		r.getStmt,
		r.listStmt,
		r.updateStmt,
		r.deleteStmt,
	}

	for _, stmt := range stmts {
		if stmt != nil {
			stmt.Close()
		}
	}
}

// package repository

// import (
// 	"context"
// 	"database/sql"
// 	"log"

// 	"github.com/google/uuid"
// 	"github.com/haakaashs/todos-backend/internal/model"
// )

// type Repository struct {
// 	db *sql.DB
// }

// func NewRepository(db *sql.DB) *Repository {
// 	return &Repository{db: db}
// }

// func (r *Repository) Create(ctx context.Context, title string) model.Todo {
// 	t := model.Todo{
// 		Id:    uuid.NewString(),
// 		Title: title,
// 	}

// 	query := `
// 		INSERT INTO todos (id, title, completed)
// 		VALUES ($1, $2, false)
// 	`

// 	r.db.ExecContext(ctx, query, t.Id, t.Title)
// 	log.Default().Println("repository: Created todo successfully:", t.Id)
// 	return t
// }

// func (r *Repository) Get(ctx context.Context, id string) (model.Todo, bool) {
// 	var t model.Todo

// 	query := `
// 		SELECT id, title, completed
// 		FROM todos
// 		WHERE id = $1
// 	`

// 	err := r.db.QueryRowContext(ctx, query, id).
// 		Scan(&t.Id, &t.Title, &t.Completed)
// 	if err == sql.ErrNoRows {
// 		return model.Todo{}, false
// 	}
// 	log.Default().Println("repository: Fetched todo successfully:", t.Id)
// 	return t, true
// }

// func (r *Repository) List(ctx context.Context) []model.Todo {
// 	query := `
// 		SELECT id, title, completed
// 		FROM todos
// 		ORDER BY created_at DESC
// 	`

// 	rows, _ := r.db.QueryContext(ctx, query)
// 	defer rows.Close()

// 	var result []model.Todo
// 	for rows.Next() {
// 		var t model.Todo
// 		rows.Scan(&t.Id, &t.Title, &t.Completed)
// 		result = append(result, t)
// 	}
// 	log.Default().Println("repository: Listed todos successfully, count:", len(result))
// 	return result
// }

// func (r *Repository) Delete(ctx context.Context, id string) {
// 	r.db.ExecContext(ctx, `DELETE FROM todos WHERE id = $1`, id)
// }

// func (r *Repository) Update(ctx context.Context, t model.Todo) model.Todo {
// 	query := `
// 		UPDATE todos
// 		SET title = $1, completed = $2
// 		WHERE id = $3
// 	`

// 	r.db.ExecContext(ctx, query, t.Title, t.Completed, t.Id)
// 	log.Default().Println("repository: Updated todo successfully:", t.Id)
// 	return t
// }

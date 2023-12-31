package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/matiaspub/todo-api/pkg/entity"
	"github.com/sirupsen/logrus"
)

type TodoListPg struct {
	db *sqlx.DB
}

func NewTodoListPg(db *sqlx.DB) *TodoListPg {
	return &TodoListPg{db: db}
}

func (r *TodoListPg) Create(userId int, todo entity.TodoList) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var listId int
	row := tx.QueryRow("INSERT INTO todo_lists (title, description) VALUES ($1, $2) RETURNING id", todo.Title, todo.Description)
	if err := row.Scan(&listId); err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	_, err = tx.Exec("INSERT INTO users_lists (user_id, list_id) VALUES ($q, $2)", userId, listId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return listId, tx.Commit()
}

func (r *TodoListPg) GetAll(userId int) ([]entity.TodoList, error) {
	var lists []entity.TodoList
	err := r.db.Select(&lists, "SELECT tl.id, title, description FROM todo_lists tl LEFT JOIN public.users_lists ul on tl.id = ul.list_id WHERE ul.user_id = $1", userId)
	return lists, err
}

func (r *TodoListPg) GetOne(userId int, listId int) (entity.TodoList, error) {
	var todoList entity.TodoList
	err := r.db.Get(todoList, "SELECT tl.id, title, description FROM todo_lists tl LEFT JOIN public.users_lists ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND tl.id = $2", userId, listId)
	return todoList, err
}

func (r *TodoListPg) Delete(userId int, listId int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	_, err = r.db.Exec("DELETE FROM users_lists WHERE user_id = $1 AND list_id = $2", userId, listId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	_, err = r.db.Exec("DELETE FROM todo_lists WHERE id = $1", listId)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func (r *TodoListPg) Update(userId int, listId int, list entity.UpdateListInput) error {
	logrus.Debugf("TodoList %d Update: %v", listId, list)
	_, err := r.db.Exec("UPDATE todo_lists SET title = COALESCE($1, title), description = COALESCE($2, description) WHERE id IN (SELECT list_id FROM users_lists WHERE user_id = $3 AND list_id = $4)", *list.Title, *list.Description, userId, listId)
	if err != nil {
		return err
	}
	return nil
}

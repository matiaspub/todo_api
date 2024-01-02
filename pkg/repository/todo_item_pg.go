package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/matiaspub/todo-api/pkg/entity"
)

type TodoItemPg struct {
	db *sqlx.DB
}

func (r *TodoItemPg) Delete(userId int, itemId int) error {
	_, err := r.db.Exec(`DELETE FROM todo_items WHERE id IN (
					SELECT ti.id FROM todo_items ti 
					    INNER JOIN public.lists_items li on ti.id = li.item_id 
					    INNER JOIN public.users_lists ul on li.list_id = ul.list_id 
				 	WHERE ul.user_id = $1 AND ti.id = $2
				)`, userId, itemId)
	return err
}

func (r *TodoItemPg) Update(userId int, itemId int, input entity.UpdateTodoItemInput) error {
	_, err := r.db.Exec(`UPDATE todo_items SET 
						title = COALESCE($1, title), 
						description = COALESCE($2, description), 
						done = COALESCE($3, done) 
					WHERE id in (
						SELECT ti.id FROM todo_items ti 
						    INNER JOIN public.lists_items li on ti.id = li.item_id 
						    INNER JOIN public.users_lists ul on li.list_id = ul.list_id 
					 	WHERE ul.user_id = $4 AND ti.id = $5
					)`,
		input.Title, input.Description, input.Done, userId, itemId)
	return err
}

func (r *TodoItemPg) GetOne(userId int, itemId int) (entity.TodoItem, error) {
	var item entity.TodoItem

	err := r.db.Get(&item, `SELECT ti.id, title, description, done 
				FROM todo_items ti
				    INNER JOIN lists_items li on ti.id = li.item_id 
				    INNER JOIN users_lists ul on li.list_id = ul.list_id 
				WHERE ti.id = $1 AND ul.user_id = $2`, itemId, userId)

	return item, err
}

func (r *TodoItemPg) GetAll(userId int, listId int) ([]entity.TodoItem, error) {
	var list []entity.TodoItem

	err := r.db.Get(&list, `SELECT ti.id, title, description, done 
			FROM todo_items ti 
			    INNER JOIN lists_items li on ti.id = li.item_id 
			    INNER JOIN users_lists ul on li.list_id = ul.list_id 
			WHERE ul.user_id = $1 
			  AND li.list_id = $2`, userId, listId)

	return list, err
}

func (r *TodoItemPg) Create(listId int, item entity.TodoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var itemId int
	row := tx.QueryRow(`INSERT INTO todo_items (title, description) VALUES ($1, $2) 
                                            RETURNING id`, item.Title, item.Description)
	err = row.Scan(&itemId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	_, err = tx.Exec("INSERT INTO lists_items (item_id, list_id) VALUES ($1, $2)", itemId, listId)
	if err != nil {
		_ = tx.Rollback()
		return 0, err
	}

	return itemId, tx.Commit()
}

func NewTodoItemPg(db *sqlx.DB) *TodoItemPg {
	return &TodoItemPg{db: db}
}

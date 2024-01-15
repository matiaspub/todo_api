package repository

import (
	"errors"
	"github.com/matiaspub/todo-api/pkg/entity"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"log"
	"testing"
)

func TestTodoItemPg_Create(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := NewTodoItemPg(db)

	type args struct {
		listId int
		item   entity.TodoItem
	}

	type mockBehaviour func(args args, id int)

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		args          args
		id            int
		wantError     bool
	}{
		{
			name: "OK",
			args: args{listId: 1, item: entity.TodoItem{Title: "test title", Description: "test description"}},
			id:   2,
			mockBehaviour: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)

				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs(args.item.Title, args.item.Description).
					WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO lists_items").
					WithArgs(id, args.listId).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "Empty Fields",
			args: args{listId: 1, item: entity.TodoItem{Title: "", Description: "test description"}},
			mockBehaviour: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)

				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs(args.item.Title, args.item.Description).
					WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantError: true,
		},
		{
			name: "2nd Insert Error",
			args: args{listId: 1, item: entity.TodoItem{Title: "test title", Description: "test description"}},
			mockBehaviour: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)

				mock.ExpectQuery("INSERT INTO todo_items").
					WithArgs(args.item.Title, args.item.Description).
					WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO lists_items").
					WithArgs(id, args.listId).
					WillReturnError(errors.New("some error"))

				mock.ExpectRollback()
			},
			wantError: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehaviour(testCase.args, testCase.id)

			itemID, err := r.Create(testCase.args.listId, testCase.args.item)
			if testCase.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, testCase.id, itemID)
			}
		})
	}
}

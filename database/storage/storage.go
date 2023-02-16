package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
)

type Storage interface {
	Save(ctx context.Context, db *sql.DB) error
	Get(ctx context.Context, db *sql.DB) (bool, []string, error)
}

type UserData struct {
	ChatID   int64
	UserName string
}

type Tasks struct {
	TaskID   int64
	TaskText string
}

func Save(s Storage, ctx context.Context, db *sql.DB) error {
	return s.Save(ctx, db)
}

func Get(s Storage, ctx context.Context, db *sql.DB) (bool, []string, error) {
	return s.Get(ctx, db)
}

func (u *UserData) Save(ctx context.Context, db *sql.DB) error {
	isID, _, err := u.Get(context.Background(), db)
	if err != nil {
		panic(err)
	}
	if isID {
		return nil
	} else {
		q := `insert into userData(chatID, userName) values(?, ?)`
		if _, err := db.ExecContext(ctx, q, u.ChatID, u.UserName); err != nil {
			return fmt.Errorf("can't add user to database: %w", err)
		}
	}
	return nil
}

func (t *Tasks) Save(ctx context.Context, db *sql.DB) error {
	q := `insert into tasks(taskID, taskText) values(?, ?)`
	if _, err := db.ExecContext(ctx, q, t.TaskID, t.TaskText); err != nil {
		return fmt.Errorf("can't add task to database: %w", err)
	}
	return nil
}

func (u *UserData) Get(ctx context.Context, db *sql.DB) (bool, []string, error) {
	q := `select chatID from userData where chatID=?`
	var resp int
	var result []string
	err := db.QueryRowContext(ctx, q, u.ChatID).Scan(&resp)
	result = append(result, strconv.Itoa(resp)) // =1 row
	if err == sql.ErrNoRows {
		return false, result, nil
	} else if err != nil {
		return false, result, fmt.Errorf("can't get data from database: %w", err)
	}

	return true, result, nil
}

func (t *Tasks) Get(ctx context.Context, db *sql.DB) (bool, []string, error) {
	q := `select taskText from tasks where taskID=?`
	var resp []string
	rows, err := db.QueryContext(ctx, q, t.TaskID) // >=1 row
	if err == sql.ErrNoRows {
		return false, []string{}, nil
	} else if err != nil {
		return false, []string{}, fmt.Errorf("can' get data from database: %w", err)
	}
	for rows.Next() { // fill a slice for selected data
		var row string
		err = rows.Scan(&row)
		if err != nil {
			return true, []string{}, err
		}
		resp = append(resp, row)
	}
	return true, resp, nil
}

func (t *Tasks) Delete(ctx context.Context, db *sql.DB, index int) (bool, error) {
	q := `delete from tasks where taskText=?`
	_, list, _ := t.Get(context.Background(), db)
	text := list[index-1]
	if result, err := db.ExecContext(ctx, q, text); err != nil {
		return false, fmt.Errorf("can't delete row from database: %w", err)
	} else if n, _ := result.RowsAffected(); n == 0 {
		return false, nil
	}
	return true, nil
}

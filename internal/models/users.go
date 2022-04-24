package models

import (
	"github.com/Masterminds/squirrel"
)

type Users struct {
	sql squirrel.StatementBuilderType
}

func UserModel(db *Database) *Users {
	return &Users{
		sql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db.db),
	}
}

type User struct {
	Name string `json:"name"`
}

func (u *Users) List() ([]User, error) {
	q := u.sql.
		Select("name").
		From("users").
		OrderBy("name ASC")

	rows, err := q.Query()
	if err != nil {
		return []User{}, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(
			&user.Name,
		)

		if err != nil {
			return []User{}, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (u *Users) Add(name string) (error) {
	q := u.sql.
		Insert("users").
		Columns("name").
		Values(name)

	rows, err := q.Query()
	if err == nil {
		defer rows.Close()
	}

	return err
}

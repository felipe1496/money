package users

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"github.com/oklog/ulid/v2"
)

type UsersRepo interface {
	ListUsers(filter UserFilter) ([]User, error)
	CreateUser(input CreateUserInput) (User, error)
}

type UsersRepoImpl struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) UsersRepo {
	return &UsersRepoImpl{db: db}
}

func (r *UsersRepoImpl) ListUsers(filter UserFilter) ([]User, error) {
	query := squirrel.
		Select("id", "name", "email", "avatar_url", "created_at", "username").
		From("users").
		PlaceholderFormat(squirrel.Dollar)

	if filter.Email != "" {
		query = query.Where(squirrel.Eq{"email": filter.Email})
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.AvatarURL,
			&user.CreatedAt,
			&user.Username,
		); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UsersRepoImpl) CreateUser(input CreateUserInput) (User, error) {
	query, args, err := squirrel.
		Insert("users").
		Columns("id", "name", "email", "avatar_url", "username", "created_at").
		Values(
			ulid.Make().String(),
			input.Name,
			input.Email,
			input.AvatarURL,
			input.Username,
			squirrel.Expr("NOW()"),
		).
		Suffix("RETURNING id, name, email, avatar_url, created_at, username").
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return User{}, err
	}

	var user User
	err = r.db.QueryRow(query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.AvatarURL,
		&user.CreatedAt,
		&user.Username,
	)

	if err != nil {
		return User{}, err
	}

	return user, nil
}

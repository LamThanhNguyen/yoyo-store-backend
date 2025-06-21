package models

import (
	"context"
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User is the type for users
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// GetUserByEmail gets a user by email address
func (m *DBModel) GetUserByEmail(email string) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	email = strings.ToLower(email)
	var u User

	row := m.DB.QueryRowContext(ctx, `
		select
			id, first_name, last_name, email, password, created_at, updated_at
		from
			users
		where email = ?`, email)

	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		return u, err
	}

	return u, nil
}

// UpdatePasswordForUser updates the password hash for a given user, by user id
func (m *DBModel) UpdatePasswordForUser(u User, hash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `update users set password = ? where id = ?`
	_, err := m.DB.ExecContext(ctx, stmt, hash, u.ID)
	if err != nil {
		return err
	}

	return nil
}

// GetAllUsersPaginated returns a slice of a subset of users
func (m *DBModel) GetAllUsersPaginated(pageSize, page int) ([]*User, int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	offset := (page - 1) * pageSize

	var users []*User

	query := `
		select
			id, last_name, first_name, email, created_at, updated_at
		from
			users
		order by
			last_name, first_name
		limit ? offset ?
	`

	rows, err := m.DB.QueryContext(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		err = rows.Scan(
			&u.ID,
			&u.LastName,
			&u.FirstName,
			&u.Email,
			&u.CreatedAt,
			&u.UpdatedAt,
		)
		if err != nil {
			return nil, 0, 0, err
		}
		users = append(users, &u)
	}

	query = `
		select
			count(id)
		from
			users
	`

	var totalRecords int
	countRow := m.DB.QueryRowContext(ctx, query)
	err = countRow.Scan(&totalRecords)
	if err != nil {
		return nil, 0, 0, err
	}

	lastPage := totalRecords / pageSize

	return users, lastPage, totalRecords, nil
}

// GetOneUser returns one user by id
func (m *DBModel) GetOneUser(id int) (User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var u User

	query := `
		select
			id, last_name, first_name, email, created_at, updated_at
		from
			users
		where id = ?`

	row := m.DB.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&u.ID,
		&u.LastName,
		&u.FirstName,
		&u.Email,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}
	return u, nil
}

// EditUser edits an existing user
func (m *DBModel) EditUser(u User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		update users set
			first_name = ?,
			last_name = ?,
			email = ?,
			updated_at = ?
		where
			id = ?`

	_, err := m.DB.ExecContext(ctx, stmt,
		u.FirstName,
		u.LastName,
		u.Email,
		time.Now(),
		u.ID,
	)

	if err != nil {
		return err
	}
	return nil
}

// AddUser inserts a user into the database
func (m *DBModel) AddUser(u User, hash string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into users (first_name, last_name, email, password, created_at, updated_at)
		values (?, ?, ?, ?, ?, ?)`

	_, err := m.DB.ExecContext(ctx, stmt,
		u.FirstName,
		u.LastName,
		u.Email,
		hash,
		time.Now(),
		time.Now(),
	)

	if err != nil {
		return err
	}
	return nil
}

// DeleteUser deletes a user by id
func (m *DBModel) DeleteUser(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `delete from users where id = ?`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	stmt = "delete from tokens where user_id = ?"
	_, err = m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

// Authenticate attempts to log a user in by comparing supplied password with password hash
func (m *DBModel) Authenticate(email, password string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = ?", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, errors.New("incorrect password")
	} else if err != nil {
		return 0, err
	}

	return id, nil
}

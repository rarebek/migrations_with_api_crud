package storage

import (
	"database/sql"
	"test_migration/models"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func connector() (*sql.Tx, error) {
	dsn := "user=postgres password=nodirbek dbname=migration sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func CreateUser(user models.User) (*models.User, error) {
	tx, err := connector()
	if err != nil {
		return nil, err
	}
	var respUser models.User
	uuid := uuid.NewString()
	query := "INSERT INTO users(uuid, name, age) VALUES ($1, $2, $3) RETURNING id, uuid, name, age"
	if err := tx.QueryRow(query, uuid, user.Name, user.Age).Scan(&respUser.Id, &respUser.Uuid, &respUser.Name, &respUser.Age); err != nil {
		tx.Rollback()
		return nil, err
	}

	query = "INSERT INTO posts(content, user_id) VALUES ($1, $2) RETURNING content"
	if err := tx.QueryRow(query, user.Content, &respUser.Id).Scan(&respUser.Content); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &respUser, nil
}

func UpdateUser(uuid string, user models.User) (*models.User, error) {
	tx, err := connector()
	if err != nil {
		return nil, err
	}
	var respUser models.User
	if err := tx.QueryRow("UPDATE users SET name = $1, age = $2 WHERE uuid = $3 RETURNING id, uuid, name, age", user.Name, user.Age, uuid).Scan(&respUser.Id, &respUser.Uuid, &respUser.Name, &respUser.Age); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.QueryRow("UPDATE posts SET content = $1 WHERE id = $2 RETURNING content", user.Content, respUser.Id).Scan(&respUser.Content); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return &respUser, nil
}

func DeleteUser(id int) (*models.User, error) {
	tx, err := connector()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	var respUser models.User

	if err := tx.QueryRow("DELETE from posts WHERE id = $1 RETURNING content", id).Scan(&respUser.Content); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.QueryRow("DELETE from users WHERE id = $1 RETURNING id, uuid, name, age", id).Scan(&respUser.Id, &respUser.Uuid, &respUser.Name, &respUser.Age); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &respUser, nil
}

func GetOne(uuid string) (*models.User, error) {
	tx, err := connector()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var respUser models.User
	if err := tx.QueryRow("SELECT u.id, u.uuid, u.name, u.age, p.content FROM users u join posts p on u.id = p.user_id WHERE uuid = $1", uuid).
		Scan(&respUser.Id, &respUser.Uuid, &respUser.Name, &respUser.Age, &respUser.Content); err != nil {
		return nil, err
	}

	return &respUser, nil
}

func GetAll(page, limit int) ([]*models.User, error) {
	tx, err := connector()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	offset := limit * (page - 1)
	query := "SELECT u.id, u.uuid, u.name, u.age, p.content FROM users u join posts p on u.id = p.user_id LIMIT $1 OFFSET $2"
	rows, err := tx.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var respUsers []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.Uuid, &user.Name, &user.Age, &user.Content); err != nil {
			return nil, err
		}
		respUsers = append(respUsers, &user)
	}

	return respUsers, nil
}

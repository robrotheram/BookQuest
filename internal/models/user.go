package models

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func GetUsers(db *bun.DB) []User {
	users := []User{}
	if err := db.NewSelect().
		Model(&users).
		Scan(context.Background()); err != nil {
		panic(err)
	}

	return users
}

func GetUser(db *bun.DB, id string) (User, error) {
	user := User{}
	err := db.NewSelect().
		Model(&user).
		Where("id = ?", id).
		Scan(context.Background())
	return user, err
}

func GetUserByEmail(db *bun.DB, email string) (User, error) {
	user := User{}
	err := db.NewSelect().
		Model(&user).
		Where("email = ?", email).
		Scan(context.Background())
	return user, err
}

func CreateUser(db *bun.DB, user User) error {
	_, err := db.NewInsert().Model(&user).Exec(context.Background())
	return err
}

func UpdateUser(db *bun.DB, user User) error {
	_, err := db.NewUpdate().Model(&user).WherePK().Exec(context.Background())
	return err
}
func UserGetLinks(db *bun.DB, id uuid.UUID) ([]Link, error) {
	user := User{}
	err := db.NewSelect().
		Model(&user).
		Where("id = ?", id).
		Relation("Links").
		Scan(context.Background())
	return user.Links, err
}

func UserGetFavorites(db *bun.DB, id uuid.UUID) ([]Link, error) {
	user := User{}
	err := db.NewSelect().
		Model(&user).
		Where("id = ?", id).
		Relation("Favourites").
		Scan(context.Background())

	return user.Favourites, err
}

package models

import (
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func GetFavourite(db *bun.DB, linkId string, userId uuid.UUID) (FavouriteLinks, error) {
	var fav FavouriteLinks
	_, err := db.NewSelect().
		Table("favourite_links").
		Where("link_id = ? AND user_id = ?", linkId, userId).
		Exec(context.Background(), &fav)
	return fav, err
}
func FavouriteExists(db *bun.DB, linkId string, userId uuid.UUID) bool {
	fav, err := GetFavourite(db, linkId, userId)
	if err != nil {
		return false
	}
	return fav.LinkId != uuid.Nil
}

func AddFavourite(db *bun.DB, linkId string, userId uuid.UUID) error {
	value := FavouriteLinks{
		LinkId: uuid.MustParse(linkId),
		UserID: (userId),
	}
	_, err := db.NewInsert().Model(&value).Exec(context.Background())
	return err
}

func RemoveFavourite(db *bun.DB, linkId string, userId uuid.UUID) error {
	_, err := db.NewDelete().
		Model((*FavouriteLinks)(nil)).
		Where("link_id = ? AND user_id = ?", linkId, userId).
		Exec(context.Background())
	return err
}

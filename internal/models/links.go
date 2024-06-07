package models

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func CreateLink(db *bun.DB, link Link, userId uuid.UUID) error {
	_, err := db.NewInsert().Model(&link).Exec(context.Background())
	if err != nil {
		return err
	}
	meta := UserToLink{
		LinkId: link.Id,
		UserID: userId,
	}
	_, err = db.NewInsert().Model(&meta).Exec(context.Background())
	return err
}

func GetLink(db *bun.DB, id string) (Link, error) {
	link := Link{}
	err := db.NewSelect().
		Model(&link).
		Where("id = ?", id).
		Scan(context.Background())
	return link, err
}

func GetLinks(db *bun.DB) []Link {
	link := []Link{}
	if err := db.NewSelect().
		Model(&link).
		Scan(context.Background()); err != nil {
		panic(err)
	}
	return link
}

func GetTopLinks(db *bun.DB) []Link {
	meta := []LinkMeta{}
	links := []Link{}
	if err := db.NewSelect().
		Model(&meta).
		Relation("Link", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("sharing = ?", PUBLIC)
		}).
		OrderExpr("clicked").
		Limit(5).
		Scan(context.Background()); err != nil {
		panic(err)
	}
	for _, m := range meta {
		links = append(links, *m.Link)
	}
	return links
}

func GetUserLinksMeta(db *bun.DB, userId uuid.UUID, limit int) ([]LinkMeta, error) {
	meta := []LinkMeta{}
	err := db.NewSelect().
		Model(&meta).
		Relation("Link").
		Relation("User", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("user_id = ?", userId)
		}).
		OrderExpr("clicked").
		Limit(limit).
		Scan(context.Background())
	return meta, err
}

func GetUserTopLinks(db *bun.DB, userId uuid.UUID) []Link {
	meta, err := GetUserLinksMeta(db, userId, 5)
	links := []Link{}
	if err != nil {
		return links
	}
	for _, m := range meta {
		links = append(links, *m.Link)
	}
	return links
}

func GetUserLinkMeta(db *bun.DB, linkId string, userId uuid.UUID) (LinkMeta, error) {
	var meta LinkMeta
	err := db.NewSelect().
		Model(&meta).
		Relation("Link", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("link_id = ?", linkId)
		}).
		Relation("User", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("user_id = ?", userId)
		}).
		Scan(context.Background())
	return meta, err
}

func CreateUserLinkMeta(db *bun.DB, linkId string, userId uuid.UUID) error {
	value := LinkMeta{
		LinkId:   uuid.MustParse(linkId),
		UserID:   userId,
		Clicked:  1,
		LastUsed: time.Now(),
	}
	if _, err := db.NewInsert().Model(&value).Exec(context.Background()); err != nil {
		return err
	}
	return nil
}

func UpdateUserLinkMeta(db *bun.DB, linkId string, userId uuid.UUID) error {
	meta, err := GetUserLinkMeta(db, linkId, userId)
	if err != nil {
		return CreateUserLinkMeta(db, linkId, userId)
	}
	meta.Clicked = meta.Clicked + 1
	meta.LastUsed = time.Now()
	_, err = db.NewUpdate().Model(&meta).WherePK().Exec(context.Background())
	return err
}

func UpdateLink(db *bun.DB, link Link) error {
	_, err := db.NewUpdate().Model(&link).WherePK().Exec(context.Background())
	return err
}

func DeleteLink(db *bun.DB, linkId string) error {
	db.NewDelete().Model((*Link)(nil)).Where("id = ?", linkId).Exec(context.Background())
	db.NewDelete().Model((*FavouriteLinks)(nil)).Where("link_id = ?", linkId).Exec(context.Background())
	db.NewDelete().Model((*TeamLink)(nil)).Where("link_id = ?", linkId).Exec(context.Background())
	db.NewDelete().Model((*UserToLink)(nil)).Where("link_id = ?", linkId).Exec(context.Background())
	return nil
}

func containsQuery(link Link, query string) bool {
	return strings.Contains(strings.ToLower(link.Title), strings.ToLower(query)) ||
		strings.Contains(strings.ToLower(link.Description), strings.ToLower(query)) ||
		strings.Contains(strings.ToLower(link.Url), strings.ToLower(query))
}

func FilterLinks(links []Link, query string) []Link {
	var filteredItems []Link
	if len(query) == 0 {
		return links
	}
	for _, item := range links {
		if containsQuery(item, query) {
			filteredItems = append(filteredItems, item)
		}
	}
	return filteredItems
}

func (l *Link) Update(link Link) {
	l.Title = link.Title
	l.Description = link.Description
	l.Tags = link.Tags
	l.Icon = link.Icon
	l.Url = link.Url
	l.Updated = time.Now()
	l.Sharing = link.Sharing
}

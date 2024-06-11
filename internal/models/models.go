package models

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func RegisterModels(db *bun.DB) {
	// Register many to many model so bun can better recognize m2m relation.
	// This should be done before you use the model for the first time.
	db.RegisterModel((*FavouriteLinks)(nil))
	db.RegisterModel((*UserToTeam)(nil))
	db.RegisterModel((*UserToLink)(nil))
	db.RegisterModel((*FavouriteLinks)(nil))
	db.RegisterModel((*TeamLink)(nil))
	db.RegisterModel((*LinkMeta)(nil))
}

func CreateSchema(ctx context.Context, db *bun.DB) error {
	if _, err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pg_trgm`); err != nil {
		log.Fatalf("failed to enable pg_trgm extension: %v", err)
	}

	models := []interface{}{
		(*Link)(nil),
		(*User)(nil),
		(*TeamLink)(nil),
		(*Team)(nil),
		(*UserToTeam)(nil),
		(*UserToLink)(nil),
		(*FavouriteLinks)(nil),
		(*LinkMeta)(nil),
	}
	for _, model := range models {
		if _, err := db.NewCreateTable().Model(model).Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}

type ShareSettings string

const PUBLIC = ShareSettings("PUBLIC")
const PRIVATE = ShareSettings("PRIVATE")
const TEAM = ShareSettings("TEAM")

type TeamPermission string

const MEMBER = TeamPermission("Member")
const OWNER = TeamPermission("Owner")

type Link struct {
	Id          uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	Title       string
	Description string
	Tags        string
	Icon        string
	Url         string
	Updated     time.Time
	Sharing     ShareSettings
}

type Team struct {
	Id          uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	Name        string
	Description string
	Memebers    []User `bun:"m2m:user_to_teams,join:Team=User"`
	Links       []Link `bun:"m2m:team_links,join:Team=Link"`
}

type TeamLink struct {
	TeamId uuid.UUID `bun:"type:uuid,pk"`
	Team   *Team     `bun:"rel:belongs-to,join:team_id=id"`
	LinkId uuid.UUID `bun:"type:uuid,pk"`
	Link   *Link     `bun:"rel:belongs-to,join:link_id=id"`
}

type UserToTeam struct {
	UserID     uuid.UUID `bun:"type:uuid,pk"`
	User       *User     `bun:"rel:belongs-to,join:user_id=id"`
	TeamId     uuid.UUID `bun:"type:uuid,pk"`
	Team       *Team     `bun:"rel:belongs-to,join:team_id=id"`
	Permission TeamPermission
}

type User struct {
	Id         uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	Username   string
	Email      string
	Picture    string
	Links      []Link `bun:"m2m:user_to_links,join:User=Link"`
	Favourites []Link `bun:"m2m:favourite_links,join:User=Link"`
}

type FavouriteLinks struct {
	LinkId uuid.UUID `bun:"type:uuid,pk"`
	Link   *Link     `bun:"rel:belongs-to,join:link_id=id"`
	UserID uuid.UUID `bun:"type:uuid,pk"`
	User   *User     `bun:"rel:belongs-to,join:user_id=id"`
}

type UserToLink struct {
	LinkId uuid.UUID `bun:"type:uuid,pk"`
	Link   *Link     `bun:"rel:belongs-to,join:link_id=id"`
	UserID uuid.UUID `bun:"type:uuid,pk"`
	User   *User     `bun:"rel:belongs-to,join:user_id=id"`
}

type LinkMeta struct {
	Id       uuid.UUID `bun:"type:uuid,pk,default:gen_random_uuid()"`
	UserID   uuid.UUID `bun:"type:uuid"`
	User     *User     `bun:"rel:belongs-to,join:user_id=id"`
	LinkId   uuid.UUID `bun:"type:uuid"`
	Link     *Link     `bun:"rel:belongs-to,join:link_id=id"`
	Clicked  int
	LastUsed time.Time
}

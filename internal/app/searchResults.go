package app

import (
	"BookQuest/internal/models"
	"context"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type SearchResult struct {
	Favorited         bool    `bun:"favorited"`
	HighestSimilarity float64 `bun:"highest_similarity"`
	models.Link
}

func ConvertLinkToSearchResult(links []models.Link) []SearchResult {
	var results = make([]SearchResult, len(links))
	for i, link := range links {
		results[i] = SearchResult{
			Link: link,
		}
	}
	return results
}

/*
*
Query using Fuzzy search using the postgress extension "pg_term"
Also Join Links table on user and teams table so that return results that have the sharing settings
Public
Private AND user matches Use_ID
Team and user user is in team
*/
func QueryLinks(db *bun.DB, searchText string, userId uuid.UUID) ([]SearchResult, error) {
	var links []SearchResult
	err := db.NewSelect().
		Table("links").
		ColumnExpr("links.*, fl.link_id IS NOT NULL AS favorited").
		ColumnExpr("GREATEST(similarity(links.title, ?), similarity(links.description, ?), similarity(links.tags, ?)) AS highest_similarity", searchText, searchText, searchText).
		Join("LEFT JOIN user_to_links utl ON utl.link_id = links.id").
		Join("LEFT JOIN user_to_teams utt ON utt.user_id = ?", userId).
		Join("LEFT JOIN team_links tl ON tl.team_id = utt.team_id").
		Join("LEFT JOIN favourite_links fl ON fl.link_id = links.id AND fl.user_id = ?", userId).
		Where("title % ? OR description % ? OR tags % ?", searchText, searchText, searchText).
		WhereGroup("AND", func(q *bun.SelectQuery) *bun.SelectQuery {
			return q.Where("links.sharing = ?", models.PUBLIC).
				WhereOr("links.sharing = ? AND utl.user_id = ?", models.PRIVATE, userId).
				WhereOr("links.sharing = ? AND tl.team_id IN (SELECT team_id FROM user_to_teams WHERE user_id = ?)", models.TEAM, userId)
		}).
		GroupExpr("links.id , fl.link_id").
		OrderExpr("highest_similarity DESC, title <-> ?, description <-> ?, tags <-> ?", searchText, searchText, searchText).
		Scan(context.Background(), &links)

	return links, err
}

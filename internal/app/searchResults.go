package app

import (
	"BookQuest/internal/models"
	"context"
	"fmt"

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
func QueryLinks(db *bun.DB, searchText string, userId string) ([]SearchResult, error) {
	var links []SearchResult
	query := `WITH user_links AS (
			SELECT 
				links.*, 
				fl.link_id IS NOT NULL AS favorited, 
				GREATEST(
					similarity(links.title, ?), 
					similarity(links.description, ?), 
					similarity(links.tags, ?)
				) AS highest_similarity 
			FROM 
				links
			LEFT JOIN favourite_links fl ON fl.link_id = links.id AND fl.user_id = ?
			LEFT JOIN user_to_links utl ON utl.link_id = links.id
			LEFT JOIN team_links tl ON tl.link_id = links.id
			LEFT JOIN teams t ON t.id = tl.team_id
			LEFT JOIN user_to_teams utt ON utt.user_id = ? 
			WHERE 
				(links.sharing = 'PUBLIC') 
				OR (links.sharing = 'PRIVATE' AND utl.user_id = ?) 
				OR (links.sharing = 'TEAM' AND (t.visability = 'PUBLIC' OR (t.visability = 'PRIVATE' AND utt.user_id IS NOT NULL)))
		)
		SELECT 
			user_links.*, 
			user_links.highest_similarity 
		FROM 
			user_links
		WHERE 
			user_links.highest_similarity > 0.1
		GROUP BY 
			user_links.id, 
			user_links.title, 
			user_links.description, 
			user_links.tags, 
			user_links.icon, 
			user_links.url, 
			user_links.updated, 
			user_links.sharing, 
			user_links.favorited, 
			user_links.highest_similarity
		ORDER BY 
			user_links.highest_similarity DESC;
	`
	err := db.NewRaw(
		query, searchText, searchText, searchText, userId, userId, userId).Scan(context.Background(), &links)

	fmt.Println(err)

	return links, err
}

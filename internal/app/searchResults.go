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
	query := `SELECT subquery.*, subquery.highest_similarity FROM (
	SELECT links.*, fl.link_id IS NOT NULL AS favorited, 
	GREATEST(
		similarity(links.title, ?), 
		similarity(links.description, ?), 
		similarity(links.tags, ?)
	) AS highest_similarity 
	FROM "links" 
	LEFT JOIN user_to_links utl ON utl.link_id = links.id 
	LEFT JOIN user_to_teams utt ON utt.user_id = ? 
	LEFT JOIN team_links tl ON tl.team_id = utt.team_id 
	LEFT JOIN favourite_links fl ON fl.link_id = links.id 
	AND fl.user_id = ? 
	WHERE (links.sharing = 'PUBLIC') OR (links.sharing = 'PRIVATE' AND utl.user_id = ?) OR 
		  (links.sharing = 'TEAM' AND tl.team_id IN (SELECT team_id FROM user_to_teams WHERE user_id = ?))
	) AS subquery
	WHERE subquery.highest_similarity > 0.1
	GROUP BY subquery.id, subquery.title, subquery.description, subquery.tags,subquery.icon,subquery.url,subquery.updated,subquery.sharing,subquery.favorited,subquery.highest_similarity
	ORDER BY subquery.highest_similarity DESC;
	`
	err := db.NewRaw(
		query, searchText, searchText, searchText, userId, userId, userId, userId).Scan(context.Background(), &links)

	fmt.Println(err)

	return links, err
}

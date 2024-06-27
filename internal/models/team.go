package models

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func GetTeam(db *bun.DB, id string) (Team, error) {
	team := Team{}
	err := db.NewSelect().
		Model(&team).
		Where("id = ?", id).
		Relation("Memebers").
		Relation("Links").
		Scan(context.Background())
	return team, err
}

func GetTeamPermissions(db *bun.DB, id string) ([]UserToTeam, error) {
	var userToTeams []UserToTeam
	err := db.NewSelect().
		Model(&userToTeams).
		Relation("User").
		Where("team_id = ?", id).
		Scan(context.Background())
	return userToTeams, err
}

func GetTeams(db *bun.DB) []Team {
	team := []Team{}
	if err := db.NewSelect().
		Model(&team).
		Scan(context.Background()); err != nil {
		panic(err)
	}
	return team
}

func GetPublicTeams(db *bun.DB) ([]Team, error) {
	teams := []Team{}
	err := db.NewSelect().
		Model(&teams).
		Where("visability = ?", PUBLIC).
		Relation("Memebers").
		Relation("Links").
		Scan(context.Background())
	return teams, err
}

func GetTeamsForUser(db *bun.DB, userID uuid.UUID) ([]Team, error) {
	var teams []Team

	// Query for public teams
	publicTeams := db.NewSelect().
		Model((*Team)(nil)).
		Where("team.visability = ?", PUBLIC)

	// Query for private teams and public teams the user is in
	userTeams := db.NewSelect().
		Model((*Team)(nil)).
		Join("JOIN user_to_teams AS ut ON ut.team_id = team.id").
		Where("ut.user_id = ?", userID).
		Where("team.visability IN (?, ?)", PRIVATE, PUBLIC)

	// Union the two queries
	err := db.NewSelect().
		With("public_teams", publicTeams).
		With("user_teams", userTeams).
		TableExpr("(SELECT * FROM public_teams UNION SELECT * FROM user_teams) AS t").
		Scan(context.Background(), &teams)

	return teams, err
}

func UpdateTeam(db *bun.DB, team Team) error {
	_, err := db.NewUpdate().Model(&team).WherePK().Exec(context.Background())
	return err
}

func CreateTeam(db *bun.DB, team Team) error {
	_, err := db.NewInsert().Model(&team).Exec(context.Background())
	return err
}

func GetTeamsByUser(db *bun.DB, userId uuid.UUID) ([]Team, error) {
	teams := []Team{}
	err := db.NewSelect().
		Model(&teams).
		Relation("Memebers").
		Join("JOIN user_to_teams AS ut ON ut.team_id = team.id").
		Where("ut.user_id = ?", userId).
		Order("team.name ASC").
		Scan(context.Background())
	return teams, err
}

func GetTeamsByLink(db *bun.DB, linkId string) ([]Team, error) {
	teams := []Team{}
	err := db.NewSelect().
		Model(&teams).
		Join("JOIN team_links AS tl ON tl.team_id = team.id").
		Where("tl.link_id = ?", linkId).
		Scan(context.Background())
	return teams, err
}

func GetLinksByTeam(db *bun.DB, id string) []Link {
	team := Team{}
	if err := db.NewSelect().
		Model(&team).
		Where("id = ?", id).
		Relation("Links").
		Scan(context.Background()); err != nil {
		panic(err)
	}
	return team.Links
}

func AddLinkToTeam(db *bun.DB, linkId, teamId string) error {
	team := TeamLink{
		TeamId: uuid.MustParse(teamId),
		LinkId: uuid.MustParse(linkId),
	}
	_, err := db.NewInsert().Model(&team).Exec(context.Background())
	return err
}

func RemoveLinksToTeam(db *bun.DB, linkId, teamId string) error {
	_, err := db.NewDelete().
		Model((*TeamLink)(nil)).
		Where("link_id = ? AND team_id = ?", linkId, teamId).
		Exec(context.Background())
	return err
}

func AddUserToTeam(db *bun.DB, user uuid.UUID, teamId string, permission TeamPermission) error {
	team := UserToTeam{
		TeamId:     uuid.MustParse(teamId),
		UserID:     user,
		Permission: permission,
	}
	_, err := db.NewInsert().Model(&team).Exec(context.Background())
	return err
}

func RemoveUserToTeam(db *bun.DB, user, teamId string) error {
	_, err := db.NewDelete().
		Model((*UserToTeam)(nil)).
		Where("user_id = ? AND team_id = ?", user, teamId).
		Exec(context.Background())
	return err
}

func ModifyUserToTeam(db *bun.DB, user string, teamId string, permission TeamPermission) error {
	team := UserToTeam{
		TeamId:     uuid.MustParse(teamId),
		UserID:     uuid.MustParse(user),
		Permission: permission,
	}
	_, err := db.NewUpdate().
		Model(&team).
		Where("user_id = ? AND team_id = ?", user, teamId).
		Exec(context.Background())
	return err
}

func (t *Team) Update(team Team) {
	t.Name = team.Name
	t.Description = team.Description
	t.Visability = team.Visability
}

func FilterTeams(teams []Team, query string) []Team {
	var filteredItems []Team
	if len(query) == 0 {
		return teams
	}
	for _, item := range teams {
		if strings.Contains(strings.ToLower(item.Name), query) ||
			strings.Contains(strings.ToLower(item.Description), query) {
			filteredItems = append(filteredItems, item)
		}
	}
	return filteredItems
}

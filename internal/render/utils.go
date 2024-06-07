package render

import "BookQuest/internal/models"

func isSelectedTeam(team models.Team, selectedTeams []models.Team) bool {
	for _, t := range selectedTeams {
		if t.Id == team.Id {
			return true
		}
	}
	return false
}

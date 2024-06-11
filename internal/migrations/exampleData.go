package migrations

import (
	"BookQuest/internal/models"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

var ctx = context.Background()

func LoadExampleData(db *bun.DB, username, email string) error {
	const ExampleUserID = "731f0bc6-23db-47fc-b6ad-ea9d307dd0ee"
	const ExampleLinkID = "f72fd161-e361-48fb-90ff-24052cb4de86"

	users := []*models.User{
		{Id: uuid.MustParse(ExampleUserID), Username: username, Email: email},
		{Id: uuid.New(), Username: "john.smith", Email: "john.smith@example.com"},
		{Id: uuid.New(), Username: "emily.brown", Email: "emily.brown@example.com"},
		{Id: uuid.New(), Username: "michael.johnson", Email: "michael.johnson@example.com"},
		{Id: uuid.New(), Username: "sarah.lee", Email: "sarah.lee@example.com"},
	}

	teams := []*models.Team{
		{Id: uuid.New(), Name: "Red Dragons", Description: "A fierce and competitive team."},
		{Id: uuid.New(), Name: "Blue Whales", Description: "Known for their strategic gameplay."},
		{Id: uuid.New(), Name: "Golden Eagles", Description: "Masters of the air, swift and precise."},
		{Id: uuid.New(), Name: "Green Giants", Description: "Strong and resilient team members."},
		{Id: uuid.New(), Name: "Silver Sharks", Description: "Fast and agile, excellent in water-based games."},
		{Id: uuid.New(), Name: "Black Panthers", Description: "Stealthy and tactical, excels in surprise attacks."},
		{Id: uuid.New(), Name: "White Wolves", Description: "Teamwork and coordination are their strengths."},
		{Id: uuid.New(), Name: "Purple Phoenix", Description: "Rises from the ashes, never gives up."},
		{Id: uuid.New(), Name: "Orange Tigers", Description: "Fierce and powerful, dominates the field."},
		{Id: uuid.New(), Name: "Yellow Hornets", Description: "Small but deadly, quick and efficient."},
	}

	members := []*models.UserToTeam{
		{UserID: users[0].Id, TeamId: teams[0].Id, Permission: models.OWNER},
		{UserID: users[1].Id, TeamId: teams[0].Id, Permission: models.MEMBER},
		{UserID: users[2].Id, TeamId: teams[0].Id, Permission: models.MEMBER},
		{UserID: users[3].Id, TeamId: teams[0].Id, Permission: models.MEMBER},
		{UserID: users[0].Id, TeamId: teams[1].Id, Permission: models.OWNER},
		{UserID: users[0].Id, TeamId: teams[2].Id, Permission: models.OWNER},
		{UserID: users[0].Id, TeamId: teams[3].Id, Permission: models.OWNER},
		{UserID: users[0].Id, TeamId: teams[4].Id, Permission: models.OWNER},
		{UserID: users[0].Id, TeamId: teams[5].Id, Permission: models.OWNER},
	}

	links := []*models.Link{
		{Id: uuid.MustParse(ExampleLinkID), Title: "Google", Description: "Search engine", Icon: "https://www.google.com/favicon.ico", Url: "https://www.google.com", Updated: time.Now(), Sharing: models.PUBLIC},
		{Id: uuid.New(), Title: "Facebook", Description: "Social media platform", Icon: "https://www.facebook.com/favicon.ico", Url: "https://www.facebook.com", Updated: time.Now(), Sharing: models.TEAM},
		{Id: uuid.New(), Title: "Twitter", Description: "Microblogging platform", Icon: "https://www.twitter.com/favicon.ico", Url: "https://www.twitter.com", Updated: time.Now(), Sharing: models.TEAM},
		{Id: uuid.New(), Title: "LinkedIn", Description: "Professional networking site", Icon: "https://www.linkedin.com/favicon.ico", Url: "https://www.linkedin.com", Updated: time.Now(), Sharing: models.PUBLIC},
		{Id: uuid.New(), Title: "GitHub", Description: "Code hosting platform", Icon: "https://www.github.com/favicon.ico", Url: "https://www.github.com", Updated: time.Now(), Sharing: models.PRIVATE},
		{Id: uuid.New(), Title: "YouTube", Description: "Video sharing site", Icon: "https://www.youtube.com/favicon.ico", Url: "https://www.youtube.com", Updated: time.Now(), Sharing: models.TEAM},
		{Id: uuid.New(), Title: "Reddit", Description: "Social news aggregation", Icon: "https://www.reddit.com/favicon.ico", Url: "https://www.reddit.com", Updated: time.Now(), Sharing: models.PUBLIC},
		{Id: uuid.New(), Title: "Amazon", Description: "Online shopping site", Icon: "https://www.amazon.com/favicon.ico", Url: "https://www.amazon.com", Updated: time.Now(), Sharing: models.PRIVATE},
		{Id: uuid.New(), Title: "Wikipedia", Description: "Free online encyclopedia", Icon: "https://www.wikipedia.org/favicon.ico", Url: "https://www.wikipedia.org", Updated: time.Now(), Sharing: models.TEAM},
		{Id: uuid.New(), Title: "Netflix", Description: "Streaming service", Icon: "https://www.netflix.com/favicon.ico", Url: "https://www.netflix.com", Updated: time.Now(), Sharing: models.PUBLIC},
		{Id: uuid.New(), Title: "Slack", Description: "Team communication tool", Icon: "https://www.slack.com/favicon.ico", Url: "https://www.slack.com", Updated: time.Now(), Sharing: models.PRIVATE},
		{Id: uuid.New(), Title: "Trello", Description: "Project management tool", Icon: "https://www.trello.com/favicon.ico", Url: "https://www.trello.com", Updated: time.Now(), Sharing: models.TEAM},
		{Id: uuid.New(), Title: "Dropbox", Description: "Cloud storage service", Icon: "https://www.dropbox.com/favicon.ico", Url: "https://www.dropbox.com", Updated: time.Now(), Sharing: models.PUBLIC},
		{Id: uuid.New(), Title: "Spotify", Description: "Music streaming service", Icon: "https://www.spotify.com/favicon.ico", Url: "https://www.spotify.com", Updated: time.Now(), Sharing: models.PRIVATE},
		{Id: uuid.New(), Title: "Pinterest", Description: "Image sharing service", Icon: "https://www.pinterest.com/favicon.ico", Url: "https://www.pinterest.com", Updated: time.Now(), Sharing: models.TEAM},
		{Id: uuid.New(), Title: "Quora", Description: "Question and answer site", Icon: "https://www.quora.com/favicon.ico", Url: "https://www.quora.com", Updated: time.Now(), Sharing: models.PUBLIC},
		{Id: uuid.New(), Title: "Instagram", Description: "Photo sharing app", Icon: "https://www.instagram.com/favicon.ico", Url: "https://www.instagram.com", Updated: time.Now(), Sharing: models.PRIVATE},
		{Id: uuid.New(), Title: "WhatsApp", Description: "Messaging app", Icon: "https://www.whatsapp.com/favicon.ico", Url: "https://www.whatsapp.com", Updated: time.Now(), Sharing: models.TEAM},
		{Id: uuid.New(), Title: "Medium", Description: "Online publishing platform", Icon: "https://www.medium.com/favicon.ico", Url: "https://www.medium.com", Updated: time.Now(), Sharing: models.PUBLIC},
		{Id: uuid.New(), Title: "Adobe", Description: "Creative software company", Icon: "https://www.adobe.com/favicon.ico", Url: "https://www.adobe.com", Updated: time.Now(), Sharing: models.PRIVATE},
	}

	teamLinks := []*models.TeamLink{
		{
			TeamId: teams[0].Id,
			LinkId: links[1].Id,
		},
		{
			TeamId: teams[1].Id,
			LinkId: links[3].Id,
		},
	}

	favouriteLinks := []*models.FavouriteLinks{
		{
			UserID: users[0].Id,
			LinkId: links[1].Id,
		},
		{
			UserID: users[0].Id,
			LinkId: links[3].Id,
		},
		{
			UserID: users[0].Id,
			LinkId: links[4].Id,
		},
		{
			UserID: users[0].Id,
			LinkId: links[5].Id,
		},
	}

	userLinks := []*models.UserToLink{
		{
			UserID: users[0].Id,
			LinkId: links[6].Id,
		},
		{
			UserID: users[0].Id,
			LinkId: links[7].Id,
		},
		{
			UserID: users[0].Id,
			LinkId: links[8].Id,
		},
		{
			UserID: users[0].Id,
			LinkId: links[9].Id,
		},
	}

	for _, value := range users {
		if _, err := db.NewInsert().Model(value).Exec(ctx); err != nil {
			return err
		}
	}
	for _, value := range teams {
		if _, err := db.NewInsert().Model(value).Exec(ctx); err != nil {
			return err
		}
	}
	for _, value := range members {
		if _, err := db.NewInsert().Model(value).Exec(ctx); err != nil {
			return err
		}
	}
	for _, value := range links {
		if _, err := db.NewInsert().Model(value).Exec(ctx); err != nil {
			return err
		}
	}
	for _, value := range teamLinks {
		if _, err := db.NewInsert().Model(value).Exec(ctx); err != nil {
			return err
		}
	}
	for _, value := range favouriteLinks {
		if _, err := db.NewInsert().Model(value).Exec(ctx); err != nil {
			return err
		}
	}
	for _, value := range userLinks {
		if _, err := db.NewInsert().Model(value).Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}

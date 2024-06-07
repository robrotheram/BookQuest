package main

import (
	BookQuest "BookQuest/internal"
	Migrations "BookQuest/internal/migrations"
	"log"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "bookquest",
		Commands: []*cli.Command{
			{
				Name:  "server",
				Usage: "start the server",
				Action: func(c *cli.Context) error {
					templateFS := os.DirFS(".")
					staticFS := os.DirFS(".")
					mux := BookQuest.NewServer(staticFS, templateFS)
					logrus.Info("Starting server on :8090")
					return http.ListenAndServe(":8090", mux)
				},
			},
			{
				Name:  "db",
				Usage: "manage database migrations",
				Subcommands: []*cli.Command{
					{
						Name:  "init",
						Usage: "create migration tables",
						Action: func(c *cli.Context) error {
							return Migrations.Create(BookQuest.DB())
						},
					},
					{
						Name:  "migrate",
						Usage: "migrate database",
						Action: func(c *cli.Context) error {
							return Migrations.Migrate(BookQuest.DB())
						},
					},
					{
						Name:  "rollback",
						Usage: "rollback the last migration group",
						Action: func(c *cli.Context) error {
							return Migrations.Rollback(BookQuest.DB())
						},
					},
					{
						Name:  "lock",
						Usage: "lock migrations",
						Action: func(c *cli.Context) error {
							return Migrations.Lock(BookQuest.DB())
						},
					},
					{
						Name:  "unlock",
						Usage: "unlock migrations",
						Action: func(c *cli.Context) error {
							return Migrations.Unlock(BookQuest.DB())
						},
					},
					{
						Name:  "create_go",
						Usage: "create Go migration",
						Action: func(c *cli.Context) error {
							return Migrations.CreateGo(BookQuest.DB(), c.Args().Slice())
						},
					},
					{
						Name:  "create_sql",
						Usage: "create up and down SQL migrations",
						Action: func(c *cli.Context) error {
							return Migrations.CreateGo(BookQuest.DB(), c.Args().Slice())
						},
					},
					{
						Name:  "status",
						Usage: "print migrations status",
						Action: func(c *cli.Context) error {
							return Migrations.Status(BookQuest.DB())
						},
					},
					{
						Name:  "mark_applied",
						Usage: "mark migrations as applied without actually running them",
						Action: func(c *cli.Context) error {
							return Migrations.MarkApplied(BookQuest.DB())
						},
					},
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	BookQuest "BookQuest/pkg"
	"embed"
	"log"
)

//go:embed views
var ThemeFS embed.FS

//go:embed static
var StaticFS embed.FS

func main() {
	err := BookQuest.LoadConfig(".")
	if err != nil {
		log.Fatalf("Config Error: %v", err)
	}
	store, err := BookQuest.NewStore(BookQuest.Configuration.DatabasePath)
	if err != nil {
		log.Fatal(err)
	}
	server := BookQuest.NewServer(store, ThemeFS, StaticFS)
	server.Start()
}

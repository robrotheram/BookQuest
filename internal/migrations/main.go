package migrations

import (
	"context"
	"embed"
	"fmt"
	"strings"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

var Migrations = migrate.NewMigrations()

//go:embed *.sql
var sqlMigrations embed.FS

func init() {
	if err := Migrations.Discover(sqlMigrations); err != nil {
		panic(err)
	}
}

func Create(db *bun.DB) error {
	migrator := migrate.NewMigrator(db, Migrations)
	return migrator.Init(context.Background())
}

func Migrate(db *bun.DB) error {
	migrator := migrate.NewMigrator(db, Migrations)

	group, err := migrator.Migrate(context.Background())
	if err != nil {
		return err
	}

	if group.ID == 0 {
		fmt.Printf("there are no new Migrations to run\n")
		return nil
	}

	fmt.Printf("migrated to %s\n", group)
	return nil
}

func Rollback(db *bun.DB) error {
	migrator := migrate.NewMigrator(db, Migrations)
	group, err := migrator.Rollback(context.Background())
	if err != nil {
		return err
	}

	if group.ID == 0 {
		fmt.Printf("there are no groups to roll back\n")
		return nil
	}

	fmt.Printf("rolled back %s\n", group)
	return nil
}

func Lock(db *bun.DB) error {
	migrator := migrate.NewMigrator(db, Migrations)
	return migrator.Lock(context.Background())
}

func Unlock(db *bun.DB) error {
	migrator := migrate.NewMigrator(db, Migrations)
	return migrator.Unlock(context.Background())
}

func CreateGo(db *bun.DB, args []string) error {

	migrator := migrate.NewMigrator(db, Migrations)

	name := strings.Join(args, "_")
	mf, err := migrator.CreateGoMigration(context.Background(), name)
	if err != nil {
		return err
	}
	fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)

	return nil
}

func CreateSQL(db *bun.DB, args []string) error {
	migrator := migrate.NewMigrator(db, Migrations)

	name := strings.Join(args, "_")
	files, err := migrator.CreateSQLMigrations(context.Background(), name)
	if err != nil {
		return err
	}

	for _, mf := range files {
		fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
	}

	return nil
}

func Status(db *bun.DB) error {

	migrator := migrate.NewMigrator(db, Migrations)

	ms, err := migrator.MigrationsWithStatus(context.Background())
	if err != nil {
		return err
	}
	fmt.Printf("Migrations: %s\n", ms)
	fmt.Printf("unapplied Migrations: %s\n", ms.Unapplied())
	fmt.Printf("last migration group: %s\n", ms.LastGroup())

	return nil
}
func MarkApplied(db *bun.DB) error {

	migrator := migrate.NewMigrator(db, Migrations)

	group, err := migrator.Migrate(context.Background(), migrate.WithNopMigration())
	if err != nil {
		return err
	}

	if group.ID == 0 {
		fmt.Printf("there are no new Migrations to mark as applied\n")
		return nil
	}

	fmt.Printf("marked as applied %s\n", group)
	return nil
}

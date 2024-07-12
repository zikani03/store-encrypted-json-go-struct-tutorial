package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func EncryptUsingPgx(dbURL string) {

	settings := Setting{
		Key: "tenant.nndi.settings.pgx",
		Data: SettingsData{
			Version:         "v0.1.0",
			AllowNewMembers: true,
		},
	}

	db, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		panic("failed to connect database")
	}

	_, err = db.Exec(context.Background(),
		`INSERT INTO settings(key, data_json) VALUES ($1, $2);`,
		settings.Key,
		settings.Data,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println("Saved settings successfully", settings)

	existing := new(Setting)

	row := db.QueryRow(context.Background(), "SELECT id, key, data_json FROM settings WHERE key = $1", settings.Key)
	if err := row.Scan(&existing.ID, &existing.Key, &existing.Data); err != nil {
		panic(err)
	}

	fmt.Println("Loaded settings successfully", existing)
}

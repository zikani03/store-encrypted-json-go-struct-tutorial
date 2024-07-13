package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/ovh/symmecrypt/keyloader"
)

func (settings SettingsData) TextValue() (pgtype.Text, error) {
	data, err := json.Marshal(settings)
	if err != nil {
		return pgtype.Text{}, fmt.Errorf("failed to marshal data %v", err)
	}

	k, err := keyloader.LoadSingleKey()
	if err != nil {
		return pgtype.Text{}, fmt.Errorf("failed to marshal data %v", err)
	}

	encrypted, err := k.Encrypt(data)
	if err != nil {
		return pgtype.Text{}, fmt.Errorf("failed to marshal data %v", err)
	}

	return pgtype.Text{
		String: base64.StdEncoding.EncodeToString(encrypted),
		Valid:  true,
	}, nil
}

func registerType(m *pgtype.Map) {
	m.RegisterType(&pgtype.Type{
		Name:  "settings",
		OID:   pgtype.TextOID,
		Codec: pgtype.TextCodec{},
	})

	m.RegisterDefaultPgType(SettingsData{}, "settings")
	valueType := reflect.TypeOf(SettingsData{})
	m.RegisterDefaultPgType(reflect.New(valueType).Interface(), "settings")
}

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
	registerType(db.TypeMap())

	row := db.QueryRow(context.Background(),
		`INSERT INTO settings(key, data_json) VALUES ($1, $2) RETURNING id;`,
		settings.Key,
		settings.Data,
	)
	if err := row.Scan(&settings.ID); err != nil {
		panic(err)
	}

	fmt.Println("Saved settings successfully", settings)

	existing := &Setting{Data: SettingsData{}}

	row = db.QueryRow(context.Background(), "SELECT id, key, data_json FROM settings WHERE id = $1", settings.ID)
	if err := row.Scan(&existing.ID, &existing.Key, &existing.Data); err != nil {
		panic(err)
	}

	fmt.Println("Loaded settings successfully", existing)
}

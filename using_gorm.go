package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/ovh/symmecrypt/keyloader"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (es SettingsData) GormValue(ctx context.Context, db *gorm.DB) (expr clause.Expr) {
	data, err := json.Marshal(es)
	if err != nil {
		db.AddError(fmt.Errorf("failed to marshal data %v", err))
	}

	k, err := keyloader.LoadSingleKey()
	if err != nil {
		db.AddError(err)
		return
	}

	encrypted, err := k.Encrypt(data)
	if err != nil {
		db.AddError(fmt.Errorf("failed to encrypt data %v", err))
		return
	}

	return clause.Expr{SQL: "?", Vars: []interface{}{base64.StdEncoding.EncodeToString(encrypted)}}
}

func EncryptUsingGorm(dbURL string) {

	settings := Setting{
		Key: "tenant.nndi.settings.gorm",
		Data: SettingsData{
			Version:         "v0.1.0",
			AllowNewMembers: true,
		},
	}

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	tx := db.Save(&settings)
	if tx.Error != nil {
		panic(tx.Error)
	}

	fmt.Println("Saved settings successfully", settings)

	existing := new(Setting)
	tx = db.Model(Setting{}).Last(existing)
	if tx.Error != nil {
		panic(tx.Error)
	}

	fmt.Println("Loaded settings successfully", existing)
}

package main

import (
	"database/sql/driver"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/ovh/symmecrypt/keyloader"
)

type Setting struct {
	ID   int64        `json:"id" gorm:""`
	Key  string       `json:"key" gorm:"key"`
	Data SettingsData `json:"data" gorm:"column:data_json;type:text"`
}

type SettingsData struct {
	Version         string `json:"version"`
	AllowNewMembers bool   `json:"allowNewMembers"`
}

func (settings *SettingsData) Scan(value interface{}) error {
	stored, ok := value.(string)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	data, err := base64.StdEncoding.DecodeString(stored)
	if err != nil {
		return err
	}

	k, err := keyloader.LoadSingleKey()
	if err != nil {
		return err
	}
	decryptedData, err := k.Decrypt(data)
	if err != nil {
		return err
	}
	result := SettingsData{}
	err = json.Unmarshal(decryptedData, &result)
	*settings = result
	return err
}

func (settings *SettingsData) Value() (driver.Value, error) {
	data, err := json.Marshal(settings)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data %v", err)
	}

	k, err := keyloader.LoadSingleKey()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data %v", err)
	}

	encrypted, err := k.Encrypt(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data %v", err)
	}

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

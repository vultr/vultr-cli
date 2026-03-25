// Package userdata provides common functionality for server user data as well
// as printer display for outputting in instance or bare metal commands
package userdata

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
)

type UserData struct {
	Data []byte
}

func NewUserDataFromFile(path string) (*UserData, error) {
	fd, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("error reading user data file: %v", err)
	}

	ud := UserData{
		Data: fd,
	}

	return &ud, nil
}

func NewUserDataFromString(ud string) *UserData {
	return &UserData{
		Data: []byte(ud),
	}
}

func (u *UserData) Base64Encode() string {
	return base64.StdEncoding.EncodeToString(u.Data)
}

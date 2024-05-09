package handler

import (
	"errors"
	"os"
)

type ZoomMeeting struct {
	Host      string `json:"host"`
	MeetingID string `json:"meeting_id"`
	Password  string `json:"password"`
}

const (
	tokenFile string = "/config/token.json"
)

func MustSaveToken(token string) {
	err := SaveToken(token)

	if err != nil {
		panic(err)
	}
}

func SaveToken(token string) error {
	cur, err := os.Getwd()

	if err != nil {
		return err
	}

	return os.WriteFile(cur+tokenFile, []byte(token), 0644)
}

func MustLoadToken() string {
	token, err := LoadToken()
	if err != nil {
		panic(err)
	}

	return token
}

func LoadToken() (string, error) {
	cur, err := os.Getwd()
	if err != nil {
		return "", err
	}

	tokenB, err := os.ReadFile(cur + tokenFile)
	if err != nil || string(tokenB) == "" {
		return "", errors.New("token is empty")
	}

	return string(tokenB), nil
}

func MustDeleteToken() {
	err := DeleteToken()
	if err != nil {
		panic(err)
	}
}

func DeleteToken() error {
	cur, err := os.Getwd()
	if err != nil {
		return err
	}

	return os.Remove(cur + tokenFile)
}

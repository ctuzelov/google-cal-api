package utils

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"
)

func GenerateToken(authCode string, config *oauth2.Config) (*oauth2.Token, error) {
	token, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		fmt.Printf("Unable to retrieve token from web: %v\n", err)
		return &oauth2.Token{}, err
	}

	return token, nil
}

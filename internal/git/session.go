package git

import (
	"context"
	"log"

	pkggit "github.com/selyusize/my-home/pkg/git"
)

func persistGitToken(client pkggit.Client, creds *CredentialStore, provider pkggit.Provider, token string) error {
	if err := client.SetAccessToken(token); err != nil {
		return err
	}
	if creds == nil {
		return nil
	}
	return creds.Save(provider, token)
}

func clearGitToken(client pkggit.Client, creds *CredentialStore, provider pkggit.Provider) error {
	client.Logout()
	if creds == nil {
		return nil
	}
	return creds.Delete(provider)
}

func restoreGitSession(ctx context.Context, client pkggit.Client, creds *CredentialStore, provider pkggit.Provider) error {
	if creds == nil {
		return nil
	}

	token, err := creds.Load(provider)
	if err != nil {
		return err
	}
	if token == "" {
		return nil
	}

	if err := client.SetAccessToken(token); err != nil {
		return err
	}
	if _, err := client.Profile(ctx); err != nil {
		log.Printf("git: stored %s token is invalid, dropping session: %v", provider, err)
		client.Logout()
		if delErr := creds.Delete(provider); delErr != nil {
			log.Printf("git: drop %s token: %v", provider, delErr)
		}
	}
	return nil
}

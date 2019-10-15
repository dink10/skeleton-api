package auth

import (
	"context"
	"errors"

	"bitbucket.org/gismart/{{Name}}/app/models"
	"bitbucket.org/gismart/{{Name}}/services/authorisation"
)

const (
	authStorage     = models.StorageKey("auth")
)

func GetFromContext(ctx context.Context) (*models.User, error) {
	if user, ok := ctx.Value(authorisation.ContextAuthUser).(*models.User); ok {
		return user, nil
	}

	return nil, errors.New("can not cast to *models.User")
}

func getStorageFromContext(ctx context.Context) (*postgres, error) {
	if repo, ok := ctx.Value(authStorage).(postgres); ok {
		return &repo, nil
	}

	return nil, errors.New("can not cast to repository")
}

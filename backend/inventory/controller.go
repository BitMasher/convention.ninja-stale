package inventory

import (
	"context"
	"errors"
)

type Controller struct {
	Repo
}

func (c Controller) GetAssets(ctx context.Context) (interface{}, error) {
	// TODO: implement
	return nil, errors.New("not implemented")
}

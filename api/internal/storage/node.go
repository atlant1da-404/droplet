package storage

import (
	"context"
	"github.com/atlant1da-404/droplet/internal/entity"
	"github.com/atlant1da-404/droplet/internal/service"
	"github.com/atlant1da-404/droplet/pkg/database"
)

type nodeStorage struct {
	*database.PostgreSQL
}

var _ service.NodeStorage = (*nodeStorage)(nil)

func NewNodeStorage(postgresql *database.PostgreSQL) service.NodeStorage {
	return &nodeStorage{postgresql}
}

func (n nodeStorage) CreateNode(ctx context.Context, node *entity.Node) (*entity.Node, error) {
	err := n.DB.WithContext(ctx).Create(node).Error
	if err != nil {
		return nil, err
	}

	return node, nil
}

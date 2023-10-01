package service

import (
	"context"
	"fmt"
	"github.com/atlant1da-404/droplet/internal/entity"
)

type nodeService struct {
	serviceContext
}

var _ NodeService = (*nodeService)(nil)

func NewNodeService(options *Options) NodeService {
	return &nodeService{
		serviceContext: serviceContext{
			storages: options.Storages,
			config:   options.Config,
			logger:   options.Logger.Named("NodeService"),
		},
	}
}

func (n nodeService) CreateNode(ctx context.Context, options *CreateNodeOptions) (*CreateNodeOutput, error) {
	logger := n.logger.
		Named("CreateNode").
		WithContext(ctx).
		With("options", options)

	node := &entity.Node{
		SenderEmail:   options.SenderEmail,
		ReceiverEmail: options.ReceiverEmail,
	}
	logger = logger.With("node", node)

	createdNode, err := n.storages.NodeStorage.CreateNode(ctx, node)
	if err != nil {
		logger.Error("failed to create new node: %w", err)
		return nil, fmt.Errorf("failed to create new node: %w", err)
	}
	logger = logger.With("createdNode", createdNode)

	logger.Info("successfully created node")
	return &CreateNodeOutput{Id: createdNode.Id}, nil
}

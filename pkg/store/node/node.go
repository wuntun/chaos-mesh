// Copyright 2020 Chaos Mesh Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package node

import (
	"context"

	"github.com/jinzhu/gorm"

	"github.com/chaos-mesh/chaos-mesh/pkg/core"
	"github.com/chaos-mesh/chaos-mesh/pkg/store/dbstore"

	ctrl "sigs.k8s.io/controller-runtime"
)

var log = ctrl.Log.WithName("store/node")

// NewStore return a new NewStore.
func NewStore(db *dbstore.DB) core.NodeStore {
	db.AutoMigrate(&core.Node{})
	//db.AutoMigrate(&core.PodRecord{})

	ns := &nodeStore{db}

	return ns
}

type nodeStore struct {
	db *dbstore.DB
}

// List returns an event list from the datastore.
func (n *nodeStore) List(ctx context.Context) ([]*core.Node, error) {
	nodeList := make([]*core.Node, 0)

	if err := n.db.Find(&nodeList).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	return nodeList, nil
}

// Find returns an event from the datastore by ID.
func (n *nodeStore) Find(ctx context.Context, name string) (*core.Node, error) {
	nodes := make([]*core.Node, 0)
	if err := n.db.Where("name = ?", name).Find(&nodes).Error; err != nil && !gorm.IsRecordNotFoundError(err) {
		return nil, err
	}

	return nodes[0], nil
}

// Create persists a new event to the datastore.
func (n *nodeStore) Create(ctx context.Context, node *core.Node) error {
	if err := n.db.Create(node).Error; err != nil {
		return err
	}

	return nil
}

// Update persists an updated event to the datastore.
func (n *nodeStore) Update(ctx context.Context, node *core.Node) error {
	return n.db.Model(core.Node{}).
		Where("name = ?", node.Name).
		Update("config", node.Config).
		Error
}

// Delete deletes a node to the datastore.
func (n *nodeStore) Delete(ctx context.Context, node *core.Node) error {
	if err := n.db.Model(core.Node{}).
		Where("name = ? ", node.Name).
		Unscoped().
		Delete(core.Node{}).Error; err != nil {
		return err
	}

	return nil
}

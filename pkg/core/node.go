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

package core

import (
	"context"
)

// EventStore defines operations for working with event.
type NodeStore interface {
	// List returns an event list from the datastore.
	List(context.Context) ([]*Node, error)

	// Find returns an event from the datastore by ID.
	Find(context.Context, string) (*Node, error)

	// Create persists a new event to the datastore.
	Create(context.Context, *Node) error

	// Update persists an updated event to the datastore.
	Update(context.Context, *Node) error

	// Delete deletes a node to the datastore.
	Delete(context.Context, *Node) error
}

// Node represents an node instance.
type Node struct {
	Name string `gorm:"primary_key" json:"name"`

	// kind means the node's kind, the value can be k8s or physic
	Kind   string `json:"kind"`
	Config string `json:"config"`
}

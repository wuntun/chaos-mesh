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
	"fmt"
	"net/http"

	"github.com/chaos-mesh/chaos-mesh/pkg/apiserver/utils"
	"github.com/chaos-mesh/chaos-mesh/pkg/clientpool"
	"github.com/chaos-mesh/chaos-mesh/pkg/config"
	"github.com/chaos-mesh/chaos-mesh/pkg/core"

	"github.com/gin-gonic/gin"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Service defines a handler service for cluster common objects.
type Service struct {
	// this kubeCli use the local token, used for list namespace of the K8s cluster
	kubeCli client.Client
	conf    *config.ChaosDashboardConfig
	node    core.NodeStore
}

// NewService returns an experiment service instance.
func NewService(
	conf *config.ChaosDashboardConfig,
	kubeCli client.Client,
	node core.NodeStore,
) *Service {
	return &Service{
		conf:    conf,
		kubeCli: kubeCli,
		node:    node,
	}
}

// Register mounts our HTTP handler on the mux.
func Register(r *gin.RouterGroup, s *Service) {
	endpoint := r.Group("/node")

	endpoint.POST("/k8s/registry/:name", s.k8sRegistry)
	endpoint.POST("/physic/registry/:name", s.physicRegistry)

	endpoint.DELETE("/delete/:name", s.delete)

	// initial k8s client saved in store
	nodes, err := s.node.List(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, node := range nodes {
		core.Nodes[node.Name] = node
		if node.Kind == "k8s" {
			// save client into poll
			_, err = clientpool.K8sClients.KubeClient(node.Name, []byte(node.Config))
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (s *Service) delete(c *gin.Context) {
	name := c.Param("name")
	fmt.Println("delete kubeconfig", name)
	err := s.node.Delete(context.Background(), &core.Node{
		Name: name,
	})
	if err != nil {
		c.Status(http.StatusInternalServerError)
		_ = c.Error(utils.ErrInternalServer.WrapWithNoMessage(err))
		return
	}

	return
}

func (s *Service) k8sRegistry(c *gin.Context) {
	name := c.Param("name")
	configBytes, err := c.GetRawData()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		_ = c.Error(utils.ErrInternalServer.WrapWithNoMessage(err))
		return
	}
	fmt.Println("name", name, "registryKubeConfig")

	// save client into poll
	_, err = clientpool.K8sClients.KubeClient(name, configBytes)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		_ = c.Error(utils.ErrInternalServer.WrapWithNoMessage(err))
		return
	}

	node := &core.Node{
		Name:   name,
		Kind:   "k8s",
		Config: string(configBytes),
	}

	err = s.node.Create(context.Background(), node)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		_ = c.Error(utils.ErrInternalServer.WrapWithNoMessage(err))
		return
	}

	return
}

func (s *Service) physicRegistry(c *gin.Context) {
	name := c.Param("name")
	configBytes, err := c.GetRawData()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		_ = c.Error(utils.ErrInternalServer.WrapWithNoMessage(err))
		return
	}
	fmt.Println("name", name, "config", string(configBytes))

	node := &core.Node{
		Name:   name,
		Kind:   "physic",
		Config: string(configBytes),
	}

	err = s.node.Create(context.Background(), node)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		_ = c.Error(utils.ErrInternalServer.WrapWithNoMessage(err))
		return
	}

	core.Nodes[name] = node

	return
}

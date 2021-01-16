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
	"encoding/base64"
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

	endpoint.POST("/registry", s.registry)
	endpoint.DELETE("/delete/:name", s.delete)
	endpoint.GET("/list", s.list)
	endpoint.GET("/get/:name", s.get)

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

func (s *Service) list(c *gin.Context) {
	nodes, err := s.node.List(context.Background())
	if err != nil {
		c.Status(http.StatusInternalServerError)
		_ = c.Error(utils.ErrInternalServer.WrapWithNoMessage(err))
		return
	}

	c.JSON(http.StatusOK, nodes)
}

func (s *Service) get(c *gin.Context) {
	name := c.Param("name")
	node, err := s.node.Find(context.Background(), name)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		_ = c.Error(utils.ErrInternalServer.WrapWithNoMessage(err))
		return
	}

	c.JSON(http.StatusOK, node)
}

func (s *Service) registry(c *gin.Context) {
	node := &core.Node{}
	if err := c.ShouldBindJSON(node); err != nil {
		c.Status(http.StatusBadRequest)
		_ = c.Error(utils.ErrInvalidRequest.WrapWithNoMessage(err))
		return
	}

	fmt.Println("registry node", node)

	configBytes, err := base64.StdEncoding.DecodeString(node.Config)
	if err != nil {
		c.Status(http.StatusBadRequest)
		_ = c.Error(utils.ErrInvalidRequest.WrapWithNoMessage(err))
		return
	}
	node.Config = string(configBytes)

	if node.Kind == "k8s" {
		// save client into poll
		_, err = clientpool.K8sClients.KubeClient(node.Name, configBytes)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			_ = c.Error(utils.ErrInternalServer.WrapWithNoMessage(err))
			return
		}
	}

	core.Nodes[node.Name] = node

	err = s.node.Create(context.Background(), node)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		_ = c.Error(utils.ErrInternalServer.WrapWithNoMessage(err))
		return
	}

	return
}

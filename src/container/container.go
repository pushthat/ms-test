package container

import (
	"errors"
	"orchestratus/src/node"
	"sync"
)

// Container is the Container definition
type Container struct {
	ID             string            `json:"id" validate:"required"`
	Name           string            `json:"name" validate:"required"`
	Image          string            `json:"image" validate:"required"`
	SchedulerHints map[string]string `json:"scheduler-hints" validate:"required"`
	Status         string            `json:"status"`
	Node           string            `json:"node"`
}

func (c Container) doesHintMatch(node node.Node) bool {
	for key, value := range c.SchedulerHints {
		if node.SchedulerHints[key] != value {
			return false
		}
	}
	return true
}

// DeleteContainer delete a container
func (c Container) DeleteContainer() (Container, error) {
	nodes := node.GetInstance()
	containers := GetInstance()

	nodeTmp := nodes[c.Node]
	nodeTmp.Capacity.Used--
	nodes[c.Node] = nodeTmp
	delete(containers, c.ID)
	SetInstance(containers)
	node.SetInstance(nodes)
	return c, nil

}

// ScheduleContainer Schedule a container to a node or return a error
func (c Container) ScheduleContainer() (Container, error) {
	nodes := node.GetInstance()
	containers := GetInstance()

	for key := range nodes {
		if c.doesHintMatch(nodes[key]) && nodes[key].Capacity.Total > nodes[key].Capacity.Used {
			nodeTmp := nodes[key]
			nodeTmp.Capacity.Used++
			nodes[key] = nodeTmp
			c.Node = nodes[key].ID
			c.Status = "running"
			containers[c.ID] = c
			SetInstance(containers)
			node.SetInstance(nodes)
			return c, nil
		}
	}
	return c, errors.New("no node could be found to schedule the container too")

}

var instancesContainers map[string]Container
var muxContainers sync.Mutex

// GetInstance Get the cache instance
func GetInstance() map[string]Container {
	if instancesContainers == nil {
		instancesContainers = make(map[string]Container)
	}
	return instancesContainers
}

// SetInstance Set the cache instance
func SetInstance(containers map[string]Container) {
	muxContainers.Lock()
	instancesContainers = containers
	muxContainers.Unlock()

}

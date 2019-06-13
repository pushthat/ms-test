package node

import "sync"

// Capacity represent the capacity definition inside the Node
type Capacity struct {
	Total int `json:"total"`
	Used  int `json:"used"`
}

// Node is the Node definition
type Node struct {
	ID             string            `json:"id" validate:"required"`
	Name           string            `json:"name" validate:"required"`
	SchedulerHints map[string]string `json:"scheduler-hints" validate:"required"`
	Capacity       Capacity          `json:"capacity" validate:"required"`
}

var instancesNodes map[string]Node
var indexLabel map[string][]string
var muxLabel sync.Mutex
var muxNodes sync.Mutex

// GetInstance Get the node instance
func GetInstance() map[string]Node {
	if instancesNodes == nil {
		instancesNodes = map[string]Node{}
	}
	return instancesNodes
}

// SetInstance Set the node instance
func SetInstance(nodes map[string]Node) {
	muxNodes.Lock()
	instancesNodes = nodes
	muxNodes.Unlock()
}

// GetInstanceLabel Get the label index instance
func GetInstanceLabel() map[string][]string {
	if indexLabel == nil {
		indexLabel = make(map[string][]string)
	}
	return indexLabel
}

// SetInstanceLabel Set the label index instance
func SetInstanceLabel(nodes map[string][]string) {
	muxLabel.Lock()
	indexLabel = nodes
	muxLabel.Unlock()
}

// ByID is used to sort map by id
type ByID []Node

func (a ByID) Len() int           { return len(a) }
func (a ByID) Less(i, j int) bool { return a[i].ID < a[j].ID }
func (a ByID) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

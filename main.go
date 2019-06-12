package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"

	"io/ioutil"
	"orchestratus/src/container"
	"orchestratus/src/node"
	objValidator "orchestratus/src/validator"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/go-playground/validator.v9"
)

func deleteContainer(c echo.Context) error {
	containersMap := container.GetInstance()

	if _, ok := containersMap[c.Param("id")]; !ok {
		return c.JSON(http.StatusNotFound, "the container could't be found")
	}
	containersMap[c.Param("id")].DeleteContainer()
	return c.JSON(http.StatusNoContent, "")

}

func scheduleContainer(c echo.Context) error {
	containerElem := &container.Container{}
	if err := c.Bind(containerElem); err != nil {
		return c.JSON(http.StatusBadRequest, "bad input parameter: "+err.Error())
	}
	if err := c.Validate(containerElem); err != nil {
		return c.JSON(http.StatusBadRequest, "bad input parameter: "+err.Error())
	}

	containersMap := container.GetInstance()

	if _, ok := containersMap[containerElem.ID]; ok {
		return c.JSON(http.StatusNotFound, "bad input parameter: container id already used")
	}
	containerScheduled, err := containerElem.ScheduleContainer()
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}
	return c.JSON(http.StatusOK, containerScheduled)
}

func listNodes(c echo.Context) error {
	nodeList := make([]node.Node, 0)
	nodeMap := node.GetInstance()

	for _, value := range nodeMap {
		nodeList = append(nodeList, value)
	}
	sort.Sort(node.ByID(nodeList))
	return c.JSON(http.StatusOK, nodeList)
}

func getNode(c echo.Context) error {
	nodeMap := node.GetInstance()
	return c.JSON(http.StatusOK, nodeMap[c.Param("id")])
}

func parseNode() error {
	nodeFile := os.Getenv("NODE_FILE")
	content, err := ioutil.ReadFile(nodeFile)
	if err != nil {
		return (err)
	}

	var dat []node.Node

	if err := json.Unmarshal(content, &dat); err != nil {
		return (err)
	}

	valid := &objValidator.NodeRequestValidator{Validator: validator.New()}

	nodes := make(map[string]node.Node)

	for _, elem := range dat {
		validationErrors := valid.Validate(elem)
		if validationErrors != nil {
			log.Println(validationErrors.Error())
		} else {
			nodes[elem.ID] = elem
		}
	}
	node.SetInstance(nodes)
	return nil
}

func main() {
	// Parsing nodes
	err := parseNode()
	if err != nil {
		log.Fatal("invalid json file :" + err.Error())
	}

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	gContainer := e.Group("")
	e.Validator = &objValidator.ContainerRequestValidator{Validator: validator.New()}
	gContainer.POST("/container/schedule", scheduleContainer)
	gContainer.DELETE("/container/:id", deleteContainer)

	gNode := e.Group("")
	gNode.GET("/nodes", listNodes)
	gNode.GET("/nodes/:id", getNode)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

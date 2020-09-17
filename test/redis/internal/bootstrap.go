package main

import (
	"fmt"
	"os"

	"github.com/gomodule/redigo/redis"
	rg "github.com/redislabs/redisgraph-go"
)

func main() {
	conn, _ := redis.Dial("tcp", ":6379")
	defer conn.Close()

	graph := rg.GraphNew("gbac", conn)
	graph.Delete()

	read := MakeAction("read", &graph)
	write := MakeAction("write", &graph)
	admin := MakeAction("admin", &graph)
	AddActionChild(admin, write, &graph)
	AddActionChild(write, read, &graph)

	userA := MakeUser("a", &graph)
	userB := MakeUser("b", &graph)
	userC := MakeUser("c", &graph)
	userD := MakeUser("d", &graph)
	userE := MakeUser("e", &graph)

	defaultContainer := rg.Node{
		Label: "container",
		Properties: map[string]interface{}{
			"id": "default",
		},
	}
	graph.AddNode(&defaultContainer)

	subContainer := rg.Node{
		Label: "container",
		Properties: map[string]interface{}{
			"id": "sub",
		},
	}
	graph.AddNode(&subContainer)
	graph.AddEdge(&rg.Edge{
		Source:      &subContainer,
		Relation:    "parent",
		Destination: &defaultContainer,
	})

	widgetA := MakeWidget("a", &defaultContainer, &graph)
	MakeWidget("b", &defaultContainer, &graph)
	MakeWidget("c", &subContainer, &graph)

	AddGeneralPerm(userA, admin, &defaultContainer, &graph)
	AddGeneralPerm(userB, write, &defaultContainer, &graph)
	AddGeneralPerm(userC, read, &defaultContainer, &graph)

	AddDirectPerm(userD, write, widgetA, &graph)

	AddResourcePerm(userE, read, "widget", &defaultContainer, &graph)

	_, err := graph.Commit()
	if err != nil {
		fmt.Println("Commit failed", err)
		os.Exit(1)
	}
	fmt.Println("Graph bootstrapped")
}

func MakeUser(id string, graph *rg.Graph) *rg.Node {
	node := &rg.Node{
		Label: "user",
		Properties: map[string]interface{}{
			"id": id,
		},
	}
	graph.AddNode(node)
	return node
}

func MakeWidget(id string, parent *rg.Node, graph *rg.Graph) *rg.Node {
	node := &rg.Node{
		Label: "widget",
		Properties: map[string]interface{}{
			"id": id,
		},
	}
	graph.AddNode(node)
	graph.AddEdge(&rg.Edge{
		Source:      node,
		Relation:    "parent",
		Destination: parent,
	})
	return node
}

func AddGeneralPerm(src *rg.Node, action *rg.Node, dest *rg.Node, graph *rg.Graph) {
	graph.AddEdge(&rg.Edge{
		Source:      src,
		Relation:    "all",
		Destination: dest,
		Properties: map[string]interface{}{
			"action": action.GetProperty("name").(string),
		},
	})
}

func AddDirectPerm(src *rg.Node, action *rg.Node, dest *rg.Node, graph *rg.Graph) {
	fmt.Println("Action", action.GetProperty("name").(string))
	graph.AddEdge(&rg.Edge{
		Source:      src,
		Relation:    dest.Label,
		Destination: dest,
		Properties: map[string]interface{}{
			"action": action.GetProperty("name").(string),
		},
	})
}

func AddResourcePerm(src *rg.Node, action *rg.Node, resource string, dest *rg.Node, graph *rg.Graph) {
	graph.AddEdge(&rg.Edge{
		Source:      src,
		Relation:    fmt.Sprintf(resource),
		Destination: dest,
		Properties: map[string]interface{}{
			"action": action.GetProperty("name").(string),
		},
	})
}

func MakeAction(action string, graph *rg.Graph) *rg.Node {
	node := &rg.Node{
		Label: "action",
		Properties: map[string]interface{}{
			"name": action,
		},
	}
	graph.AddNode(node)
	return node
}

func AddActionChild(parent *rg.Node, child *rg.Node, graph *rg.Graph) {
	graph.AddEdge(&rg.Edge{
		Source:      parent,
		Relation:    "can",
		Destination: child,
	})
}

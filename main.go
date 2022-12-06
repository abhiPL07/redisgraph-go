package main

import (
	"github.com/gomodule/redigo/redis"
	rg "github.com/redislabs/redisgraph-go"
)

type Node struct {
	Subject string
	Marks   map[string]interface{}
}

func NodeMake() (*Node, *Node) {
	m1 := make(map[string]interface{})
	m1["semOne"] = 93
	m1["semTwo"] = 95
	node1 := Node{
		Subject: "english",
		Marks:   m1,
	}
	m2 := make(map[string]interface{})
	m2["semOne"] = 91
	m2["semTwo"] = 96
	node2 := Node{
		Subject: "mathematics",
		Marks:   m2,
	}
	return &node1, &node2
}

func main() {
	conn, _ := redis.Dial("tcp", "0.0.0.0:6379")
	defer conn.Close()
	graph := rg.GraphNew("classTwo", conn)

	e, m := NodeMake()
	eng := rg.Node{
		Label:      e.Subject,
		Properties: e.Marks,
	}
	graph.AddNode(&eng)

	math := rg.Node{
		Label:      m.Subject,
		Properties: m.Marks,
	}
	graph.AddNode(&math)

	edge := rg.Edge{
		Source:      &eng,
		Relation:    "passed",
		Destination: &math,
		Properties: map[string]interface{}{
			"success": "Yes",
		},
	}
	graph.AddEdge(&edge)

	graph.Commit()

	query := `MATCH (e:english)-[p:passed]->(m:mathematics)
	       RETURN e.semOne, e.semTwo, p.success, m.semOne, m.semTwo`
	rs, _ := graph.Query(query)
	rs.PrettyPrint()
}

package main

import (
	"sync"

	"github.com/gomodule/redigo/redis"
	rg "github.com/redislabs/redisgraph-go"
)

type Node struct {
	Subject string
	Marks   map[string]interface{}
}

var wg sync.WaitGroup
var e *Node
var m *Node

func NodeMake(in string) {
	if in == "e" {
		m1 := map[string]interface{}{
			"semOne": 93,
			"semTwo": 95,
		}
		e = &Node{
			Subject: "english",
			Marks:   m1,
		}
	} else {
		m2 := map[string]interface{}{
			"semOne": 91,
			"semTwo": 96,
		}
		m = &Node{
			Subject: "mathematics",
			Marks:   m2,
		}
	}
	wg.Done()
}

func main() {
	conn, _ := redis.Dial("tcp", "0.0.0.0:6379")
	defer conn.Close()
	graph := rg.GraphNew("classConc", conn)

	wg.Add(2)
	go NodeMake("e")
	go NodeMake("")
	wg.Wait()
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
	wg.Wait()

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
	wg.Wait()
}

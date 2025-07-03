/*
Copyright 2013 Alessandro Frossi

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dijkstra

import (
	"testing"
)

var (
	graph *testGraph
)

func init() {
	nodes := map[string]interface{}{
		"A": struct{}{},
		"B": struct{}{},
		"C": struct{}{},
		"D": struct{}{},
		"E": struct{}{},
		"F": struct{}{},
		"G": struct{}{},
		"S": struct{}{},
		"T": struct{}{},
		"U": struct{}{},
	}
	edges := map[string]map[string]interface{}{
		"A": {
			"B": struct{}{},
			"C": struct{}{},
		},
		"B": {
			"C": struct{}{},
			"D": struct{}{},
		},
		"C": {
			"E": struct{}{},
			"G": struct{}{},
		},
		"D": {
			"C": struct{}{},
		},
		"E": {
			"F": struct{}{},
		},
		"F": {
			"G": struct{}{},
		},
		"G": {
			"T": struct{}{},
		},
		"S": {
			"A": struct{}{},
			// "B": struct{}{},
		},
	}
	graph = newTestGraph()
	graph.nodes = nodes
	// revEdges := make(map[string]map[string]interface{})
	graph.edges = edges
	for k1, v := range edges {
		for k2, _ := range v {
			if _, ok := graph.reverseEdges[k2]; !ok {
				graph.reverseEdges[k2] = make(map[string]interface{})
			}
			graph.reverseEdges[k2][k1] = struct{}{}
		}
	}
}

func TestSinglePath(t *testing.T) {
	path, valid := SearchPath(graph, "S", "T", BIDIR)
	if !valid {
		t.Fatal("Validity error.")
	}
	if len(path.Path) <= 2 {
		t.Fatal("Path length error.")
	}
	expPath := []string{"S", "A", "C", "G", "T"}
	for i, v := range path.Path {
		if v.Node != expPath[i] {
			t.Fatal("Wrong path.")
		}
	}
	if int(path.Weight) != len(expPath)-1 {
		t.Logf("Path: %#v\n", path)
		t.Fatalf("Wrong path weight:\nExpected: %d\nGot: %f\n", len(expPath)-1, path.Weight)
	}
}

func TestShortcut(t *testing.T) {
	// adding fake edge
	bk := graph.edges["C"]
	graph.edges["C"]["T"] = struct{}{}

	revBk := graph.reverseEdges["T"]
	graph.reverseEdges["T"]["C"] = struct{}{}

	path, valid := SearchPath(graph, "S", "T", BIDIR)
	if !valid {
		t.Fatal("Validity error.")
	}
	if len(path.Path) <= 2 {
		t.Fatal("Path length error.")
	}
	expPath := []string{"S", "A", "C", "T"}
	for i, v := range path.Path {
		if v.Node != expPath[i] {
			t.Fatal("Wrong path.")
		}
	}
	if int(path.Weight) != len(expPath)-1 {
		t.Logf("Path: %#v\n", path)
		t.Fatalf("Wrong path weight:\nExpected: %d\nGot: %f\n", len(expPath)-1, path.Weight)
	}

	graph.edges["C"] = bk
	graph.reverseEdges["T"] = revBk
}

func TestUnreachable(t *testing.T) {
	_, valid1 := SearchPath(graph, "S", "U", VANILLA)
	_, valid2 := SearchPath(graph, "S", "U", BIDIR)

	if valid1 || valid2 {
		t.Fatalf("A path was found in an unconnected graph.")
	}
}

func TestEquality(t *testing.T) {
	path1, valid1 := SearchPath(graph, "S", "T", VANILLA)
	path2, valid2 := SearchPath(graph, "S", "T", BIDIR)

	if !valid1 || !valid2 {
		t.Fatal("A path search failed.")
	}
	if !path1.IsEqual(path2) {
		t.Fatal("The algorithms yield different paths.")
	}
}

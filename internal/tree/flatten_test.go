package tree

import "testing"

func TestFlattenSingleRootObject(t *testing.T) {
	data := []any{
		map[string]any{
			"name": "Ana",
			"age":  30.0,
			"tags": []any{"go"},
		},
	}

	nodes := Flatten(data)
	if len(nodes) != 7 {
		t.Fatalf("expected 7 nodes, got %d", len(nodes))
	}

	if nodes[0].Type != ObjectOpen || nodes[0].Path != "." {
		t.Fatalf("unexpected root node: %+v", nodes[0])
	}

	if nodes[1].Type != KeyValue || nodes[1].Path != ".age" || nodes[1].Key != "age" {
		t.Fatalf("unexpected age node: %+v", nodes[1])
	}

	if nodes[2].Type != KeyValue || nodes[2].Path != ".name" || nodes[2].Key != "name" {
		t.Fatalf("unexpected name node: %+v", nodes[2])
	}

	if nodes[3].Type != ArrayOpen || nodes[3].Path != ".tags" || nodes[3].Key != "tags" {
		t.Fatalf("unexpected tags open node: %+v", nodes[3])
	}

	if nodes[4].Type != ArrayElement || nodes[4].Path != ".tags[0]" {
		t.Fatalf("unexpected array element node: %+v", nodes[4])
	}

	if nodes[5].Type != ArrayClose || nodes[5].Path != ".tags" {
		t.Fatalf("unexpected tags close node: %+v", nodes[5])
	}

	if nodes[6].Type != ObjectClose || nodes[6].Path != "." {
		t.Fatalf("unexpected object close node: %+v", nodes[6])
	}
}

func TestFlattenMultipleTopLevelValues(t *testing.T) {
	data := []any{"a", 2.0}
	nodes := Flatten(data)

	if len(nodes) != 4 {
		t.Fatalf("expected 4 nodes, got %d", len(nodes))
	}

	if nodes[0].Type != ArrayOpen || nodes[0].Path != "." {
		t.Fatalf("unexpected synthetic root open node: %+v", nodes[0])
	}

	if nodes[1].Type != ArrayElement || nodes[1].Path != ".[0]" || nodes[1].IsLast {
		t.Fatalf("unexpected first element node: %+v", nodes[1])
	}

	if nodes[2].Type != ArrayElement || nodes[2].Path != ".[1]" || !nodes[2].IsLast {
		t.Fatalf("unexpected second element node: %+v", nodes[2])
	}

	if nodes[3].Type != ArrayClose || nodes[3].Path != "." {
		t.Fatalf("unexpected synthetic root close node: %+v", nodes[3])
	}
}

func TestFlattenUsesBracketNotationForNonIdentifiers(t *testing.T) {
	data := []any{
		map[string]any{
			"weird-key": map[string]any{
				"x y": true,
			},
		},
	}

	nodes := Flatten(data)
	var found bool
	for _, node := range nodes {
		if node.Type != KeyValue {
			continue
		}

		if node.Path == `.["weird-key"]["x y"]` {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("expected bracket notation path to be present")
	}
}

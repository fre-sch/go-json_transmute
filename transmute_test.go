package transmute

import (
	"reflect"
	"testing"
)

func TestStringNormal(t *testing.T) {
	expr := "string value"
	context := map[string]any(nil)
	result, err := Transmute(expr, context)
	if err != nil {
		t.Fatalf("failed transmute with: %#+v", err)
	}
	expected := "string value"
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("\nexpected: %#+v\nreceived %#+v\n", expected, result)
	}
}

func TestStringPath(t *testing.T) {
	expr := "$.nested.context.key"
	context := map[string]any{
		"nested": map[string]any{
			"context": map[string]any{
				"key": "expected value",
			},
		},
	}
	result, err := Transmute(expr, context)
	if err != nil {
		t.Fatalf("failed transmute with: %#+v", err)
	}
	expected := "expected value"
	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("\nexpected: %#+v\nreceived %#+v\n", expected, result)
	}
}

func TestFormat(t *testing.T) {
	expr := map[string]any{
		"#format": `Hello {{index . "key"}}!`,
	}
	context := map[string]string{
		"key": "world",
	}

	result, err := Transmute(expr, context)

	if err != nil {
		t.Fatalf("failed transmute with: %#+v", err)
	}

	expected := "Hello world!"

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("\nexpected: %#+v\nreceived %#+v\n", expected, result)
	}
}

func TestFormatNested(t *testing.T) {
	expr := map[string]any{
		"first": map[string]any{
			"#format": `Hello {{index . "key"}}!`,
		},
		"second": map[string]any{
			"nested": map[string]any{
				"#format": `Hello {{index . "key"}}!`,
			},
		},
		"third": []any{
			"untouched",
			map[string]any{
				"#format": `Hello {{index . "key"}}!`,
			},
		},
	}
	context := map[string]string{
		"key": "world",
	}

	result, err := Transmute(expr, context)

	if err != nil {
		t.Fatalf("failed transmute with: %#+v", err)
	}

	expected := map[string]any{
		"first": "Hello world!",
		"second": map[string]any{
			"nested": "Hello world!",
		},
		"third": []any{
			"untouched",
			"Hello world!",
		},
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("\nexpected: %#+v\nreceived %#+v\n", expected, result)
	}
}

func TestEachSimple(t *testing.T) {
	expr := map[string]any{
		"#each": []any{
			"one", "two", "three",
		},
		"key": "value",
	}
	context := any(nil)
	result, err := Transmute(expr, context)

	if err != nil {
		t.Fatalf("failed transmute with: %#+v", err)
	}

	expected := []any{
		map[string]any{
			"key": "value",
		},
		map[string]any{
			"key": "value",
		},
		map[string]any{
			"key": "value",
		},
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("\nexpected: %#+v\nreceived %#+v\n", expected, result)
	}
}

func TestEachItem(t *testing.T) {
	expr := map[string]any{
		"#each": []any{
			"one", "two", "three",
		},
		"key": "$.it",
	}
	context := any(nil)
	result, err := Transmute(expr, context)

	if err != nil {
		t.Fatalf("failed transmute with: %#+v", err)
	}

	expected := []any{
		map[string]any{
			"key": "one",
		},
		map[string]any{
			"key": "two",
		},
		map[string]any{
			"key": "three",
		},
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("\nexpected: %#+v\nreceived %#+v\n", expected, result)
	}
}

func TestEachItemFormat(t *testing.T) {
	expr := map[string]any{
		"#each": []any{
			"one", "two", "three",
		},
		"key": map[string]any{
			"#format": `nested {{path "$.it"}}`,
		},
	}
	context := any(nil)
	result, err := Transmute(expr, context)

	if err != nil {
		t.Fatalf("failed transmute with: %#+v", err)
	}

	expected := []any{
		map[string]any{
			"key": "nested one",
		},
		map[string]any{
			"key": "nested two",
		},
		map[string]any{
			"key": "nested three",
		},
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("\nexpected: %#+v\nreceived %#+v\n", expected, result)
	}
}

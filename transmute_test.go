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

func TestMapSimple(t *testing.T) {
	expr := map[string]any{
		"#map": []any{
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

func TestMapItem(t *testing.T) {
	expr := map[string]any{
		"#map": []any{
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

func TestMapItemAndParent(t *testing.T) {
	expr := map[string]any{
		"#map":  "$.tags",
		"title": "$.it",
		"price": "$.parent.price",
	}
	context := map[string]any{
		"tags": []any{
			"one", "two", "three",
		},
		"price": 1337,
	}
	result, err := Transmute(expr, context)

	if err != nil {
		t.Fatalf("failed transmute with: %#+v", err)
	}

	expected := []any{
		map[string]any{
			"title": "one",
			"price": 1337,
		},
		map[string]any{
			"title": "two",
			"price": 1337,
		},
		map[string]any{
			"title": "three",
			"price": 1337,
		},
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("\nexpected: %#+v\nreceived %#+v\n", expected, result)
	}
}

func TestMapItemFormat(t *testing.T) {
	expr := map[string]any{
		"#map": []any{
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

func TestFirstSlice(t *testing.T) {
	expr := map[string]any{
		"#first": []any{
			"one", "two", "three",
		},
	}
	context := any(nil)
	result, err := Transmute(expr, context)

	if err != nil {
		t.Fatalf("failed transmute with: %#+v", err)
	}

	var expected any = "one"

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("\nexpected: %#+v\nreceived %#+v\n", expected, result)
	}
}

func TestFirstPath(t *testing.T) {
	expr := map[string]any{
		"#first": "$.items",
	}
	context := map[string]any{
		"items": []any{
			"one", "two", "three",
		},
	}
	result, err := Transmute(expr, context)

	if err != nil {
		t.Fatalf("failed transmute with: %#+v", err)
	}

	var expected any = "one"

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("\nexpected: %#+v\nreceived %#+v\n", expected, result)
	}
}

func TestFirstString(t *testing.T) {
	expr := map[string]any{
		"#first": "not a slice",
	}
	context := any(nil)
	result, err := Transmute(expr, context)

	if err != nil {
		t.Fatalf("failed transmute with: %#+v", err)
	}

	var expected any = "not a slice"

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("\nexpected: %#+v\nreceived %#+v\n", expected, result)
	}
}

func TestFirstMap(t *testing.T) {
	expr := map[string]any{
		"#first": map[string]any{
			"not a slice": "not a slice",
		},
	}
	context := any(nil)
	result, err := Transmute(expr, context)

	if err != nil {
		t.Fatalf("failed transmute with: %#+v", err)
	}

	var expected any = map[string]any{
		"not a slice": "not a slice",
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("\nexpected: %#+v\nreceived %#+v\n", expected, result)
	}
}

func TestTransmute(t *testing.T) {
	expr := map[string]any{
		"#map": map[string]any{
			"#transmute": "$.var",
		},
		"label": "$.it.title",
	}
	context := map[string]any{
		"var": "$.items",
		"items": []any{
			map[string]any{
				"title": "one",
			},
			map[string]any{
				"title": "two",
			},
		},
	}
	result, err := Transmute(expr, context)

	if err != nil {
		t.Fatalf("failed transmute with: %#+v", err)
	}

	var expected any = []any{
		map[string]any{
			"label": "one",
		},
		map[string]any{
			"label": "two",
		},
	}

	if !reflect.DeepEqual(expected, result) {
		t.Fatalf("\nexpected: %#+v\nreceived %#+v\n", expected, result)
	}
}

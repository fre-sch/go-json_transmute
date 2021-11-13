package transmute

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	"github.com/oliveagle/jsonpath"
)

type any = interface{}

const OperatorFormat = "#format"
const OperatorMap = "#map"
const OperatorFirst = "#first"
const OperatorTransmute = "#transmute"

func Transmute(value any, context any) (any, error) {
	switch valueTyped := value.(type) {
	case string:
		return transmuteString(valueTyped, context)
	case map[string]any:
		return transmuteMap(valueTyped, context)
	case []any:
		return transmuteSlice(valueTyped, context)
	default:
		return value, nil
	}
}

func transmuteString(value string, context any) (result any, err error) {
	if result, err = jsonpath.JsonPathLookup(context, value); err == nil {
		return result, nil
	}
	return value, nil
}

func transmuteSlice(value []any, context any) (result []any, err error) {
	for _, item := range value {
		var resultItem any
		if resultItem, err = Transmute(item, context); err != nil {
			return
		}
		result = append(result, resultItem)
	}
	return
}

func transmuteMap(value map[string]any, context any) (any, error) {
	if _, ok := value[OperatorFormat]; ok {
		return transmuteOpFormat(value, context)
	}
	if _, ok := value[OperatorMap]; ok {
		return transmuteOpMap(value, context)
	}
	if _, ok := value[OperatorFirst]; ok {
		return transmuteOpFirst(value, context)
	}
	if item, ok := value[OperatorTransmute]; ok {
		res, err := Transmute(item, context)
		if err == nil {
			return Transmute(res, context)
		}
		return nil, err
	}

	return transmuteMapItems(value, context)
}

func transmuteMapItems(value map[string]any, context any) (result map[string]any, err error) {
	result = make(map[string]any)
	for key, item := range value {
		var itemTransmuted any
		itemTransmuted, err = Transmute(item, context)
		if err != nil {
			return
		}
		result[key] = itemTransmuted
	}
	return
}

func newTplJsonPathLookup(context any) func(string) any {
	return func(path string) (result any) {
		if result, err := jsonpath.JsonPathLookup(context, path); err == nil {
			return result
		}
		return nil
	}
}

func transmuteOpFormat(value map[string]any, context any) (result string, err error) {
	switch formatString := value[OperatorFormat].(type) {
	case string:
		buf := &bytes.Buffer{}
		var tpl *template.Template
		if tpl, err = template.New("").Funcs(template.FuncMap{
			"path": newTplJsonPathLookup(context),
		}).Parse(formatString); err != nil {
			return
		}
		if err = tpl.Execute(buf, context); err != nil {
			return
		}
		return buf.String(), nil
	default:
		err = errors.New(`value for key "#format" must be string`)
		return
	}
}

func transmuteOpMap(value map[string]any, context any) (result []any, err error) {
	var mapTransmuted any
	var mapAsSlice []any
	var ok bool

	if mapTransmuted, err = Transmute(value[OperatorMap], context); err != nil {
		return
	}
	if mapAsSlice, ok = mapTransmuted.([]any); !ok {
		err = errors.New(fmt.Sprintf(
			"#map expected to evaluate to []interface{}, actual %#+v",
			mapTransmuted))
		return
	}

	rest := restMap(value, OperatorMap)

	for _, item := range mapAsSlice {
		itemContext := map[string]any{
			"parent": context,
			"it":     item,
		}
		var resultIt any
		if resultIt, err = Transmute(rest, itemContext); err != nil {
			return
		}
		result = append(result, resultIt)
	}
	return
}

func transmuteOpFirst(value map[string]any, context any) (result any, err error) {
	var firstTransmuted any
	var firstSlice []any
	var ok bool

	if firstTransmuted, err = Transmute(value[OperatorFirst], context); err != nil {
		return
	}

	if firstSlice, ok = firstTransmuted.([]any); ok {
		for _, val := range firstSlice {
			result = val
			return
		}
	}
	result = firstTransmuted
	return
}

func restMap(value map[string]any, omitKeys ...string) (result map[string]any) {
	result = make(map[string]any)
	for k, v := range value {
		result[k] = v
	}
	for _, key := range omitKeys {
		delete(result, key)
	}
	return
}

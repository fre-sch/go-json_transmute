package transmute

import (
	"bytes"
	"errors"
	"fmt"
	"math/big"
	"text/template"

	"github.com/oliveagle/jsonpath"
)

type any = interface{}

const OperatorExtend = "#extend"
const OperatorFirst = "#first"
const OperatorFormat = "#format"
const OperatorMap = "#map"
const OperatorSum = "#sum"
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
	if item, ok := value[OperatorTransmute]; ok {
		res, err := Transmute(item, context)
		if err == nil {
			return Transmute(res, context)
		}
		return nil, err
	}
	if _, ok := value[OperatorExtend]; ok {
		return transmuteOpExtend(value, context)
	}
	if _, ok := value[OperatorFirst]; ok {
		return transmuteOpFirst(value, context)
	}
	if _, ok := value[OperatorFormat]; ok {
		return transmuteOpFormat(value, context)
	}
	if _, ok := value[OperatorMap]; ok {
		return transmuteOpMap(value, context)
	}
	if _, ok := value[OperatorSum]; ok {
		return transmuteOpSum(value, context)
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

func transmuteOpSum(value map[string]any, context any) (result float64, err error) {
	var sumTransmuted any
	var sumSlice []any
	var ok bool

	sumTransmuted, err = Transmute(value[OperatorSum], context)
	if err != nil {
		return
	}

	if sumSlice, ok = sumTransmuted.([]any); !ok {
		err = errors.New(fmt.Sprintf(
			"#sum expected to evaluate to []interface{}, actual %#+v",
			sumTransmuted))
		return
	}

	result = sumValues(sumSlice...)
	return result, nil
}

func sumValues(values ...any) float64 {
	var x, y big.Float

	for _, v := range values {
		if _, _, err := x.Parse(fmt.Sprint(v), 10); err == nil {
			y.Add(&y, &x)
		}
	}
	t, _ := y.Float64()
	return t
}

func transmuteOpExtend(value map[string]any, context any) (result any, err error) {
	var baseAny any
	var resultMap map[string]any
	var ok bool

	if baseAny, err = Transmute(value[OperatorExtend], context); err != nil {
		return
	}
	if resultMap, ok = baseAny.(map[string]any); !ok {
		return baseAny, nil
	}

	restMap := restMap(value, OperatorExtend)
	var currentValue any
	for key, value := range restMap {
		if currentValue, err = Transmute(value, context); err == nil {
			resultMap[key] = currentValue
		}
	}

	return resultMap, nil
}

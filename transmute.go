package transmute

import (
	"bytes"
	"errors"
	"fmt"
	"text/template"

	"github.com/oliveagle/jsonpath"
)

type any = interface{}

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
	if _, ok := value["#format"]; ok {
		return transmuteFormat(value, context)
	}
	if _, ok := value["#each"]; ok {
		return transmuteEach(value, context)
	}

	resultMap := make(map[string]any)
	for key, item := range value {
		itemResult, err := Transmute(item, context)
		if err != nil {
			return nil, err
		}
		resultMap[key] = itemResult
	}

	return resultMap, nil
}

func newTplJsonPathLookup(context any) func(string) any {
	return func(path string) (result any) {
		if result, err := jsonpath.JsonPathLookup(context, path); err == nil {
			return result
		}
		return nil
	}
}

func transmuteFormat(value map[string]any, context any) (result string, err error) {
	switch formatString := value["#format"].(type) {
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

func transmuteEach(value map[string]any, context any) (result []any, err error) {
	var exprAny any
	var exprSlice []any
	var ok bool

	if exprAny, err = Transmute(value["#each"], context); err != nil {
		return
	}
	if exprSlice, ok = exprAny.([]any); !ok {
		err = errors.New(fmt.Sprintf(
			"#each expected to evaluate to []interface{}, actual %#+v",
			exprAny))
		return
	}

	rest := restMap(value, "#each")

	for _, item := range exprSlice {
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

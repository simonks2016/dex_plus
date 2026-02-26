package payload

import (
	"fmt"
	"reflect"
	"strconv"
	"sync"

	"github.com/goccy/go-json"
)

type fieldInfo struct {
	index []int
	kind  reflect.Kind
	typ   reflect.Type
}

type typeInfo struct {
	fields map[string]fieldInfo // key = binance tag
}

var typeCache sync.Map // map[reflect.Type]*typeInfo

func getTypeInfo(t reflect.Type) *typeInfo {
	// 只支持 struct 或 *struct
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return &typeInfo{fields: map[string]fieldInfo{}}
	}

	if v, ok := typeCache.Load(t); ok {
		return v.(*typeInfo)
	}

	ti := &typeInfo{fields: make(map[string]fieldInfo, t.NumField())}
	for i := 0; i < t.NumField(); i++ {
		sf := t.Field(i)

		// 跳过未导出字段
		if sf.PkgPath != "" {
			continue
		}

		tag := sf.Tag.Get("binance")
		if tag == "" || tag == "-" {
			continue
		}

		ti.fields[tag] = fieldInfo{
			index: sf.Index,
			kind:  sf.Type.Kind(),
			typ:   sf.Type,
		}
	}

	typeCache.Store(t, ti)
	return ti
}

func setSlice(fv reflect.Value, destTyp reflect.Type, raw any) error {
	rawVal := reflect.ValueOf(raw)
	if rawVal.Kind() != reflect.Slice && rawVal.Kind() != reflect.Array {
		return fmt.Errorf("raw data is not a slice: %T", raw)
	}

	n := rawVal.Len()
	// 创建对应类型的 Slice: 例如 []string 或 [][]string
	elemTyp := destTyp.Elem()
	slice := reflect.MakeSlice(destTyp, n, n)

	for i := 0; i < n; i++ {
		item := rawVal.Index(i).Interface()
		targetElem := slice.Index(i)

		// 递归处理元素
		// 这里构造一个临时的 fieldInfo 来复用 setValueFast
		itemFi := fieldInfo{
			kind: elemTyp.Kind(),
			typ:  elemTyp,
		}

		if err := setValueFast(targetElem, itemFi, item); err != nil {
			return fmt.Errorf("slice index %d: %w", i, err)
		}
	}

	fv.Set(slice)
	return nil
}

func setValueFast(fv reflect.Value, fi fieldInfo, raw any) error {
	// 1. 处理切片/数组类型
	if fi.kind == reflect.Slice {
		return setSlice(fv, fi.typ, raw)
	}

	// 2. 原有的基础类型逻辑
	switch fi.kind {
	case reflect.String:
		s, err := toString(raw)
		if err == nil {
			fv.SetString(s)
		}
		return err
	case reflect.Bool:
		b, err := toBool(raw)
		if err == nil {
			fv.SetBool(b)
		}
		return err
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := toInt64(raw)
		if err == nil {
			fv.SetInt(n)
		}
		return err
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n, err := toUint64(raw)
		if err == nil {
			fv.SetUint(n)
		}
		return err
	case reflect.Float32, reflect.Float64:
		f, err := toFloat64(raw)
		if err == nil {
			fv.SetFloat(f)
		}
		return err
	}

	// 3. 兜底逻辑
	rvv := reflect.ValueOf(raw)
	if rvv.IsValid() && rvv.Type().AssignableTo(fi.typ) {
		fv.Set(rvv)
		return nil
	}
	return fmt.Errorf("unsupported kind %v for raw %T", fi.kind, raw)
}

func toString(v any) (string, error) {
	switch x := v.(type) {
	case string:
		return x, nil
	case []byte:
		return string(x), nil
	case json.Number:
		return x.String(), nil
	case float64:
		// JSON 默认 number -> float64
		return strconv.FormatFloat(x, 'f', -1, 64), nil
	case float32:
		return strconv.FormatFloat(float64(x), 'f', -1, 32), nil
	case int:
		return strconv.Itoa(x), nil
	case int64:
		return strconv.FormatInt(x, 10), nil
	case uint64:
		return strconv.FormatUint(x, 10), nil
	case bool:
		if x {
			return "true", nil
		}
		return "false", nil
	default:
		return "", fmt.Errorf("cannot convert %T to string", v)
	}
}

func toBool(v any) (bool, error) {
	switch x := v.(type) {
	case bool:
		return x, nil
	case string:
		// "true"/"false"/"1"/"0"
		if x == "1" || x == "true" || x == "TRUE" {
			return true, nil
		}
		if x == "0" || x == "false" || x == "FALSE" {
			return false, nil
		}
		return false, fmt.Errorf("invalid bool string: %q", x)
	case float64:
		return x != 0, nil
	case int:
		return x != 0, nil
	case int64:
		return x != 0, nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", v)
	}
}

func toInt64(v any) (int64, error) {
	switch x := v.(type) {
	case int:
		return int64(x), nil
	case int64:
		return x, nil
	case float64:
		return int64(x), nil
	case json.Number:
		return x.Int64()
	case string:
		// 允许 "123" / "123.45"（后者会截断）
		if i, err := strconv.ParseInt(x, 10, 64); err == nil {
			return i, nil
		}
		f, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return 0, err
		}
		return int64(f), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", v)
	}
}

func toUint64(v any) (uint64, error) {
	switch x := v.(type) {
	case uint64:
		return x, nil
	case int:
		if x < 0 {
			return 0, fmt.Errorf("negative int %d", x)
		}
		return uint64(x), nil
	case int64:
		if x < 0 {
			return 0, fmt.Errorf("negative int64 %d", x)
		}
		return uint64(x), nil
	case float64:
		if x < 0 {
			return 0, fmt.Errorf("negative float %f", x)
		}
		return uint64(x), nil
	case json.Number:
		i, err := x.Int64()
		if err != nil {
			return 0, err
		}
		if i < 0 {
			return 0, fmt.Errorf("negative number %d", i)
		}
		return uint64(i), nil
	case string:
		if u, err := strconv.ParseUint(x, 10, 64); err == nil {
			return u, nil
		}
		f, err := strconv.ParseFloat(x, 64)
		if err != nil {
			return 0, err
		}
		if f < 0 {
			return 0, fmt.Errorf("negative float %f", f)
		}
		return uint64(f), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to uint64", v)
	}
}

func toFloat64(v any) (float64, error) {
	switch x := v.(type) {
	case float64:
		return x, nil
	case float32:
		return float64(x), nil
	case int:
		return float64(x), nil
	case int64:
		return float64(x), nil
	case json.Number:
		return x.Float64()
	case string:
		return strconv.ParseFloat(x, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", v)
	}
}

type BinancePayloadType interface {
	AggTrade | Trade | OrderBookDelta | OrderBookSnapshot
}

// DecodeBinanceMap 将 map 中按 binance tag 的字段写入 out（不存在则忽略）
// DecodeBinanceMap 将 map 中按 binance tag 的字段写入并返回实例 T
func DecodeBinanceMap[T BinancePayloadType](m map[string]any) (T, error) {
	var t T
	rv := reflect.ValueOf(&t).Elem() // 获取 T 的可寻址 Value

	// 处理 T 是指针的情况，例如 *Trade
	structVal := rv
	if rv.Kind() == reflect.Pointer {
		// 如果 T 是指针类型，需要初始化它指向的结构体
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
		}
		structVal = rv.Elem()
	}

	if structVal.Kind() != reflect.Struct {
		return t, fmt.Errorf("target type %T must be a struct or struct pointer", t)
	}

	ti := getTypeInfo(structVal.Type())

	// 遍历缓存的字段信息
	for tag, fi := range ti.fields {
		raw, ok := m[tag]
		if !ok || raw == nil {
			continue
		}

		fv := structVal.FieldByIndex(fi.index)
		if !fv.CanSet() {
			continue
		}

		// 核心赋值逻辑
		if err := setValueFast(fv, fi, raw); err != nil {
			// 根据业务需求，可以选择跳过错误字段或直接返回
			// return t, fmt.Errorf("field %s: %w", tag, err)
			continue
		}
	}

	return t, nil
}

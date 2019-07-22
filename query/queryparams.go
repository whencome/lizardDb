package query

import (
	"strconv"
	"fmt"
	"strings"
	"bytes"
)

// 查询参数
type QueryParams struct {
	params  		[]interface{}
}

func NewQueryParams() *QueryParams {
	return &QueryParams {
		params : make([]interface{}, 0),
	}
}

// 添加参数
func (this *QueryParams) AddBatch(values... interface{}) {
	if len(values) == 0 {
		return
	}
	for _, v := range values {
		this.Add(v)
	}
}

// 添加参数
func (this *QueryParams) Add(value interface{}) {
	v := ""
	switch value.(type) {
	case int,int8,int16,int32,int64,uint,uint8,uint16,uint32,uint64:
		intValue := value.(int64)
		v = strconv.FormatInt(intValue, 10)
	case float32, float64:
		floatValue := value.(float64)
		v = strconv.FormatFloat(floatValue, 'f', -1, 64)
	case string:
		v = value.(string)
	case bool:
		v = strconv.FormatBool(value.(bool))
	case []int:
		intArr := value.([]int)
		v = this.formatIntArr(intArr)
	case []float64:
		intArr := value.([]float64)
		v = this.formatFloatArr(intArr)
	case []rune:
		v = string(value.([]rune))
	case []byte:
		v = string(value.([]byte))
	case []string:
		v = "(\"" + strings.Join(value.([]string), "\",\"") + "\")"
	default:
		v = fmt.Sprintf("%s", value)
	}
	this.params = append(this.params, v)
}

func (this *QueryParams) formatIntArr(value []int) string {
	if len(value) < 1 {
		return "(-1)"
	}
	buf := bytes.Buffer{}
	buf.WriteString("(")
	for i, v := range value {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.FormatInt(int64(v), 10))
	}
	buf.WriteString(")")
	return buf.String()
}

func (this *QueryParams) formatFloatArr(value []float64) string {
	if len(value) < 1 {
		return "(-1)"
	}
	buf := bytes.Buffer{}
	buf.WriteString("(")
	for i, v := range value {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.FormatFloat(v, 'f', -1, 64))
	}
	buf.WriteString(")")
	return buf.String()
}

// 获取参数列表
func (this *QueryParams) Params() []interface{} {
	return this.params
}

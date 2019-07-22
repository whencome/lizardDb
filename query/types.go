package query

import (
	"fmt"
	"bytes"
	"strings"
)

type Number struct {
	value 		interface{}
	strValue 	string 		 // 保存string值
}

type Numbers struct {
	values 		[]*Number
}

type Value struct {
	data 		interface{}
	size 		int
	strValue 	string 		 // 保存string值
	isString 	bool
	isNumber	bool
	isNil		bool
	isArray     bool
}

func NewNumber(v interface{}) *Number {
	num := &Number{}
	switch v.(type) {
	case int:
		num.value = v
		num.strValue = fmt.Sprintf("%d", v.(int))
	case int8:
		num.value = v
		num.strValue = fmt.Sprintf("%d", v.(int8))
	case int16:
		num.value = v
		num.strValue = fmt.Sprintf("%d", v.(int16))
	case int32:
		num.value = v
		num.strValue = fmt.Sprintf("%d", v.(int32))
	case int64:
		num.value = v
		num.strValue = fmt.Sprintf("%d", v.(int64))
	case uint:
		num.value = v
		num.strValue = fmt.Sprintf("%d", v.(uint))
	case uint8:
		num.value = v
		num.strValue = fmt.Sprintf("%d", v.(uint8))
	case uint16:
		num.value = v
		num.strValue = fmt.Sprintf("%d", v.(uint16))
	case uint32:
		num.value = v
		num.strValue = fmt.Sprintf("%d", v.(uint32))
	case uint64:
		num.value = v
		num.strValue = fmt.Sprintf("%d", v.(uint64))
	case float32:
		num.value = v
		num.strValue = fmt.Sprintf("%f", v.(float32))
	case float64:
		num.value = v
		num.strValue = fmt.Sprintf("%f", v.(float64))
	default:
		// 其他情况不处理
	}
	return num
}

func (num *Number) ToString() string {
	if num == nil {
		return ""
	}
	return num.strValue
}

func NewNumbers(v interface{}) *Numbers {
	var nums = &Numbers{
		values : make([]*Number, 0),
	}
	switch v.(type) {
	case int,int8,int16,int32,int64,uint,uint8,uint16,uint32,uint64,float32,float64:
		nums.Add(NewNumber(v))
	case []int:
		values := v.([]int)
		for _, n := range values {
			nums.Add(NewNumber(n))
		}
	case []int8:
		values := v.([]int8)
		for _, n := range values {
			nums.Add(NewNumber(n))
		}
	case []int16:
		values := v.([]int16)
		for _, n := range values {
			nums.Add(NewNumber(n))
		}
	case []int32:
		values := v.([]int32)
		for _, n := range values {
			nums.Add(NewNumber(n))
		}
	case []int64:
		values := v.([]int64)
		for _, n := range values {
			nums.Add(NewNumber(n))
		}
	case []uint:
		values := v.([]uint)
		for _, n := range values {
			nums.Add(NewNumber(n))
		}
	case []uint8:
		values := v.([]uint8)
		for _, n := range values {
			nums.Add(NewNumber(n))
		}
	case []uint16:
		values := v.([]uint16)
		for _, n := range values {
			nums.Add(NewNumber(n))
		}
	case []uint32:
		values := v.([]uint32)
		for _, n := range values {
			nums.Add(NewNumber(n))
		}
	case []uint64:
		values := v.([]uint64)
		for _, n := range values {
			nums.Add(NewNumber(n))
		}
	case []float32:
		values := v.([]float32)
		for _, n := range values {
			nums.Add(NewNumber(n))
		}
	case []float64:
		values := v.([]float64)
		for _, n := range values {
			nums.Add(NewNumber(n))
		}
	default:
	}
	return nums
}

// add a number to numbers
func (nums *Numbers) Add(n *Number) {
	if n == nil {
		return
	}
	nums.values = append(nums.values, n)
}

func (nums *Numbers) Size() int {
	if nums == nil {
		return 0
	}
	return len(nums.values)
}

// format numbers to string， join with “,”
func (nums *Numbers) ToString() string {
	if nums.Size() == 0 {
		return ""
	}
	buf := bytes.Buffer{}
	for idx, number := range nums.values {
		if idx > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(number.ToString())
	}
	return buf.String()
}

func NewValue(d interface{}) *Value {
	// 单独处理值为nil的情况
	if d == nil {
		return &Value{
			data : d,
			isNil:true,
			strValue:"NULL",
		}
	}
	// 其他情况处理
	value := &Value{}
	switch d.(type) {
	case int,int8,int16,int32,int64,uint,uint8,uint16,uint32,uint64,float32,float64:
		num := NewNumber(d)
		value.size = 1
		value.data = num
		value.isNumber = true
		value.strValue = num.ToString()
	case string:
		value.size = 1
		value.data = d.(string)
		value.isString = true
		value.strValue = "\"" + d.(string) + "\""
	case []int,[]int8,[]int16,[]int32,[]int64,[]uint,[]uint8,[]uint16,[]uint32,[]uint64,[]float32,[]float64:
		nums := NewNumbers(d)
		value.data = nums
		value.size = nums.Size()
		value.isNumber = true
		value.isArray = true
		value.strValue = nums.ToString()
	case []string:
		strs := d.([]string)
		value.data = strs
		value.size = len(strs)
		value.isString = true
		value.isArray = true
		value.strValue = "\"" + strings.Join(strs, "\",\"") + "\""
	default:
	}
	return value
}

func (v *Value) ToString() string {
	return v.strValue
}

func (v *Value) GetUnquoteString() string {
	// 非数组
	if !v.isArray {
		if v.isString {
			return v.data.(string)
		}
		return v.strValue
	}
	// 数组
	if v.isString {
		return strings.Join(v.data.([]string), ",")
	}
	return v.strValue
}

func (v *Value) GetString(pos int) string {
	if pos > v.size {
		return ""
	}
	if !v.isArray {
		return v.strValue
	}
	if v.isNumber {
		vals, ok := v.data.(*Numbers)
		if !ok {
			return ""
		}
		if len(vals.values) < pos {
			return ""
		}
		return vals.values[0].ToString()
	} else {
		return v.data.([]string)[0]
	}
}
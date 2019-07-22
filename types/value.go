package types

import "strconv"

type Value struct {
	data  interface{}
}

func NewValue(d interface{}) *Value {
	return &Value{
		data : d,
	}
}

func (v *Value) String() string {
	switch v.data.(type) {
	case string:
		return v.data.(string)
	case int64:
		tmp, ok := v.data.(int64)
		if !ok {
			return ""
		}
		return strconv.FormatInt(tmp, 10)
	case []uint8:
		tmp, ok := v.data.([]uint8)
		if !ok {
			return ""
		}
		return string(tmp)
	}
	return ""
}

func (v *Value) Int() int64 {
	strV := v.String()
	if strV == "" {
		return 0
	}
	iv, err := strconv.ParseInt(strV, 10, 64)
	if err != nil {
		return 0
	}
	return iv
}

func (v *Value) Uint() uint64 {
	strV := v.String()
	if strV == "" {
		return 0
	}
	iv, err := strconv.ParseUint(strV, 10, 64)
	if err != nil {
		return 0
	}
	return iv
}

func (v *Value) Bool() bool {
	strV := v.String()
	if strV == "" {
		return false
	}
	bv, err := strconv.ParseBool(strV)
	if err != nil {
		return false
	}
	return bv
}

func (v *Value) Float() float64 {
	strV := v.String()
	if strV == "" {
		return 0.00
	}
	fv, err := strconv.ParseFloat(strV, 64)
	if err != nil {
		return 0.00
	}
	return fv
}

package query

import (
	"bytes"
	"strconv"
)

func getIntList(value []int64) string {
	if len(value) < 1 {
		return "()"
	}
	buf := bytes.Buffer{}
	buf.WriteString("(")
	for i, v := range value {
		if i > 0 {
			buf.WriteString(",")
		}
		buf.WriteString(strconv.FormatInt(v, 10))
	}
	buf.WriteString(")")
	return buf.String()
}
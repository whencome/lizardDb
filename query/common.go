package query

import "errors"

var (
	// 条件为空
	ErrEmptyCondition = errors.New("empty query condition")
)

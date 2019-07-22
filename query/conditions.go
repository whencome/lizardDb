package query

import (
	"bytes"
	"strings"
	"regexp"
	"errors"
)

// 定义通用接口
type SqlCondition interface {
	Build() (string, error)
}

// 查询条件
type Conditions struct {
	Conds		[]*Condition
}

// 基本条件单元
type Condition struct {
	Logic	string  // 逻辑关系，and，or
	Value   string  // 条件值
}

// 普通文本条件，即已经构造好了的查询条件
type PlainCondition struct {
	Value   string  // 条件值
}

// 指定字段的查询条件
type FieldCondition struct {
	Field   string
	Operator 	string 	// 逻辑关系,=,in,not in,between
	Values  interface{}
}

// 创建一个新的Conditions对象
func NewConditions() *Conditions {
	return &Conditions{
		Conds : make([]*Condition, 0),
	}
}

// 创建一个新的PlainCondition对象
func NewPlainCondition(cond string) *PlainCondition {
	return &PlainCondition{
		Value : cond,
	}
}

// 创建一个新的FieldCondition对象
func NewFieldCondition(field, op string, vals interface{}) *FieldCondition {
	return &FieldCondition{
		Field : field,
		Operator:op,
		Values:vals,
	}
}

func (fc *FieldCondition) Build() (string, error) {
	var buf bytes.Buffer
	buf.WriteString(fc.Field)
	buf.WriteString(" ")
	re,_ := regexp.Compile("\\s+")
	var logic = strings.ToUpper(re.ReplaceAllString(fc.Operator, " "))
	switch logic {
	case "=","!=","<>",">",">=","<","<=":
		v := NewValue(fc.Values)
		buf.WriteString(v.ToString())
	case "IN": fallthrough
	case "NOT IN":
		buf.WriteString(logic)
		buf.WriteString(" ")
		buf.WriteString("(")
		v := NewValue(fc.Values)
		buf.WriteString(v.ToString())
		buf.WriteString(")")
	case "BETWEEN":
		vals := NewValue(fc.Values)
		buf.WriteString(" BETWEEN ")
		buf.WriteString(vals.GetString(0))
		buf.WriteString(" AND ")
		buf.WriteString(vals.GetString(1))
	case "IS", "IS NOT":
		v := NewValue(fc.Values)
		buf.WriteString(v.ToString())
	}
	return buf.String(), nil
}

func (pc *PlainCondition) Build() (string, error) {
	return pc.Value, nil
}

func (c *Conditions) addPlain(logic string, cond string) error {
	condition := &Condition{
		Logic : logic,
		Value : cond,
	}
	c.Conds = append(c.Conds, condition)
	return nil
}

func (c *Conditions) add(logic string, sc SqlCondition) error {
	patch, err := sc.Build()
	if err != nil {
		return err
	}
	if patch == "" {
		return errors.New("build condition failed")
	}
	condition := &Condition{
		Logic : logic,
		Value : patch,
	}
	c.Conds = append(c.Conds, condition)
	return nil
}

func (c *Conditions) And(sc SqlCondition) error {
	return c.add("and", sc)
}

func (c *Conditions) Or(sc SqlCondition) error {
	return c.add("or", sc)
}

func (c *Conditions) AndPlain(where string) error {
	return c.addPlain("and", where)
}

func (c *Conditions) OrPlain(where string) error {
	return c.addPlain("or", where)
}

func (c *Conditions) IsEmpty() bool {
	if c == nil || len(c.Conds) == 0 {
		return true
	}
	return false
}

func (c *Conditions) Build() (string, error) {
	if len(c.Conds) == 0 {
		return "",ErrEmptyCondition
	}
	var buf bytes.Buffer
	for i, cond := range c.Conds {
		if i > 0 {
			buf.WriteString(strings.ToUpper(cond.Logic))
		}
		buf.WriteString(" ( ")
		buf.WriteString(cond.Value)
		buf.WriteString(" ) ")
	}
	return buf.String(), nil
}
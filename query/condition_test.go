package query

import (
	"testing"
)

func TestConditions_Build(t *testing.T) {
	var conditions = NewConditions()
	err := conditions.AndPlain("id > 100")
	if err != nil {
		t.Error(err)
	}
	err = conditions.And(NewPlainCondition("score > 60"))
	if err != nil {
		t.Error(err)
	}
	err = conditions.And(NewFieldCondition("uid", "in", []int{100001,100002,100003}))
	if err != nil {
		t.Error(err)
	}
	err = conditions.And(NewFieldCondition("class", "not in", []string{"classA","classB","classC"}))
	if err != nil {
		t.Error(err)
	}
	rs, err := conditions.Build()
	if err != nil {
		t.Error(err)
	}
	t.Log(rs)
}

func TestConditions_BuildOrLogic(t *testing.T) {
	var conditions = NewConditions()
	conds1 := NewConditions()
	conds2 := NewConditions()
	err := conds1.AndPlain("id > 100")
	if err != nil {
		t.Error(err)
	}
	err = conds1.And(NewPlainCondition("score > 60"))
	if err != nil {
		t.Error(err)
	}
	err = conds1.And( NewFieldCondition("uid", "in", []int{100001,100002,100003}))
	if err != nil {
		t.Error(err)
	}
	err = conds2.And(NewFieldCondition("class", "not in", []string{"classA","classB","classC"}))
	if err != nil {
		t.Error(err)
	}

	conditions.And(conds1)
	conditions.Or(conds2)
	rs, err := conditions.Build()
	if err != nil {
		t.Error(err)
	}
	t.Log(rs)
}
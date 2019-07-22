package query

import "bytes"

const (
	LeftJoin  = "LEFT JOIN"
	RightJoin = "RIGHT JOIN"
	InnerJoin = "INNER JOIN"
)

type JoinTable struct {
	joinType 		string   	  // join方式，left，right，inner
	tableName       string        // 表名
	conds			*Conditions   // 连表条件
}

type JoinTables struct {
	joins			[]*JoinTable
}

func NewJoinTables() *JoinTables {
	return &JoinTables{
		joins : make([]*JoinTable, 0),
	}
}

func NewJoinTable(joinType, table string, cond *Conditions) *JoinTable {
	return &JoinTable{
		joinType 	: joinType,
		tableName	: table,
		conds 		: cond,
	}
}

func NewInnerJoinTable(table string, cond *Conditions) *JoinTable {
	return NewJoinTable("inner", table, cond)
}

func NewLeftJoinTable(table string, cond *Conditions) *JoinTable {
	return NewJoinTable("left", table, cond)
}
func NewRightJoinTable(table string, cond *Conditions) *JoinTable {
	return NewJoinTable("right", table, cond)
}

func (this *JoinTable) Build() (string, error) {
	var buf bytes.Buffer
	switch this.joinType {
	case "left":
		buf.WriteString("LEFT JOIN")
	case "right":
		buf.WriteString("RIGHT JOIN")
	default:
		buf.WriteString("INNER JOIN")
	}
	buf.WriteString(" ")
	buf.WriteString(this.tableName)
	buf.WriteString(" ON")
	onCond, err := this.conds.Build()
	if err != nil {
		return "", err
	}
	buf.WriteString(onCond)
	return buf.String(), nil
}

func (jts *JoinTables) Add(joinTable *JoinTable) {
	jts.joins = append(jts.joins, joinTable)
}

func (jts *JoinTables) Join(table string, cond *Conditions) {
	jt := NewInnerJoinTable(table, cond)
	jts.joins = append(jts.joins, jt)
}

func (jts *JoinTables) LeftJoin(table string, cond *Conditions) {
	jt := NewLeftJoinTable(table, cond)
	jts.joins = append(jts.joins, jt)
}

func (jts *JoinTables) RightJoin(table string, cond *Conditions) {
	jt := NewRightJoinTable(table, cond)
	jts.joins = append(jts.joins, jt)
}

func (jts *JoinTables) Build() (string, error) {
	if len(jts.joins) == 0 {
		return "", nil
	}
	var buf bytes.Buffer
	for _, jt := range jts.joins {
		joinStr, err := jt.Build()
		if err != nil {
			return "", err
		}
		buf.WriteString(joinStr)
		buf.WriteString(" ")
	}
	return buf.String(), nil
}
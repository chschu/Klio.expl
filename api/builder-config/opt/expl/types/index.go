package types

import "fmt"

type Index interface {
	fmt.Stringer
	SqlCondition(cmp IndexComparison) (sqlCondition string, params []any)
}

type IndexComparison int

const (
	IndexStartingWith IndexComparison = iota
	IndexMatching
	IndexEndingWith
)

var ascendingIndexComparisonOperator = map[IndexComparison]string{
	IndexStartingWith: ">=",
	IndexMatching:     "=",
	IndexEndingWith:   "<=",
}

var descendingIndexComparisonOperator = map[IndexComparison]string{
	IndexStartingWith: "<=",
	IndexMatching:     "=",
	IndexEndingWith:   ">=",
}

type HeadIndex uint

func (i HeadIndex) SqlCondition(cmp IndexComparison) (string, []any) {
	return "head_index " + ascendingIndexComparisonOperator[cmp] + " ?", []any{uint(i)}
}

func (i HeadIndex) String() string {
	return fmt.Sprintf("%d", i)
}

type TailIndex uint

func (i TailIndex) SqlCondition(cmp IndexComparison) (string, []any) {
	return "tail_index " + descendingIndexComparisonOperator[cmp] + " ?", []any{uint(i)}
}

func (i TailIndex) String() string {
	return fmt.Sprintf("-%d", i)
}

type PermanentIndex uint

func (i PermanentIndex) SqlCondition(cmp IndexComparison) (string, []any) {
	return "permanent_index " + ascendingIndexComparisonOperator[cmp] + " ?", []any{uint(i)}
}

func (i PermanentIndex) String() string {
	return fmt.Sprintf("p%d", i)
}

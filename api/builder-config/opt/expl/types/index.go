package types

import "fmt"

type Index interface {
	fmt.Stringer
	SqlCondition(cmp IndexComparison) (sqlCondition string, params []any)
}

type IndexComparison string

const (
	IndexStartingWith IndexComparison = ">="
	IndexEndingWith   IndexComparison = "<="
)

type HeadIndex uint

func (i HeadIndex) SqlCondition(cmp IndexComparison) (string, []any) {
	return "head_index " + string(cmp) + " ?", []any{uint(i)}
}

func (i HeadIndex) String() string {
	return fmt.Sprintf("%d", i)
}

type TailIndex uint

func (i TailIndex) SqlCondition(cmp IndexComparison) (string, []any) {
	return "0-tail_index " + string(cmp) + " 0-?", []any{uint(i)}
}

func (i TailIndex) String() string {
	return fmt.Sprintf("-%d", i)
}

type PermanentIndex uint

func (i PermanentIndex) SqlCondition(cmp IndexComparison) (string, []any) {
	return "permanent_index " + string(cmp) + " ?", []any{uint(i)}
}

func (i PermanentIndex) String() string {
	return fmt.Sprintf("p%d", i)
}

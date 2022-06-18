package types

import "fmt"

type Index interface {
	Value() uint
	Descending() bool // true iff newer entries have smaller index values
	DatabaseColumn() string
}

type IndexRange struct {
	From Index
	To   Index
}

type IndexSpec []IndexRange

type HeadIndex uint

func (i HeadIndex) Value() uint {
	return uint(i)
}

func (i HeadIndex) Descending() bool {
	return false
}

func (i HeadIndex) DatabaseColumn() string {
	return "head_index"
}

func (i HeadIndex) String() string {
	return fmt.Sprintf("%d", i)
}

type TailIndex uint

func (i TailIndex) Value() uint {
	return uint(i)
}

func (i TailIndex) Descending() bool {
	return true
}

func (i TailIndex) DatabaseColumn() string {
	return "tail_index"
}

func (i TailIndex) String() string {
	return fmt.Sprintf("-%d", i)
}

type PermanentIndex uint

func (i PermanentIndex) Value() uint {
	return uint(i)
}

func (i PermanentIndex) Descending() bool {
	return false
}

func (i PermanentIndex) DatabaseColumn() string {
	return "permanent_index"
}

func (i PermanentIndex) String() string {
	return fmt.Sprintf("p%d", i)
}

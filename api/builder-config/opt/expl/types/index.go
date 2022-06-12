package types

import "fmt"

type HeadIndex uint

type TailIndex uint

type PermanentIndex uint

func (i HeadIndex) String() string {
	return fmt.Sprintf("%d", i)
}

func (i TailIndex) String() string {
	return fmt.Sprintf("-%d", i)
}

func (i PermanentIndex) String() string {
	return fmt.Sprintf("p%d", i)
}

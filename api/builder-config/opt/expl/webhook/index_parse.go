package webhook

import (
	"errors"
	"fmt"
	"klio/expl/types"
	"regexp"
	"strconv"
)

func parseIndex(index string) (types.Index, error) {
	sep := regexp.MustCompile("^(?P<Prefix>|-|p)(?P<N>[1-9][0-9]*)$")
	match := sep.FindStringSubmatch(index)
	if match == nil {
		return nil, errors.New(fmt.Sprintf("invalid index: %s", index))
	}
	prefix := match[sep.SubexpIndex("Prefix")]
	nStr := match[sep.SubexpIndex("N")]

	n64, err := strconv.ParseUint(nStr, 10, 32)
	if err != nil {
		return nil, err
	}
	n := uint(n64)

	switch prefix {
	case "":
		return types.HeadIndex(n), nil
	case "-":
		return types.TailIndex(n), nil
	case "p":
		return types.PermanentIndex(n), nil
	}

	return nil, errors.New(fmt.Sprintf("invalid index: %s", index))
}

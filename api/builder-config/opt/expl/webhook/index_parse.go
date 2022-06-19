package webhook

import (
	"errors"
	"fmt"
	"klio/expl/types"
	"regexp"
	"strconv"
	"strings"
)

func parseIndexSpec(indexSpec string) (types.IndexSpec, error) {
	var out = types.IndexSpec{}
	sep := regexp.MustCompile("\\pZ+")
	indexRanges := sep.Split(indexSpec, -1)
	for _, indexRange := range indexRanges {
		ir, err := parseIndexRange(indexRange)
		if err != nil {
			return nil, err
		}
		out = append(out, *ir)
	}
	return out, nil
}

func parseIndexRange(indexRange string) (*types.IndexRange, error) {
	indexes := strings.SplitN(indexRange, ":", 3)
	l := len(indexes)

	if l < 1 || l > 2 {
		return nil, errors.New(fmt.Sprintf("invalid index range: %s", indexRange))
	}

	from, err := parseIndex(indexes[0])
	if err != nil {
		return nil, err
	}

	var to types.Index
	if l == 2 {
		to, err = parseIndex(indexes[1])
		if err != nil {
			return nil, err
		}
	} else {
		to = from
	}

	return &types.IndexRange{
		From: from,
		To:   to,
	}, nil
}

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

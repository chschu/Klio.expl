package types

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type IndexSpecParser interface {
	ParseIndexSpec(s string) (IndexSpec, error)
	ParseIndexRange(s string) (IndexRange, error)
	ParseIndex(s string) (Index, error)
}

func NewIndexSpecParser() IndexSpecParser {
	return &indexSpecParser{}
}

type indexSpecParser struct{}

func (p *indexSpecParser) ParseIndexSpec(s string) (IndexSpec, error) {
	var out = indexSpec{}
	sep := regexp.MustCompile("\\pZ+")
	indexRanges := sep.Split(s, -1)
	for _, indexRange := range indexRanges {
		ir, err := p.ParseIndexRange(indexRange)
		if err != nil {
			return nil, err
		}
		out = append(out, ir)
	}
	return out, nil
}

func (p *indexSpecParser) ParseIndexRange(s string) (IndexRange, error) {
	indexes := strings.SplitN(s, ":", 3)
	l := len(indexes)

	if l < 1 || l > 2 {
		return nil, fmt.Errorf("invalid index range: %s", s)
	}

	from, err := p.ParseIndex(indexes[0])
	if err != nil {
		return nil, err
	}

	var to Index
	if l == 2 {
		to, err = p.ParseIndex(indexes[1])
		if err != nil {
			return nil, err
		}
	} else {
		to = from
	}

	return &indexRange{
		from: from,
		to:   to,
	}, nil
}

func (p *indexSpecParser) ParseIndex(s string) (Index, error) {
	sep := regexp.MustCompile("^(?P<Prefix>|-|p)(?P<N>[1-9][0-9]*)$")
	match := sep.FindStringSubmatch(s)
	if match == nil {
		return nil, fmt.Errorf("invalid index: %s", s)
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
		return HeadIndex(n), nil
	case "-":
		return TailIndex(n), nil
	case "p":
		return PermanentIndex(n), nil
	}

	return nil, fmt.Errorf("invalid index: %s", s)
}

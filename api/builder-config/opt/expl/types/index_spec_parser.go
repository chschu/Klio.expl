package types

import (
	"fmt"
	"regexp"
	"strconv"
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
	var ranges []IndexRange
	sep := regexp.MustCompile("\\pZ+")
	indexRanges := sep.Split(s, -1)
	for _, indexRange := range indexRanges {
		if len(indexRange) > 0 {
			ir, err := p.ParseIndexRange(indexRange)
			if err != nil {
				return nil, err
			}
			ranges = append(ranges, ir)
		}
	}
	return NewIndexSpec(ranges...), nil
}

func (p *indexSpecParser) ParseIndexRange(s string) (IndexRange, error) {
	sep := regexp.MustCompile("^\\pZ*(?P<From>\\PZ+?)(:(?P<To>\\PZ+?))?\\pZ*$")
	match := sep.FindStringSubmatch(s)
	if match == nil {
		return nil, fmt.Errorf("invalid index range: %s", s)
	}

	fromStr := match[sep.SubexpIndex("From")]
	from, err := p.ParseIndex(fromStr)
	if err != nil {
		return nil, err
	}

	toStr := match[sep.SubexpIndex("To")]
	if toStr == "" {
		return NewIndexRange(from, from), nil
	}
	to, err := p.ParseIndex(toStr)
	if err != nil {
		return nil, err
	}

	return NewIndexRange(from, to), nil
}

func (p *indexSpecParser) ParseIndex(s string) (Index, error) {
	sep := regexp.MustCompile("^\\pZ*(?P<Prefix>|-|p)(?P<N>[1-9]\\d*)\\pZ*$")
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

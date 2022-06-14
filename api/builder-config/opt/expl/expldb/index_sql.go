package expldb

import (
	"fmt"
	"klio/expl/types"
)

func indexRangesSqlCondition(indexRanges []types.IndexRange, params []any) (string, []any) {
	nextParamPart := func(value any) string {
		params = append(params, value)
		return fmt.Sprintf("$%d", len(params))
	}

	indexPart := func(cmp string, index types.Index) string {
		var prefix string
		if index.Descending() {
			prefix = "0-" // convert to ascending values for comparison
		} else {
			prefix = ""
		}
		return prefix + index.DatabaseColumn() + cmp + prefix + nextParamPart(index.Value())
	}

	indexRangePart := func(indexRange types.IndexRange) string {
		return indexPart(">=", indexRange.From) + " AND " + indexPart("<=", indexRange.To)
	}

	sqlCondition := "false"
	for _, indexRange := range indexRanges {
		sqlCondition += " OR (" + indexRangePart(indexRange) + ")"
	}

	return sqlCondition, params
}

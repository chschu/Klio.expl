package expldb

import (
	"klio/expl/types"
)

func indexSpecSqlCondition(indexSpec types.IndexSpec) (sqlCondition string, params []any) {
	indexPart := func(cmp string, index types.Index) string {
		var prefix string
		if index.Descending() {
			prefix = "0-" // convert to ascending values for comparison
		} else {
			prefix = ""
		}
		params = append(params, index.Value())
		return prefix + index.DatabaseColumn() + cmp + prefix + "?"
	}

	indexRangePart := func(indexRange types.IndexRange) string {
		return indexPart(">=", indexRange.From) + " AND " + indexPart("<=", indexRange.To)
	}

	sqlCondition = "false"
	for _, indexRange := range indexSpec {
		sqlCondition += " OR (" + indexRangePart(indexRange) + ")"
	}
	return sqlCondition, params
}

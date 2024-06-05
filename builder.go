package gobuilder

import (
	"fmt"
	"strings"

	"golang.org/x/exp/maps"
)

type GoBuilder struct {
	sql string
}

func (gb *GoBuilder) Select(table string, columns ...string) *GoBuilder {
	if len(columns) == 0 {
		columns = append(columns, "*")
	}
	gb.sql = fmt.Sprintf("SELECT %v FROM %v", strings.Join(columns, ", "), table)
	return gb
}

func (gb *GoBuilder) SelectDistinct(table string, columns ...string) *GoBuilder {
	if len(columns) == 0 {
		columns = append(columns, "*")
	}
	gb.sql = fmt.Sprintf("SELECT DISTINCT %v FROM %v", strings.Join(columns, ", "), table)
	return gb
}

func (gb *GoBuilder) Insert(table string, args map[string]any) *GoBuilder {
	if len(args) != 0 {
		keys := maps.Keys(args)
		var values []string
		for _, v := range maps.Values(args) {
			values = append(values, fmt.Sprintf("'%v'", v))
		}
		gb.sql = fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", table, strings.Join(keys, ", "), strings.Join(values, ", "))
	}
	return gb
}

func (gb *GoBuilder) Update(table string, args map[string]any) *GoBuilder {
	if len(args) != 0 {
		var setClauses []string
		for key, value := range args {
			setClauses = append(setClauses, fmt.Sprintf("%v = '%v'", key, value))
		}
		gb.sql = fmt.Sprintf("UPDATE %v SET %v", table, strings.Join(setClauses, ", "))
	}
	return gb
}

func (gb *GoBuilder) Delete(table string) *GoBuilder {
	gb.sql = fmt.Sprintf("DELETE FROM %v", table)
	return gb
}

func (gb *GoBuilder) Where(key, opt string, val any) *GoBuilder {
	return gb.where("AND", key, opt, val)
}

func (gb *GoBuilder) OrWhere(key, opt string, val any) *GoBuilder {
	return gb.where("OR", key, opt, val)
}

func (gb *GoBuilder) where(OP, key, opt string, val any) *GoBuilder {
	var clause string

	switch v := val.(type) {
	case int, int64, float32, float64:
		clause = fmt.Sprintf("%v %v %v", key, opt, v)
	case *GoBuilder:
		clause = fmt.Sprintf("%v %v (%v)", key, opt, v.ToSql())
	default:
		clause = fmt.Sprintf("%v %v '%v'", key, opt, v)
	}

	if strings.Contains(gb.sql, "WHERE") {
		gb.sql = fmt.Sprintf("%v %v %v", gb.sql, OP, clause)
	} else {
		gb.sql = fmt.Sprintf("%v WHERE %v", gb.sql, clause)
	}
	return gb
}

func (gb *GoBuilder) In(column string, args ...any) *GoBuilder {
	return gb.in("AND", column, args...)
}

func (gb *GoBuilder) OrIn(column string, args ...any) *GoBuilder {
	return gb.in("OR", column, args...)
}

func (gb *GoBuilder) in(OP, column string, args ...any) *GoBuilder {
	if len(args) > 0 {
		var values []string
		for _, arg := range args {
			switch v := arg.(type) {
			case int, int64, float32, float64:
				values = append(values, fmt.Sprintf("%v", v))
			default:
				values = append(values, fmt.Sprintf("'%v'", v))
			}
		}
		clause := fmt.Sprintf("%v IN (%v)", column, strings.Join(values, ", "))
		if strings.Contains(gb.sql, "WHERE") {
			gb.sql = fmt.Sprintf("%v %v %v", gb.sql, OP, clause)
		} else {
			gb.sql = fmt.Sprintf("%v WHERE %v", gb.sql, clause)
		}
	}
	return gb
}

func (gb *GoBuilder) Between(column string, args ...any) *GoBuilder {
	return gb.between("AND", column, args...)
}

func (gb *GoBuilder) OrBetween(column string, args ...any) *GoBuilder {
	return gb.between("OR", column, args...)
}

func (gb *GoBuilder) between(OP, column string, args ...any) *GoBuilder {
	if len(args) == 2 {
		var clause string
		switch args[0].(type) {
		case int, int64, float32, float64:
			clause = fmt.Sprintf("%v BETWEEN %v AND %v", column, args[0], args[1])
		default:
			clause = fmt.Sprintf("%v BETWEEN '%v' AND '%v'", column, args[0], args[1])
		}
		if strings.Contains(gb.sql, "WHERE") {
			gb.sql = fmt.Sprintf("%v %v %v", gb.sql, OP, clause)
		} else {
			gb.sql = fmt.Sprintf("%v WHERE %v", gb.sql, clause)
		}
	}
	return gb
}

func (gb *GoBuilder) IsNull(column string) *GoBuilder {
	clause := fmt.Sprintf("%v IS NULL", column)
	if strings.Contains(gb.sql, "WHERE") {
		gb.sql = fmt.Sprintf("%v AND %v", gb.sql, clause)
	} else {
		gb.sql = fmt.Sprintf("%v WHERE %v", gb.sql, clause)
	}
	return gb
}

func (gb *GoBuilder) OrIsNull(column string) *GoBuilder {
	clause := fmt.Sprintf("%v IS NULL", column)
	if strings.Contains(gb.sql, "WHERE") {
		gb.sql = fmt.Sprintf("%v OR %v", gb.sql, clause)
	} else {
		gb.sql = fmt.Sprintf("%v WHERE %v", gb.sql, clause)
	}
	return gb
}

func (gb *GoBuilder) IsNotNull(column string) *GoBuilder {
	clause := fmt.Sprintf("%v IS NOT NULL", column)
	if strings.Contains(gb.sql, "WHERE") {
		gb.sql = fmt.Sprintf("%v AND %v", gb.sql, clause)
	} else {
		gb.sql = fmt.Sprintf("%v WHERE %v", gb.sql, clause)
	}
	return gb
}

func (gb *GoBuilder) OrIsNotNull(column string) *GoBuilder {
	clause := fmt.Sprintf("%v IS NOT NULL", column)
	if strings.Contains(gb.sql, "WHERE") {
		gb.sql = fmt.Sprintf("%v OR %v", gb.sql, clause)
	} else {
		gb.sql = fmt.Sprintf("%v WHERE %v", gb.sql, clause)
	}
	return gb
}

func (gb *GoBuilder) Having(condition string) *GoBuilder {
	if strings.Contains(gb.sql, "HAVING") {
		gb.sql = fmt.Sprintf("%v AND %v", gb.sql, condition)
	} else {
		gb.sql = fmt.Sprintf("%v HAVING %v", gb.sql, condition)
	}
	return gb
}

func (gb *GoBuilder) OrHaving(condition string) *GoBuilder {
	if strings.Contains(gb.sql, "HAVING") {
		gb.sql = fmt.Sprintf("%v OR %v", gb.sql, condition)
	} else {
		gb.sql = fmt.Sprintf("%v HAVING %v", gb.sql, condition)
	}
	return gb
}

func (gb *GoBuilder) Join(joinType, table, onCondition string) *GoBuilder {
	gb.sql = fmt.Sprintf("%v %v JOIN %v ON %v", gb.sql, joinType, table, onCondition)
	return gb
}

func (gb *GoBuilder) Limit(start, limit int) *GoBuilder {
	gb.sql = fmt.Sprintf("%v LIMIT %v, %v", gb.sql, start, limit)
	return gb
}

func (gb *GoBuilder) GroupBy(columns ...string) *GoBuilder {
	gb.sql = fmt.Sprintf("%v GROUP BY %v", gb.sql, strings.Join(columns, ", "))
	return gb
}

func (gb *GoBuilder) OrderBy(columns ...string) *GoBuilder {
	gb.sql = fmt.Sprintf("%v ORDER BY %v ASC", gb.sql, strings.Join(columns, ", "))
	return gb
}

func (gb *GoBuilder) OrderByDesc(columns ...string) *GoBuilder {
	gb.sql = fmt.Sprintf("%v ORDER BY %v DESC", gb.sql, strings.Join(columns, ", "))
	return gb
}

func (gb *GoBuilder) Union(sql string) *GoBuilder {
	gb.sql = fmt.Sprintf("%v UNION %v", gb.sql, sql)
	return gb
}

func (gb *GoBuilder) ToSql() string {
	return gb.sql
}

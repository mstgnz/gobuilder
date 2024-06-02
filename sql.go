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
	gb.sql = fmt.Sprintf("SELECT %v FROM %v", strings.Join(columns, ","), table)
	return gb
}

func (gb *GoBuilder) Insert(table string, args map[string]string) *GoBuilder {
	if len(args) != 0 {
		keys := maps.Keys(args)
		values := ""
		for _, v := range maps.Values(args) {
			values += fmt.Sprintf("'%v',", v)
		}
		values = strings.Trim(values, ",")
		gb.sql = fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", table, strings.Join(keys, ","), values)
	}
	return gb
}

func (gb *GoBuilder) Update(table string, args map[string]string) *GoBuilder {
	if len(args) != 0 {
		set := ""
		for _, key := range maps.Keys(args) {
			set += fmt.Sprintf("%v='%v', ", key, args[key])
		}
		set = strings.Trim(set, ", ")
		gb.sql = fmt.Sprintf("UPDATE %v SET %v", table, set)
	}
	return gb
}

func (gb *GoBuilder) Delete(table string) *GoBuilder {
	gb.sql = fmt.Sprintf("DELETE FROM %v", table)
	return gb
}

func (gb *GoBuilder) Where(key, opt, val string) *GoBuilder {
	if strings.Contains(gb.sql, "WHERE") {
		gb.sql = fmt.Sprintf("%v AND %v%v'%v'", gb.sql, key, opt, val)
	} else {
		gb.sql = fmt.Sprintf("%v WHERE %v%v'%v'", gb.sql, key, opt, val)
	}
	return gb
}

func (gb *GoBuilder) OrWhere(key, opt, val string) *GoBuilder {
	if strings.Contains(gb.sql, "WHERE") {
		gb.sql = fmt.Sprintf("%v OR %v%v'%v'", gb.sql, key, opt, val)
	} else {
		gb.sql = fmt.Sprintf("%v WHERE %v%v'%v'", gb.sql, key, opt, val)
	}
	return gb
}

func (gb *GoBuilder) In(column string, args ...string) *GoBuilder {
	if len(args) > 0 {
		var values string
		for _, v := range args {
			values += fmt.Sprintf("'%v', ", v)
		}
		values = fmt.Sprintf("(%v)", strings.Trim(values, ", "))
		if strings.Contains(gb.sql, "WHERE") {
			gb.sql = fmt.Sprintf("%v AND %v", gb.sql, values)
		} else {
			gb.sql = fmt.Sprintf("%v WHERE %v IN %v", gb.sql, column, values)
		}
	}
	return gb
}

func (gb *GoBuilder) OrIn(column string, args ...string) *GoBuilder {
	if len(args) > 0 {
		var values string
		for _, v := range args {
			values += fmt.Sprintf("'%v', ", v)
		}
		values = fmt.Sprintf("(%v)", strings.Trim(values, ", "))
		if strings.Contains(gb.sql, "WHERE") {
			gb.sql = fmt.Sprintf("%v OR %v", gb.sql, values)
		} else {
			gb.sql = fmt.Sprintf("%v WHERE %v IN %v", gb.sql, column, values)
		}
	}
	return gb
}

func (gb *GoBuilder) Between(column string, vals ...string) *GoBuilder {
	if len(vals) == 2 {
		if strings.Contains(gb.sql, "WHERE") {
			gb.sql = fmt.Sprintf("%v AND %v BETWEEN '%v' AND '%v'", gb.sql, column, vals[0], vals[1])
		} else {
			gb.sql = fmt.Sprintf("%v WHERE %v BETWEEN '%v' AND '%v'", gb.sql, column, vals[0], vals[1])
		}
	}
	return gb
}

func (gb *GoBuilder) OrBetween(column string, vals ...string) *GoBuilder {
	if len(vals) == 2 {
		if strings.Contains(gb.sql, "WHERE") {
			gb.sql = fmt.Sprintf("%v OR %v BETWEEN '%v' AND '%v'", gb.sql, column, vals[0], vals[1])
		} else {
			gb.sql = fmt.Sprintf("%v WHERE %v BETWEEN '%v' AND '%v'", gb.sql, column, vals[0], vals[1])
		}
	}
	return gb
}

func (gb *GoBuilder) Join(joinName, table, equal string) *GoBuilder {
	gb.sql = fmt.Sprintf("%v %v JOIN %v ON %v", gb.sql, joinName, table, equal)
	return gb
}

func (gb *GoBuilder) Limit(start, limit int) *GoBuilder {
	gb.sql = fmt.Sprintf("%v LIMIT %v,%v", gb.sql, start, limit)
	return gb
}

func (gb *GoBuilder) GroupBy(column string) *GoBuilder {
	gb.sql = fmt.Sprintf("%v GROUP BY %v", gb.sql, column)
	return gb
}

func (gb *GoBuilder) OrderBy(column, sortName string) *GoBuilder {
	gb.sql = fmt.Sprintf("%v ORDER BY %v %v", gb.sql, column, sortName)
	return gb
}

func (gb *GoBuilder) Union(sql string) *GoBuilder {
	gb.sql = gb.sql + " UNION " + sql
	return gb
}

func (gb *GoBuilder) Sql() string {
	return gb.sql
}

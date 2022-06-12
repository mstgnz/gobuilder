package main

import (
	"fmt"
	"strings"

	"golang.org/x/exp/maps"
)

type Sql struct {
	sql string
}

func (s *Sql) Select(table string, columns []string) *Sql {
	if len(columns) == 0 {
		columns = append(columns, "*")
	}
	s.sql = fmt.Sprintf("SELECT %v FROM %v", strings.Join(columns, ","), table)
	return s
}

func (s *Sql) Insert(table string, args map[string]string) *Sql {
	if len(args) != 0 {
		keys := maps.Keys(args)
		values := ""
		for _, v := range maps.Values(args) {
			values += fmt.Sprintf("'%v',", v)
		}
		values = strings.Trim(values, ",")
		s.sql = fmt.Sprintf("INSERT INTO %v (%v) VALUES (%v)", table, strings.Join(keys, ","), values)
	}
	return s
}

func (s *Sql) Update(table string, args map[string]string) *Sql {
	if len(args) != 0 {
		set := ""
		for _, key := range maps.Keys(args) {
			set += fmt.Sprintf("%v='%v', ", key, args[key])
		}
		set = strings.Trim(set, ", ")
		s.sql = fmt.Sprintf("UPDATE %v SET %v", table, set)
	}
	return s
}

func (s *Sql) Delete(table string) *Sql {
	s.sql = fmt.Sprintf("DELETE FROM %v", table)
	return s
}

func (s *Sql) Where(key, opt, val string) *Sql {
	if strings.Contains(s.sql, "WHERE") {
		s.sql = fmt.Sprintf("%v AND %v%v'%v'", s.sql, key, opt, val)
	} else {
		s.sql = fmt.Sprintf("%v WHERE %v%v'%v'", s.sql, key, opt, val)
	}
	return s
}

func (s *Sql) OrWhere(key, opt, val string) *Sql {
	if strings.Contains(s.sql, "WHERE") {
		s.sql = fmt.Sprintf("%v OR %v%v'%v'", s.sql, key, opt, val)
	} else {
		s.sql = fmt.Sprintf("%v WHERE %v%v'%v'", s.sql, key, opt, val)
	}
	return s
}

func (s *Sql) In(column string, args []string) *Sql {
	if len(args) > 0 {
		var values string
		for _, v := range args {
			values += fmt.Sprintf("'%v', ", v)
		}
		values = fmt.Sprintf("(%v)", strings.Trim(values, ", "))
		if strings.Contains(s.sql, "WHERE") {
			s.sql = fmt.Sprintf("%v AND %v", s.sql, values)
		} else {
			s.sql = fmt.Sprintf("%v WHERE %v IN %v", s.sql, column, values)
		}
	}
	return s
}

func (s *Sql) OrIn(column string, args []string) *Sql {
	if len(args) > 0 {
		var values string
		for _, v := range args {
			values += fmt.Sprintf("'%v', ", v)
		}
		values = fmt.Sprintf("(%v)", strings.Trim(values, ", "))
		if strings.Contains(s.sql, "WHERE") {
			s.sql = fmt.Sprintf("%v OR %v", s.sql, values)
		} else {
			s.sql = fmt.Sprintf("%v WHERE %v IN %v", s.sql, column, values)
		}
	}
	return s
}

func (s *Sql) Between(column, val1, val2 string) *Sql {
	if strings.Contains(s.sql, "WHERE") {
		s.sql = fmt.Sprintf("%v AND %v BETWEEN '%v' AND '%v'", s.sql, column, val1, val2)
	} else {
		s.sql = fmt.Sprintf("%v WHERE %v BETWEEN '%v' AND '%v'", s.sql, column, val1, val2)
	}
	return s
}

func (s *Sql) OrBetween(column, val1, val2 string) *Sql {
	if strings.Contains(s.sql, "WHERE") {
		s.sql = fmt.Sprintf("%v OR %v BETWEEN '%v' AND '%v'", s.sql, column, val1, val2)
	} else {
		s.sql = fmt.Sprintf("%v WHERE %v BETWEEN '%v' AND '%v'", s.sql, column, val1, val2)
	}
	return s
}

func (s *Sql) Join(joinName, table, equal string) *Sql {
	s.sql = fmt.Sprintf("%v %v JOIN %v ON %v", s.sql, joinName, table, equal)
	return s
}

func (s *Sql) Limit(start, limit int) *Sql {
	s.sql = fmt.Sprintf("%v LIMIT %v,%v", s.sql, start, limit)
	return s
}

func (s *Sql) GroupBy(column string) *Sql {
	s.sql = fmt.Sprintf("%v GROUP BY %v", s.sql, column)
	return s
}

func (s *Sql) OrderBy(column, sortName string) *Sql {
	s.sql = fmt.Sprintf("%v ORDER BY %v %v", s.sql, column, sortName)
	return s
}

func (s *Sql) Union(sql string) *Sql {
	s.sql = s.sql +" UNION "+ sql
	return s
}

func (s *Sql) Get() string {
	return s.sql
}

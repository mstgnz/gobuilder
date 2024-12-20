package gobuilder

import (
	"reflect"
	"testing"
)

var (
	query          string
	params         []any
	queryExpected  string
	paramsExpected []any
	gb             = NewGoBuilder(Postgres)
)

func TestSql_Select(t *testing.T) {
	queryExpected = "SELECT * FROM users"
	paramsExpected = []any{}
	query, params = gb.Table("users").Select().Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "SELECT * FROM users"
	query = gb.Table("users").Select().Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_Select_With_Columns(t *testing.T) {
	queryExpected = "SELECT firstname, lastname FROM users"
	paramsExpected = []any{}
	query, params = gb.Table("users").Select("firstname", "lastname").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "SELECT firstname, lastname FROM users"
	query = gb.Table("users").Select("firstname", "lastname").Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_Distinct(t *testing.T) {
	queryExpected = "SELECT DISTINCT name, age FROM users"
	paramsExpected := []any{}
	query, params = gb.Table("users").SelectDistinct("name", "age").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "SELECT DISTINCT name, age FROM users"
	query = gb.Table("users").SelectDistinct("name", "age").Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_Insert(t *testing.T) {
	queryExpected = "INSERT INTO users (firstname, lastname) VALUES ($1, $2)"
	paramsExpected := []any{"Mesut", "GENEZ"}
	args := map[string]any{"firstname": "Mesut", "lastname": "GENEZ"}
	query, params = gb.Table("users").Insert(args).Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "INSERT INTO users (firstname, lastname) VALUES ('Mesut', 'GENEZ')"
	query = gb.Table("users").Insert(args).Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_Update(t *testing.T) {
	// sometimes it may not pass the test because the map is unordered
	args := map[string]any{"firstname": "Mesut", "lastname": "GENEZ"}
	queryExpected = "UPDATE users SET firstname = $1, lastname = $2"
	paramsExpected := []any{"Mesut", "GENEZ"}
	query, params = gb.Table("users").Update(args).Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "UPDATE users SET firstname = 'Mesut', lastname = 'GENEZ'"
	query = gb.Table("users").Update(args).Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_Delete(t *testing.T) {
	queryExpected = "DELETE FROM users"
	paramsExpected := []any{}
	query, params = gb.Table("users").Delete().Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "DELETE FROM users"
	query = gb.Table("users").Delete().Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_Where(t *testing.T) {
	queryExpected = "WHERE firstname = $1"
	paramsExpected := []any{"Mesut"}
	query, params = gb.Where("firstname", "=", "Mesut").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "WHERE firstname = 'Mesut'"
	query = gb.Where("firstname", "=", "Mesut").Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_WhereWithInt(t *testing.T) {
	queryExpected = "WHERE id = $1"
	paramsExpected := []any{55}
	query, params = gb.Where("id", "=", 55).Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "WHERE id = 55"
	query = gb.Where("id", "=", 55).Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_WhereWithFloat(t *testing.T) {
	queryExpected = "WHERE amount = $1"
	paramsExpected := []any{55.5}
	query, params = gb.Where("amount", "=", 55.5).Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "WHERE amount = 55.5"
	query = gb.Where("amount", "=", 55.5).Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_Where_With_And(t *testing.T) {
	queryExpected = "WHERE firstname = $1 AND lastname = $2"
	paramsExpected := []any{"Mesut", "GENEZ"}
	query, params = gb.Where("firstname", "=", "Mesut").Where("lastname", "=", "GENEZ").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "WHERE firstname = 'Mesut' AND lastname = 'GENEZ'"
	query = gb.Where("firstname", "=", "Mesut").Where("lastname", "=", "GENEZ").Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_OrWhere(t *testing.T) {
	queryExpected = "WHERE firstname = $1"
	paramsExpected := []any{"Mesut"}
	query, params = gb.OrWhere("firstname", "=", "Mesut").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "WHERE firstname = 'Mesut'"
	query = gb.OrWhere("firstname", "=", "Mesut").Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_OrWhere_With_And(t *testing.T) {
	queryExpected = "WHERE firstname = $1 OR lastname = $2"
	paramsExpected := []any{"Mesut", "GENEZ"}
	query, params = gb.OrWhere("firstname", "=", "Mesut").OrWhere("lastname", "=", "GENEZ").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "WHERE firstname = 'Mesut' OR lastname = 'GENEZ'"
	query = gb.OrWhere("firstname", "=", "Mesut").OrWhere("lastname", "=", "GENEZ").Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_In(t *testing.T) {
	queryExpected = "WHERE firstname IN ($1, $2)"
	paramsExpected := []any{"Mesut", 33}
	query, params = gb.In("firstname", "Mesut", 33).Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "WHERE firstname IN ('Mesut', 33)"
	query = gb.In("firstname", "Mesut", 33).Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_InWithInt(t *testing.T) {
	queryExpected = "WHERE id IN ($1, $2, $3)"
	paramsExpected := []any{12, 34, 55}
	query, params = gb.In("id", 12, 34, 55).Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "WHERE id IN (12, 34, 55)"
	query = gb.In("id", 12, 34, 55).Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_InAnd(t *testing.T) {
	queryExpected = "WHERE firstname IN ($1) AND lastname IN ($2)"
	paramsExpected := []any{"Mesut", "GENEZ"}
	query, params = gb.In("firstname", "Mesut").In("lastname", "GENEZ").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "WHERE firstname IN ('Mesut') AND lastname IN ('GENEZ')"
	query = gb.In("firstname", "Mesut").In("lastname", "GENEZ").Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_OrIn(t *testing.T) {
	queryExpected = "WHERE firstname IN ($1)"
	paramsExpected := []any{"Mesut"}
	query, params = gb.OrIn("firstname", "Mesut").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "WHERE firstname IN ('Mesut')"
	query = gb.OrIn("firstname", "Mesut").Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_OrInAnd(t *testing.T) {
	queryExpected = "WHERE firstname IN ($1, $2) OR lastname IN ($3)"
	paramsExpected := []any{"Mesut", "GENEZ", "GENEZ"}
	query, params = gb.OrIn("firstname", "Mesut", "GENEZ").OrIn("lastname", "GENEZ").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "WHERE firstname IN ('Mesut', 'GENEZ') OR lastname IN ('GENEZ')"
	query = gb.OrIn("firstname", "Mesut", "GENEZ").OrIn("lastname", "GENEZ").Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_Between(t *testing.T) {
	queryExpected = "WHERE firstname BETWEEN $1 AND $2"
	paramsExpected := []any{"Mesut", "GENEZ"}
	query, params = gb.Between("firstname", "Mesut", "GENEZ").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "WHERE firstname BETWEEN 'Mesut' AND 'GENEZ'"
	query = gb.Between("firstname", "Mesut", "GENEZ").Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_BetweenWithInt(t *testing.T) {
	queryExpected = "WHERE id BETWEEN $1 AND $2"
	paramsExpected := []any{12, 55}
	query, params = gb.Between("id", 12, 55).Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "WHERE id BETWEEN 12 AND 55"
	query = gb.Between("id", 12, 55).Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_BetweenWithWhere(t *testing.T) {
	queryExpected = "WHERE firstname BETWEEN $1 AND $2 AND lastname BETWEEN $3 AND $4"
	paramsExpected := []any{"Mesut", "GENEZ", "Mesut", "GENEZ"}
	query, params = gb.Between("firstname", "Mesut", "GENEZ").Between("lastname", "Mesut", "GENEZ").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "WHERE firstname BETWEEN 'Mesut' AND 'GENEZ' AND lastname BETWEEN 'Mesut' AND 'GENEZ'"
	query = gb.Between("firstname", "Mesut", "GENEZ").Between("lastname", "Mesut", "GENEZ").Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_OrBetween(t *testing.T) {
	queryExpected = "WHERE firstname BETWEEN $1 AND $2"
	paramsExpected := []any{"Mesut", "GENEZ"}
	query, params = gb.OrBetween("firstname", "Mesut", "GENEZ").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "WHERE firstname BETWEEN 'Mesut' AND 'GENEZ'"
	query = gb.OrBetween("firstname", "Mesut", "GENEZ").Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_OrBetweenWithWhere(t *testing.T) {
	queryExpected = "WHERE firstname BETWEEN $1 AND $2 OR lastname BETWEEN $3 AND $4"
	paramsExpected := []any{"Mesut", "GENEZ", "Mesut", "GENEZ"}
	query, params = gb.OrBetween("firstname", "Mesut", "GENEZ").OrBetween("lastname", "Mesut", "GENEZ").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "WHERE firstname BETWEEN 'Mesut' AND 'GENEZ' OR lastname BETWEEN 'Mesut' AND 'GENEZ'"
	query = gb.OrBetween("firstname", "Mesut", "GENEZ").OrBetween("lastname", "Mesut", "GENEZ").Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_IsNull(t *testing.T) {
	queryExpected = "SELECT * FROM users WHERE age = $1 AND name IS NULL"
	paramsExpected := []any{30}
	query, params = gb.Table("users").Select().Where("age", "=", 30).IsNull("name").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "SELECT * FROM users WHERE age = 30 AND name IS NULL"
	query = gb.Table("users").Select().Where("age", "=", 30).IsNull("name").Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_IsNotNull(t *testing.T) {
	queryExpected = "SELECT * FROM users WHERE age = $1 AND name IS NOT NULL"
	paramsExpected := []any{30}
	query, params = gb.Table("users").Select().Where("age", "=", 30).IsNotNull("name").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "SELECT * FROM users WHERE age = 30 AND name IS NOT NULL"
	query = gb.Table("users").Select().Where("age", "=", 30).IsNotNull("name").Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_Having(t *testing.T) {
	queryExpected = "SELECT * FROM orders GROUP BY user_id HAVING COUNT(*) > 1"
	paramsExpected := []any{}
	query, params = gb.Table("orders").Select().GroupBy("user_id").Having("COUNT(*) > 1").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "SELECT * FROM orders GROUP BY user_id HAVING COUNT(*) > 1"
	query = gb.Table("orders").Select().GroupBy("user_id").Having("COUNT(*) > 1").Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_Join(t *testing.T) {
	queryExpected = "INNER JOIN users ON roles"
	paramsExpected := []any{}
	query, params = gb.Join("INNER", "users", "roles").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "INNER JOIN users ON roles"
	query = gb.Join("INNER", "users", "roles").Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_Limit(t *testing.T) {
	queryExpected = "LIMIT 1, 5"
	paramsExpected := []any{}
	query, params = gb.Limit(1, 5).Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "LIMIT 1, 5"
	query = gb.Limit(1, 5).Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_Group(t *testing.T) {
	queryExpected = "GROUP BY firstname"
	paramsExpected := []any{}
	query, params = gb.GroupBy("firstname").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_GroupMultiple(t *testing.T) {
	queryExpected = "GROUP BY firstname, lastname, email"
	paramsExpected := []any{}
	query, params = gb.GroupBy("firstname", "lastname", "email").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_OrderBy(t *testing.T) {
	queryExpected = "ORDER BY firstname ASC"
	paramsExpected := []any{}
	query, params = gb.OrderBy("firstname").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_OrderByMultiple(t *testing.T) {
	queryExpected = "ORDER BY firstname, lastname ASC"
	paramsExpected := []any{}
	query, params = gb.OrderBy("firstname", "lastname").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_OrderByDesc(t *testing.T) {
	queryExpected = "ORDER BY firstname DESC"
	paramsExpected := []any{}
	query, params = gb.OrderByDesc("firstname").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_OrderByDescMultiple(t *testing.T) {
	queryExpected = "ORDER BY firstname, lastname DESC"
	paramsExpected := []any{}
	query, params = gb.OrderByDesc("firstname", "lastname").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_Union(t *testing.T) {
	queryExpected = "UNION select * from companies"
	paramsExpected := []any{}
	query, params = gb.Union("select * from companies").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_SubQuery(t *testing.T) {
	paramsExpected := []any{"SELECT id FROM users WHERE age > $1"}
	mainBuilder := NewGoBuilder(Postgres)
	subBuilder := NewGoBuilder(Postgres)

	query, params = subBuilder.Table("users").Select("id").Where("age", ">", 30).Prepare()
	query, params = mainBuilder.Table("orders").Select("order_id", "user_id").In("user_id", query).Prepare()

	queryExpected = "SELECT order_id, user_id FROM orders WHERE user_id IN ($1)"

	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

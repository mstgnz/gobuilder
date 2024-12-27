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

func TestSql_Create(t *testing.T) {
	queryExpected = "INSERT INTO users (firstname, lastname) VALUES ($1, $2)"
	paramsExpected := []any{"Mesut", "GENEZ"}
	args := map[string]any{"firstname": "Mesut", "lastname": "GENEZ"}
	query, params = gb.Table("users").Create(args).Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	queryExpected = "INSERT INTO users (firstname, lastname) VALUES ('Mesut', 'GENEZ')"
	query = gb.Table("users").Create(args).Sql()
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
	query = gb.Table("orders").Select().GroupBy("user_id").Having("COUNT(*) > 1").Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_Join(t *testing.T) {
	queryExpected = "INNER JOIN users ON users.id = roles.user_id"
	paramsExpected := []any{}
	query, params = gb.Join("users", "users.id", "=", "roles.user_id").Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
	query = gb.Join("users", "users.id", "=", "roles.user_id").Sql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
}

func TestSql_Limit(t *testing.T) {
	queryExpected = "OFFSET 1 LIMIT 5"
	paramsExpected := []any{}
	query, params = gb.Limit(1, 5).Prepare()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
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

func TestSql_ErrorHandling(t *testing.T) {
	gb := NewGoBuilder(Postgres)

	// Table olmadan sorgu oluşturma denemesi
	gb.Select("id").Where("age", ">", 30)
	if gb.Error() == nil {
		t.Error("Table olmadan sorgu oluşturulduğunda hata vermeli")
	}
}

func TestSql_JsonOperations(t *testing.T) {
	queryExpected := "SELECT * FROM users WHERE data->>'name' = $1"
	paramsExpected := []any{"John"}

	query, params := gb.Table("users").Select().Where("data->>'name'", "=", "John").Prepare()

	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_AggregateFunctions(t *testing.T) {
	testCases := []struct {
		name     string
		builder  func() (string, []any)
		expected string
		params   []any
	}{
		{
			name: "Count",
			builder: func() (string, []any) {
				return gb.Table("users").Select("COUNT(*) as total").Prepare()
			},
			expected: "SELECT COUNT(*) as total FROM users",
			params:   []any{},
		},
		{
			name: "Sum",
			builder: func() (string, []any) {
				return gb.Table("orders").Select("SUM(amount) as total_amount").Prepare()
			},
			expected: "SELECT SUM(amount) as total_amount FROM orders",
			params:   []any{},
		},
		{
			name: "Multiple Aggregates",
			builder: func() (string, []any) {
				return gb.Table("orders").
					Select(
						"customer_id",
						"COUNT(*) as total_orders",
						"SUM(amount) as total_amount",
						"AVG(amount) as avg_amount",
						"MIN(amount) as min_amount",
						"MAX(amount) as max_amount",
					).
					GroupBy("customer_id").
					Having("COUNT(*) > ?", 5).
					OrderBy("total_amount").
					Prepare()
			},
			expected: "SELECT customer_id, COUNT(*) as total_orders, SUM(amount) as total_amount, AVG(amount) as avg_amount, MIN(amount) as min_amount, MAX(amount) as max_amount FROM orders GROUP BY customer_id HAVING COUNT(*) > $1 ORDER BY total_amount ASC",
			params:   []any{5},
		},
		{
			name: "Aggregate with CASE",
			builder: func() (string, []any) {
				return gb.Table("orders").
					Select(
						"status",
						"COUNT(*) as total",
						"SUM(CASE WHEN amount > 1000 THEN 1 ELSE 0 END) as high_value_orders",
					).
					GroupBy("status").
					Prepare()
			},
			expected: "SELECT status, COUNT(*) as total, SUM(CASE WHEN amount > 1000 THEN 1 ELSE 0 END) as high_value_orders FROM orders GROUP BY status",
			params:   []any{},
		},
		{
			name: "Aggregate with Having and Complex Conditions",
			builder: func() (string, []any) {
				return gb.Table("sales").
					Select(
						"department",
						"COUNT(*) as total_sales",
						"AVG(amount) as avg_sale",
					).
					GroupBy("department").
					Having("COUNT(*) > ? AND AVG(amount) > ?", 10, 1000).
					OrderByDesc("avg_sale").
					Prepare()
			},
			expected: "SELECT department, COUNT(*) as total_sales, AVG(amount) as avg_sale FROM sales GROUP BY department HAVING COUNT(*) > $1 AND AVG(amount) > $2 ORDER BY avg_sale DESC",
			params:   []any{10, 1000},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, params := tc.builder()
			if query != tc.expected {
				t.Errorf("expected query %v, got %v", tc.expected, query)
			}
			if !reflect.DeepEqual(params, tc.params) {
				t.Errorf("expected params %v, got %v", tc.params, params)
			}
		})
	}
}

func TestSql_WithCTE(t *testing.T) {
	subQuery := NewGoBuilder(Postgres).Table("orders").Select("user_id", "COUNT(*) as order_count").GroupBy("user_id")
	queryExpected := "WITH user_orders AS (SELECT user_id, COUNT(*) as order_count FROM orders GROUP BY user_id) SELECT * FROM user_orders WHERE order_count > $1"
	paramsExpected := []any{5}

	query, params := gb.With("user_orders", subQuery).Table("user_orders").Select().Where("order_count", ">", 5).Prepare()

	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_BatchInsert(t *testing.T) {
	records := []map[string]any{
		{"name": "John", "age": 30},
		{"name": "Jane", "age": 25},
	}

	queryExpected := "INSERT INTO users (age, name) VALUES ($1, $2), ($3, $4)"
	paramsExpected := []any{30, "John", 25, "Jane"}

	query, params := gb.Table("users").CreateBatch(records).Prepare()

	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_LockMechanism(t *testing.T) {
	queryExpected := "SELECT * FROM users FOR UPDATE"
	paramsExpected := []any{}

	query, params := gb.Table("users").Select().Lock("FOR UPDATE").Prepare()

	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_ConditionalClauses(t *testing.T) {
	age := 30
	name := "John"

	queryExpected := "SELECT * FROM users WHERE age > $1 AND name = $2"
	paramsExpected := []any{30, "John"}

	query, params := gb.Table("users").Select().WhenThen(age > 0, func(b *GoBuilder) *GoBuilder { return b.Where("age", ">", age) }, nil).WhenThen(name != "", func(b *GoBuilder) *GoBuilder { return b.Where("name", "=", name) }, nil).Prepare()

	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_DatabaseSpecificFeatures(t *testing.T) {
	testCases := []struct {
		name     string
		dialect  SQLDialect
		builder  func(gb *GoBuilder) (string, []any)
		expected string
		params   []any
	}{
		{
			name:    "PostgreSQL - RETURNING Clause",
			dialect: Postgres,
			builder: func(gb *GoBuilder) (string, []any) {
				return gb.Table("users").Create(map[string]any{"name": "John"}, "id").Prepare()
			},
			expected: "INSERT INTO users (name) VALUES ($1) RETURNING id",
			params:   []any{"John"},
		},
		{
			name:    "MySQL - ON DUPLICATE KEY UPDATE",
			dialect: MySQL,
			builder: func(gb *GoBuilder) (string, []any) {
				return gb.Table("users").Create(map[string]any{"id": 1, "name": "John"}).OnDuplicateKeyUpdate(map[string]any{"name": "John"}).Prepare()
			},
			expected: "INSERT INTO users (id, name) VALUES (?, ?) ON DUPLICATE KEY UPDATE name = ?",
			params:   []any{1, "John", "John"},
		},
		{
			name:    "SQLServer - TOP Clause",
			dialect: SQLServer,
			builder: func(gb *GoBuilder) (string, []any) {
				return gb.Table("users").Select().Top(10).Prepare()
			},
			expected: "SELECT TOP 10 * FROM users",
			params:   []any{},
		},
		{
			name:    "Oracle - ROWNUM",
			dialect: Oracle,
			builder: func(gb *GoBuilder) (string, []any) {
				return gb.Table("users").Select().Where("ROWNUM", "<=", 10).Prepare()
			},
			expected: "SELECT * FROM users WHERE ROWNUM <= :1",
			params:   []any{10},
		},
		{
			name:    "SQLite - PRAGMA",
			dialect: SQLite,
			builder: func(gb *GoBuilder) (string, []any) {
				return gb.Pragma("foreign_keys", "ON").Prepare()
			},
			expected: "PRAGMA foreign_keys = ON",
			params:   []any{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gb := NewGoBuilder(tc.dialect)
			query, params := tc.builder(gb)

			if query != tc.expected {
				t.Errorf("expected query %v, got %v", tc.expected, query)
			}
			if !reflect.DeepEqual(params, tc.params) {
				t.Errorf("expected params %v, got %v", tc.params, params)
			}
		})
	}
}

func TestSql_ComplexJoins(t *testing.T) {
	testCases := []struct {
		name     string
		builder  func() (string, []any)
		expected string
		params   []any
	}{
		{
			name: "Multiple Inner Joins",
			builder: func() (string, []any) {
				return gb.Table("orders").
					Select("orders.id", "customers.name", "products.title").
					Join("customers", "customers.id", "=", "orders.customer_id").
					Join("products", "products.id", "=", "orders.product_id").
					Where("orders.status", "=", "pending").
					Prepare()
			},
			expected: "SELECT orders.id, customers.name, products.title FROM orders INNER JOIN customers ON customers.id = orders.customer_id INNER JOIN products ON products.id = orders.product_id WHERE orders.status = $1",
			params:   []any{"pending"},
		},
		{
			name: "Mixed Join Types",
			builder: func() (string, []any) {
				return gb.Table("users").
					Select("users.name", "orders.total", "addresses.city").
					LeftJoin("orders", "orders.user_id", "=", "users.id").
					RightJoin("addresses", "addresses.user_id", "=", "users.id").
					Where("users.status", "=", "active").
					Prepare()
			},
			expected: "SELECT users.name, orders.total, addresses.city FROM users LEFT JOIN orders ON orders.user_id = users.id RIGHT JOIN addresses ON addresses.user_id = users.id WHERE users.status = $1",
			params:   []any{"active"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, params := tc.builder()
			if query != tc.expected {
				t.Errorf("expected query %v, got %v", tc.expected, query)
			}
			if !reflect.DeepEqual(params, tc.params) {
				t.Errorf("expected params %v, got %v", tc.params, params)
			}
		})
	}
}

func TestSql_ComplexSubqueries(t *testing.T) {
	testCases := []struct {
		name     string
		builder  func() (string, []any)
		expected string
		params   []any
	}{
		{
			name: "Subquery in WHERE",
			builder: func() (string, []any) {
				subQuery := NewGoBuilder(Postgres).
					Table("orders").
					Select("customer_id").
					Where("total", ">", 1000)
				return gb.Table("customers").
					Select("name", "email").
					Where("id", "IN", subQuery).
					Prepare()
			},
			expected: "SELECT name, email FROM customers WHERE id IN (SELECT customer_id FROM orders WHERE total > $1)",
			params:   []any{1000},
		},
		{
			name: "Subquery in FROM",
			builder: func() (string, []any) {
				subQuery := NewGoBuilder(Postgres).
					Table("orders").
					Select("customer_id", "COUNT(*) as order_count").
					GroupBy("customer_id").
					Having("COUNT(*) > ?", 5)
				return gb.Table("("+subQuery.Sql()+") as order_stats").
					Select("customer_id", "order_count").
					OrderBy("order_count").
					Prepare()
			},
			expected: "SELECT customer_id, order_count FROM (SELECT customer_id, COUNT(*) as order_count FROM orders GROUP BY customer_id HAVING COUNT(*) > 5) as order_stats ORDER BY order_count ASC",
			params:   []any{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, params := tc.builder()
			if query != tc.expected {
				t.Errorf("expected query %v, got %v", tc.expected, query)
			}
			if !reflect.DeepEqual(params, tc.params) {
				t.Errorf("expected params %v, got %v", tc.params, params)
			}
		})
	}
}

func TestSql_WindowFunctions(t *testing.T) {
	testCases := []struct {
		name     string
		builder  func() (string, []any)
		expected string
		params   []any
	}{
		{
			name: "ROW_NUMBER",
			builder: func() (string, []any) {
				return gb.Table("sales").
					Select(
						"employee_id",
						"amount",
						"ROW_NUMBER() OVER (PARTITION BY employee_id ORDER BY amount DESC) as rank",
					).
					Where("amount", ">", 1000).
					Prepare()
			},
			expected: "SELECT employee_id, amount, ROW_NUMBER() OVER (PARTITION BY employee_id ORDER BY amount DESC) as rank FROM sales WHERE amount > $1",
			params:   []any{1000},
		},
		{
			name: "Multiple Window Functions",
			builder: func() (string, []any) {
				return gb.Table("products").
					Select(
						"category_id",
						"price",
						"ROW_NUMBER() OVER (PARTITION BY category_id ORDER BY price DESC) as price_rank",
						"RANK() OVER (PARTITION BY category_id ORDER BY price DESC) as overall_rank",
						"DENSE_RANK() OVER (PARTITION BY category_id ORDER BY price DESC) as dense_rank",
					).
					Where("price", ">", 100).
					Prepare()
			},
			expected: "SELECT category_id, price, ROW_NUMBER() OVER (PARTITION BY category_id ORDER BY price DESC) as price_rank, RANK() OVER (PARTITION BY category_id ORDER BY price DESC) as overall_rank, DENSE_RANK() OVER (PARTITION BY category_id ORDER BY price DESC) as dense_rank FROM products WHERE price > $1",
			params:   []any{100},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, params := tc.builder()
			if query != tc.expected {
				t.Errorf("expected query %v, got %v", tc.expected, query)
			}
			if !reflect.DeepEqual(params, tc.params) {
				t.Errorf("expected params %v, got %v", tc.params, params)
			}
		})
	}
}

func TestSql_CaseExpressions(t *testing.T) {
	testCases := []struct {
		name     string
		builder  func() (string, []any)
		expected string
		params   []any
	}{
		{
			name: "Simple CASE",
			builder: func() (string, []any) {
				return gb.Table("orders").
					Select(
						"id",
						"amount",
						"CASE status WHEN 'pending' THEN 'In Progress' WHEN 'completed' THEN 'Done' ELSE 'Unknown' END as status_text",
					).
					Where("amount", ">", 100).
					Prepare()
			},
			expected: "SELECT id, amount, CASE status WHEN 'pending' THEN 'In Progress' WHEN 'completed' THEN 'Done' ELSE 'Unknown' END as status_text FROM orders WHERE amount > $1",
			params:   []any{100},
		},
		{
			name: "Searched CASE",
			builder: func() (string, []any) {
				return gb.Table("products").
					Select(
						"id",
						"name",
						"price",
						"CASE WHEN price < 100 THEN 'Budget' WHEN price < 500 THEN 'Regular' ELSE 'Premium' END as category",
					).
					OrderBy("price").
					Prepare()
			},
			expected: "SELECT id, name, price, CASE WHEN price < 100 THEN 'Budget' WHEN price < 500 THEN 'Regular' ELSE 'Premium' END as category FROM products ORDER BY price ASC",
			params:   []any{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			query, params := tc.builder()
			if query != tc.expected {
				t.Errorf("expected query %v, got %v", tc.expected, query)
			}
			if !reflect.DeepEqual(params, tc.params) {
				t.Errorf("expected params %v, got %v", tc.params, params)
			}
		})
	}
}

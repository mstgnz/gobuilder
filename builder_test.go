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
	s              = NewGoBuilder(Postgres)
)

func TestSql_Select(t *testing.T) {
	queryExpected = "SELECT * FROM users"
	paramsExpected = []any{}
	query, params = s.Table("users").Select().ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_Select_With_Columns(t *testing.T) {
	queryExpected = "SELECT firstname, lastname FROM users"
	paramsExpected = []any{}
	query, params = s.Table("users").Select("firstname", "lastname").ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_Distinct(t *testing.T) {
	queryExpected = "SELECT DISTINCT name, age FROM users"
	paramsExpected := []any{}
	query, params = s.Table("users").SelectDistinct("name", "age").ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_Insert(t *testing.T) {
	queryExpected = "INSERT INTO users (firstname, lastname) VALUES ($1, $2)"
	paramsExpected := []any{"Mesut", "GENEZ"}
	args := map[string]any{"firstname": "Mesut", "lastname": "GENEZ"}
	query, params = s.Table("users").Insert(args).ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_Update(t *testing.T) {
	// sometimes it may not pass the test because the map is unordered
	queryExpected = "UPDATE users SET firstname = $1, lastname = $2"
	paramsExpected := []any{"Mesut", "GENEZ"}
	query, params = s.Table("users").Update(map[string]any{"firstname": "Mesut", "lastname": "GENEZ"}).ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_Delete(t *testing.T) {
	queryExpected = "DELETE FROM users"
	paramsExpected := []any{}
	query, params = s.Table("users").Delete().ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_Where(t *testing.T) {
	queryExpected = "WHERE firstname = $1"
	paramsExpected := []any{"Mesut"}
	query, params = s.Where("firstname", "=", "Mesut").ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_WhereWithInt(t *testing.T) {
	queryExpected = "WHERE id = $1"
	paramsExpected := []any{55}
	query, params = s.Where("id", "=", 55).ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_WhereWithFloat(t *testing.T) {
	queryExpected = "WHERE amount = $1"
	paramsExpected := []any{55.5}
	query, params = s.Where("amount", "=", 55.5).ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_Where_With_And(t *testing.T) {
	queryExpected = "WHERE firstname = $1 AND lastname = $2"
	paramsExpected := []any{"Mesut", "GENEZ"}
	query, params = s.Where("firstname", "=", "Mesut").Where("lastname", "=", "GENEZ").ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_OrWhere(t *testing.T) {
	queryExpected = "WHERE firstname = $1"
	paramsExpected := []any{"Mesut"}
	query, params = s.OrWhere("firstname", "=", "Mesut").ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_OrWhere_With_And(t *testing.T) {
	queryExpected = "WHERE firstname = $1 OR lastname = $2"
	paramsExpected := []any{"Mesut", "GENEZ"}
	query, params = s.OrWhere("firstname", "=", "Mesut").OrWhere("lastname", "=", "GENEZ").ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_In(t *testing.T) {
	queryExpected = "WHERE firstname IN ($1, $2)"
	paramsExpected := []any{"Mesut", 33}
	query, params = s.In("firstname", "Mesut", 33).ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_InWithInt(t *testing.T) {
	queryExpected = "WHERE id IN ($1, $2, $3)"
	paramsExpected := []any{12, 34, 55}
	query, params = s.In("id", 12, 34, 55).ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_InAnd(t *testing.T) {
	queryExpected = "WHERE firstname IN ($1) AND lastname IN ($2)"
	paramsExpected := []any{"Mesut", "GENEZ"}
	query, params = s.In("firstname", "Mesut").In("lastname", "GENEZ").ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_OrIn(t *testing.T) {
	queryExpected = "WHERE firstname IN ($1)"
	paramsExpected := []any{"Mesut"}
	query, params = s.OrIn("firstname", "Mesut").ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_OrInAnd(t *testing.T) {
	queryExpected = "WHERE firstname IN ($1, $2) OR lastname IN ($3)"
	paramsExpected := []any{"Mesut", "GENEZ", "GENEZ"}
	query, params = s.OrIn("firstname", "Mesut", "GENEZ").OrIn("lastname", "GENEZ").ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_Between(t *testing.T) {
	queryExpected = "WHERE firstname BETWEEN $1 AND $2"
	paramsExpected := []any{"Mesut", "GENEZ"}
	query, params = s.Between("firstname", "Mesut", "GENEZ").ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_BetweenWithInt(t *testing.T) {
	queryExpected = "WHERE id BETWEEN $1 AND $2"
	paramsExpected := []any{12, 55}
	query, params = s.Between("id", 12, 55).ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_BetweenWithWhere(t *testing.T) {
	queryExpected = "WHERE firstname BETWEEN $1 AND $2 AND lastname BETWEEN $3 AND $4"
	paramsExpected := []any{"Mesut", "GENEZ", "Mesut", "GENEZ"}
	query, params = s.Between("firstname", "Mesut", "GENEZ").Between("lastname", "Mesut", "GENEZ").ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_OrBetween(t *testing.T) {
	queryExpected = "WHERE firstname BETWEEN $1 AND $2"
	paramsExpected := []any{"Mesut", "GENEZ"}
	query, params = s.OrBetween("firstname", "Mesut", "GENEZ").ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_OrBetweenWithWhere(t *testing.T) {
	queryExpected = "WHERE firstname BETWEEN $1 AND $2 OR lastname BETWEEN $3 AND $4"
	paramsExpected := []any{"Mesut", "GENEZ", "Mesut", "GENEZ"}
	query, params = s.OrBetween("firstname", "Mesut", "GENEZ").OrBetween("lastname", "Mesut", "GENEZ").ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_IsNull(t *testing.T) {
	queryExpected = "SELECT * FROM users WHERE age = $1 AND name IS NULL"
	paramsExpected := []any{30}
	query, params = s.Table("users").Select().Where("age", "=", 30).IsNull("name").ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_IsNotNull(t *testing.T) {
	queryExpected = "SELECT * FROM users WHERE age = $1 AND name IS NOT NULL"
	paramsExpected := []any{30}
	query, params = s.Table("users").Select().Where("age", "=", 30).IsNotNull("name").ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_Having(t *testing.T) {
	queryExpected = "SELECT * FROM orders GROUP BY user_id HAVING COUNT(*) > 1"
	paramsExpected := []any{}
	query, params = s.Table("orders").Select().GroupBy("user_id").Having("COUNT(*) > 1").ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_Join(t *testing.T) {
	queryExpected = "INNER JOIN users ON roles"
	paramsExpected := []any{}
	query, params = s.Join("INNER", "users", "roles").ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_Limit(t *testing.T) {
	queryExpected = "LIMIT 1, 5"
	paramsExpected := []any{}
	query, params = s.Limit(1, 5).ToSql()
	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

func TestSql_Group(t *testing.T) {
	queryExpected = "GROUP BY firstname"
	paramsExpected := []any{}
	query, params = s.GroupBy("firstname").ToSql()
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
	query, params = s.GroupBy("firstname", "lastname", "email").ToSql()
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
	query, params = s.OrderBy("firstname").ToSql()
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
	query, params = s.OrderBy("firstname", "lastname").ToSql()
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
	query, params = s.OrderByDesc("firstname").ToSql()
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
	query, params = s.OrderByDesc("firstname", "lastname").ToSql()
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
	query, params = s.Union("select * from companies").ToSql()
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

	query, params = subBuilder.Table("users").Select("id").Where("age", ">", 30).ToSql()
	query, params = mainBuilder.Table("orders").Select("order_id", "user_id").In("user_id", query).ToSql()

	queryExpected = "SELECT order_id, user_id FROM orders WHERE user_id IN ($1)"

	if !reflect.DeepEqual(queryExpected, query) {
		t.Errorf("queryExpected = %v, query %v", queryExpected, query)
	}
	if !reflect.DeepEqual(paramsExpected, params) {
		t.Errorf("paramsExpected = %v, params %v", paramsExpected, params)
	}
}

package gobuilder

import (
	"reflect"
	"testing"
)

var (
	got      string
	expected string
	s        GoBuilder
)

func TestSql_Select(t *testing.T) {
	expected = "SELECT * FROM users"
	got = s.Table("users").Select().ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Select_With_Columns(t *testing.T) {
	expected = "SELECT firstname, lastname FROM users"
	got = s.Table("users").Select("firstname", "lastname").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Distinct(t *testing.T) {
	expected = "SELECT DISTINCT name, age FROM users"
	got = s.Table("users").SelectDistinct("name", "age").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Insert(t *testing.T) {
	insert := map[string]any{"firstname": "Mesut", "lastname": "GENEZ"}
	expected = "INSERT INTO users (firstname, lastname) VALUES ('Mesut', 'GENEZ')"
	got = s.Table("users").Insert(insert).ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Update(t *testing.T) {
	// sometimes it may not pass the test because the map is unordered
	expected = "UPDATE users SET firstname = 'Mesut', lastname = 'GENEZ'"
	got = s.Table("users").Update(map[string]any{"firstname": "Mesut", "lastname": "GENEZ"}).ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Delete(t *testing.T) {
	expected = "DELETE FROM users"
	got = s.Table("users").Delete().ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Where(t *testing.T) {
	expected = "WHERE firstname = 'Mesut'"
	got = s.Where("firstname", "=", "Mesut").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_WhereWithInt(t *testing.T) {
	expected = "WHERE id = 55"
	got = s.Where("id", "=", 55).ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_WhereWithFloat(t *testing.T) {
	expected = "WHERE amount = 55.5"
	got = s.Where("amount", "=", 55.5).ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Where_With_And(t *testing.T) {
	expected = "WHERE firstname = 'Mesut' AND lastname = 'GENEZ'"
	got = s.Where("firstname", "=", "Mesut").Where("lastname", "=", "GENEZ").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_OrWhere(t *testing.T) {
	expected = "WHERE firstname = 'Mesut'"
	got = s.OrWhere("firstname", "=", "Mesut").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_OrWhere_With_And(t *testing.T) {
	expected = "WHERE firstname = 'Mesut' OR lastname = 'GENEZ'"
	got = s.OrWhere("firstname", "=", "Mesut").OrWhere("lastname", "=", "GENEZ").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_In(t *testing.T) {
	expected = "WHERE firstname IN ('Mesut', 33)"
	got = s.In("firstname", "Mesut", 33).ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_InWithInt(t *testing.T) {
	expected = "WHERE id IN (12, 34, 55)"
	got = s.In("id", 12, 34, 55).ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_InAnd(t *testing.T) {
	expected = "WHERE firstname IN ('Mesut') AND lastname IN ('GENEZ')"
	got = s.In("firstname", "Mesut").In("lastname", "GENEZ").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_OrIn(t *testing.T) {
	expected = "WHERE firstname IN ('Mesut')"
	got = s.OrIn("firstname", "Mesut").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_OrInAnd(t *testing.T) {
	expected = "WHERE firstname IN ('Mesut', 'GENEZ') OR lastname IN ('GENEZ')"
	got = s.OrIn("firstname", "Mesut", "GENEZ").OrIn("lastname", "GENEZ").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Between(t *testing.T) {
	expected = "WHERE firstname BETWEEN 'Mesut' AND 'GENEZ'"
	got = s.Between("firstname", "Mesut", "GENEZ").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_BetweenWithInt(t *testing.T) {
	expected = "WHERE id BETWEEN 12 AND 55"
	got = s.Between("id", 12, 55).ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_BetweenWithWhere(t *testing.T) {
	expected = "WHERE firstname BETWEEN 'Mesut' AND 'GENEZ' AND lastname BETWEEN 'Mesut' AND 'GENEZ'"
	got = s.Between("firstname", "Mesut", "GENEZ").Between("lastname", "Mesut", "GENEZ").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_OrBetween(t *testing.T) {
	expected = "WHERE firstname BETWEEN 'Mesut' AND 'GENEZ'"
	got = s.OrBetween("firstname", "Mesut", "GENEZ").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_OrBetweenWithWhere(t *testing.T) {
	expected = "WHERE firstname BETWEEN 'Mesut' AND 'GENEZ' OR lastname BETWEEN 'Mesut' AND 'GENEZ'"
	got = s.OrBetween("firstname", "Mesut", "GENEZ").OrBetween("lastname", "Mesut", "GENEZ").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_IsNull(t *testing.T) {
	expected = "SELECT * FROM users WHERE age = 30 AND name IS NULL"
	got = s.Table("users").Select().Where("age", "=", 30).IsNull("name").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_IsNotNull(t *testing.T) {
	expected = "SELECT * FROM users WHERE age = 30 AND name IS NOT NULL"
	got = s.Table("users").Select().Where("age", "=", 30).IsNotNull("name").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Having(t *testing.T) {
	expected = "SELECT * FROM orders GROUP BY user_id HAVING COUNT(*) > 1"
	got = s.Table("orders").Select().GroupBy("user_id").Having("COUNT(*) > 1").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_OrHaving(t *testing.T) {
	expected = "SELECT * FROM orders GROUP BY user_id HAVING COUNT(*) > 1 OR COUNT(*) < 5"
	got = s.Table("orders").Select().GroupBy("user_id").OrHaving("COUNT(*) > 1").OrHaving("COUNT(*) < 5").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Join(t *testing.T) {
	expected = "INNER JOIN users ON roles"
	got = s.Join("INNER", "users", "roles").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Limit(t *testing.T) {
	expected = "LIMIT 1, 5"
	got = s.Limit(1, 5).ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Group(t *testing.T) {
	expected = "GROUP BY firstname"
	got = s.GroupBy("firstname").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_GroupMultiple(t *testing.T) {
	expected = "GROUP BY firstname, lastname, email"
	got = s.GroupBy("firstname", "lastname", "email").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_OrderBy(t *testing.T) {
	expected = "ORDER BY firstname ASC"
	got = s.OrderBy("firstname").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_OrderByMultiple(t *testing.T) {
	expected = "ORDER BY firstname, lastname ASC"
	got = s.OrderBy("firstname", "lastname").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_OrderByDesc(t *testing.T) {
	expected = "ORDER BY firstname DESC"
	got = s.OrderByDesc("firstname").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_OrderByDescMultiple(t *testing.T) {
	expected = "ORDER BY firstname, lastname DESC"
	got = s.OrderByDesc("firstname", "lastname").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Union(t *testing.T) {
	expected = "UNION select * from companies"
	got = s.Union("select * from companies").ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_SubQuery(t *testing.T) {
	mainBuilder := &GoBuilder{}
	subBuilder := &GoBuilder{}

	subBuilder.Table("users").Select("id").Where("age", ">", 30)
	mainBuilder.Table("orders").Select("order_id", "user_id").Where("user_id", "IN", subBuilder)

	expected = "SELECT order_id, user_id FROM orders WHERE user_id IN (SELECT id FROM users WHERE age > 30)"
	got = mainBuilder.ToSql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

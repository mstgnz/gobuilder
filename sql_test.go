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
	s.sql = ""
	expected = "SELECT * FROM users"
	got = s.Select("users").Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Select_With_Columns(t *testing.T) {
	s.sql = ""
	expected = "SELECT firstname,lastname FROM users"
	got = s.Select("users", "firstname", "lastname").Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Insert(t *testing.T) {
	s.sql = ""
	insert := map[string]string{"firstname": "Mesut", "lastname": "GENEZ"}
	expected = "INSERT INTO users (firstname,lastname) VALUES ('Mesut','GENEZ')"
	got = s.Insert("users", insert).Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Update(t *testing.T) {
	s.sql = ""
	expected = "UPDATE users SET firstname='Mesut', lastname='GENEZ'"
	got = s.Update("users", map[string]string{"firstname": "Mesut", "lastname": "GENEZ"}).Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Delete(t *testing.T) {
	s.sql = ""
	expected = "DELETE FROM users"
	got = s.Delete("users").Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Where(t *testing.T) {
	s.sql = ""
	expected = " WHERE firstname='Mesut'"
	got = s.Where("firstname", "=", "Mesut").Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Where_With_And(t *testing.T) {
	s.sql = ""
	expected = " WHERE firstname='Mesut' AND lastname='GENEZ'"
	got = s.Where("firstname", "=", "Mesut").Where("lastname", "=", "GENEZ").Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_OrWhere(t *testing.T) {
	s.sql = ""
	expected = " WHERE firstname='Mesut'"
	got = s.OrWhere("firstname", "=", "Mesut").Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_OrWhere_With_And(t *testing.T) {
	s.sql = ""
	expected = " WHERE firstname='Mesut' OR lastname='GENEZ'"
	got = s.OrWhere("firstname", "=", "Mesut").OrWhere("lastname", "=", "GENEZ").Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_In(t *testing.T) {
	s.sql = ""
	expected = " WHERE firstname IN ('Mesut')"
	got = s.In("firstname", "Mesut").Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_InAnd(t *testing.T) {
	s.sql = ""
	expected = " WHERE firstname IN ('Mesut') AND ('GENEZ')"
	got = s.In("firstname", "Mesut").In("lastname", "GENEZ").Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_OrIn(t *testing.T) {
	s.sql = ""
	expected = " WHERE firstname IN ('Mesut')"
	got = s.OrIn("firstname", "Mesut").Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_OrInAnd(t *testing.T) {
	s.sql = ""
	expected = " WHERE firstname IN ('Mesut') OR ('GENEZ')"
	got = s.OrIn("firstname", "Mesut").OrIn("lastname", "GENEZ").Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Between(t *testing.T) {
	s.sql = ""
	expected = " WHERE firstname BETWEEN 'Mesut' AND 'GENEZ'"
	got = s.Between("firstname", "Mesut", "GENEZ").Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_BetweenWithWhere(t *testing.T) {
	s.sql = ""
	expected = " WHERE firstname BETWEEN 'Mesut' AND 'GENEZ' AND lastname BETWEEN 'Mesut' AND 'GENEZ'"
	got = s.Between("firstname", "Mesut", "GENEZ").Between("lastname", "Mesut", "GENEZ").Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_OrBetween(t *testing.T) {
	s.sql = ""
	expected = " WHERE firstname BETWEEN 'Mesut' AND 'GENEZ'"
	got = s.OrBetween("firstname", "Mesut", "GENEZ").Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_OrBetweenWithWhere(t *testing.T) {
	s.sql = ""
	expected = " WHERE firstname BETWEEN 'Mesut' AND 'GENEZ' OR lastname BETWEEN 'Mesut' AND 'GENEZ'"
	got = s.OrBetween("firstname", "Mesut", "GENEZ").OrBetween("lastname", "Mesut", "GENEZ").Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Join(t *testing.T) {
	s.sql = ""
	expected = " INNER JOIN users ON roles"
	got = s.Join("INNER", "users", "roles").Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Limit(t *testing.T) {
	s.sql = ""
	expected = " LIMIT 1,5"
	got = s.Limit(1, 5).Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Group(t *testing.T) {
	s.sql = ""
	expected = " GROUP BY firstname"
	got = s.GroupBy("firstname").Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Order(t *testing.T) {
	s.sql = ""
	expected = " ORDER BY firstname ASC"
	got = s.OrderBy("firstname", "ASC").Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

func TestSql_Union(t *testing.T) {
	s.sql = ""
	expected = " UNION select * from companies"
	got = s.Union("select * from companies").Sql()
	if !reflect.DeepEqual(expected, got) {
		t.Errorf("expected = %v, got %v", expected, got)
	}
}

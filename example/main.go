package main

import (
	"fmt"

	query "github.com/mstgnz/gobuilder"
)

var (
	gb  query.GoBuilder
	sql string
)

func main() {

	sql = gb.Select("users").Where("id", "=", "1").Sql()
	fmt.Printf("All Columns: \n%s\n", sql)

	sql = gb.Select("users", "firstname", "lastname", "create_date").
		Where("id", "=", "1").
		Sql()
	fmt.Printf("Filter Columns: \n%s\n", sql)

	sql = gb.Select("users").
		Where("id", "=", "1").
		OrWhere("email", "=", "loremipsum@lrmpsm.com").
		Sql()
	fmt.Printf("Where Or Where: \n%s\n", sql)

	sql = gb.Select("users as u", "u.firstname", "u.lastname", "a.address").
		Join("INNER", "address as a", "a.user_id=u.id").
		Where("u.email", "=", "loremipsum@lrmpsm.com").
		Sql()
	fmt.Printf("Join: \n%s\n", sql)

	sql = gb.Select("users").
		Where("id", "=", "1").
		Between("create_date", "2021-01-01", "2021-03-16").
		Sql()
	fmt.Printf("Between: \n%s\n", sql)

	sql = gb.Select("users").
		Where("id", "=", "1").
		Between("create_date", "2021-01-01", "2021-03-16").
		Limit(1, 5).
		Sql()
	fmt.Printf("Limit: \n%s\n", sql)

	sql = gb.Select("users").
		Where("id", "=", "1").
		Between("create_date", "2021-01-01", "2021-03-16").
		GroupBy("lastname").
		Sql()
	fmt.Printf("Group By: \n%s\n", sql)

	sql = gb.Select("users").
		Where("id", "=", "1").
		Between("create_date", "2021-01-01", "2021-03-16").
		GroupBy("lastname").
		OrderBy("id", "DESC").
		Sql()
	fmt.Printf("Order By: \n%s\n", sql)

	sql = gb.Select("users").Where("lastname", "=", "lorem").Sql()
	sql = gb.Select("users").Where("lastname", "=", "ipsum").Union(sql).Sql()
	fmt.Printf("Union: \n%s\n", sql)

}

package main

import (
	"fmt"

	"github.com/mstgnz/gobuilder"
)

var (
	gb     = gobuilder.NewGoBuilder(gobuilder.Postgres)
	query  string
	params []any
)

func main() {

	query, params = gb.Table("users").Select().Where("id", "=", "1").Prepare()
	fmt.Printf("All Columns Prepare: \n%s\t%v\n", query, params)
	query = gb.Table("users").Select().Where("id", "=", "1").Sql()
	fmt.Printf("All Columns Sql: \n%s\n\n", query)

	query, params = gb.Table("users").Select("firstname", "lastname", "created_at").Where("id", "=", 1).Prepare()
	fmt.Printf("Filter Columns Prepare: \n%s\t%v\n", query, params)
	query = gb.Table("users").Select("firstname", "lastname", "created_at").Where("id", "=", 1).Sql()
	fmt.Printf("Filter Columns Sql: \n%s\n\n", query)

	query, params = gb.Table("users").Select().Where("id", "=", "1").OrWhere("email", "=", "loremipsum@lrmpsm.com").Prepare()
	fmt.Printf("Where Or Where Prepare: \n%s\t%v\n\n", query, params)
	query = gb.Table("users").Select().Where("id", "=", "1").OrWhere("email", "=", "loremipsum@lrmpsm.com").Sql()
	fmt.Printf("Where Or Where Sql: \n%s\n\n", query)

	query, params = gb.Table("users as u").Select("u.firstname", "u.lastname", "a.address").
		Join("address as a", "a.user_id", "=", "u.id").
		Where("u.email", "=", "loremipsum@lrmpsm.com").Prepare()
	fmt.Printf("Join Prepare: \n%s\t%v\n", query, params)
	query = gb.Table("users as u").Select("u.firstname", "u.lastname", "a.address").
		Join("address as a", "a.user_id", "=", "u.id").
		Where("u.email", "=", "loremipsum@lrmpsm.com").Sql()
	fmt.Printf("Join Sql: \n%s\n\n", query)

	query, params = gb.Table("users").Select().
		Where("id", "=", "1").
		Between("created_at", "2021-01-01", "2021-03-16").Prepare()
	fmt.Printf("Between Prepare: \n%s\t%v\n", query, params)
	query = gb.Table("users").Select().
		Where("id", "=", "1").
		Between("created_at", "2021-01-01", "2021-03-16").Sql()
	fmt.Printf("Between Sql: \n%s\n\n", query)

	query, params = gb.Table("users").Select().
		Where("id", "=", "1").
		Between("created_at", "2021-01-01", "2021-03-16").
		Limit(1, 5).Prepare()
	fmt.Printf("Limit Prepare: \n%s\t%v\n", query, params)
	query = gb.Table("users").Select().
		Where("id", "=", "1").
		Between("created_at", "2021-01-01", "2021-03-16").
		Limit(1, 5).Sql()
	fmt.Printf("Limit Sql: \n%s\n\n", query)

	query, params = gb.Table("users").Select().
		Where("id", "=", "1").
		Between("created_at", "2021-01-01", "2021-03-16").
		GroupBy("lastname").Prepare()
	fmt.Printf("Group By Prepare: \n%s\t%v\n", query, params)
	query = gb.Table("users").Select().
		Where("id", "=", "1").
		Between("created_at", "2021-01-01", "2021-03-16").
		GroupBy("lastname").Sql()
	fmt.Printf("Group By Sql: \n%s\n\n", query)

	query, params = gb.Table("users").Select().
		Where("id", "=", "1").
		Between("created_at", "2021-01-01", "2021-03-16").
		GroupBy("lastname").
		OrderBy("id").Prepare()
	fmt.Printf("Order By Prepare: \n%s\t%v\n", query, params)
	query = gb.Table("users").Select().
		Where("id", "=", "1").
		Between("created_at", "2021-01-01", "2021-03-16").
		GroupBy("lastname").
		OrderBy("id").Sql()
	fmt.Printf("Order By Sql: \n%s\n\n", query)

	query, params = gb.Table("users").Delete().Where("email", "=", "loremipsum@lrmpsm.com").Prepare()
	fmt.Printf("Delete Prepare: \n%s\t%v\n", query, params)
	query = gb.Table("users").Delete().Where("email", "=", "loremipsum@lrmpsm.com").Sql()
	fmt.Printf("Delete Sql: \n%s\n\n", query)

	query, params = gb.Table("users").Select().Where("lastname", "=", "lorem").Prepare()
	query, params = gb.Table("users").Select().Where("lastname", "=", "ipsum").Union(query).Prepare()
	fmt.Printf("Union: \n%s\n%v\n\n", query, params)

	// example subquery
	mainBuilder := gobuilder.NewGoBuilder(gobuilder.Postgres)
	subBuilder := gobuilder.NewGoBuilder(gobuilder.Postgres)

	query, params = subBuilder.Table("users").Select("id").Where("age", ">", 30).Prepare()
	query, params = mainBuilder.Table("orders").Select("order_id", "user_id").Where("user_id", "IN", query).Prepare()

	fmt.Printf("SubQuery: \n%s\n%v\n\n", query, params)
}

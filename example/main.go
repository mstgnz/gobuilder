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

	query, params = gb.Table("users").Select().Where("id", "=", "1").ToSql()
	fmt.Printf("All Columns: \n%s\n%v\n\n", query, params)

	query, params = gb.Table("users").Select("firstname", "lastname", "created_at").
		Where("id", "=", 1).
		ToSql()
	fmt.Printf("Filter Columns: \n%s\n%v\n\n", query, params)

	query, params = gb.Table("users").Select().
		Where("id", "=", "1").
		OrWhere("email", "=", "loremipsum@lrmpsm.com").
		ToSql()
	fmt.Printf("Where Or Where: \n%s\n%v\n\n", query, params)

	query, params = gb.Table("users as u").Select("u.firstname", "u.lastname", "a.address").
		Join("INNER", "address as a", "a.user_id=u.id").
		Where("u.email", "=", "loremipsum@lrmpsm.com").
		ToSql()
	fmt.Printf("Join: \n%s\n%v\n\n", query, params)

	query, params = gb.Table("users").Select().
		Where("id", "=", "1").
		Between("created_at", "2021-01-01", "2021-03-16").
		ToSql()
	fmt.Printf("Between: \n%s\n%v\n\n", query, params)

	query, params = gb.Table("users").Select().
		Where("id", "=", "1").
		Between("created_at", "2021-01-01", "2021-03-16").
		Limit(1, 5).
		ToSql()
	fmt.Printf("Limit: \n%s\n%v\n\n", query, params)

	query, params = gb.Table("users").Select().
		Where("id", "=", "1").
		Between("created_at", "2021-01-01", "2021-03-16").
		GroupBy("lastname").
		ToSql()
	fmt.Printf("Group By: \n%s\n%v\n\n", query, params)

	query, params = gb.Table("users").Select().
		Where("id", "=", "1").
		Between("created_at", "2021-01-01", "2021-03-16").
		GroupBy("lastname").
		OrderBy("id").
		ToSql()
	fmt.Printf("Order By: \n%s\n%v\n\n", query, params)

	query, params = gb.Table("users").Delete().Where("email", "=", "loremipsum@lrmpsm.com").ToSql()
	fmt.Printf("Delete: \n%s\n%v\n\n", query, params)

	query, params = gb.Table("users").Select().Where("lastname", "=", "lorem").ToSql()
	query, params = gb.Table("users").Select().Where("lastname", "=", "ipsum").Union(query).ToSql()
	fmt.Printf("Union: \n%s\n%v\n\n", query, params)

	// example subquery
	mainBuilder := gobuilder.NewGoBuilder(gobuilder.Postgres)
	subBuilder := gobuilder.NewGoBuilder(gobuilder.Postgres)

	query, params = subBuilder.Table("users").Select("id").Where("age", ">", 30).ToSql()
	query, params = mainBuilder.Table("orders").Select("order_id", "user_id").Where("user_id", "IN", query).ToSql()

	fmt.Printf("SubQuery: \n%s\n%v\n\n", query, params)
}

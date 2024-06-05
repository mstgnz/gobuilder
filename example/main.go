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

	sql = gb.Select("users").Where("id", "=", "1").ToSql()
	fmt.Printf("All Columns: \n%s\n", sql)

	sql = gb.Select("users", "firstname", "lastname", "create_date").
		Where("id", "=", "1").
		ToSql()
	fmt.Printf("Filter Columns: \n%s\n", sql)

	sql = gb.Select("users").
		Where("id", "=", "1").
		OrWhere("email", "=", "loremipsum@lrmpsm.com").
		ToSql()
	fmt.Printf("Where Or Where: \n%s\n", sql)

	sql = gb.Select("users as u", "u.firstname", "u.lastname", "a.address").
		Join("INNER", "address as a", "a.user_id=u.id").
		Where("u.email", "=", "loremipsum@lrmpsm.com").
		ToSql()
	fmt.Printf("Join: \n%s\n", sql)

	sql = gb.Select("users").
		Where("id", "=", "1").
		Between("create_date", "2021-01-01", "2021-03-16").
		ToSql()
	fmt.Printf("Between: \n%s\n", sql)

	sql = gb.Select("users").
		Where("id", "=", "1").
		Between("create_date", "2021-01-01", "2021-03-16").
		Limit(1, 5).
		ToSql()
	fmt.Printf("Limit: \n%s\n", sql)

	sql = gb.Select("users").
		Where("id", "=", "1").
		Between("create_date", "2021-01-01", "2021-03-16").
		GroupBy("lastname").
		ToSql()
	fmt.Printf("Group By: \n%s\n", sql)

	sql = gb.Select("users").
		Where("id", "=", "1").
		Between("create_date", "2021-01-01", "2021-03-16").
		GroupBy("lastname").
		OrderBy("id").
		ToSql()
	fmt.Printf("Order By: \n%s\n", sql)

	sql = gb.Select("users").Where("lastname", "=", "lorem").ToSql()
	sql = gb.Select("users").Where("lastname", "=", "ipsum").Union(sql).ToSql()
	fmt.Printf("Union: \n%s\n", sql)

	// example subquery
	mainBuilder := &query.GoBuilder{}
	subBuilder := &query.GoBuilder{}

	subBuilder.Select("users", "id").Where("age", ">", 30)
	mainBuilder.Select("orders", "order_id", "user_id").Where("user_id", "IN", subBuilder)

	fmt.Println(mainBuilder.ToSql())
	// SELECT order_id, user_id FROM orders WHERE user_id IN (SELECT id FROM users WHERE age > '30')

}

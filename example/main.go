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

	sql = gb.Table("users").Select().Where("id", "=", "1").ToSql()
	fmt.Printf("All Columns: \n%s\n\n", sql)

	sql = gb.Table("users").Select("firstname", "lastname", "created_at").
		Where("id", "=", 1).
		ToSql()
	fmt.Printf("Filter Columns: \n%s\n\n", sql)

	sql = gb.Table("users").Select().
		Where("id", "=", "1").
		OrWhere("email", "=", "loremipsum@lrmpsm.com").
		ToSql()
	fmt.Printf("Where Or Where: \n%s\n\n", sql)

	sql = gb.Table("users as u").Select("u.firstname", "u.lastname", "a.address").
		Join("INNER", "address as a", "a.user_id=u.id").
		Where("u.email", "=", "loremipsum@lrmpsm.com").
		ToSql()
	fmt.Printf("Join: \n%s\n\n", sql)

	sql = gb.Table("users").Select().
		Where("id", "=", "1").
		Between("created_at", "2021-01-01", "2021-03-16").
		ToSql()
	fmt.Printf("Between: \n%s\n\n", sql)

	sql = gb.Table("users").Select().
		Where("id", "=", "1").
		Between("created_at", "2021-01-01", "2021-03-16").
		Limit(1, 5).
		ToSql()
	fmt.Printf("Limit: \n%s\n\n", sql)

	sql = gb.Table("users").Select().
		Where("id", "=", "1").
		Between("created_at", "2021-01-01", "2021-03-16").
		GroupBy("lastname").
		ToSql()
	fmt.Printf("Group By: \n%s\n\n", sql)

	sql = gb.Table("users").Select().
		Where("id", "=", "1").
		Between("created_at", "2021-01-01", "2021-03-16").
		GroupBy("lastname").
		OrderBy("id").
		ToSql()
	fmt.Printf("Order By: \n%s\n\n", sql)

	sql = gb.Table("users").Select().Where("lastname", "=", "lorem").ToSql()
	sql = gb.Table("users").Select().Where("lastname", "=", "ipsum").Union(sql).ToSql()
	fmt.Printf("Union: \n%s\n\n", sql)

	// example subquery
	mainBuilder := &query.GoBuilder{}
	subBuilder := &query.GoBuilder{}

	subBuilder.Table("users").Select("id").Where("age", ">", 30)
	mainBuilder.Table("orders").Select("order_id", "user_id").Where("user_id", "IN", subBuilder)

	fmt.Println(mainBuilder.ToSql())
	// SELECT order_id, user_id FROM orders WHERE user_id IN (SELECT id FROM users WHERE age > '30')

}

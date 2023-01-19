package main

import (
	"fmt"

	generator "github.com/mstgnz/golang-sql-generator"
)

var (
	s   generator.Sql
	sql string
)

func main() {

	sql = s.Select("users", nil).Where("id", "=", "1").Get()
	fmt.Printf("All Columns: \n%s\n", sql)

	sql = s.Select("users", []string{"firstname", "lastname", "create_date"}).
		Where("id", "=", "1").
		Get()
	fmt.Printf("Filter Columns: \n%s\n", sql)

	sql = s.Select("users", nil).
		Where("id", "=", "1").
		OrWhere("email", "=", "loremipsum@lrmpsm.com").
		Get()
	fmt.Printf("Where Or Where: \n%s\n", sql)

	sql = s.Select("users as u", []string{"u.firstname", "u.lastname", "a.address"}).
		Join("INNER", "address as a", "a.user_id=u.id").
		Where("u.email", "=", "loremipsum@lrmpsm.com").
		Get()
	fmt.Printf("Join: \n%s\n", sql)

	sql = s.Select("users", nil).
		Where("id", "=", "1").
		Between("create_date", "2021-01-01", "2021-03-16").
		Get()
	fmt.Printf("Between: \n%s\n", sql)

	sql = s.Select("users", nil).
		Where("id", "=", "1").
		Between("create_date", "2021-01-01", "2021-03-16").
		Limit(1, 5).
		Get()
	fmt.Printf("Limit: \n%s\n", sql)

	sql = s.Select("users", nil).
		Where("id", "=", "1").
		Between("create_date", "2021-01-01", "2021-03-16").
		GroupBy("lastname").
		Get()
	fmt.Printf("Group By: \n%s\n", sql)

	sql = s.Select("users", nil).
		Where("id", "=", "1").
		Between("create_date", "2021-01-01", "2021-03-16").
		GroupBy("lastname").
		OrderBy("id", "DESC").
		Get()
	fmt.Printf("Order By: \n%s\n", sql)

	sql = s.Select("users", nil).Where("lastname", "=", "lorem").Get()
	sql = s.Select("users", nil).Where("lastname", "=", "ipsum").Union(sql).Get()
	fmt.Printf("Union: \n%s\n", sql)

}

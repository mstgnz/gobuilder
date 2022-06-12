package main

import (
	"fmt"
)

var (
	s   Sql
	sql string
)

func main() {

	sql = s.Delete("users").Where("email", "=", "loremipsum@lrmpsm.com").Get()

	fmt.Print(sql)
}

# Golang Query Builder

This package is designed for constructing SQL queries in Go. It does not execute the queries; you need to use a database connection (e.g., MySQL, PostgreSQL) to run them.

## Installation

To install the package, use:

```bash
go get github.com/mstgnz/gobuilder
```

## Usage

Import the package and create a new builder instance:

```go
import "github.com/mstgnz/gobuilder"

var gb = gobuilder.NewGoBuilder(gobuilder.Postgres)
```

## Examples

### Select Queries

#### Select All Columns
```go
gb.Table("users").Select().Where("id", "=", 1).Sql()
```
SQL Output:
```sql
SELECT * FROM users WHERE id = 1
```

#### Select Specific Columns
```go
gb.Table("users").Select("firstname", "lastname", "created_at").Where("id", "=", 1).Sql()
```
SQL Output:
```sql
SELECT firstname, lastname, created_at FROM users WHERE id = 1
```

#### Select with Conditions
```go
gb.Table("users").Select("name", "email").Where("status", "=", "active").Where("age", ">", 18).Sql()
```
SQL Output:
```sql
SELECT name, email FROM users WHERE status = 'active' AND age > 18
```

### Insert Query
```go
args := map[string]any{"firstname": "John", "lastname": "Doe"}
gb.Table("users").Create(args).Sql()
```
SQL Output:
```sql
INSERT INTO users (firstname, lastname) VALUES ('John', 'Doe')
```

### Update Query
```go
args := map[string]any{"firstname": "Jane"}
gb.Table("users").Update(args).Where("id", "=", 1).Sql()
```
SQL Output:
```sql
UPDATE users SET firstname = 'Jane' WHERE id = 1
```

### Delete Query
```go
gb.Table("users").Delete().Where("id", "=", 1).Sql()
```
SQL Output:
```sql
DELETE FROM users WHERE id = 1
```

### Raw SQL
```go
gb.Raw("SELECT * FROM users WHERE id = ?", 1).Sql()
```
SQL Output:
```sql
SELECT * FROM users WHERE id = 1
```

### SQL Injection Prevention
```go
gb.Table("users").Where("username", "=", "admin' OR '1'='1").Prepare()
```
SQL Output:
```sql
SELECT * FROM users WHERE username = $1
Params: ['admin'' OR ''1''=''1']
```

### Join
```go
gb.Table("orders").Join("users", "users.id", "=", "orders.user_id").Select("orders.id", "users.name").Sql()
```
SQL Output:
```sql
SELECT orders.id, users.name FROM orders INNER JOIN users ON users.id = orders.user_id
```

### Aggregate Functions
```go
gb.Table("orders").Select("COUNT(*) as total").Sql()
```
SQL Output:
```sql
SELECT COUNT(*) as total FROM orders
```

### Subquery
```go
subQuery := gb.Table("orders").Select("customer_id").Where("total", ">", 1000)
gb.Table("customers").Select("name").Where("id", "IN", subQuery).Sql()
```
SQL Output:
```sql
SELECT name FROM customers WHERE id IN (SELECT customer_id FROM orders WHERE total > 1000)
```

### Complex Conditions
```go
age := 30
name := "John"
gb.Table("users").Select().WhenThen(age > 0, func(b *GoBuilder) *GoBuilder { return b.Where("age", ">", age) }, nil).WhenThen(name != "", func(b *GoBuilder) *GoBuilder { return b.Where("name", "=", name) }, nil).Sql()
```
SQL Output:
```sql
SELECT * FROM users WHERE age > 30 AND name = 'John'
```

## Benchmark Results

The following benchmark results provide an overview of the performance of various SQL operations using the query builder:

```
goos: darwin
goarch: arm64
cpu: Apple M1
BenchmarkSelect-8                 381279              3140 ns/op
BenchmarkInsert-8                 315511              3654 ns/op
BenchmarkUpdate-8                 300626              3990 ns/op
BenchmarkDelete-8                 375622              2948 ns/op
```

These results were obtained on an Apple M1 CPU and may vary based on hardware and system configuration.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
# Golang Query Builder
This package is only for creating sql text. To run the created sql text, you must create a database connection. (Ex: mysql, postgresql).

# Examples

## Select
```go
import "github.com/mstgnz/gobuilder"

var gb = gobuilder.NewGoBuilder(gobuilder.Postgres)
```

```go
// all columns
gb.Table("users").Select().Where("id","=","1").Sql()
gb.Table("users").Select().Where("id","=","1").Prepare()

// filter columns
gb.Table("users").Select("firstname", "lastname", "created_at").
    Where("id", "=", 1).
    Prepare()
```
```sql
Result: SELECT * FROM users WHERE id = $1
Result: SELECT * FROM users WHERE id = 1
Result: SELECT firstname,lastname,created_at FROM users WHERE id = $1
Params: [1]
```
### where orWhere
```go
gb.Table("users").Select().
    Where("id", "=", "1").
    OrWhere("email", "=", "loremipsum@lrmpsm.com").
    Prepare()
```
```sql
Result: SELECT * FROM users WHERE id=$1 OR email=$2
Params: ["1", "loremipsum@lrmpsm.com"]
```
### join
```go
gb.Table("users as u").Select("u.firstname", "u.lastname", "a.address").
    Join("address as a", "a.user_id","=","u.id").
    Where("u.email", "=", "loremipsum@lrmpsm.com").
    Prepare()
gb.Table("users as u").Select("u.firstname", "u.lastname", "a.address").
    Join("address as a", "a.user_id","=","u.id").
    Where("u.email", "=", "loremipsum@lrmpsm.com").
    Sql()
```
```sql
Result: SELECT u.firstname,u.lastname,a.address FROM users as u INNER JOIN address as a ON a.user_id=u.id WHERE u.email=$1
Params: ["loremipsum@lrmpsm.com"]
Result: SELECT u.firstname,u.lastname,a.address FROM users as u INNER JOIN address as a ON a.user_id=u.id WHERE u.email='loremipsum@lrmpsm.com'
```
### between
```go
gb.Table("users").Select().
	Where("id", "=", "1").
	Between("created_at", "2021-01-01", "2021-03-16").
	Prepare()
gb.Table("users").Select().
	Where("id", "=", "1").
	Between("created_at", "2021-01-01", "2021-03-16").
	Sql()
```
```sql
Result: SELECT * FROM users WHERE id=$1 AND created_at BETWEEN $2 AND $3
Params: [1, "2021-01-01", "2021-03-16"]
Result: SELECT * FROM users WHERE id=1 AND created_at BETWEEN '2021-01-01' AND '2021-03-16'
```
### limit
```go
gb.Table("users").Select().
    Where("id", "=", "1").
    Between("created_at", "2021-01-01", "2021-03-16").
    Limit(1, 5).
    Prepare()
```
```sql
Result: SELECT * FROM users WHERE id=$1 AND created_at BETWEEN $2 AND $3 LIMIT 1,5
Params: [1, "2021-01-01", "2021-03-16"]
```
### group by
```go
gb.Table("users").Select().
	Where("id", "=", "1").
	Between("created_at", "2021-01-01", "2021-03-16").
	GroupBy("lastname").
	Prepare()
```
```sql
Result: SELECT * FROM users WHERE id=$1 AND created_at BETWEEN $2 AND $3 GROUP BY lastname
Params: [1, "2021-01-01", "2021-03-16"]
```
### order by
```go
gb.Table("users").Select().
	Where("id", "=", "1").
	Between("created_at", "2021-01-01", "2021-03-16").
	GroupBy("lastname").
	OrderBy("id", "DESC").
	Prepare()
```
```sql
Result: SELECT * FROM users WHERE id=$1 AND created_at BETWEEN $2 AND $3 GROUP BY lastname ORDER BY id DESC
Params: [1, "2021-01-01", "2021-03-16"]
```
### union
```go
s1, _ := gb.Table("users").Select().Where("lastname", "=", "lorem").Prepare()
s2, _ := gb.Table("users").Select().Where("lastname", "=", "ipsum").Union(s1).Prepare()
```
```sql
Result: SELECT * FROM users WHERE lastname=$1 UNION SELECT * FROM users WHERE lastname=$2
Params: ["lorem", "ipsum"]
```

## Create
```go
args := map[string]any{
"firstname": "Lorem",
"lastname":  "IPSUM",
}
gb.Table("users").Create(args).Prepare()
```
```sql
Result: INSERT INTO users (lastname,firstname) VALUES ($1,$2)
Params: ["Lorem", "IPSUM"]
```

## Update
```go
args := map[string]any{
"firstname": "Lorem",
"lastname":  "IPSUM",
}
gb.Table("users").Update(args).Where("email", "=", "loremipsum@lrmpsm.com").Prepare()
```
```sql
Result: UPDATE users SET firstname=$1, lastname=$2 WHERE email=$3
Params: ["Lorem", "IPSUM", "loremipsum@lrmpsm.com"]
```

## Delete
```go
gb.Table("users").Delete().Where("email", "=", "loremipsum@lrmpsm.com").Prepare()
```
```sql
Result: DELETE FROM users WHERE email=$1
Params: ["loremipsum@lrmpsm.com"]
```
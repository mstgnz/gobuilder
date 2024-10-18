# Golang Query Builder
This package is only for creating sql text. To run the created sql text, you must create a database connection. (Ex: mysql, postgresql).

# Examples

## Select
```go
import "github.com/mstgnz/gobuilder"

var gb GoBuilder
```

```go
// all columns
gb.Table("users").Select().Where("id","=","1").ToSql()

// filter columns
gb.Table("users").Select("firstname", "lastname", "created_at").
    Where("id", "=", 1).
    ToSql()
```
```sql
Result: SELECT * FROM users WHERE id = '1'
Result: SELECT firstname,lastname,created_at FROM users WHERE id = 1
```
### where orWhere
```go
gb.Table("users").Select().
    Where("id", "=", "1").
    OrWhere("email", "=", "loremipsum@lrmpsm.com").
    ToSql()
```
```sql
Result: SELECT * FROM users WHERE id='1' OR email='loremipsum@lrmpsm.com'
```
### join
```go
gb.Table("users as u").Select("u.firstname", "u.lastname", "a.address").
    Join("INNER", "address as a", "a.user_id=u.id").
    Where("u.email", "=", "loremipsum@lrmpsm.com").
    ToSql()
```
```sql
Result: SELECT u.firstname,u.lastname,a.address FROM users as u INNER JOIN address as a ON a.user_id=u.id WHERE u.email='loremipsum@lrmpsm.com'
```
### between
```go
gb.Table("users").Select().
	Where("id", "=", "1").
	Between("created_at", "2021-01-01", "2021-03-16").
	ToSql()
```
```sql
Result: SELECT * FROM users WHERE id='1' AND created_at BETWEEN '2021-01-01' AND '2021-03-16'
```
### limit
```go
gb.Table("users").Select().
    Where("id", "=", "1").
    Between("created_at", "2021-01-01", "2021-03-16").
    Limit(1, 5).
    ToSql()
```
```sql
Result: SELECT * FROM users WHERE id='1' AND created_at BETWEEN '2021-01-01' AND '2021-03-16' LIMIT 1,5
```
### group by
```go
gb.Table("users").Select().
	Where("id", "=", "1").
	Between("created_at", "2021-01-01", "2021-03-16").
	GroupBy("lastname").
	ToSql()
```
```sql
Result: SELECT * FROM users WHERE id='1' AND created_at BETWEEN '2021-01-01' AND '2021-03-16' GROUP BY lastname
```
### order by
```go
gb.Table("users").Select().
	Where("id", "=", "1").
	Between("created_at", "2021-01-01", "2021-03-16").
	GroupBy("lastname").
	OrderBy("id", "DESC").
	ToSql()
```
```sql
Result: SELECT * FROM users WHERE id='1' AND created_at BETWEEN '2021-01-01' AND '2021-03-16' GROUP BY lastname ORDER BY id DESC
```
### union
```go
s1 := gb.Table("users").Select().Where("lastname", "=", "lorem").ToSql()
s2 := gb.Table("users").Select().Where("lastname", "=", "ipsum").Union(s1).ToSql()
```
```sql
Result: SELECT * FROM users WHERE lastname='ipsum' UNION SELECT * FROM users WHERE lastname='lorem'
```

## Insert
```go
args := map[string]string{
"firstname": "Lorem",
"lastname":  "IPSUM",
}
gb.Table("users").Insert(args).ToSql()
```
```sql
Result : INSERT INTO users (lastname,firstname) VALUES ('Lorem','IPSUM')
```

## Update
```go
args := map[string]string{
"firstname": "Lorem",
"lastname":  "IPSUM",
}
gb.Table("users").Update(args).Where("email", "=", "loremipsum@lrmpsm.com").ToSql()
```
```sql
Result: UPDATE users SET firstname='Lorem', lastname='IPSUM' WHERE email='loremipsum@lrmpsm.com'
```

## Delete
```go
gb.Table("users").Delete().Where("email", "=", "loremipsum@lrmpsm.com").ToSql()
```
```sql
Result: DELETE FROM users WHERE email='loremipsum@lrmpsm.com'
```
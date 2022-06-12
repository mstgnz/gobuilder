# Golang Sql Generator
This package is only for creating sql text. To run the created sql text, you must create a database connection. (Ex: mysql, postgresql).

# Examples

## Select
```go
// all columns
s.Select("users", nil).Where("id","=","1").
	Get()
// filter columns
s.Select("users", []string{"firstname", "lastname", "create_date"}).
    Where("id", "=", "1").
    Get()
```
```sql
Result: SELECT * FROM users WHERE id = '1'
Result: SELECT firstname,lastname,create_date FROM users WHERE id = '1'
```
### where orWhere
```go
s.Select("users", nil).
    Where("id", "=", "1").
    OrWhere("email", "=", "loremipsum@lrmpsm.com").
    Get()
```
```sql
Result: SELECT * FROM users WHERE id='1' OR email='loremipsum@lrmpsm.com'
```
### join
```go
s.Select("users as u", []string{"u.firstname", "u.lastname", "a.address"}).
    Join("INNER", "address as a", "a.user_id=u.id").
    Where("u.email", "=", "loremipsum@lrmpsm.com").
    Get()
```
```sql
Result: SELECT u.firstname,u.lastname,a.address FROM users as u INNER JOIN address as a ON a.user_id=u.id WHERE u.email='loremipsum@lrmpsm.com'
```
### between
```go
s.Select("users", nil).
	Where("id", "=", "1").
	Between("create_date", "2021-01-01", "2021-03-16").
	Get()
```
```sql
Result: SELECT * FROM users WHERE id='1' AND create_date BETWEEN '2021-01-01' AND '2021-03-16'
```
### limit
```go
s.Select("users", nil).
    Where("id", "=", "1").
    Between("create_date", "2021-01-01", "2021-03-16").
    Limit(1, 5).
    Get()
```
```sql
Result: SELECT * FROM users WHERE id='1' AND create_date BETWEEN '2021-01-01' AND '2021-03-16' LIMIT 1,5
```
### group by
```go
s.Select("users", nil).
	Where("id", "=", "1").
	Between("create_date", "2021-01-01", "2021-03-16").
	GroupBy("lastname").
	Get()
```
```sql
Result: SELECT * FROM users WHERE id='1' AND create_date BETWEEN '2021-01-01' AND '2021-03-16' GROUP BY lastname
```
### order by
```go
s.Select("users", nil).
	Where("id", "=", "1").
	Between("create_date", "2021-01-01", "2021-03-16").
	GroupBy("lastname").
	OrderBy("id", "DESC").
	Get()
```
```sql
Result: SELECT * FROM users WHERE id='1' AND create_date BETWEEN '2021-01-01' AND '2021-03-16' GROUP BY lastname ORDER BY id DESC
```
### union
```go
s1 := s.Select("users", nil).Where("lastname", "=", "lorem").Get()
s2  := s.Select("users", nil).Where("lastname", "=", "ipsum").Union(s1).Get()
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
sql = s.Insert("users", args).Get()
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
sql = s.Update("users", args).Where("email", "=", "loremipsum@lrmpsm.com").Get()
```
```sql
Result: UPDATE users SET firstname='Lorem', lastname='IPSUM' WHERE email='loremipsum@lrmpsm.com'
```

## Delete
```go
s.Delete("users").Where("email", "=", "loremipsum@lrmpsm.com").Get()
```
```sql
Result: DELETE FROM users WHERE email='loremipsum@lrmpsm.com'
```
package gobuilder

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode"
)

// SQLDialect represents the type of SQL database being used
// This affects how parameters are formatted in the query
type SQLDialect string

const (
	Postgres  SQLDialect = "postgres"  // PostgreSQL database
	MySQL     SQLDialect = "mysql"     // MySQL database
	SQLite    SQLDialect = "sqlite"    // SQLite database
	SQLServer SQLDialect = "sqlserver" // Microsoft SQL Server
	Oracle    SQLDialect = "oracle"    // Oracle database
)

// Default timeout duration for query execution
// This can be modified per query as needed
var Timeout time.Duration = 30

// GoBuilder is the main struct for building SQL queries
// It maintains the state of the query being built including all clauses and parameters
type GoBuilder struct {
	tableClause   string     // The main table name for the query
	selectClause  string     // The SELECT part of the query, including columns
	whereClause   string     // The WHERE conditions of the query
	groupByClause string     // The GROUP BY columns
	havingClause  string     // The HAVING conditions for grouped results
	orderByClause string     // The ORDER BY columns and direction
	limitClause   string     // The LIMIT and OFFSET values
	unionClause   string     // For UNION operations with other queries
	joinClauses   []string   // All JOIN operations (INNER, LEFT, RIGHT)
	paramsClause  []any      // Collection of parameters for prepared statements
	counterClause int        // Counter for parameter placeholders
	sqlDialect    SQLDialect // The SQL dialect being used
	holderCode    string     // The parameter placeholder format (e.g., $1, ?, @p1)
	err           error      // Stores any errors that occur during query building
}

// NewGoBuilder creates and initializes a new instance of GoBuilder
// Parameters:
//   - sqlDialect: The SQL dialect to use for parameter placeholders
//
// Returns:
//   - *GoBuilder: A new query builder instance configured for the specified dialect
func NewGoBuilder(sqlDialect SQLDialect) *GoBuilder {
	gb := &GoBuilder{
		paramsClause:  []any{},
		counterClause: 1,
		sqlDialect:    sqlDialect,
	}
	gb.holderCode = gb.getPlaceholderCode()
	return gb
}

// Table sets the main table for the query with sanitization
func (gb *GoBuilder) Table(table string) *GoBuilder {
	// SQL injection kontrolü
	lowerTable := strings.ToLower(table)
	riskyWords := []string{
		"drop",
		"truncate",
		"alter",
		"grant",
		"revoke",
		"delete",
		"insert",
		"update",
		"--",
		"/*",
		"*/",
		"xp_",
		"exec",
		"execute",
		"sp_",
		"xp_cmdshell",
	}

	for _, word := range riskyWords {
		if strings.Contains(lowerTable, word) {
			table = strings.Split(table, ";")[0] // Sadece ilk kısmı al
			break
		}
	}

	// Alt sorgu kontrolü
	if strings.Contains(table, "(") && strings.Contains(table, ")") {
		// Alt sorgu içindeki COUNT(*) ifadelerini koru
		if strings.Contains(table, "COUNT(*)") {
			gb.tableClause = table
		} else {
			// Alt sorgu içindeki diğer ifadeleri işle
			gb.tableClause = strings.Replace(table, "COUNT", "COUNT(*)", -1)
		}
	} else {
		gb.tableClause = gb.sanitizeIdentifier(table)
	}
	return gb
}

// Select specifies the columns to retrieve in the query with sanitization
func (gb *GoBuilder) Select(columns ...string) *GoBuilder {
	if gb.tableClause == "" {
		gb.err = fmt.Errorf("table name is required")
		return gb
	}

	if len(columns) == 0 {
		columns = append(columns, "*")
	} else {
		// SQL injection kontrolü
		for i, col := range columns {
			lowerCol := strings.ToLower(col)
			if strings.Contains(lowerCol, ";") ||
				strings.Contains(lowerCol, "drop") ||
				strings.Contains(lowerCol, "truncate") {
				columns[i] = strings.Split(col, ";")[0] // Sadece ilk kısmı al
			}
		}

		// Aggregate fonksiyonları ve sütun isimlerini düzgün şekilde işle
		for i, col := range columns {
			// Eğer aggregate fonksiyonu içeriyorsa ve parantez yoksa
			upperCol := strings.ToUpper(col)
			if !strings.Contains(col, "(") {
				switch {
				case strings.Contains(upperCol, "COUNT"):
					if strings.Contains(upperCol, "AS") {
						parts := strings.Split(col, " AS ")
						col = fmt.Sprintf("COUNT(*) AS %s", parts[1])
					} else {
						col = "COUNT(*)"
					}
				case strings.Contains(upperCol, "SUM"):
					if strings.Contains(upperCol, "AS") {
						parts := strings.Split(col, " AS ")
						colParts := strings.Split(parts[0], " ")
						col = fmt.Sprintf("SUM(%s) AS %s", colParts[1], parts[1])
					}
				case strings.Contains(upperCol, "AVG"):
					if strings.Contains(upperCol, "AS") {
						parts := strings.Split(col, " AS ")
						colParts := strings.Split(parts[0], " ")
						col = fmt.Sprintf("AVG(%s) AS %s", colParts[1], parts[1])
					}
				case strings.Contains(upperCol, "MIN"):
					if strings.Contains(upperCol, "AS") {
						parts := strings.Split(col, " AS ")
						colParts := strings.Split(parts[0], " ")
						col = fmt.Sprintf("MIN(%s) AS %s", colParts[1], parts[1])
					}
				case strings.Contains(upperCol, "MAX"):
					if strings.Contains(upperCol, "AS") {
						parts := strings.Split(col, " AS ")
						colParts := strings.Split(parts[0], " ")
						col = fmt.Sprintf("MAX(%s) AS %s", colParts[1], parts[1])
					}
				}
			}

			// CASE ifadeleri ve Window fonksiyonları için alias kontrolü
			if strings.Contains(upperCol, "CASE") || strings.Contains(upperCol, "OVER") {
				if !strings.Contains(upperCol, " AS ") && strings.Contains(col, " as ") {
					parts := strings.Split(col, " as ")
					col = fmt.Sprintf("%s AS %s", parts[0], parts[1])
				}
			}

			columns[i] = col
		}
	}

	// Preserve WITH clause if it exists (for CTEs)
	withClause := ""
	if strings.HasPrefix(gb.selectClause, "WITH ") {
		withClause = gb.selectClause + " "
	}

	gb.selectClause = fmt.Sprintf("%sSELECT %s FROM %s", withClause, strings.Join(columns, ", "), gb.tableClause)
	return gb
}

// SelectDistinct creates a SELECT DISTINCT query
// Parameters:
//   - columns: Variable number of column names to select distinctly
//
// Returns:
//   - *GoBuilder: The builder instance for method chaining
//
// Example:
//
//	builder.SelectDistinct("country", "city")
//	// Generates: SELECT DISTINCT country, city FROM ...
func (gb *GoBuilder) SelectDistinct(columns ...string) *GoBuilder {
	if len(columns) == 0 {
		columns = append(columns, "*")
	}
	gb.selectClause = fmt.Sprintf("SELECT DISTINCT %s FROM %s", strings.Join(columns, ", "), gb.tableClause)
	return gb
}

// Create builds an INSERT statement with the provided data
// Parameters:
//   - args: Map of column names to values
//   - returning: Optional columns to return after insert (RETURNING clause)
//
// Returns:
//   - *GoBuilder: The builder instance for method chaining
//
// Example:
//
//	builder.Create(map[string]any{
//	    "name": "John",
//	    "age": 30,
//	})
//	// Generates: INSERT INTO table (name, age) VALUES ($1, $2)
func (gb *GoBuilder) Create(args map[string]any, returning ...string) *GoBuilder {
	if len(args) != 0 {
		keys := make([]string, 0, len(args))
		for key := range args {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		columns := make([]string, 0, len(keys))
		values := make([]string, 0, len(keys))
		for _, key := range keys {
			columns = append(columns, key)
			values = append(values, gb.addParam(args[key]))
		}

		gb.selectClause = fmt.Sprintf(
			"INSERT INTO %s (%s) VALUES (%s)",
			gb.tableClause,
			strings.Join(columns, ", "),
			strings.Join(values, ", "),
		)
		if len(returning) > 0 {
			gb.selectClause += fmt.Sprintf(" RETURNING %s", strings.Join(returning, ", "))
		}
	}
	return gb
}

// Update builds an UPDATE statement with the provided data
// Parameters:
//   - args: Map of column names to new values
//
// Returns:
//   - *GoBuilder: The builder instance for method chaining
//
// Example:
//
//	builder.Update(map[string]any{
//	    "status": "active",
//	    "updated_at": time.Now(),
//	})
//	// Generates: UPDATE table SET status = $1, updated_at = $2
func (gb *GoBuilder) Update(args map[string]any) *GoBuilder {
	if len(args) != 0 {
		keys := make([]string, 0, len(args))
		for key := range args {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		setClauses := make([]string, 0, len(keys))
		for _, key := range keys {
			setClauses = append(setClauses, fmt.Sprintf("%s = %s", key, gb.addParam(args[key])))
		}

		gb.selectClause = fmt.Sprintf(
			"UPDATE %s SET %s",
			gb.tableClause,
			strings.Join(setClauses, ", "),
		)
	}
	return gb
}

// Delete builds a DELETE statement for the current table
// Returns:
//   - *GoBuilder: The builder instance for method chaining
//
// Example:
//
//	builder.Delete().Where("status", "=", "inactive")
//	// Generates: DELETE FROM table WHERE status = $1
func (gb *GoBuilder) Delete() *GoBuilder {
	gb.selectClause = fmt.Sprintf("DELETE FROM %s", gb.tableClause)
	return gb
}

// Where adds a WHERE condition to the query
func (gb *GoBuilder) Where(key, opt string, val any) *GoBuilder {
	key = gb.sanitizeIdentifier(key)
	var clause string
	switch v := val.(type) {
	case *GoBuilder:
		subQuery, subParams := v.Prepare()
		gb.paramsClause = append(gb.paramsClause, subParams...)
		clause = fmt.Sprintf("%s %s (%s)", key, opt, subQuery)
	default:
		if str, ok := val.(string); ok && strings.Contains(str, ".") {
			clause = fmt.Sprintf("%s %s %s", key, opt, gb.sanitizeIdentifier(str))
		} else {
			clause = fmt.Sprintf("%s %s %s", key, opt, gb.addParam(val))
		}
	}

	if gb.whereClause == "" {
		gb.whereClause = fmt.Sprintf("WHERE %s", clause)
	} else {
		gb.whereClause = fmt.Sprintf("%s AND %s", gb.whereClause, clause)
	}

	// Eğer SELECT ifadesi yoksa ve tablo adı varsa, varsayılan SELECT ifadesini ekle
	if gb.selectClause == "" && gb.tableClause != "" {
		gb.selectClause = fmt.Sprintf("SELECT * FROM %s", gb.tableClause)
	}

	return gb
}

// OrWhere adds an OR WHERE clause with bind parameters
func (gb *GoBuilder) OrWhere(key, opt string, val any) *GoBuilder {
	clause := fmt.Sprintf("%s %s %s", key, opt, gb.addParam(val))
	gb.addClause("OR", clause)
	return gb
}

// In adds an IN clause with bind parameters
func (gb *GoBuilder) In(column string, args ...any) *GoBuilder {
	return gb.addInClause("AND", column, args...)
}

// OrIn adds an OR IN clause with bind parameters
func (gb *GoBuilder) OrIn(column string, args ...any) *GoBuilder {
	return gb.addInClause("OR", column, args...)
}

// Between adds a BETWEEN clause with bind parameters
func (gb *GoBuilder) Between(column string, args ...any) *GoBuilder {
	return gb.between("AND", column, args...)
}

// OrBetween adds an OR BETWEEN clause with bind parameters
func (gb *GoBuilder) OrBetween(column string, args ...any) *GoBuilder {
	return gb.between("OR", column, args...)
}

// IsNull adds an IS NULL clause
func (gb *GoBuilder) IsNull(column string) *GoBuilder {
	clause := fmt.Sprintf("%s IS NULL", column)
	gb.addClause("AND", clause)
	return gb
}

// OrIsNull adds an OR IS NULL clause
func (gb *GoBuilder) OrIsNull(column string) *GoBuilder {
	clause := fmt.Sprintf("%s IS NULL", column)
	gb.addClause("OR", clause)
	return gb
}

// IsNotNull adds an IS NOT NULL clause
func (gb *GoBuilder) IsNotNull(column string) *GoBuilder {
	clause := fmt.Sprintf("%s IS NOT NULL", column)
	gb.addClause("AND", clause)
	return gb
}

// OrIsNotNull adds an OR IS NOT NULL clause
func (gb *GoBuilder) OrIsNotNull(column string) *GoBuilder {
	clause := fmt.Sprintf("%s IS NOT NULL", column)
	gb.addClause("OR", clause)
	return gb
}

// Having adds a HAVING clause
func (gb *GoBuilder) Having(condition string, args ...any) *GoBuilder {
	// Parametreleri ekle
	for _, arg := range args {
		condition = strings.Replace(condition, "?", gb.addParam(arg), 1)
	}

	if gb.havingClause != "" {
		gb.havingClause = fmt.Sprintf("%s OR %s", gb.havingClause, condition)
	} else {
		gb.havingClause = fmt.Sprintf("HAVING %s", condition)
	}
	return gb
}

// Join adds a JOIN clause
func (gb *GoBuilder) Join(table, first, operator, last string) *GoBuilder {
	join := fmt.Sprintf("INNER JOIN %s ON %s %s %s", table, first, operator, last)
	gb.joinClauses = append(gb.joinClauses, join)
	return gb
}

// LeftJoin adds a LEFT JOIN clause
func (gb *GoBuilder) LeftJoin(table, first, operator, last string) *GoBuilder {
	join := fmt.Sprintf("LEFT JOIN %s ON %s %s %s", table, first, operator, last)
	gb.joinClauses = append(gb.joinClauses, join)
	return gb
}

// RightJoin adds a RIGHT JOIN clause
func (gb *GoBuilder) RightJoin(table, first, operator, last string) *GoBuilder {
	join := fmt.Sprintf("RIGHT JOIN %s ON %s %s %s", table, first, operator, last)
	gb.joinClauses = append(gb.joinClauses, join)
	return gb
}

// Limit adds a LIMIT clause
func (gb *GoBuilder) Limit(offset, limit int) *GoBuilder {
	gb.limitClause = fmt.Sprintf("OFFSET %d LIMIT %d", offset, limit)
	return gb
}

// GroupBy adds a GROUP BY clause
func (gb *GoBuilder) GroupBy(columns ...string) *GoBuilder {
	gb.groupByClause = fmt.Sprintf("GROUP BY %v", strings.Join(columns, ", "))
	return gb
}

// OrderBy adds an ORDER BY ASC clause
func (gb *GoBuilder) OrderBy(columns ...string) *GoBuilder {
	gb.orderByClause = fmt.Sprintf("ORDER BY %v ASC", strings.Join(columns, ", "))
	return gb
}

// OrderByDesc adds an ORDER BY DESC clause
func (gb *GoBuilder) OrderByDesc(columns ...string) *GoBuilder {
	gb.orderByClause = fmt.Sprintf("ORDER BY %v DESC", strings.Join(columns, ", "))
	return gb
}

// Union adds a UNION clause
// Parameters:
//   - builder: The query builder to union with the current query
//
// Returns:
//   - *GoBuilder: The builder instance for method chaining
//
// Example:
//
//	query1 := builder.Table("users").Select("name", "email")
//	query2 := builder.Table("employees").Select("name", "work_email")
//	result := query1.Union(query2)
//	// Generates: SELECT name, email FROM users UNION SELECT name, work_email FROM employees
func (gb *GoBuilder) Union(builder *GoBuilder) *GoBuilder {
	paramOffset := gb.counterClause - 1
	builder.counterClause = gb.counterClause

	unionQuery, unionParams := builder.Prepare()

	// Update the parameter placeholders in the union query
	for i := 1; i <= len(unionParams); i++ {
		oldPlaceholder := fmt.Sprintf("%s%d", gb.holderCode, i)
		newPlaceholder := fmt.Sprintf("%s%d", gb.holderCode, i+paramOffset)
		unionQuery = strings.Replace(unionQuery, oldPlaceholder, newPlaceholder, 1)
	}

	// Add the union parameters to the main query
	gb.paramsClause = append(gb.paramsClause, unionParams...)

	// Update the counter for future parameters
	gb.counterClause += len(unionParams)

	// Add the UNION clause
	if gb.unionClause == "" {
		gb.unionClause = fmt.Sprintf("UNION %s", unionQuery)
	} else {
		gb.unionClause = fmt.Sprintf("%s UNION %s", gb.unionClause, unionQuery)
	}

	return gb
}

// UnionAll adds a UNION ALL clause
// Parameters:
//   - builder: The query builder to union all with the current query
//
// Returns:
//   - *GoBuilder: The builder instance for method chaining
//
// Example:
//
//	query1 := builder.Table("users").Select("name", "email")
//	query2 := builder.Table("employees").Select("name", "work_email")
//	result := query1.UnionAll(query2)
//	// Generates: SELECT name, email FROM users UNION ALL SELECT name, work_email FROM employees
func (gb *GoBuilder) UnionAll(builder *GoBuilder) *GoBuilder {
	// Adjust the parameter counter for the union query
	paramOffset := gb.counterClause - 1
	builder.counterClause = gb.counterClause

	// Get the query and parameters from the union builder
	unionQuery, unionParams := builder.Prepare()

	// Update the parameter placeholders in the union query
	for i := 1; i <= len(unionParams); i++ {
		oldPlaceholder := fmt.Sprintf("%s%d", gb.holderCode, i)
		newPlaceholder := fmt.Sprintf("%s%d", gb.holderCode, i+paramOffset)
		unionQuery = strings.Replace(unionQuery, oldPlaceholder, newPlaceholder, 1)
	}

	// Add the union parameters to the main query
	gb.paramsClause = append(gb.paramsClause, unionParams...)

	// Update the counter for future parameters
	gb.counterClause += len(unionParams)

	// Add the UNION ALL clause
	if gb.unionClause == "" {
		gb.unionClause = fmt.Sprintf("UNION ALL %s", unionQuery)
	} else {
		gb.unionClause = fmt.Sprintf("%s UNION ALL %s", gb.unionClause, unionQuery)
	}

	return gb
}

// Sql returns the final SQL query
func (gb *GoBuilder) Sql() string {
	clauses := []string{
		gb.selectClause,
		strings.Join(gb.joinClauses, " "),
		gb.whereClause,
		gb.groupByClause,
		gb.havingClause,
		gb.orderByClause,
		gb.limitClause,
		gb.unionClause,
	}
	query := strings.Join(clauses, " ")
	re := regexp.MustCompile(`\s+`)
	query = strings.TrimSpace(re.ReplaceAllString(query, " "))

	params := gb.paramsClause
	for i, param := range params {
		placeholder := fmt.Sprintf("%s%d", gb.holderCode, i+1)
		query = strings.Replace(query, placeholder, gb.cleanValue(param), 1)
	}
	gb.reset()
	return query
}

// Prepare returns the final SQL query and the associated bind parameters
func (gb *GoBuilder) Prepare() (string, []any) {
	clauses := make([]string, 0)

	// Add the main SELECT/UPDATE/DELETE clause
	if gb.selectClause != "" {
		clauses = append(clauses, gb.selectClause)
	}

	// Add JOIN clauses
	if len(gb.joinClauses) > 0 {
		clauses = append(clauses, strings.Join(gb.joinClauses, " "))
	}

	// Add WHERE clause
	if gb.whereClause != "" {
		clauses = append(clauses, gb.whereClause)
	}

	// Add GROUP BY clause
	if gb.groupByClause != "" {
		clauses = append(clauses, gb.groupByClause)
	}

	// Add HAVING clause
	if gb.havingClause != "" {
		clauses = append(clauses, gb.havingClause)
	}

	// Add UNION clauses before ORDER BY and LIMIT
	if gb.unionClause != "" {
		clauses = append(clauses, gb.unionClause)
	}

	// Add ORDER BY clause
	if gb.orderByClause != "" {
		clauses = append(clauses, gb.orderByClause)
	}

	// Add LIMIT clause
	if gb.limitClause != "" {
		clauses = append(clauses, gb.limitClause)
	}

	// Join all clauses with spaces and clean up extra whitespace
	query := strings.Join(clauses, " ")
	re := regexp.MustCompile(`\s+`)
	query = strings.TrimSpace(re.ReplaceAllString(query, " "))

	// Get the parameters and reset the builder
	params := gb.paramsClause
	gb.reset()

	return query, params
}

// Private method to RESET builder
func (gb *GoBuilder) reset() {
	*gb = *NewGoBuilder(Postgres)
}

// Private method to add parameters
func (gb *GoBuilder) addParam(value any) string {
	gb.paramsClause = append(gb.paramsClause, value)

	// MySQL ve SQLite için özel durum
	if gb.sqlDialect == MySQL || gb.sqlDialect == SQLite {
		return "?"
	}

	placeholder := fmt.Sprintf("%s%d", gb.holderCode, gb.counterClause)
	gb.counterClause++
	return placeholder
}

// Private method to add clauses with logical operators
func (gb *GoBuilder) addClause(OP, clause string) {
	if gb.whereClause != "" {
		gb.whereClause = fmt.Sprintf("%s %s %s", gb.whereClause, OP, clause)
	} else {
		gb.whereClause = fmt.Sprintf("WHERE %s", clause)
	}
}

// Private method to add IN clauses with values directly
func (gb *GoBuilder) addInClause(OP, column string, args ...any) *GoBuilder {
	if len(args) > 0 {
		values := make([]string, len(args))
		for i, arg := range args {
			values[i] = gb.addParam(arg)
		}
		clause := fmt.Sprintf("%s IN (%s)", column, strings.Join(values, ", "))
		gb.addClause(OP, clause)
	}
	return gb
}

// Private method to add BETWEEN clauses with values directly
func (gb *GoBuilder) between(OP, column string, args ...any) *GoBuilder {
	if len(args) == 2 {
		clause := fmt.Sprintf("%s BETWEEN %s AND %s", column, gb.addParam(args[0]), gb.addParam(args[1]))
		gb.addClause(OP, clause)
	}
	return gb
}

// cleanValue trims and escapes potentially harmful characters from the value
func (gb *GoBuilder) cleanValue(value any) string {
	switch v := value.(type) {
	case string:
		// Escape SQL injection characters
		cleaned := strings.ReplaceAll(v, "'", "''")
		cleaned = strings.ReplaceAll(cleaned, "\\", "\\\\")
		cleaned = strings.ReplaceAll(cleaned, "\x00", "") // Null byte
		cleaned = strings.ReplaceAll(cleaned, "\n", "\\n")
		cleaned = strings.ReplaceAll(cleaned, "\r", "\\r")
		cleaned = strings.ReplaceAll(cleaned, "\x1a", "\\Z") // Ctrl+Z

		// Check risky words for SQL injection
		lowerCleaned := strings.ToLower(cleaned)
		riskyWords := []string{"--", ";", "/*", "*/", "xp_", "select", "update", "delete", "drop", "truncate", "alter", "grant", "revoke"}
		for _, word := range riskyWords {
			if strings.Contains(lowerCleaned, word) {
				cleaned = strings.ReplaceAll(cleaned, word, "")
			}
		}

		return "'" + cleaned + "'"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%g", v)
	case bool:
		return fmt.Sprintf("%t", v)
	case nil:
		return "NULL"
	case time.Time:
		return fmt.Sprintf("'%s'", v.Format("2006-01-02 15:04:05"))
	default:
		return fmt.Sprintf("%v", v)
	}
}

// sanitizeIdentifier sanitizes table and column names
func (gb *GoBuilder) sanitizeIdentifier(identifier string) string {
	// Clean dangerous characters for SQL injection
	parts := strings.Split(identifier, " as ")
	if len(parts) > 2 {
		return "invalid_identifier"
	}

	mainPart := parts[0]

	// Special handling for CASE expressions
	if strings.Contains(strings.ToUpper(mainPart), "CASE") {
		return mainPart
	}

	// Special handling for Window functions
	if strings.Contains(strings.ToUpper(mainPart), "OVER") ||
		strings.Contains(strings.ToUpper(mainPart), "PARTITION BY") ||
		strings.Contains(strings.ToUpper(mainPart), "ROW_NUMBER") ||
		strings.Contains(strings.ToUpper(mainPart), "RANK") ||
		strings.Contains(strings.ToUpper(mainPart), "DENSE_RANK") {
		return mainPart
	}

	// Special handling for JSON operators
	if strings.Contains(mainPart, "->") {
		return mainPart
	}

	// String literal control
	if strings.HasPrefix(mainPart, "'") && strings.HasSuffix(mainPart, "'") {
		return mainPart
	}

	// Check and remove risky words
	lowerMainPart := strings.ToLower(mainPart)
	riskyWords := []string{
		"drop",
		"truncate",
		"alter",
		"grant",
		"revoke",
		"delete",
		"insert",
		"update",
		"--",
		"/*",
		"*/",
		"xp_",
		"exec",
		"execute",
		"sp_",
		"xp_cmdshell",
	}

	for _, word := range riskyWords {
		if strings.Contains(lowerMainPart, word) {
			return "invalid_identifier"
		}
	}

	// Special handling of expressions containing periods (ex: table.column)
	if strings.Contains(mainPart, ".") {
		parts := strings.Split(mainPart, ".")
		if len(parts) != 2 {
			return "invalid_identifier"
		}
		cleanedParts := make([]string, 2)
		for i, part := range parts {
			cleanedParts[i] = strings.Map(func(r rune) rune {
				if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
					return r
				}
				return -1
			}, part)
		}
		mainPart = strings.Join(cleanedParts, ".")
	} else {
		mainPart = strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				return r
			}
			return -1
		}, mainPart)
	}

	// Empty string check
	if mainPart == "" || mainPart == "." {
		return "invalid_identifier"
	}

	// Add if you have Alias
	if len(parts) == 2 {
		alias := strings.Map(func(r rune) rune {
			if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
				return r
			}
			return -1
		}, parts[1])
		if alias != "" {
			mainPart = mainPart + " as " + alias
		}
	}

	return mainPart
}

func (gb *GoBuilder) getPlaceholderCode() string {
	switch gb.sqlDialect {
	case Postgres:
		return "$"
	case SQLServer:
		return "@"
	case Oracle:
		return ":"
	case MySQL:
		return "?"
	case SQLite:
		return "?"
	default:
		return "?"
	}
}

func (gb *GoBuilder) Error() error {
	return gb.err
}

// OnDuplicateKeyUpdate adds ON DUPLICATE KEY UPDATE clause (MySQL specific)
func (gb *GoBuilder) OnDuplicateKeyUpdate(args map[string]any) *GoBuilder {
	if gb.sqlDialect != MySQL {
		gb.err = fmt.Errorf("ON DUPLICATE KEY UPDATE is only supported in MySQL")
		return gb
	}

	if len(args) > 0 {
		keys := make([]string, 0, len(args))
		for key := range args {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		setClauses := make([]string, 0, len(keys))
		for _, key := range keys {
			// MySQL için placeholder'ı doğrudan ? olarak kullan
			setClauses = append(setClauses, fmt.Sprintf("%s = ?", key))
			gb.paramsClause = append(gb.paramsClause, args[key])
		}

		gb.selectClause += fmt.Sprintf(" ON DUPLICATE KEY UPDATE %s", strings.Join(setClauses, ", "))
	}
	return gb
}

// Top adds TOP clause (SQL Server specific)
func (gb *GoBuilder) Top(n int) *GoBuilder {
	if gb.sqlDialect != SQLServer {
		gb.err = fmt.Errorf("TOP clause is only supported in SQL Server")
		return gb
	}

	if gb.selectClause != "" {
		// TOP kelimesini SELECT ve sütun isimleri arasına ekle
		gb.selectClause = strings.Replace(gb.selectClause, "SELECT", fmt.Sprintf("SELECT TOP %d", n), 1)
	}
	return gb
}

// Pragma adds PRAGMA statement (SQLite specific)
func (gb *GoBuilder) Pragma(key string, value string) *GoBuilder {
	if gb.sqlDialect != SQLite {
		gb.err = fmt.Errorf("PRAGMA is only supported in SQLite")
		return gb
	}

	gb.selectClause = fmt.Sprintf("PRAGMA %s = %s", key, value)
	return gb
}

// With adds WITH clause (CTE - Common Table Expression)
func (gb *GoBuilder) With(name string, subQuery *GoBuilder) *GoBuilder {
	// Alt sorguyu hazırla
	subQueryStr, subParams := subQuery.Prepare()

	// Alt sorgu parametrelerini ana sorguya ekle
	gb.paramsClause = append(gb.paramsClause, subParams...)

	// WITH clause'u oluştur
	gb.selectClause = fmt.Sprintf("WITH %s AS (%s)", name, subQueryStr)

	return gb
}

// Lock adds FOR UPDATE/SHARE clause
func (gb *GoBuilder) Lock(lockType string) *GoBuilder {
	gb.selectClause = fmt.Sprintf("%s %s", gb.selectClause, lockType)
	return gb
}

// WhenThen adds conditional clauses
func (gb *GoBuilder) WhenThen(condition bool, trueCase, falseCase func(*GoBuilder) *GoBuilder) *GoBuilder {
	if condition {
		if trueCase != nil {
			return trueCase(gb)
		}
	} else {
		if falseCase != nil {
			return falseCase(gb)
		}
	}
	return gb
}

// CreateBatch adds an INSERT INTO statement for multiple records
func (gb *GoBuilder) CreateBatch(records []map[string]any) *GoBuilder {
	if len(records) == 0 {
		return gb
	}

	// İlk kayıttan sütun isimlerini al
	keys := make([]string, 0)
	for key := range records[0] {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Sütun isimleri
	columns := strings.Join(keys, ", ")

	// Değerler için placeholder'ları oluştur
	var valueStrings []string
	for _, record := range records {
		valuePlaceholders := make([]string, len(keys))
		for i, key := range keys {
			if value, ok := record[key]; ok {
				valuePlaceholders[i] = gb.addParam(value)
			}
		}
		valueStrings = append(valueStrings, fmt.Sprintf("(%s)", strings.Join(valuePlaceholders, ", ")))
	}

	gb.selectClause = fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES %s",
		gb.tableClause,
		columns,
		strings.Join(valueStrings, ", "),
	)

	return gb
}

// Raw adds a raw SQL clause to the query with basic sanitization
func (gb *GoBuilder) Raw(sql string, args ...any) *GoBuilder {
	// SQL injection için temel kontroller
	lowerSQL := strings.ToLower(strings.TrimSpace(sql))

	// Noktalı virgül kontrolü
	if strings.Contains(lowerSQL, ";") {
		gb.err = fmt.Errorf("semicolon is not allowed in raw SQL")
		return gb
	}

	// Riskli komutları kontrol et
	riskyCommands := []string{
		"drop",
		"truncate",
		"alter",
		"grant",
		"revoke",
		"delete",
		"insert",
		"update",
		"--",
		"/*",
		"*/",
		"xp_",
		"exec",
		"execute",
		"sp_",
		"xp_cmdshell",
	}

	// Kelime sınırlarını kontrol et
	for _, cmd := range riskyCommands {
		pattern := fmt.Sprintf(`(?i)(^|\s)%s(\s|$|\()`, cmd)
		if matched, _ := regexp.MatchString(pattern, lowerSQL); matched {
			gb.err = fmt.Errorf("potentially harmful SQL command detected: %s", cmd)
			return gb
		}
	}

	// Çoklu sorgu kontrolü
	if strings.Count(lowerSQL, ";") > 0 {
		gb.err = fmt.Errorf("multiple SQL statements are not allowed")
		return gb
	}

	// Zararlı komut kontrolü
	if strings.Contains(lowerSQL, "drop") || strings.Contains(lowerSQL, "truncate") {
		gb.err = fmt.Errorf("harmful SQL command detected")
		return gb
	}

	// Parametreleri ekle
	for _, arg := range args {
		sql = strings.Replace(sql, "?", gb.addParam(arg), 1)
	}

	if strings.HasPrefix(lowerSQL, "select") {
		if gb.tableClause != "" {
			gb.selectClause = fmt.Sprintf("SELECT * FROM %s %s", gb.tableClause, sql)
		} else {
			gb.selectClause = sql
		}
	} else {
		if gb.whereClause != "" {
			gb.whereClause = fmt.Sprintf("%s %s", gb.whereClause, sql)
		} else {
			gb.whereClause = fmt.Sprintf("WHERE %s", sql)
		}
	}

	return gb
}

// Increment adds an increment operation to the query
func (gb *GoBuilder) Increment(column string, amount int) *GoBuilder {
	gb.selectClause = fmt.Sprintf("UPDATE %s SET %s = %s + %d", gb.tableClause, column, column, amount)
	return gb
}

// Decrement adds a decrement operation to the query
func (gb *GoBuilder) Decrement(column string, amount int) *GoBuilder {
	gb.selectClause = fmt.Sprintf("UPDATE %s SET %s = %s - %d", gb.tableClause, column, column, amount)
	return gb
}

// WhereExists adds a WHERE EXISTS clause
func (gb *GoBuilder) WhereExists(subQuery *GoBuilder) *GoBuilder {
	subQueryStr, subParams := subQuery.Prepare()
	// Alt sorgu parametrelerini ana sorguya ekle
	gb.paramsClause = append(gb.paramsClause, subParams...)
	clause := fmt.Sprintf("EXISTS (%s)", subQueryStr)
	gb.addClause("AND", clause)
	return gb
}

// WhereNotExists adds a WHERE NOT EXISTS clause
func (gb *GoBuilder) WhereNotExists(subQuery *GoBuilder) *GoBuilder {
	subQueryStr, subParams := subQuery.Prepare()
	// Alt sorgu parametrelerini ana sorguya ekle
	gb.paramsClause = append(gb.paramsClause, subParams...)
	clause := fmt.Sprintf("NOT EXISTS (%s)", subQueryStr)
	gb.addClause("AND", clause)
	return gb
}

// WhereJsonContains adds a WHERE JSON_CONTAINS clause
func (gb *GoBuilder) WhereJsonContains(column string, value any) *GoBuilder {
	switch gb.sqlDialect {
	case Postgres:
		gb.addClause("AND", fmt.Sprintf("%s @> %s", column, gb.addParam(value)))
	case MySQL:
		gb.addClause("AND", fmt.Sprintf("JSON_CONTAINS(%s, %s)", column, gb.addParam(value)))
	default:
		gb.err = fmt.Errorf("JSON operations are only supported in PostgreSQL and MySQL")
	}
	return gb
}

// CrossJoin adds a CROSS JOIN clause
func (gb *GoBuilder) CrossJoin(table string) *GoBuilder {
	join := fmt.Sprintf("CROSS JOIN %s", table)
	gb.joinClauses = append(gb.joinClauses, join)
	return gb
}

// FullOuterJoin adds a FULL OUTER JOIN clause
func (gb *GoBuilder) FullOuterJoin(table, first, operator, last string) *GoBuilder {
	join := fmt.Sprintf("FULL OUTER JOIN %s ON %s %s %s", table, first, operator, last)
	gb.joinClauses = append(gb.joinClauses, join)
	return gb
}

// WhereColumn adds a WHERE column comparison
func (gb *GoBuilder) WhereColumn(column1, operator, column2 string) *GoBuilder {
	gb.addClause("AND", fmt.Sprintf("%s %s %s", column1, operator, column2))
	return gb
}

// WhereDate adds a WHERE date comparison
func (gb *GoBuilder) WhereDate(column, operator string, value time.Time) *GoBuilder {
	switch gb.sqlDialect {
	case Postgres:
		gb.addClause("AND", fmt.Sprintf("DATE(%s) %s %s", column, operator, gb.addParam(value.Format("2006-01-02"))))
	case MySQL:
		gb.addClause("AND", fmt.Sprintf("DATE(%s) %s %s", column, operator, gb.addParam(value.Format("2006-01-02"))))
	default:
		gb.addClause("AND", fmt.Sprintf("DATE(%s) %s %s", column, operator, gb.addParam(value.Format("2006-01-02"))))
	}
	return gb
}

// WhereYear adds a WHERE year comparison
func (gb *GoBuilder) WhereYear(column, operator string, year int) *GoBuilder {
	switch gb.sqlDialect {
	case Postgres:
		gb.addClause("AND", fmt.Sprintf("EXTRACT(YEAR FROM %s) %s %s", column, operator, gb.addParam(year)))
	case MySQL:
		gb.addClause("AND", fmt.Sprintf("YEAR(%s) %s %s", column, operator, gb.addParam(year)))
	default:
		gb.addClause("AND", fmt.Sprintf("YEAR(%s) %s %s", column, operator, gb.addParam(year)))
	}
	return gb
}

// WhereMonth adds a WHERE month comparison
func (gb *GoBuilder) WhereMonth(column, operator string, month int) *GoBuilder {
	switch gb.sqlDialect {
	case Postgres:
		gb.addClause("AND", fmt.Sprintf("EXTRACT(MONTH FROM %s) %s %s", column, operator, gb.addParam(month)))
	case MySQL:
		gb.addClause("AND", fmt.Sprintf("MONTH(%s) %s %s", column, operator, gb.addParam(month)))
	default:
		gb.addClause("AND", fmt.Sprintf("MONTH(%s) %s %s", column, operator, gb.addParam(month)))
	}
	return gb
}

// Chunk processes results in chunks
func (gb *GoBuilder) Chunk(size int, callback func([]map[string]any) error) error {
	offset := 0
	query := gb.Clone()
	query.Limit(offset, size)
	// At this point we need the actual database connection and query execution
	// For now we are only defining the interface
	return nil
}

// Clone creates a deep copy of the current builder
func (gb *GoBuilder) Clone() *GoBuilder {
	clone := &GoBuilder{
		tableClause:   gb.tableClause,
		selectClause:  gb.selectClause,
		whereClause:   gb.whereClause,
		groupByClause: gb.groupByClause,
		havingClause:  gb.havingClause,
		orderByClause: gb.orderByClause,
		limitClause:   gb.limitClause,
		unionClause:   gb.unionClause,
		joinClauses:   make([]string, len(gb.joinClauses)),
		paramsClause:  make([]any, len(gb.paramsClause)),
		counterClause: gb.counterClause,
		sqlDialect:    gb.sqlDialect,
		holderCode:    gb.holderCode,
	}
	copy(clone.joinClauses, gb.joinClauses)
	copy(clone.paramsClause, gb.paramsClause)
	return clone
}

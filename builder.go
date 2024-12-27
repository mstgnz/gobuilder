package gobuilder

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

type SQLDialect string

const (
	Postgres  SQLDialect = "postgres"
	MySQL     SQLDialect = "mysql"
	SQLite    SQLDialect = "sqlite"
	SQLServer SQLDialect = "sqlserver"
	Oracle    SQLDialect = "oracle"
)

// default time for query statute of limitations, you can change this value for each query as you need.
var Timeout time.Duration = 30

type GoBuilder struct {
	tableClause   string
	selectClause  string
	whereClause   string
	groupByClause string
	havingClause  string
	orderByClause string
	limitClause   string
	unionClause   string
	joinClauses   []string
	paramsClause  []any
	counterClause int
	holderClause  SQLDialect
	holderCode    string
	err           error
}

// NewGoBuilder initializes a new instance of GoBuilder
func NewGoBuilder(holderClause SQLDialect) *GoBuilder {
	gb := &GoBuilder{
		paramsClause:  []any{},
		counterClause: 1,
		holderClause:  holderClause,
	}
	gb.holderCode = gb.getPlaceholderCode()
	return gb
}

// Table specifies the table name for the query
func (gb *GoBuilder) Table(table string) *GoBuilder {
	gb.tableClause = table
	return gb
}

// Select defines the columns to be selected in the query
func (gb *GoBuilder) Select(columns ...string) *GoBuilder {
	if gb.tableClause == "" {
		gb.err = fmt.Errorf("table name is required")
		return gb
	}
	if len(columns) == 0 {
		columns = append(columns, "*")
	}

	// Eğer WITH clause varsa (CTE), onu koru
	withClause := ""
	if strings.HasPrefix(gb.selectClause, "WITH ") {
		withClause = gb.selectClause + " "
	}

	gb.selectClause = fmt.Sprintf("%sSELECT %s FROM %s", withClause, strings.Join(columns, ", "), gb.tableClause)
	return gb
}

// SelectDistinct defines the distinct columns to be selected
func (gb *GoBuilder) SelectDistinct(columns ...string) *GoBuilder {
	if len(columns) == 0 {
		columns = append(columns, "*")
	}
	gb.selectClause = fmt.Sprintf("SELECT DISTINCT %s FROM %s", strings.Join(columns, ", "), gb.tableClause)
	return gb
}

// Create adds an INSERT INTO statement to the query with bind parameters
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

// Update builds an UPDATE statement with bind parameters
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

// Delete builds a DELETE statement
func (gb *GoBuilder) Delete() *GoBuilder {
	gb.selectClause = fmt.Sprintf("DELETE FROM %s", gb.tableClause)
	return gb
}

// Where adds a WHERE clause with bind parameters
func (gb *GoBuilder) Where(key, opt string, val any) *GoBuilder {
	var clause string
	switch v := val.(type) {
	case *GoBuilder:
		subQuery, subParams := v.Prepare()
		for _, param := range subParams {
			gb.addParam(param)
		}
		clause = fmt.Sprintf("%s %s (%s)", key, opt, subQuery)
	default:
		clause = fmt.Sprintf("%s %s %s", key, opt, gb.addParam(val))
	}
	gb.addClause("AND", clause)
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
func (gb *GoBuilder) Union(sql string) *GoBuilder {
	if gb.unionClause == "" {
		gb.unionClause = fmt.Sprintf("UNION %v", sql)
	} else {
		gb.unionClause = fmt.Sprintf("%v UNION %v", gb.unionClause, sql)
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

	// WITH clause'u varsa, onu ilk sıraya ekle
	if strings.HasPrefix(gb.selectClause, "WITH ") {
		parts := strings.SplitN(gb.selectClause, " SELECT ", 2)
		if len(parts) > 1 {
			clauses = append(clauses, parts[0])
			gb.selectClause = "SELECT " + parts[1]
		}
	}

	// Diğer clause'ları ekle
	if gb.selectClause != "" {
		clauses = append(clauses, gb.selectClause)
	}
	if len(gb.joinClauses) > 0 {
		clauses = append(clauses, strings.Join(gb.joinClauses, " "))
	}
	if gb.whereClause != "" {
		clauses = append(clauses, gb.whereClause)
	}
	if gb.groupByClause != "" {
		clauses = append(clauses, gb.groupByClause)
	}
	if gb.havingClause != "" {
		clauses = append(clauses, gb.havingClause)
	}
	if gb.orderByClause != "" {
		clauses = append(clauses, gb.orderByClause)
	}
	if gb.limitClause != "" {
		clauses = append(clauses, gb.limitClause)
	}
	if gb.unionClause != "" {
		clauses = append(clauses, gb.unionClause)
	}

	query := strings.Join(clauses, " ")
	re := regexp.MustCompile(`\s+`)
	query = strings.TrimSpace(re.ReplaceAllString(query, " "))
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
	if gb.holderClause == MySQL || gb.holderClause == SQLite {
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
		cleaned := strings.ReplaceAll(v, "'", "''")
		return "'" + cleaned + "'"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (gb *GoBuilder) getPlaceholderCode() string {
	switch gb.holderClause {
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
	if gb.holderClause != MySQL {
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
	if gb.holderClause != SQLServer {
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
	if gb.holderClause != SQLite {
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
	for _, param := range subParams {
		gb.addParam(param)
	}

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

package query

import (
	"database/sql"
	"fmt"
	"imohamedsheta/gocrud/database"
	"imohamedsheta/gocrud/pkg/logger"
	"strings"
)

type QueryBuilder struct {
	tableName      string
	joins          []Join
	fields         []string
	where          []string
	orderBy        []string
	groupBy        []string
	having         []string
	limit          int
	offset         int
	values         []any
	updated_fields []string
}

type Join struct {
	Table    string
	First    string
	Operator string
	Second   string
	JoinType string // "INNER", "LEFT", etc.
}

func Table(tableName string) *QueryBuilder {
	return &QueryBuilder{
		tableName: tableName,
		joins:     []Join{},
		fields:    []string{"*"},
		where:     []string{},
		groupBy:   []string{},
		having:    []string{},
		orderBy:   []string{},
		limit:     0,
		offset:    0,
		values:    []any{},
	}
}

func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
	qb.fields = fields
	return qb
}

func (qb *QueryBuilder) Where(field, operator string, value interface{}) *QueryBuilder {
	qb.where = append(qb.where, fmt.Sprintf("%s %s ?", field, operator))
	qb.values = append(qb.values, value)
	return qb
}

func (qb *QueryBuilder) OrderBy(orderBy []string) *QueryBuilder {
	qb.orderBy = orderBy
	return qb
}

func (qb *QueryBuilder) Limit(limit int) *QueryBuilder {
	qb.limit = limit
	return qb
}

func (qb *QueryBuilder) Offset(offset int) *QueryBuilder {
	qb.offset = offset
	return qb
}

func (qb *QueryBuilder) GroupBy(groupBy []string) *QueryBuilder {
	qb.groupBy = groupBy
	return qb
}

func (qb *QueryBuilder) Having(having []string) *QueryBuilder {
	qb.having = having
	return qb
}

func (qb *QueryBuilder) Join(table, first, operator, second string) *QueryBuilder {
	qb.joins = append(qb.joins, Join{
		Table:    table,
		First:    first,
		Operator: operator,
		Second:   second,
		JoinType: "INNER",
	})
	return qb
}

func (qb *QueryBuilder) LeftJoin(table, first, operator, second string) *QueryBuilder {
	qb.joins = append(qb.joins, Join{
		Table:    table,
		First:    first,
		Operator: operator,
		Second:   second,
		JoinType: "LEFT",
	})
	return qb
}

func (qb *QueryBuilder) Count() ([]map[string]interface{}, error) {
	qb.fields = []string{"COUNT(*) AS count"}
	return qb.Execute()
}

func (qb *QueryBuilder) ToSql() string {
	return qb.Build()
}

func (qb *QueryBuilder) Paginate(page, perPage int) ([]map[string]interface{}, error) {
	if page < 1 {
		page = 1
	}
	qb.limit = perPage
	qb.offset = (page - 1) * perPage

	return qb.Execute()
}

func (qb *QueryBuilder) First() (map[string]interface{}, error) {
	results, err := qb.Limit(1).Execute()
	if results != nil {
		return results[0], err
	}
	return nil, err
}

func (qb *QueryBuilder) Get() ([]map[string]interface{}, error) {
	return qb.Execute()
}

func (qb *QueryBuilder) GetValues() []any {
	return qb.values
}

func (qb *QueryBuilder) InsertSql(data []map[string]interface{}) (string, []any, error) {
	if len(data) == 0 {
		return fmt.Sprintf("INSERT INTO %s DEFAULT VALUES", qb.tableName), nil, nil
	}

	// Get fields from the first record
	fields := []string{}
	for key := range data[0] {
		fields = append(fields, key)
	}

	// Prepare values and placeholders
	var values []any
	placeholdersList := []string{}
	for _, record := range data {
		placeholders := []string{}
		for _, field := range fields {
			// Directly use the value from the record
			values = append(values, record[field])
			placeholders = append(placeholders, "?")
		}
		placeholdersList = append(placeholdersList, "("+strings.Join(placeholders, ", ")+")")
	}

	// Build the final query
	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		qb.tableName,
		strings.Join(fields, ", "),
		strings.Join(placeholdersList, ", "),
	)

	return query, values, nil
}

func (qb *QueryBuilder) Insert(data []map[string]interface{}) (sql.Result, error) {
	db := database.DB()
	query, values, err := qb.InsertSql(data)
	if err != nil {
		return nil, err
	}

	result, err := db.Exec(query, values...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (qb *QueryBuilder) UpdateSql(data map[string]any) (string, []any, error) {
	if len(data) == 0 {
		return "", nil, fmt.Errorf("no data provided for update")
	}

	query := qb.UpdateBuild(data)

	return query, qb.values, nil
}

func (qb *QueryBuilder) Update(data map[string]any) (sql.Result, error) {
	db := database.DB()
	query, values, err := qb.UpdateSql(data)

	if err != nil {
		return nil, err
	}

	result, err := db.Exec(query, values...)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (qb *QueryBuilder) Execute() ([]map[string]any, error) {
	db := database.DB() // Get the database connection

	query := qb.Build() // Build the SQL query string with bound values

	// Execute the query with parameters
	rows, err := db.Query(query, qb.values...)
	if err != nil {
		logger.Log().Error(err.Error()) // Log the error if query fails
		return nil, err                 // Return the error to the caller
	}
	defer rows.Close() // Ensure that the rows are closed once the function exits

	// all result rows
	var results []map[string]any

	// Get the column names from the result set (used for mapping results)
	columns, err := rows.Columns()
	if err != nil {
		logger.Log().Error(err.Error())
		return nil, err
	}

	// Prepare for scanning values into the result
	values := make([]any, len(columns))    // Actual values will be stored here
	valuePtrs := make([]any, len(columns)) // Need for rows.Scan to manipulate the values
	for i := range values {
		valuePtrs[i] = &values[i] // Point each valuePtr to returned value
	}

	// Iterate through all rows
	for rows.Next() {
		// Scan the row values need pointers to write the values to
		if err := rows.Scan(valuePtrs...); err != nil {
			logger.Log().Error(err.Error())
			return nil, err
		}

		// the result is a map of column names to their values
		result := make(map[string]any)

		for i, col := range columns {
			currentValue := values[i] // Get the value for the current column

			// If the value is a byte slice (e.g., text or varchar), convert it to string
			if currentInBytes, ok := currentValue.([]byte); ok {
				result[col] = string(currentInBytes)
			} else {
				result[col] = currentValue // If not []byte, store the value as it is
			}
		}

		results = append(results, result)
	}

	return results, nil
}

func (qb *QueryBuilder) Values() []any {
	return qb.values
}

func (qb *QueryBuilder) Build() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("SELECT %s FROM %s", strings.Join(qb.fields, ", "), qb.tableName))

	for _, join := range qb.joins {
		sb.WriteString(fmt.Sprintf(" %s JOIN %s ON %s %s %s", join.JoinType, join.Table, join.First, join.Operator, join.Second))
	}

	if len(qb.where) > 0 {
		sb.WriteString(" WHERE " + strings.Join(qb.where, " AND "))
	}
	if len(qb.groupBy) > 0 {
		sb.WriteString(" GROUP BY " + strings.Join(qb.groupBy, ", "))
	}
	if len(qb.having) > 0 {
		sb.WriteString(" HAVING " + strings.Join(qb.having, " AND "))
	}
	if len(qb.orderBy) > 0 {
		sb.WriteString(" ORDER BY " + strings.Join(qb.orderBy, ", "))
	}
	if qb.limit != 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", qb.limit))
	}
	if qb.offset != 0 {
		sb.WriteString(fmt.Sprintf(" OFFSET %d", qb.offset))
	}

	return sb.String()
}

// Update builds the SQL query for an UPDATE statement
func (qb *QueryBuilder) UpdateBuild(data map[string]any) string {
	var sb strings.Builder

	var updatedFields []string
	var values []any

	for field, value := range data {
		updatedFields = append(updatedFields, fmt.Sprintf("%s = ?", field))
		values = append(values, value)
	}

	qb.values = append(values, qb.values...)

	sb.WriteString(fmt.Sprintf("UPDATE %s SET %s ", qb.tableName, strings.Join(updatedFields, ", ")))

	if len(qb.where) > 0 {
		sb.WriteString(" WHERE " + strings.Join(qb.where, " AND "))
	}

	if qb.limit != 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", qb.limit))
	}

	return sb.String()
}

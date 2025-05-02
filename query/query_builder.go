package query

import (
	"database/sql"
	"fmt"
	"imohamedsheta/gocrud/database"
	"imohamedsheta/gocrud/pkg/logger"
	"math"
	"strings"

	"github.com/jmoiron/sqlx"
)

type QueryBuilder struct {
	tableName string
	joins     []Join
	fields    []string
	where     []string
	orderBy   []string
	groupBy   []string
	having    []string
	limit     int
	offset    int
	values    []any
	tx        *sqlx.Tx // Transaction field
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

func (qb *QueryBuilder) WhereIn(field string, values []interface{}) *QueryBuilder {
	placeholders := make([]string, len(values))
	for i := range values {
		placeholders[i] = "?"
	}

	qb.where = append(qb.where, fmt.Sprintf("%s IN (%s)", field, strings.Join(placeholders, ", ")))
	qb.values = append(qb.values, values...)
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

func (qb *QueryBuilder) Paginate(perPage int, page int, withCount bool) ([]map[string]any, map[string]any, error) {

	if page < 1 {
		page = 1
	}
	qb.limit = perPage
	qb.offset = (page - 1) * perPage

	// First, execute the paginated query to get the results
	results, err := qb.Execute()
	if err != nil {
		return nil, nil, err
	}

	meta := map[string]any{
		"per_page":     perPage,
		"current_page": page,
	}

	if withCount {
		// Use the Count() method to get the total count of records
		countResults, err := qb.Count()
		if err != nil {
			return nil, nil, err
		}
		if len(countResults) > 0 {
			count := countResults[0]["count"]

			meta["total_count"] = count
			meta["last_page"] = int(math.Ceil(float64(count.(int64)) / float64(perPage)))
			return results, meta, nil
		}
	}

	// If no count is needed, just return the results
	return results, meta, nil
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

func (qb *QueryBuilder) Insert(data []map[string]any) (sql.Result, error) {
	db := database.DB()
	query, values, err := qb.InsertSql(data)
	if err != nil {
		return nil, err
	}

	var result sql.Result
	if qb.tx != nil {
		result, err = qb.tx.Exec(query, values...) // Execute within the transaction
	} else {
		result, err = db.Exec(query, values...) // Normal execution
	}

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
	query, _, err := qb.UpdateSql(data)

	if err != nil {
		return nil, err
	}

	var result sql.Result
	if qb.tx != nil {
		result, err = qb.tx.Exec(query, qb.values...) // Execute within the transaction
	} else {
		result, err = db.Exec(query, qb.values...) // Normal execution
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (qb *QueryBuilder) Execute() ([]map[string]any, error) {
	db := database.DB()
	var tx *sqlx.Tx

	if qb.tx != nil {
		tx = qb.tx
	}

	query := qb.Build()

	// Execute the query with parameters
	rows, err := executeQuery(tx, db, query, qb.values...)
	if err != nil {
		logger.Log().Error(err.Error()) // Log the error if query fails
		return nil, err                 // Return the error to the caller
	}
	defer rows.Close() // Ensure that the rows are closed once the function exits

	var results []map[string]any
	columns, err := rows.Columns()
	if err != nil {
		logger.Log().Error(err.Error())
		return nil, err
	}

	values := make([]any, len(columns))
	valuePtrs := make([]any, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		if err := rows.Scan(valuePtrs...); err != nil {
			logger.Log().Error(err.Error())
			return nil, err
		}

		result := make(map[string]any)
		for i, col := range columns {
			currentValue := values[i]
			if currentInBytes, ok := currentValue.([]byte); ok {
				result[col] = string(currentInBytes)
			} else {
				result[col] = currentValue
			}
		}

		results = append(results, result)
	}

	return results, nil
}

func executeQuery(tx *sqlx.Tx, db *sqlx.DB, query string, args ...any) (*sql.Rows, error) {
	if tx != nil {
		return tx.Query(query, args...)
	}
	return db.Query(query, args...)
}

func (qb *QueryBuilder) Values() []any {
	return qb.values
}

// StartTransaction starts a new transaction
func (qb *QueryBuilder) BeginTransaction() error {
	db := database.DB()
	tx, err := db.Beginx() // Start a transaction
	if err != nil {
		return err
	}
	qb.tx = tx
	return nil
}

// CommitTransaction commits the current transaction
func (qb *QueryBuilder) Commit() error {
	if qb.tx == nil {
		return fmt.Errorf("no transaction to commit")
	}
	return qb.tx.Commit()
}

// RollbackTransaction rolls back the current transaction
func (qb *QueryBuilder) Rollback() error {
	if qb.tx == nil {
		return fmt.Errorf("no transaction to rollback")
	}
	return qb.tx.Rollback()
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

func (qb *QueryBuilder) DeleteBuild() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("DELETE FROM %s", qb.tableName))

	if len(qb.where) > 0 {
		sb.WriteString(" WHERE " + strings.Join(qb.where, " AND "))
	}
	if len(qb.orderBy) > 0 {
		sb.WriteString(" ORDER BY " + strings.Join(qb.orderBy, ", "))
	}
	if qb.limit != 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", qb.limit))
	}

	return sb.String()
}

func (qb *QueryBuilder) Delete() (sql.Result, error) {
	db := database.DB()
	query := qb.DeleteBuild()

	var result sql.Result
	var err error

	if qb.tx != nil {
		result, err = qb.tx.Exec(query, qb.values...) // Execute within the transaction
	} else {
		result, err = db.Exec(query, qb.values...) // Normal execution
	}

	if err != nil {
		return nil, err
	}

	return result, nil
}

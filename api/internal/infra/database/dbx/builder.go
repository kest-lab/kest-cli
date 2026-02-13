package dbx

import (
	"context"
	"database/sql"
	"strings"

	"gorm.io/gorm"
)

// Builder wraps gorm.DB to provide a fluent API for building queries.
// Use the Use() function to create a builder from an existing *gorm.DB instance.
type Builder struct {
	db *gorm.DB
}

// Use creates a builder from an existing gorm.DB session.
// This is the primary way to create a Builder - pass in your DI-injected *gorm.DB.
func Use(db *gorm.DB) *Builder {
	return &Builder{db: db.Session(&gorm.Session{})}
}

// DB returns the underlying gorm.DB reference.
func (b *Builder) DB() *gorm.DB {
	return b.db
}

// WithContext attaches a context to the builder session.
func (b *Builder) WithContext(ctx context.Context) *Builder {
	b.db = b.db.WithContext(ctx)
	return b
}

// Table sets the table name for the query.
func (b *Builder) Table(name string) *Builder {
	b.db = b.db.Table(name)
	return b
}

// Model sets the model for the query.
func (b *Builder) Model(value any) *Builder {
	b.db = b.db.Model(value)
	return b
}

// Select specifies fields to retrieve.
func (b *Builder) Select(fields ...string) *Builder {
	if len(fields) == 0 {
		return b
	}
	b.db = b.db.Select(strings.Join(fields, ", "))
	return b
}

// Where adds a where clause to the query.
func (b *Builder) Where(query any, args ...any) *Builder {
	b.db = b.db.Where(query, args...)
	return b
}

// OrWhere adds an OR where clause to the query.
func (b *Builder) OrWhere(query any, args ...any) *Builder {
	b.db = b.db.Or(query, args...)
	return b
}

// WhereIn adds a WHERE IN clause to the query.
func (b *Builder) WhereIn(column string, values any) *Builder {
	b.db = b.db.Where(column+" IN ?", values)
	return b
}

// WhereNotIn adds a WHERE NOT IN clause to the query.
func (b *Builder) WhereNotIn(column string, values any) *Builder {
	b.db = b.db.Where(column+" NOT IN ?", values)
	return b
}

// Order adds an order by clause.
func (b *Builder) Order(value string) *Builder {
	b.db = b.db.Order(value)
	return b
}

// Group adds a group by clause.
func (b *Builder) Group(columns string) *Builder {
	b.db = b.db.Group(columns)
	return b
}

// Having adds a having clause.
func (b *Builder) Having(query any, args ...any) *Builder {
	b.db = b.db.Having(query, args...)
	return b
}

// Limit limits the number of records returned.
func (b *Builder) Limit(limit int) *Builder {
	b.db = b.db.Limit(limit)
	return b
}

// Offset sets the offset for records returned.
func (b *Builder) Offset(offset int) *Builder {
	b.db = b.db.Offset(offset)
	return b
}

// Join adds an inner join clause.
func (b *Builder) Join(expr string, args ...any) *Builder {
	b.db = b.db.Joins(expr, args...)
	return b
}

// LeftJoin adds a left join clause.
func (b *Builder) LeftJoin(expr string, args ...any) *Builder {
	b.db = b.db.Joins("LEFT JOIN "+expr, args...)
	return b
}

// RightJoin adds a right join clause.
func (b *Builder) RightJoin(expr string, args ...any) *Builder {
	b.db = b.db.Joins("RIGHT JOIN "+expr, args...)
	return b
}

// First fetches the first record into dest.
func (b *Builder) First(dest any) error {
	return b.db.First(dest).Error
}

// Get fetches all matching records into dest.
func (b *Builder) Get(dest any) error {
	return b.db.Find(dest).Error
}

// Value fetches a single column into dest.
func (b *Builder) Value(column string, dest any) error {
	return b.db.Limit(1).Pluck(column, dest).Error
}

// Count returns the total number of matching records.
func (b *Builder) Count() (int64, error) {
	var total int64
	if err := b.db.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

// Exists returns true if any record matches the current query.
func (b *Builder) Exists() (bool, error) {
	rows, err := b.db.Session(&gorm.Session{}).Select("1").Limit(1).Rows()
	if err != nil {
		return false, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	return rows.Next(), nil
}

// Paginate executes the query with built-in pagination semantics.
type Pagination struct {
	Page    int
	PerPage int
	Total   int64
}

// Paginate fetches the requested page into dest and returns pagination metadata.
func (b *Builder) Paginate(page, perPage int, dest any) (*Pagination, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 20
	}

	countDB := b.db.Session(&gorm.Session{})
	var total int64
	if err := countDB.Count(&total).Error; err != nil {
		return nil, err
	}

	queryDB := b.db.Session(&gorm.Session{})
	if err := queryDB.Offset((page - 1) * perPage).Limit(perPage).Find(dest).Error; err != nil {
		return nil, err
	}

	return &Pagination{
		Page:    page,
		PerPage: perPage,
		Total:   total,
	}, nil
}

// Create inserts value using the underlying query builder context.
func (b *Builder) Create(value any) error {
	return b.db.Create(value).Error
}

// Update updates selected columns on the matched rows.
func (b *Builder) Update(column string, value any) error {
	return b.db.Update(column, value).Error
}

// Updates updates multiple columns on the matched rows.
func (b *Builder) Updates(values any) error {
	return b.db.Updates(values).Error
}

// Delete removes matched rows.
func (b *Builder) Delete(value any, conds ...any) error {
	return b.db.Delete(value, conds...).Error
}

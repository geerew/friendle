package dao

import (
	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/utils/pagination"
)

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Options defines a limited set of database query options
//
// This will passed to the internal DAO builder options struct for use in the query
// building process
type Options struct {
	// OrderBy can be used to order the results
	//
	// Example: []string{"id DESC", "display_name ASC"}
	OrderBy []string

	// OrderByClause can be used to set a custom ORDER BY clause, for example
	// when using a case expression
	OrderByClause squirrel.Sqlizer

	// Where is any valid squirrel WHERE expression
	//
	// Examples:
	//
	//   EQ:   squirrel.Eq{"id": "123"}
	//   IN:   squirrel.Eq{"id": []string{"123", "456"}}
	//   OR:   squirrel.Or{squirrel.Expr("id = ?", "123"), squirrel.Expr("id = ?", "456")}
	//   AND:  squirrel.And{squirrel.Eq{"id": "123"}, squirrel.Eq{"display_name": "Alice"}}
	//   LIKE: squirrel.Like{"display_name": "%ali%"}
	//   NOT:  squirrel.NotEq{"id": "123"}
	Where squirrel.Sqlizer

	// Pagination applies OFFSET/LIMIT to list queries. When set, listGeneric runs a COUNT
	// query and updates the same instance via SetCount before selecting the page.
	Pagination *pagination.Pagination

	// ApiQuery is the list `q` query string from an HTTP request
	ApiQuery string
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// NewOptions creates an empty Options builder
func NewOptions() *Options {
	return &Options{}
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// WithOrderBy sets the ORDER BY
//
// Calling multiple times will override the previous WithOrderBy call
func (o *Options) WithOrderBy(fields ...string) *Options {
	o.OrderBy = append([]string(nil), fields...)
	return o
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// WithOrderByClause sets a custom ORDER BY clause
//
// Calling multiple times will override the previous WithOrderByClause call
func (o *Options) WithOrderByClause(clause squirrel.Sqlizer) *Options {
	o.OrderByClause = clause
	return o
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// WithWhere sets the WHERE clause using a squirrel.Sqlizer
//
// Calling multiple times will override the previous WithWhere call
func (o *Options) WithWhere(pred squirrel.Sqlizer) *Options {
	o.Where = pred
	return o
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// WithPagination sets the pagination options
//
// Calling multiple times will override the previous WithPagination call
func (o *Options) WithPagination(p *pagination.Pagination) *Options {
	o.Pagination = p
	return o
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// WithApiQuery sets the API query for LIST operations, which is typically passed from
// c.Query("q", ""). An empty string is the same as no query
//
// Calling multiple times will override the previous WithApiQuery call
func (o *Options) WithApiQuery(q string) *Options {
	o.ApiQuery = q
	return o
}

package dao

import (
	"github.com/Masterminds/squirrel"
	"github.com/geerew/friendle/models"
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

	// IncludeMembers includes group members when querying groups
	//
	// Valid when querying groups
	IncludeMembers bool

	// IncludeJoinRequests includes pending join requests when querying groups
	//
	// Valid when querying groups
	IncludeJoinRequests bool

	// IncludeParticipations includes round participations when querying rounds
	//
	// Valid when querying rounds
	IncludeParticipations bool

	// IncludeGuesses includes guess rows on round participations
	//
	// Valid when querying rounds
	IncludeGuesses bool

	// IncludeUsers includes related user rows on group members, round pickers, and
	// round participations
	//
	// Valid when querying groups or rounds
	IncludeUsers bool

	// IncludeMemberCount sets MemberCount on groups loaded via groups or group member queries
	//
	// Valid when querying groups or group members with WithGroup()
	IncludeMemberCount bool

	// IncludeGroup includes the parent group when querying group members
	//
	// Valid when querying group members
	IncludeGroup bool
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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// WithMembers enables group member inclusion in queries
//
// Can be used when querying groups
func (o *Options) WithMembers() *Options {
	o.IncludeMembers = true
	return o
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// WithJoinRequests enables pending join request inclusion in queries
//
// Can be used when querying groups
func (o *Options) WithJoinRequests() *Options {
	o.IncludeJoinRequests = true
	return o
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// WithParticipations enables round participation inclusion in queries
//
// Can be used when querying rounds
func (o *Options) WithParticipations() *Options {
	o.IncludeParticipations = true
	return o
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// WithGuesses enables guess rows on round participations
//
// Can be used when querying rounds
func (o *Options) WithGuesses() *Options {
	o.IncludeGuesses = true
	return o
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// WithUsers enables related user rows on loaded group members, round pickers, and
// round participations
//
// Can be used when querying groups or rounds
func (o *Options) WithUsers() *Options {
	o.IncludeUsers = true
	return o
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// WithMemberCount enables member count inclusion on loaded groups
//
// Can be used when querying groups or group members with WithGroup()
func (o *Options) WithMemberCount() *Options {
	o.IncludeMemberCount = true
	return o
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// WithGroup enables parent group inclusion when querying group members
func (o *Options) WithGroup() *Options {
	o.IncludeGroup = true
	return o
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// WithGroupNameSearch applies a case-insensitive name search with relevance ordering
//
// Can be used when querying groups
func (o *Options) WithGroupNameSearch(q string) *Options {
	g := models.GROUP_TABLE
	like := "%" + q + "%"

	o.WithWhere(squirrel.Like{"LOWER(" + g + ".name)": like})
	o.WithOrderByClause(groupNameSearchOrder(g, q))

	return o
}

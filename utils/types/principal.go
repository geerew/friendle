package types

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Principal identifies the authenticated caller attached to a request context
type Principal struct {
	UserID   string
	SiteRole SiteRole
}

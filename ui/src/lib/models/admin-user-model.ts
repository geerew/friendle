import { array, object, picklist, string, type InferOutput } from 'valibot';
import { BasePaginationSchema, type PaginationReqParams } from './pagination-model';

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// SiteRoleSchema represents the role of a site user
const SiteRoleSchema = picklist(['site_admin', 'site_user']);

export type SiteRole = InferOutput<typeof SiteRoleSchema>;

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// SelectSiteRoles represents the select options for site roles
export const SelectSiteRoles = [
	{ value: 'site_user', label: 'User' },
	{ value: 'site_admin', label: 'Admin' }
];

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// formatSiteRole formats a site role for display
export function formatSiteRole(role: SiteRole): string {
	return SelectSiteRoles.find((item) => item.value === role)?.label ?? role;
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// siteRoleBadgeLabel returns the badge label for a site role
export function siteRoleBadgeLabel(role: SiteRole): string {
	return role === 'site_admin' ? 'admin' : 'user';
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// AdminUserSchema represents a user
export const AdminUserSchema = object({
	id: string(),
	username: string(),
	displayName: string(),
	siteRole: SiteRoleSchema
});

export type AdminUserModel = InferOutput<typeof AdminUserSchema>;

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// AdminUserCreateSchema represents a request to create a user
export const AdminUserCreateSchema = object({
	username: string(),
	displayName: string(),
	password: string(),
	siteRole: SiteRoleSchema
});

export type AdminUserCreateModel = InferOutput<typeof AdminUserCreateSchema>;

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// AdminUserPaginationSchema represents a pagination of users
export const AdminUserPaginationSchema = object({
	...BasePaginationSchema.entries,
	items: array(AdminUserSchema)
});

export type AdminUserPaginationModel = InferOutput<typeof AdminUserPaginationSchema>;

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// AdminUserReqParams represents the parameters for listing users
export type AdminUserReqParams = PaginationReqParams;

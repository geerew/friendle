import { array, object, picklist, string, type InferOutput } from 'valibot';
import { BasePaginationSchema, type PaginationReqParams } from './pagination-model';

const SiteRoleSchema = picklist(['site_admin', 'site_user']);

export type SiteRole = InferOutput<typeof SiteRoleSchema>;

export const SelectSiteRoles = [
	{ value: 'site_user', label: 'User' },
	{ value: 'site_admin', label: 'Admin' }
];

export function formatSiteRole(role: SiteRole): string {
	return SelectSiteRoles.find((item) => item.value === role)?.label ?? role;
}

export const AdminUserSchema = object({
	id: string(),
	username: string(),
	displayName: string(),
	siteRole: SiteRoleSchema
});

export type AdminUserModel = InferOutput<typeof AdminUserSchema>;

export const AdminUserCreateSchema = object({
	username: string(),
	displayName: string(),
	password: string(),
	siteRole: SiteRoleSchema
});

export type AdminUserCreateModel = InferOutput<typeof AdminUserCreateSchema>;

export const AdminUserPaginationSchema = object({
	...BasePaginationSchema.entries,
	items: array(AdminUserSchema)
});

export type AdminUserPaginationModel = InferOutput<typeof AdminUserPaginationSchema>;

export type AdminUserReqParams = PaginationReqParams;

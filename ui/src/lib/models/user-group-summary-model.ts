import { number, object, optional, picklist, string, type InferOutput } from 'valibot';

const GroupRoleSchema = picklist(['group_admin', 'group_user']);

export const UserGroupSummarySchema = object({
	id: string(),
	name: string(),
	memberCount: number(),
	groupRole: optional(GroupRoleSchema)
});

export type UserGroupSummaryModel = InferOutput<typeof UserGroupSummarySchema>;

import { array, number, object, string, type InferOutput } from 'valibot';
import { BasePaginationSchema, type PaginationReqParams } from './pagination-model';
import { GroupRoleSchema } from './group-model';

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupMemberSchema represents a group member
export const GroupMemberSchema = object({
	userId: string(),
	displayName: string(),
	groupRole: GroupRoleSchema,
	timesPicked: number(),
	pickerSkips: number()
});

export type GroupMemberModel = InferOutput<typeof GroupMemberSchema>;

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupMemberPaginationSchema represents a paginated list of group members
export const GroupMemberPaginationSchema = object({
	...BasePaginationSchema.entries,
	items: array(GroupMemberSchema)
});

export type GroupMemberPaginationModel = InferOutput<typeof GroupMemberPaginationSchema>;

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListGroupMembersParams represents the parameters for listing group members
export type ListGroupMembersParams = PaginationReqParams;

// SelectGroupRoles represents radio options for group member roles
export const SelectGroupRoles = [
	{ value: 'group_admin', label: 'Admin' },
	{ value: 'group_user', label: 'Member' }
] as const;

// UpdateGroupMemberRoleRequest represents a group member role update request
export type UpdateGroupMemberRoleRequest = Pick<GroupMemberModel, 'groupRole'>;

import { array, object, string, type InferOutput } from 'valibot';
import { BasePaginationSchema, type PaginationReqParams } from './pagination-model';
import { GroupRoleSchema } from './group-model';

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupMemberSchema represents a group member
export const GroupMemberSchema = object({
	userId: string(),
	displayName: string(),
	groupRole: GroupRoleSchema
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

import {
	array,
	boolean,
	number,
	object,
	optional,
	picklist,
	string,
	type InferOutput
} from 'valibot';
import { BasePaginationSchema, type PaginationReqParams } from './pagination-model';

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupRoleSchema represents the role of a group member
export const GroupRoleSchema = picklist(['group_admin', 'group_user']);

export type GroupRole = InferOutput<typeof GroupRoleSchema>;

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// JoinRequestStatusSchema represents the status of a join request
export const JoinRequestStatusSchema = picklist(['pending', 'rejected']);

export type JoinRequestStatus = InferOutput<typeof JoinRequestStatusSchema>;

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupAdminSummarySchema represents the summary of a group admin
export const GroupAdminSummarySchema = object({
	pendingJoinRequestCount: number(),
	rejectedJoinRequestCount: number()
});

export type GroupAdminSummary = InferOutput<typeof GroupAdminSummarySchema>;

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupSchema represents a group
export const GroupSchema = object({
	id: string(),
	createdAt: string(),
	updatedAt: string(),
	name: string(),
	createdBy: string(),
	memberCount: number(),
	memberThresholdMet: boolean(),
	groupRole: optional(GroupRoleSchema),
	joinRequestStatus: optional(JoinRequestStatusSchema),
	adminSummary: optional(GroupAdminSummarySchema)
});

export type GroupModel = InferOutput<typeof GroupSchema>;

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupPaginationSchema represents a pagination of groups
export const GroupPaginationSchema = object({
	...BasePaginationSchema.entries,
	items: array(GroupSchema)
});

export type GroupPaginationModel = InferOutput<typeof GroupPaginationSchema>;

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// CreateGroupRequest represents a request to create a group
export type CreateGroupRequest = {
	name: string;
};

// UpdateGroupRequest represents a request to update a group
export type UpdateGroupRequest = {
	name: string;
};

// MAX_GROUP_NAME_LENGTH matches the server-side group name limit
export const MAX_GROUP_NAME_LENGTH = 64;

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListGroupsParams represents the parameters for listing groups
export type ListGroupsParams = PaginationReqParams;

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListSelfGroupsParams represents the parameters for listing groups the authenticated user belongs to
export type ListSelfGroupsParams = PaginationReqParams;

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// SearchGroupsParams represents the parameters for searching groups
export type SearchGroupsParams = PaginationReqParams & {
	name: string;
};

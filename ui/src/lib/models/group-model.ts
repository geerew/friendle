import { array, number, object, optional, picklist, string, type InferOutput } from 'valibot';
import { BasePaginationSchema, type PaginationReqParams } from './pagination-model';

export const GroupRoleSchema = picklist(['group_admin', 'group_user']);

export type GroupRole = InferOutput<typeof GroupRoleSchema>;

export const JoinRequestStatusSchema = picklist(['pending', 'rejected']);

export type JoinRequestStatus = InferOutput<typeof JoinRequestStatusSchema>;

export const GroupSchema = object({
	id: string(),
	createdAt: string(),
	updatedAt: string(),
	name: string(),
	createdBy: string(),
	memberCount: number(),
	groupRole: optional(GroupRoleSchema),
	joinRequestStatus: optional(JoinRequestStatusSchema)
});

export type GroupModel = InferOutput<typeof GroupSchema>;

export const GroupPaginationSchema = object({
	...BasePaginationSchema.entries,
	items: array(GroupSchema)
});

export type GroupPaginationModel = InferOutput<typeof GroupPaginationSchema>;

export type CreateGroupRequest = {
	name: string;
};

export type ListGroupsParams = PaginationReqParams;

export type ListSelfGroupsParams = PaginationReqParams;

export type SearchGroupsParams = PaginationReqParams & {
	name: string;
};

import { array, object, string, type InferOutput } from 'valibot';
import { BasePaginationSchema, type PaginationReqParams } from './pagination-model';

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupJoinRequestSchema represents a pending join request
export const GroupJoinRequestSchema = object({
	userId: string(),
	displayName: string()
});

export type GroupJoinRequestModel = InferOutput<typeof GroupJoinRequestSchema>;

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// GroupJoinRequestPaginationSchema represents a paginated list of join requests
export const GroupJoinRequestPaginationSchema = object({
	...BasePaginationSchema.entries,
	items: array(GroupJoinRequestSchema)
});

export type GroupJoinRequestPaginationModel = InferOutput<typeof GroupJoinRequestPaginationSchema>;

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// ListGroupPendingJoinRequestsParams represents the parameters for listing pending join requests
export type ListGroupPendingJoinRequestsParams = PaginationReqParams;

import { array, number, object, string, type InferOutput } from 'valibot';
import { BasePaginationSchema, type PaginationReqParams } from './pagination-model';

export const GroupSchema = object({
	id: string(),
	createdAt: string(),
	updatedAt: string(),
	name: string(),
	createdBy: string(),
	memberCount: number()
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
	q: string;
};

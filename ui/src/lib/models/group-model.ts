import { array, object, optional, string, type InferOutput } from 'valibot';
import { BasePaginationSchema, type PaginationReqParams } from './pagination-model';

export const GroupSchema = object({
	id: string(),
	name: string(),
	createdBy: optional(string())
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

export type ListSelfGroupsParams = PaginationReqParams;

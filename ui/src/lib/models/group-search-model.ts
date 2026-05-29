import { array, boolean, number, object, string, type InferOutput } from 'valibot';
import { BasePaginationSchema, type PaginationReqParams } from './pagination-model';

export const GroupSearchResultSchema = object({
	id: string(),
	name: string(),
	memberCount: number(),
	isMember: boolean(),
	joinPending: boolean()
});

export type GroupSearchResultModel = InferOutput<typeof GroupSearchResultSchema>;

export const GroupSearchPaginationSchema = object({
	...BasePaginationSchema.entries,
	items: array(GroupSearchResultSchema)
});

export type GroupSearchPaginationModel = InferOutput<typeof GroupSearchPaginationSchema>;

export type GroupSearchReqParams = PaginationReqParams & {
	q: string;
};

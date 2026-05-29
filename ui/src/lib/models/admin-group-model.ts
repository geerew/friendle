import { array, number, object, string, type InferOutput } from 'valibot';
import { BasePaginationSchema, type PaginationReqParams } from './pagination-model';

export const AdminGroupSchema = object({
	id: string(),
	name: string(),
	memberCount: number()
});

export type AdminGroupModel = InferOutput<typeof AdminGroupSchema>;

export const AdminGroupPaginationSchema = object({
	...BasePaginationSchema.entries,
	items: array(AdminGroupSchema)
});

export type AdminGroupPaginationModel = InferOutput<typeof AdminGroupPaginationSchema>;

export type AdminGroupReqParams = PaginationReqParams & {
	orderBy?: string;
};

import { number, object } from 'valibot';

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// SelectPaginationPerPage represents the select options for pagination per page
export const SelectPaginationPerPage = [
	{ value: '10', label: '10' },
	{ value: '25', label: '25' },
	{ value: '50', label: '50' },
	{ value: '100', label: '100' }
];

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// BasePaginationSchema represents the base pagination schema
export const BasePaginationSchema = object({
	page: number(),
	perPage: number(),
	totalPages: number(),
	totalItems: number()
});

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// PaginationReqParams represents the parameters for pagination
export type PaginationReqParams = {
	page?: number;
	perPage?: number;
};

import { ApiError, apiFetch, parseJson } from './fetch';
import { buildQueryString } from '$lib/utils';
import { safeParse } from 'valibot';
import {
	GroupPaginationSchema,
	type CreateGroupRequest,
	type GroupModel,
	type GroupPaginationModel,
	type ListSelfGroupsParams
} from '$lib/models/group-model';

export async function listSelfGroups(
	params?: ListSelfGroupsParams
): Promise<GroupPaginationModel> {
	const qs = params ? buildQueryString(params) : '';
	const response = await apiFetch('/api/groups/self' + (qs ? `?${qs}` : ''));

	if (response.ok) {
		const data = await response.json();
		const result = safeParse(GroupPaginationSchema, data);

		if (!result.success) {
			throw new ApiError('Invalid response from the server', response.status);
		}

		return result.output;
	}

	const data = (await response.json()) as { message?: string };
	throw new ApiError(data.message || 'Request failed', response.status);
}

export async function createGroup(data: CreateGroupRequest): Promise<GroupModel> {
	const response = await apiFetch('/api/groups/', {
		method: 'POST',
		body: JSON.stringify(data)
	});
	return parseJson(response);
}

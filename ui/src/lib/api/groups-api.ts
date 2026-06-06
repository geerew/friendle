import { ApiError, apiFetch, parseJson } from './fetch';
import { buildQueryString } from '$lib/utils';
import { safeParse } from 'valibot';
import {
	GroupPaginationSchema,
	GroupSchema,
	type CreateGroupRequest,
	type GroupModel,
	type GroupPaginationModel,
	type ListGroupsParams,
	type ListSelfGroupsParams,
	type SearchGroupsParams
} from '$lib/models/group-model';

export async function listGroups(params?: ListGroupsParams): Promise<GroupPaginationModel> {
	const qs = params ? buildQueryString(params) : '';
	const response = await apiFetch('/api/groups/' + (qs ? `?${qs}` : ''));

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

export async function searchGroups(params: SearchGroupsParams): Promise<GroupPaginationModel> {
	const qs = buildQueryString(params);
	const response = await apiFetch('/api/groups/search' + (qs ? `?${qs}` : ''));

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

export async function getGroup(id: string): Promise<GroupModel> {
	const response = await apiFetch(`/api/groups/${id}`);

	if (response.ok) {
		const data = await response.json();
		const result = safeParse(GroupSchema, data);

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

export async function deleteGroup(groupId: string): Promise<void> {
	const response = await apiFetch(`/api/groups/${groupId}`, { method: 'DELETE' });

	if (response.ok || response.status === 204) {
		return;
	}

	const data = (await response.json()) as { message?: string };
	throw new ApiError(data.message || 'Request failed', response.status);
}

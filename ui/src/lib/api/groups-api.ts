import { ApiError, apiFetch, parseJson } from './fetch';
import { buildQueryString } from '$lib/utils';
import { array, safeParse } from 'valibot';
import {
	GroupSearchPaginationSchema,
	type GroupSearchPaginationModel,
	type GroupSearchReqParams
} from '$lib/models/group-search-model';
import {
	UserGroupSummarySchema,
	type UserGroupSummaryModel
} from '$lib/models/user-group-summary-model';
import type {
	CreateGroupRequest,
	Group,
	GroupDetail,
	JoinRequest,
	UpdateGroupRequest
} from '$lib/types/group';

export async function listMyGroups(): Promise<UserGroupSummaryModel[]> {
	const response = await apiFetch('/api/groups/mine');
	const data = await parseJson<unknown>(response);
	const result = safeParse(array(UserGroupSummarySchema), data);

	if (!result.success) {
		throw new ApiError('Invalid response from the server', response.status);
	}

	return result.output;
}

export async function searchGroups(params: GroupSearchReqParams): Promise<GroupSearchPaginationModel> {
	const qs = buildQueryString(params);
	const response = await apiFetch(`/api/groups/search?${qs}`);
	const data = await parseJson<unknown>(response);
	const result = safeParse(GroupSearchPaginationSchema, data);

	if (!result.success) {
		throw new ApiError('Invalid response from the server', response.status);
	}

	return result.output;
}

export async function getGroup(id: string): Promise<GroupDetail> {
	const response = await apiFetch(`/api/groups/${id}`);
	return parseJson(response);
}

export async function createGroup(data: CreateGroupRequest): Promise<Group> {
	const response = await apiFetch('/api/groups/', {
		method: 'POST',
		body: JSON.stringify(data)
	});
	return parseJson(response);
}

export async function updateGroup(id: string, data: UpdateGroupRequest): Promise<Group> {
	const response = await apiFetch(`/api/groups/${id}`, {
		method: 'PUT',
		body: JSON.stringify(data)
	});
	return parseJson(response);
}

export async function deleteGroup(id: string): Promise<void> {
	const response = await apiFetch(`/api/groups/${id}`, { method: 'DELETE' });
	await parseJson(response);
}

export async function joinGroup(id: string): Promise<JoinRequest> {
	const response = await apiFetch(`/api/groups/${id}/join-requests`, { method: 'POST' });

	if (response.ok) {
		return parseJson(response);
	}

	const data = (await response.json()) as { message?: string };
	throw new ApiError(data.message || 'Request failed', response.status);
}

export async function cancelJoinRequest(groupId: string): Promise<void> {
	const response = await apiFetch(`/api/groups/${groupId}/join-requests/me`, { method: 'DELETE' });

	if (response.ok || response.status === 204) {
		return;
	}

	const data = (await response.json()) as { message?: string };
	throw new ApiError(data.message || 'Request failed', response.status);
}

export async function leaveGroup(id: string): Promise<void> {
	const response = await apiFetch(`/api/groups/${id}/leave`, { method: 'POST' });
	await parseJson(response);
}

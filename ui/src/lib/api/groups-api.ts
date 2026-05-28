import { ApiError, apiFetch, parseJson } from './fetch';
import { array, safeParse } from 'valibot';
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
	const response = await apiFetch('/api/groups/');
	const data = await parseJson<unknown>(response);
	const result = safeParse(array(UserGroupSummarySchema), data);

	if (!result.success) {
		throw new ApiError('Invalid response from the server', response.status);
	}

	return result.output;
}

export async function searchGroups(query: string): Promise<Group[]> {
	const params = new URLSearchParams({ q: query });
	const response = await apiFetch(`/api/groups/search?${params.toString()}`);
	const data = await parseJson<{ items?: Group[] } | Group[] | null>(response);
	if (data == null) return [];
	return Array.isArray(data) ? data : (data.items ?? []);
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

export async function joinGroup(id: string): Promise<JoinRequest | Group> {
	const response = await apiFetch(`/api/groups/${id}/join`, { method: 'POST' });
	return parseJson(response);
}

export async function leaveGroup(id: string): Promise<void> {
	const response = await apiFetch(`/api/groups/${id}/leave`, { method: 'POST' });
	await parseJson(response);
}

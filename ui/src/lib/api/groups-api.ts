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
import {
	GroupMemberPaginationSchema,
	GroupMemberSchema,
	type GroupMemberModel,
	type GroupMemberPaginationModel,
	type ListGroupMembersParams,
	type UpdateGroupMemberRoleRequest
} from '$lib/models/group-member-model';
import {
	GroupJoinRequestPaginationSchema,
	type GroupJoinRequestPaginationModel,
	type ListGroupPendingJoinRequestsParams,
	type ListGroupRejectedJoinRequestsParams
} from '$lib/models/group-join-request-model';

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Query a list of groups (paginated)
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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Query a list of groups the authenticated user belongs to (paginated)
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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Search for groups (paginated)
export async function searchGroups(params: SearchGroupsParams): Promise<GroupPaginationModel> {
	const qs = buildQueryString(params);
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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Query a group by ID
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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Create a group
export async function createGroup(data: CreateGroupRequest): Promise<GroupModel> {
	const response = await apiFetch('/api/groups/', {
		method: 'POST',
		body: JSON.stringify(data)
	});
	return parseJson(response);
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Request to join a group
export async function requestGroupJoin(groupId: string): Promise<GroupModel> {
	const response = await apiFetch(`/api/groups/${groupId}/join`, {
		method: 'POST'
	});

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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Query group members (paginated)
export async function listGroupMembers(
	groupId: string,
	params?: ListGroupMembersParams
): Promise<GroupMemberPaginationModel> {
	const qs = params ? buildQueryString(params) : '';
	const response = await apiFetch(`/api/groups/${groupId}/members` + (qs ? `?${qs}` : ''));

	if (response.ok) {
		const data = await response.json();
		const result = safeParse(GroupMemberPaginationSchema, data);

		if (!result.success) {
			throw new ApiError('Invalid response from the server', response.status);
		}

		return result.output;
	}

	const data = (await response.json()) as { message?: string };
	throw new ApiError(data.message || 'Request failed', response.status);
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Update a group member role
export async function updateGroupMemberRole(
	groupId: string,
	userId: string,
	body: UpdateGroupMemberRoleRequest
): Promise<GroupMemberModel> {
	const response = await apiFetch(`/api/groups/${groupId}/members/${userId}`, {
		method: 'PATCH',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(body)
	});

	if (response.ok) {
		const data = await response.json();
		const result = safeParse(GroupMemberSchema, data);

		if (!result.success) {
			throw new ApiError('Invalid response from the server', response.status);
		}

		return result.output;
	}

	const data = (await response.json()) as { message?: string };
	throw new ApiError(data.message || 'Request failed', response.status);
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Query pending group join requests (paginated)
export async function listGroupPendingJoinRequests(
	groupId: string,
	params?: ListGroupPendingJoinRequestsParams
): Promise<GroupJoinRequestPaginationModel> {
	const qs = params ? buildQueryString(params) : '';
	const response = await apiFetch(`/api/groups/${groupId}/pending` + (qs ? `?${qs}` : ''));

	if (response.ok) {
		const data = await response.json();
		const result = safeParse(GroupJoinRequestPaginationSchema, data);

		if (!result.success) {
			throw new ApiError('Invalid response from the server', response.status);
		}

		return result.output;
	}

	const data = (await response.json()) as { message?: string };
	throw new ApiError(data.message || 'Request failed', response.status);
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Query rejected group join requests (paginated)
export async function listGroupRejectedJoinRequests(
	groupId: string,
	params?: ListGroupRejectedJoinRequestsParams
): Promise<GroupJoinRequestPaginationModel> {
	const qs = params ? buildQueryString(params) : '';
	const response = await apiFetch(`/api/groups/${groupId}/rejected` + (qs ? `?${qs}` : ''));

	if (response.ok) {
		const data = await response.json();
		const result = safeParse(GroupJoinRequestPaginationSchema, data);

		if (!result.success) {
			throw new ApiError('Invalid response from the server', response.status);
		}

		return result.output;
	}

	const data = (await response.json()) as { message?: string };
	throw new ApiError(data.message || 'Request failed', response.status);
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Approve a pending group join request
export async function approveGroupJoinRequest(groupId: string, userId: string): Promise<void> {
	const response = await apiFetch(`/api/groups/${groupId}/pending/${userId}/approve`, {
		method: 'POST'
	});

	if (response.ok || response.status === 204) {
		return;
	}

	const data = (await response.json()) as { message?: string };
	throw new ApiError(data.message || 'Request failed', response.status);
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Decline a pending group join request
export async function declineGroupJoinRequest(groupId: string, userId: string): Promise<void> {
	const response = await apiFetch(`/api/groups/${groupId}/pending/${userId}/decline`, {
		method: 'POST'
	});

	if (response.ok || response.status === 204) {
		return;
	}

	const data = (await response.json()) as { message?: string };
	throw new ApiError(data.message || 'Request failed', response.status);
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Delete a group
export async function deleteGroup(groupId: string): Promise<void> {
	const response = await apiFetch(`/api/groups/${groupId}`, { method: 'DELETE' });

	if (response.ok || response.status === 204) {
		return;
	}

	const data = (await response.json()) as { message?: string };
	throw new ApiError(data.message || 'Request failed', response.status);
}

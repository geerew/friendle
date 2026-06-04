import { ApiError, apiFetch } from './fetch';
import { buildQueryString } from '$lib/utils';
import { safeParse } from 'valibot';
import {
	AdminGroupPaginationSchema,
	type AdminGroupPaginationModel,
	type AdminGroupReqParams
} from '$lib/models/admin-group-model';
import {
	AdminUserPaginationSchema,
	type AdminUserCreateModel,
	type AdminUserPaginationModel,
	type AdminUserReqParams
} from '$lib/models/admin-user-model';

export async function listUsers(params?: AdminUserReqParams): Promise<AdminUserPaginationModel> {
	const qs = params ? buildQueryString(params) : '';
	const response = await apiFetch('/api/users' + (qs ? `?${qs}` : ''));

	if (response.ok) {
		const data = await response.json();
		const result = safeParse(AdminUserPaginationSchema, data);

		if (!result.success) {
			throw new ApiError('Invalid response from the server', response.status);
		}

		return result.output;
	}

	const data = (await response.json()) as { message?: string };
	throw new ApiError(data.message || 'Request failed', response.status);
}

export async function createUser(data: AdminUserCreateModel): Promise<void> {
	const response = await apiFetch('/api/users', {
		method: 'POST',
		body: JSON.stringify(data)
	});

	if (response.ok || response.status === 201) {
		return;
	}

	const body = (await response.json()) as { message?: string };
	throw new ApiError(body.message || 'Request failed', response.status);
}

export async function listGroups(params?: AdminGroupReqParams): Promise<AdminGroupPaginationModel> {
	const qs = params ? buildQueryString(params) : '';
	const response = await apiFetch('/api/groups' + (qs ? `?${qs}` : ''));

	if (response.ok) {
		const data = await response.json();
		const result = safeParse(AdminGroupPaginationSchema, data);

		if (!result.success) {
			throw new ApiError('Invalid response from the server', response.status);
		}

		return result.output;
	}

	const data = (await response.json()) as { message?: string };
	throw new ApiError(data.message || 'Request failed', response.status);
}

export async function deleteUser(userId: string): Promise<void> {
	const response = await apiFetch(`/api/users/${userId}`, { method: 'DELETE' });

	if (response.ok || response.status === 204) {
		return;
	}

	const data = (await response.json()) as { message?: string };
	throw new ApiError(data.message || 'Request failed', response.status);
}

export async function deleteAdminGroup(groupId: string): Promise<void> {
	const response = await apiFetch(`/api/groups/${groupId}`, { method: 'DELETE' });

	if (response.ok || response.status === 204) {
		return;
	}

	const data = (await response.json()) as { message?: string };
	throw new ApiError(data.message || 'Request failed', response.status);
}

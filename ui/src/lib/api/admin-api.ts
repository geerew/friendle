import { ApiError, apiFetch } from './fetch';
import { buildQueryString } from '$lib/utils';
import { safeParse } from 'valibot';
import {
	AdminUserPaginationSchema,
	type AdminUserCreateModel,
	type AdminUserPaginationModel,
	type AdminUserReqParams
} from '$lib/models/admin-user-model';

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Query a list of users (paginated)
export async function listUsers(params?: AdminUserReqParams): Promise<AdminUserPaginationModel> {
	const qs = params ? buildQueryString(params) : '';
	const response = await apiFetch('/api/admin/users/' + (qs ? `?${qs}` : ''));

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

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Create a user
export async function createUser(data: AdminUserCreateModel): Promise<void> {
	const response = await apiFetch('/api/admin/users/', {
		method: 'POST',
		body: JSON.stringify(data)
	});

	if (response.ok || response.status === 201) {
		return;
	}

	const body = (await response.json()) as { message?: string };
	throw new ApiError(body.message || 'Request failed', response.status);
}

// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

// Delete a user
export async function deleteUser(userId: string): Promise<void> {
	const response = await apiFetch(`/api/admin/users/${userId}/`, { method: 'DELETE' });

	if (response.ok || response.status === 204) {
		return;
	}

	const data = (await response.json()) as { message?: string };
	throw new ApiError(data.message || 'Request failed', response.status);
}

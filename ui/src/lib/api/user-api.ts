import { ApiError, apiFetch } from './fetch';
import type { AdminUserCreateModel } from '$lib/models/admin-user-model';

export async function createUser(data: AdminUserCreateModel): Promise<void> {
	const response = await apiFetch('/api/users/', {
		method: 'POST',
		body: JSON.stringify(data)
	});

	if (response.ok || response.status === 201) {
		return;
	}

	const body = (await response.json()) as { message?: string };
	throw new ApiError(body.message || 'Request failed', response.status);
}

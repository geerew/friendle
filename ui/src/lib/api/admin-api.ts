import { apiFetch, parseJson } from './fetch';
import type { PaginatedGroups, PaginatedUsers } from '$lib/types/admin';

export async function listUsers(): Promise<PaginatedUsers> {
	const response = await apiFetch('/api/admin/users');
	return parseJson(response);
}

export async function listGroups(): Promise<PaginatedGroups> {
	const response = await apiFetch('/api/admin/groups');
	return parseJson(response);
}

export async function deleteUser(userId: string): Promise<void> {
	const response = await apiFetch(`/api/admin/users/${userId}`, { method: 'DELETE' });
	await parseJson(response);
}

export async function deleteAdminGroup(groupId: string): Promise<void> {
	const response = await apiFetch(`/api/admin/groups/${groupId}`, { method: 'DELETE' });
	await parseJson(response);
}

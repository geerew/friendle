import { apiFetch, parseJson } from './fetch';
import type { CreateGroupRequest, GroupModel } from '$lib/models/group-model';

export async function createGroup(data: CreateGroupRequest): Promise<GroupModel> {
	const response = await apiFetch('/api/groups/', {
		method: 'POST',
		body: JSON.stringify(data)
	});
	return parseJson(response);
}

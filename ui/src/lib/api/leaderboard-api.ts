import { apiFetch, parseJson } from './fetch';
import type { Leaderboard } from '$lib/types/leaderboard';

export async function getLeaderboard(groupId: string): Promise<Leaderboard> {
	const response = await apiFetch(`/api/groups/${groupId}/leaderboard`);
	return parseJson(response);
}

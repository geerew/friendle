import { apiFetch, parseJson } from './fetch';
import type {
	Round,
	RevealResponse,
	SubmitGuessRequest,
	SubmitGuessResponse,
	SubmitWordRequest
} from '$lib/types/round';

export async function getCurrentRound(groupId: string): Promise<Round> {
	const response = await apiFetch(`/api/groups/${groupId}/rounds/current`);
	return parseJson(response);
}

export async function submitWord(groupId: string, data: SubmitWordRequest): Promise<void> {
	await apiFetch(`/api/groups/${groupId}/rounds/current/word`, {
		method: 'POST',
		body: JSON.stringify(data)
	});
}

export async function submitGuess(
	groupId: string,
	data: SubmitGuessRequest
): Promise<SubmitGuessResponse> {
	const response = await apiFetch(`/api/groups/${groupId}/rounds/current/guesses`, {
		method: 'POST',
		body: JSON.stringify(data)
	});
	return parseJson(response);
}

export async function revealWord(groupId: string): Promise<RevealResponse> {
	const response = await apiFetch(`/api/groups/${groupId}/rounds/current/reveal`);
	return parseJson(response);
}

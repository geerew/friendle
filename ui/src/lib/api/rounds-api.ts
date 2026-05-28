import { apiFetch, parseJson } from './fetch';
import type { Guess, RevealResponse, Round, SubmitGuessRequest, SubmitGuessResponse, SubmitWordRequest } from '$lib/types/round';

export async function getCurrentRound(groupId: string): Promise<Round | null> {
	const response = await apiFetch(`/api/groups/${groupId}/rounds/current`);
	if (response.status === 404) return null;
	return parseJson(response);
}

export async function submitWord(groupId: string, data: SubmitWordRequest): Promise<Round> {
	const response = await apiFetch(`/api/groups/${groupId}/rounds/current/word`, {
		method: 'POST',
		body: JSON.stringify(data)
	});
	return parseJson(response);
}

export async function getMyGuess(groupId: string): Promise<Guess | null> {
	const response = await apiFetch(`/api/groups/${groupId}/rounds/current/guesses`);
	if (response.status === 404) return null;
	return parseJson(response);
}

export async function submitGuess(groupId: string, data: SubmitGuessRequest): Promise<SubmitGuessResponse> {
	const response = await apiFetch(`/api/groups/${groupId}/rounds/current/guesses`, {
		method: 'POST',
		body: JSON.stringify(data)
	});
	return parseJson(response);
}

export async function revealWord(groupId: string): Promise<RevealResponse> {
	const response = await apiFetch(`/api/groups/${groupId}/rounds/current/reveal`, {
		method: 'POST'
	});
	return parseJson(response);
}

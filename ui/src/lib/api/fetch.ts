import { auth } from '$lib/auth.svelte';

export class ApiError extends Error {
	status: number;

	constructor(message: string, status = 400) {
		super(message);
		this.name = 'ApiError';
		this.status = status;
	}
}

export async function apiFetch(input: RequestInfo, init?: RequestInit): Promise<Response> {
	const response = await fetch(input, {
		credentials: 'include',
		...init,
		headers: {
			...(init?.body ? { 'Content-Type': 'application/json' } : {}),
			...init?.headers
		}
	});

	if (response.status === 401 || response.status === 403) {
		const isAuthRoute =
			typeof input === 'string' &&
			(input.includes('/api/auth/login') ||
				input.includes('/api/auth/register') ||
				input.includes('/api/auth/signup-status') ||
				input.includes('/api/auth/bootstrap/'));

		if (!isAuthRoute) {
			auth.clear();
			if (typeof window !== 'undefined') {
				window.location.href = '/auth/login/';
			}
		}
	}

	return response;
}

export async function parseJson<T>(response: Response): Promise<T> {
	if (response.ok) {
		if (response.status === 204 || response.status === 201 && response.headers.get('content-length') === '0') {
			return undefined as T;
		}

		const text = await response.text();
		if (!text) return undefined as T;

		return JSON.parse(text) as T;
	}

	let message = 'Request failed';

	try {
		const data = (await response.json()) as { message?: string };
		if (data.message) message = data.message;
	} catch {
		// ignore parse errors
	}

	throw new ApiError(message, response.status);
}

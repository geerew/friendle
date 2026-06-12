import { auth } from '$lib/auth.svelte';

/** Minimum time mutation requests take so button loaders do not flicker */
export const MIN_LOADING_MS = 500;

const mutationMethods = new Set(['POST', 'PUT', 'PATCH', 'DELETE']);

export type ApiFetchOptions = {
	/** Override the default minimum duration for this request; pass false to disable */
	minDurationMs?: number | false;
};

export class ApiError extends Error {
	status: number;

	constructor(message: string, status = 400) {
		super(message);
		this.name = 'ApiError';
		this.status = status;
	}
}

// isApiError reports whether err is an ApiError, including across duplicate module instances
export function isApiError(err: unknown): err is ApiError {
	return (
		err instanceof ApiError ||
		(typeof err === 'object' &&
			err !== null &&
			(err as ApiError).name === 'ApiError' &&
			typeof (err as ApiError).status === 'number' &&
			typeof (err as ApiError).message === 'string')
	);
}

// apiErrorFromResponse builds an ApiError from a non-OK fetch response body
export async function apiErrorFromResponse(response: Response): Promise<ApiError> {
	let message = 'Request failed';

	try {
		const text = await response.text();

		if (text) {
			const data = JSON.parse(text) as { message?: string };

			if (data.message) {
				message = data.message;
			}
		}
	} catch {
		// ignore parse errors
	}

	return new ApiError(message, response.status);
}

// delay resolves after at least ms milliseconds
function delay(ms: number): Promise<void> {
	return new Promise((resolve) => setTimeout(resolve, ms));
}

// minDurationMs returns the minimum wait for this request, if any
function minDurationMs(
	init: RequestInit | undefined,
	options?: ApiFetchOptions
): number | undefined {
	if (options?.minDurationMs === false) {
		return undefined;
	}

	if (options?.minDurationMs !== undefined) {
		return options.minDurationMs;
	}

	const method = (init?.method ?? 'GET').toUpperCase();
	if (!mutationMethods.has(method)) {
		return undefined;
	}

	return MIN_LOADING_MS;
}

// runFetch performs the HTTP request and applies session handling
async function runFetch(input: RequestInfo, init?: RequestInit): Promise<Response> {
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
				input.includes('/api/auth/me') ||
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

export async function apiFetch(
	input: RequestInfo,
	init?: RequestInit,
	options?: ApiFetchOptions
): Promise<Response> {
	const minMs = minDurationMs(init, options);

	if (!minMs) {
		return runFetch(input, init);
	}

	const [response] = await Promise.all([runFetch(input, init), delay(minMs)]);

	return response;
}

export async function parseJson<T>(response: Response): Promise<T> {
	if (response.ok) {
		if (
			response.status === 204 ||
			(response.status === 201 && response.headers.get('content-length') === '0')
		) {
			return undefined as T;
		}

		const text = await response.text();
		if (!text) return undefined as T;

		return JSON.parse(text) as T;
	}

	let message = 'Request failed';

	try {
		const text = await response.text();

		if (text) {
			const data = JSON.parse(text) as { message?: string };

			if (data.message) {
				message = data.message;
			}
		}
	} catch {
		// ignore parse errors
	}

	throw new ApiError(message, response.status);
}

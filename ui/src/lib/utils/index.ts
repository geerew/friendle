import { isApiError } from '$lib/api/fetch';
import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

// cn merges class values into a single string
export function cn(...inputs: ClassValue[]): string {
	return twMerge(clsx(inputs));
}

// MIN_LOADING_MS is the minimum delay when loading
export const MIN_LOADING_MS = 150;

// withMinLoadingDelay ensures button spinners do not flicker on fast responses
export async function withMinLoadingDelay<T>(
	task: Promise<T>,
	ms: number = MIN_LOADING_MS
): Promise<T> {
	const [result] = await Promise.all([
		task,
		new Promise<void>((resolve) => setTimeout(resolve, ms))
	]);

	return result;
}

// buildQueryString builds a query string from a record of parameters
export function buildQueryString(
	params: Record<string, string | number | boolean | undefined>
): string {
	const searchParams = new URLSearchParams();

	for (const [key, value] of Object.entries(params)) {
		if (value !== undefined) {
			searchParams.append(key, value.toString());
		}
	}

	return searchParams.toString();
}

// apiErrorMessage formats an API error message
export function apiErrorMessage(err: unknown, fallback: string): string {
	return isApiError(err) ? err.message : fallback;
}

// isJoinRequestNotFound reports whether approve/decline failed because the pending request
// is no longer valid. Those endpoints only return 404 in that scenario
export function isJoinRequestNotFound(err: unknown): boolean {
	return isApiError(err) && err.status === 404;
}

// isPasswordFieldError reports whether an error message is related to password fields
export function isPasswordFieldError(error: string | null): boolean {
	if (!error) return false;

	const message = error.toLowerCase();

	return (
		message.includes('passwords do not match') ||
		message.includes('password must be at least') ||
		message.includes('all password fields are required')
	);
}

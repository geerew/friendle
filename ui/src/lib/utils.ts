import { ApiError } from '$lib/api/fetch';
import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export const MIN_BUTTON_LOADING_MS = 150;

export function cn(...inputs: ClassValue[]): string {
	return twMerge(clsx(inputs));
}

// withMinLoadingDelay ensures button spinners do not flicker on fast responses
export async function withMinLoadingDelay<T>(
	task: Promise<T>,
	ms: number = MIN_BUTTON_LOADING_MS
): Promise<T> {
	const [result] = await Promise.all([
		task,
		new Promise<void>((resolve) => setTimeout(resolve, ms))
	]);

	return result;
}

export function buildQueryString(params: Record<string, string | number | boolean | undefined>): string {
	const searchParams = new URLSearchParams();

	for (const [key, value] of Object.entries(params)) {
		if (value !== undefined) {
			searchParams.append(key, value.toString());
		}
	}

	return searchParams.toString();
}

export function apiErrorMessage(err: unknown, fallback: string): string {
	return err instanceof ApiError ? err.message : fallback;
}

export function isPasswordFieldError(error: string | null): boolean {
	if (!error) return false;

	const message = error.toLowerCase();

	return (
		message.includes('passwords do not match') ||
		message.includes('password must be at least') ||
		message.includes('all password fields are required')
	);
}


import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

export function cn(...inputs: ClassValue[]): string {
	return twMerge(clsx(inputs));
}

export function buildQueryString(params: Record<string, string | number | undefined>): string {
	const searchParams = new URLSearchParams();

	for (const [key, value] of Object.entries(params)) {
		if (value !== undefined) {
			searchParams.append(key, value.toString());
		}
	}

	return searchParams.toString();
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

export function formatMemberCount(count: number): string {
	return count === 1 ? '1 member' : `${count} members`;
}

export function formatGroupCount(count: number): string {
	return count === 1 ? '1 group' : `${count} groups`;
}

type ClassValue = string | false | null | undefined;

export function cn(...inputs: ClassValue[]): string {
	return inputs.filter(Boolean).join(' ');
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

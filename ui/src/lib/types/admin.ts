import type { Group } from '$lib/types/group';
import type { User } from '$lib/types/auth';

export type AdminUser = User & {
	createdAt?: string;
};

export type AdminGroup = Group & {
	createdAt?: string;
	memberCount: number;
};

export type PaginatedUsers = {
	items: AdminUser[];
	total: number;
};

export type PaginatedGroups = {
	items: AdminGroup[];
	total: number;
};

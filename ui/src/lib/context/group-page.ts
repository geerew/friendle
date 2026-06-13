import type { GroupModel } from '$lib/models/group-model';

export const GROUP_PAGE_KEY = Symbol('group-page');

export type GroupPageContext = {
	group: GroupModel | null;
	loading: boolean;
	reloadGroup: (options?: { silent?: boolean }) => Promise<void>;
};

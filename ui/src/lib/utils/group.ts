import type { BreadcrumbItem } from '$lib/components/breadcrumb.svelte';
import type { GroupModel } from '$lib/models/group-model';

// groupHomeBreadcrumb builds breadcrumb items for the group detail page
export function groupHomeBreadcrumb(): BreadcrumbItem[] {
	return [{ label: 'Group' }];
}

// groupChildBreadcrumb builds breadcrumb items for a group sub-page
export function groupChildBreadcrumb(groupId: string, pageLabel: string): BreadcrumbItem[] {
	return [
		{ label: 'Group', href: `/groups/${groupId}/` },
		{ label: pageLabel }
	];
}

// isGroupMember reports whether the viewer belongs to the group
export function isGroupMember(group: Pick<GroupModel, 'groupRole'>): boolean {
	return group.groupRole === 'group_admin' || group.groupRole === 'group_user';
}

// isGroupAdmin reports whether the viewer is a group admin
export function isGroupAdmin(group: Pick<GroupModel, 'groupRole'>): boolean {
	return group.groupRole === 'group_admin';
}

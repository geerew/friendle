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

// groupSettingsChildBreadcrumb builds breadcrumb items for a group settings sub-page
export function groupSettingsChildBreadcrumb(groupId: string, pageLabel: string): BreadcrumbItem[] {
	return [
		{ label: 'Group', href: `/groups/${groupId}/` },
		{ label: 'Settings', href: `/groups/${groupId}/settings/` },
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

// joinRequestResolvedNotice returns toast content when a pending join request was already resolved
export function joinRequestResolvedNotice(
	group: Pick<GroupModel, 'groupRole' | 'joinRequestStatus'>
): { variant: 'success' | 'error'; message: string } {
	if (isGroupMember(group)) {
		return { variant: 'success', message: 'You were added to the group' };
	}

	if (group.joinRequestStatus === 'rejected') {
		return { variant: 'error', message: 'Join request was declined' };
	}

	return { variant: 'error', message: 'Join request is no longer pending' };
}

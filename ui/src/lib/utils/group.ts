import type { BreadcrumbItem } from '$lib/components/breadcrumb.svelte';
import type { GroupModel } from '$lib/models/group-model';

const groupBreadcrumbNameMaxLength = 16;

// truncateGroupBreadcrumbName shortens a group name for breadcrumb display
export function truncateGroupBreadcrumbName(
	name: string,
	maxLength: number = groupBreadcrumbNameMaxLength
): string {
	if (name.length <= maxLength) {
		return name;
	}

	return `${name.slice(0, maxLength)}…`;
}

// groupHomeBreadcrumb builds breadcrumb items for the group detail page
export function groupHomeBreadcrumb(groupName: string): BreadcrumbItem[] {
	const truncatedName = truncateGroupBreadcrumbName(groupName);

	return [
		{
			label: 'Group',
			accentLabel: truncatedName,
			title: truncatedName !== groupName ? groupName : undefined
		}
	];
}

// groupChildBreadcrumb builds breadcrumb items for a group sub-page
export function groupChildBreadcrumb(
	groupId: string,
	groupName: string,
	pageLabel: string
): BreadcrumbItem[] {
	const truncatedName = truncateGroupBreadcrumbName(groupName);

	return [
		{
			label: 'Group',
			accentLabel: truncatedName,
			href: `/groups/${groupId}/`,
			title: truncatedName !== groupName ? groupName : undefined
		},
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

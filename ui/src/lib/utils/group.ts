import type { GroupModel } from '$lib/models/group-model';

// isGroupMember reports whether the viewer belongs to the group
export function isGroupMember(group: Pick<GroupModel, 'groupRole'>): boolean {
	return group.groupRole === 'group_admin' || group.groupRole === 'group_user';
}

// isGroupAdmin reports whether the viewer is a group admin
export function isGroupAdmin(group: Pick<GroupModel, 'groupRole'>): boolean {
	return group.groupRole === 'group_admin';
}

export type GroupRole = 'group_admin' | 'group_user';

export type RoundStatus = 'awaiting_word' | 'active' | 'completed' | 'skipped';

export type Group = {
	id: string;
	name: string;
	createdBy?: string;
	memberCount?: number;
	myRole?: GroupRole;
};

export type GroupDetail = Group & {
	groupRole?: GroupRole;
	members?: GroupMember[];
	joinRequests?: JoinRequest[];
};

export type GroupMember = {
	id: string;
	userId: string;
	displayName: string;
	groupRole: GroupRole;
	timesPicked: number;
};

export type CreateGroupRequest = {
	name: string;
};

export type UpdateGroupRequest = {
	name?: string;
};

export type JoinRequestStatus = 'pending' | 'approved' | 'rejected';

export type JoinRequest = {
	id: string;
	status: JoinRequestStatus;
};

export type RoundSummary = {
	id: string;
	roundDate: string;
	status: RoundStatus;
	pickerUserId?: string;
	word?: string;
};

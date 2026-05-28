export type GroupRole = 'group_admin' | 'group_user';

export type RoundStatus = 'awaiting_word' | 'active' | 'completed' | 'skipped';

export type Group = {
	id: string;
	name: string;
	createdBy: string;
	intervalHours: number;
	timezone: string;
	memberCount?: number;
	myRole?: GroupRole;
};

export type GroupDetail = Group & {
	members?: GroupMember[];
	currentRound?: RoundSummary | null;
};

export type GroupMember = {
	id: string;
	userId: string;
	username: string;
	displayName: string;
	groupRole: GroupRole;
	timesPicked: number;
	pickerSkips: number;
};

export type RoundSummary = {
	id: string;
	groupId: string;
	roundDate: string;
	pickerUserId: string;
	pickerUsername?: string;
	status: RoundStatus;
	isPicker?: boolean;
};

export type CreateGroupRequest = {
	name: string;
	intervalHours?: number;
	timezone?: string;
};

export type UpdateGroupRequest = {
	name?: string;
	intervalHours?: number;
	timezone?: string;
};

export type JoinRequestStatus = 'pending' | 'approved' | 'rejected';

export type JoinRequest = {
	id: string;
	groupId: string;
	userId: string;
	username?: string;
	status: JoinRequestStatus;
};

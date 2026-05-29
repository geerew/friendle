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
	groupRole?: GroupRole;
	members?: GroupMember[];
	joinRequests?: JoinRequest[];
	leaderboard?: LeaderboardEntry[];
	previousRounds?: RoundSummary[];
	round?: {
		status: RoundStatus | 'none';
		yourRole?: string;
		canReveal?: boolean;
		groupRole?: GroupRole;
	};
};

export type LeaderboardEntry = {
	userId: string;
	displayName: string;
	totalScore: number;
	roundsPlayed: number;
};

export type RoundSummary = {
	id: string;
	roundDate: string;
	status: RoundStatus;
	pickerUserId?: string;
	word?: string;
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

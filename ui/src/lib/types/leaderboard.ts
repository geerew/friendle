export type LeaderboardEntry = {
	userId: string;
	username: string;
	displayName: string;
	totalScore: number;
	roundsPlayed: number;
	roundsWon: number;
	averageAttempts: number;
	rank: number;
};

export type Leaderboard = {
	groupId: string;
	entries: LeaderboardEntry[];
};

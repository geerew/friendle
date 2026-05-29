import type { GroupRole, JoinRequest, RoundStatus } from './group';

export type TileState = 'empty' | 'correct' | 'present' | 'absent' | 'tbd';

export type GuessOutcome = 'correct' | 'partial' | 'incorrect';

export type GuessRow = {
	letters: string;
	states: TileState[];
};

export type ApiGuessRow = {
	word: string;
	result: TileState[];
	outcome: GuessOutcome;
};

export type RoundParticipation = {
	id: string;
	userId: string;
	displayName?: string;
	solved: boolean;
	finished: boolean;
	score: number;
	firstGuessAt?: string;
	completedAt?: string;
	guesses?: ApiGuessRow[];
};

export type Round = {
	id?: string;
	groupId?: string;
	roundDate?: string;
	status: RoundStatus | 'none';
	pickerUserId?: string;
	word?: string;
	yourRole?: string;
	participations?: RoundParticipation[];
};

export type Guess = {
	id: string;
	roundId: string;
	userId: string;
	attempt: number;
	word: string;
	outcome: GuessOutcome;
	result: TileState[];
};

export type SubmitGuessRequest = {
	word: string;
};

export type SubmitGuessResponse = {
	result: TileState[];
	attempt: number;
	won: boolean;
	finished: boolean;
};

export type SubmitWordRequest = {
	word: string;
};

export type RevealResponse = {
	status: RoundStatus;
	pickerUserId: string;
	pickerDisplayName?: string;
	word?: string;
	guesses: Array<{
		userId: string;
		displayName: string;
		attemptsUsed: number;
		solved: boolean;
		score: number;
	}>;
};

export type RoundSummary = {
	id: string;
	roundDate: string;
	status: RoundStatus;
	pickerUserId?: string;
	word?: string;
};

export function mapApiRows(apiRows: ApiGuessRow[] | undefined): GuessRow[] {
	if (!apiRows) {
		return [];
	}

	return apiRows.map((row) => ({
		letters: row.word,
		states: row.result
	}));
}

export function participationFromRound(round: Round | null, userId: string): RoundParticipation | null {
	const participation = round?.participations?.find((p) => p.userId === userId);
	if (!participation) {
		return null;
	}

	return {
		...participation,
		guesses: participation.guesses
	};
}

export function participationToBoard(participation: RoundParticipation | null): {
	rows: GuessRow[];
	finished: boolean;
	solved: boolean;
} {
	return {
		rows: mapApiRows(participation?.guesses),
		finished: participation?.finished ?? false,
		solved: participation?.solved ?? false
	};
}

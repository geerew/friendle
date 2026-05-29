import type { RoundStatus } from './group';

export type TileState = 'empty' | 'correct' | 'present' | 'absent' | 'tbd';

export type GuessRow = {
	letters: string;
	states: TileState[];
};

export type ApiGuessRow = {
	word: string;
	result: TileState[];
};

export type CurrentRound = {
	roundId?: string;
	status: RoundStatus | 'none';
	yourRole?: string;
	attemptsUsed?: number;
	finished?: boolean;
	solved?: boolean;
	rows?: ApiGuessRow[];
};

export type RoundParticipation = {
	attemptsUsed: number;
	solved: boolean;
	finished: boolean;
	score?: number;
	rows: GuessRow[];
};

export type Guess = {
	id: string;
	roundId: string;
	userId: string;
	attemptNumber: number;
	word: string;
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

export type LeaderboardEntry = {
	userId: string;
	displayName: string;
	totalScore: number;
	roundsPlayed: number;
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

export function currentRoundToParticipation(round: CurrentRound | null): RoundParticipation | null {
	if (!round || round.status === 'none') {
		return null;
	}

	return {
		attemptsUsed: round.attemptsUsed ?? 0,
		solved: round.solved ?? false,
		finished: round.finished ?? false,
		rows: mapApiRows(round.rows)
	};
}

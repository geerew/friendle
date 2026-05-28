import type { RoundStatus } from './group';

export type TileState = 'empty' | 'correct' | 'present' | 'absent' | 'tbd';

export type GuessRow = {
	letters: string;
	states: TileState[];
};

export type Round = {
	id: string;
	groupId: string;
	roundDate: string;
	pickerUserId: string;
	pickerUsername?: string;
	status: RoundStatus;
	isPicker?: boolean;
};

export type Guess = {
	id: string;
	roundId: string;
	userId: string;
	attemptsUsed: number;
	solved: boolean;
	rows: GuessRow[];
	score: number;
	finished: boolean;
};

export type SubmitGuessRequest = {
	word: string;
};

export type SubmitGuessResponse = {
	guess: Guess;
	roundFinished?: boolean;
};

export type SubmitWordRequest = {
	word: string;
};

export type RevealResponse = {
	word: string;
};

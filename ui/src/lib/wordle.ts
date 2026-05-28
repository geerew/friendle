import type { GuessRow, TileState } from '$lib/types/round';

const rank: Record<TileState, number> = {
	empty: 0,
	tbd: 0,
	absent: 1,
	present: 2,
	correct: 3
};

export function buildLetterStates(rows: GuessRow[]): Record<string, TileState> {
	const states: Record<string, TileState> = {};

	for (const row of rows) {
		for (let i = 0; i < row.letters.length; i++) {
			const letter = row.letters[i]?.toUpperCase();
			const state = row.states[i] ?? 'empty';

			if (!letter || state === 'empty' || state === 'tbd') continue;

			const current = states[letter] ?? 'empty';
			if (rank[state] > rank[current]) {
				states[letter] = state;
			}
		}
	}

	return states;
}

export function normalizeWord(value: string): string {
	return value.toUpperCase().replace(/[^A-Z]/g, '').slice(0, 5);
}

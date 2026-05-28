<script lang="ts">
	import type { GuessRow, TileState } from '$lib/types/round';

	type Props = {
		rows?: GuessRow[];
		currentWord?: string;
		currentRow?: number;
		shake?: boolean;
		revealed?: boolean;
	};

	let {
		rows = [],
		currentWord = '',
		currentRow = 0,
		shake = false,
		revealed = false
	}: Props = $props();

	const maxRows = 6;
	const wordLength = 5;

	function tileClass(state: TileState, filled: boolean): string {
		const base =
			'flex h-14 w-14 items-center justify-center border-2 text-2xl font-bold uppercase sm:h-16 sm:w-16';

		if (state === 'correct') return `${base} border-transparent bg-tile-correct text-white`;
		if (state === 'present') return `${base} border-transparent bg-tile-present text-white`;
		if (state === 'absent') return `${base} border-transparent bg-tile-absent text-white`;
		if (state === 'tbd') return `${base} border-border bg-tile-empty text-white`;
		if (filled) return `${base} border-text-muted bg-tile-empty text-white`;
		return `${base} border-border bg-transparent text-white`;
	}

	function displayRows(): Array<{ letters: string; states: TileState[] }> {
		const result: Array<{ letters: string; states: TileState[] }> = [];

		for (let i = 0; i < maxRows; i++) {
			if (i < rows.length) {
				result.push(rows[i]);
				continue;
			}

			if (i === currentRow && !revealed) {
				const letters = currentWord.padEnd(wordLength, ' ').slice(0, wordLength);
				const states = Array.from({ length: wordLength }, (_, idx) =>
					letters[idx]?.trim() ? ('tbd' as TileState) : ('empty' as TileState)
				);
				result.push({ letters: letters.replace(/ /g, ''), states });
				continue;
			}

			result.push({
				letters: '',
				states: Array.from({ length: wordLength }, () => 'empty' as TileState)
			});
		}

		return result;
	}
</script>

<div class="mx-auto flex w-full flex-col items-center gap-1.5 {shake ? 'animate-shake' : ''}">
	{#each displayRows() as row, rowIndex (rowIndex)}
		<div class="flex gap-1.5">
			{#each Array.from({ length: wordLength }) as _, colIndex (colIndex)}
				{@const letter = row.letters[colIndex] ?? ''}
				{@const state = row.states[colIndex] ?? 'empty'}
				<div class={tileClass(state, Boolean(letter))}>
					{letter}
				</div>
			{/each}
		</div>
	{/each}
</div>

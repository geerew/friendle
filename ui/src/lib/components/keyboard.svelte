<script lang="ts">
	import type { TileState } from '$lib/types/round';

	type Props = {
		letterStates?: Record<string, TileState>;
		disabled?: boolean;
		onKey?: (key: string) => void;
		onEnter?: () => void;
		onBackspace?: () => void;
	};

	let {
		letterStates = {},
		disabled = false,
		onKey,
		onEnter,
		onBackspace
	}: Props = $props();

	const rows = ['QWERTYUIOP', 'ASDFGHJKL', 'ZXCVBNM'];

	function keyClass(letter: string): string {
		const state = letterStates[letter] ?? 'empty';
		const base =
			'flex h-14 min-w-0 flex-1 items-center justify-center rounded text-sm font-bold uppercase sm:h-16 sm:text-base';

		if (state === 'correct') return `${base} bg-key-correct text-white`;
		if (state === 'present') return `${base} bg-key-present text-white`;
		if (state === 'absent') return `${base} bg-key-absent text-white`;
		return `${base} bg-key-bg text-white`;
	}

	function handleClick(key: string): void {
		if (disabled) return;
		onKey?.(key);
	}

	function handleEnter(): void {
		if (disabled) return;
		onEnter?.();
	}

	function handleBackspace(): void {
		if (disabled) return;
		onBackspace?.();
	}
</script>

<div class="mx-auto flex w-full max-w-md flex-col gap-2 px-1 pb-2">
	{#each rows as row, rowIndex (rowIndex)}
		<div class="flex justify-center gap-1.5">
			{#if rowIndex === 2}
				<button
					type="button"
					class="flex h-14 flex-[1.5] items-center justify-center rounded bg-key-bg px-2 text-xs font-bold uppercase sm:h-16"
					onclick={handleEnter}
					{disabled}
				>
					Enter
				</button>
			{/if}

			{#each row.split('') as letter (letter)}
				<button type="button" class={keyClass(letter)} onclick={() => handleClick(letter)} {disabled}>
					{letter}
				</button>
			{/each}

			{#if rowIndex === 2}
				<button
					type="button"
					class="flex h-14 flex-[1.5] items-center justify-center rounded bg-key-bg px-2 text-xs font-bold uppercase sm:h-16"
					onclick={handleBackspace}
					{disabled}
					aria-label="Backspace"
				>
					⌫
				</button>
			{/if}
		</div>
	{/each}
</div>

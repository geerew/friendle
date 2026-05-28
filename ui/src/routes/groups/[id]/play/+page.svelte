<script lang="ts">
	import { page } from '$app/state';
	import { ApiError } from '$lib/api';
	import { getGroup } from '$lib/api/groups-api';
	import { getCurrentRound, getMyGuess, revealWord, submitGuess } from '$lib/api/rounds-api';
	import { AppShell, Keyboard, TileGrid } from '$lib/components';
	import type { Guess, GuessRow } from '$lib/types/round';
	import { buildLetterStates, normalizeWord } from '$lib/wordle';

	const groupId = $derived(page.params.id ?? '');

	let rows = $state<GuessRow[]>([]);
	let currentWord = $state('');
	let currentRow = $state(0);
	let finished = $state(false);
	let solved = $state(false);
	let loading = $state(true);
	let submitting = $state(false);
	let shake = $state(false);
	let error = $state<string | null>(null);
	let reveal = $state<string | null>(null);
	let roundActive = $state(false);
	let groupName = $state('Group');

	const letterStates = $derived(buildLetterStates(rows));

	$effect(() => {
		loadGame();
	});

	async function loadGame(): Promise<void> {
		loading = true;
		error = null;

		try {
			const [round, group] = await Promise.all([getCurrentRound(groupId), getGroup(groupId)]);
			groupName = group.name;
			roundActive = round?.status === 'active';

			if (!roundActive) {
				error = round ? `Round is ${round.status.replace('_', ' ')}.` : 'No active round.';
				return;
			}

			const guess = await getMyGuess(groupId);
			applyGuess(guess);
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to load game';
		} finally {
			loading = false;
		}
	}

	function applyGuess(guess: Guess | null): void {
		const rawRows = guess?.rows ?? (guess as { rowsJson?: GuessRow[] | string } | null)?.rowsJson;
		if (typeof rawRows === 'string') {
			try {
				rows = JSON.parse(rawRows) as GuessRow[];
			} catch {
				rows = [];
			}
		} else {
			rows = rawRows ?? [];
		}

		currentRow = rows.length;
		finished = guess?.finished ?? false;
		solved = guess?.solved ?? false;
		currentWord = '';
	}

	function handleKey(key: string): void {
		if (finished || submitting) return;
		currentWord = normalizeWord(currentWord + key);
	}

	function handleBackspace(): void {
		if (finished || submitting) return;
		currentWord = currentWord.slice(0, -1);
	}

	async function triggerShake(): Promise<void> {
		shake = true;
		await new Promise((resolve) => setTimeout(resolve, 500));
		shake = false;
	}

	async function handleEnter(): Promise<void> {
		if (finished || submitting) return;

		if (currentWord.length !== 5) {
			await triggerShake();
			error = 'Not enough letters';
			return;
		}

		submitting = true;
		error = null;

		try {
			const response = await submitGuess(groupId, { word: currentWord });
			applyGuess(response.guess);

			if (response.guess.finished && !response.guess.solved) {
				try {
					const revealed = await revealWord(groupId);
					reveal = revealed.word;
				} catch {
					// reveal may be restricted until round ends
				}
			}
		} catch (err) {
			await triggerShake();
			error = err instanceof ApiError ? err.message : 'Guess rejected';
		} finally {
			submitting = false;
		}
	}

	function handlePhysicalKeyboard(event: KeyboardEvent): void {
		if (loading || !roundActive || finished) return;

		if (event.key === 'Enter') {
			event.preventDefault();
			void handleEnter();
			return;
		}

		if (event.key === 'Backspace') {
			event.preventDefault();
			handleBackspace();
			return;
		}

		if (/^[a-zA-Z]$/.test(event.key)) {
			event.preventDefault();
			handleKey(event.key.toUpperCase());
		}
	}
</script>

<svelte:window onkeydown={handlePhysicalKeyboard} />

<AppShell
	breadcrumb={[
		{ label: groupName, href: `/groups/${groupId}/` },
		{ label: 'Play' }
	]}
>
	{#if loading}
		<p class="text-text-muted">Loading…</p>
	{:else if !roundActive}
		<p class="text-text-muted">{error ?? 'No active round.'}</p>
	{:else}
		<div class="flex flex-1 flex-col gap-4">
			<TileGrid {rows} {currentWord} {currentRow} {shake} revealed={finished} />

			{#if solved}
				<p class="text-center text-tile-correct">You got it!</p>
			{:else if finished}
				<p class="text-center text-text-muted">
					Out of guesses.{#if reveal} The word was {reveal}.{/if}
				</p>
			{/if}

			{#if error}
				<p class="text-center text-sm text-error">{error}</p>
			{/if}

			<div class="mt-auto">
				<Keyboard
					{letterStates}
					disabled={finished || submitting}
					onKey={handleKey}
					onEnter={handleEnter}
					onBackspace={handleBackspace}
				/>
			</div>
		</div>
	{/if}
</AppShell>

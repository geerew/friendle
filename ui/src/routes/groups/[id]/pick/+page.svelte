<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { ApiError } from '$lib/api';
	import { getCurrentRound, submitWord } from '$lib/api/rounds-api';
	import { normalizeWord } from '$lib/wordle';
	import { AppShell } from '$lib/components';
	import { Button, Field, Input } from '$lib/components/ui';

	const groupId = $derived(page.params.id ?? '');

	let word = $state('');
	let loading = $state(true);
	let submitting = $state(false);
	let error = $state<string | null>(null);
	let canPick = $state(false);

	$effect(() => {
		loadRound();
	});

	async function loadRound(): Promise<void> {
		loading = true;
		error = null;

		try {
			const round = await getCurrentRound(groupId);
			canPick = round?.status === 'awaiting_word' && Boolean(round.isPicker);
			if (round && !canPick) {
				error = 'You are not the picker for this round.';
			}
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to load round';
		} finally {
			loading = false;
		}
	}

	function handleInput(event: Event): void {
		const target = event.currentTarget as HTMLInputElement;
		word = normalizeWord(target.value);
		target.value = word;
	}

	async function handleSubmit(event: SubmitEvent): Promise<void> {
		event.preventDefault();

		if (word.length !== 5) {
			error = 'Word must be exactly 5 letters';
			return;
		}

		submitting = true;
		error = null;

		try {
			await submitWord(groupId, { word });
			await goto(`/groups/${groupId}/`);
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to submit word';
		} finally {
			submitting = false;
		}
	}
</script>

<AppShell title="Pick Word" backHref="/groups/{groupId}/">
	{#if loading}
		<p class="text-text-muted">Loading…</p>
	{:else if !canPick}
		<p class="text-text-muted">{error ?? 'Waiting for your turn to pick.'}</p>
	{:else}
		<p class="text-sm text-text-muted">
			Choose a secret 5-letter word for today's round. Other members will try to guess it.
		</p>

		<form class="flex flex-col gap-4" onsubmit={handleSubmit}>
			<Field label="Secret word">
				<Input
					class="text-center text-2xl tracking-[0.5em] uppercase"
					value={word}
					oninput={handleInput}
					maxlength={5}
					autocomplete="off"
					autocapitalize="characters"
					required
				/>
			</Field>

			{#if error}
				<p class="text-sm text-error">{error}</p>
			{/if}

			<Button type="submit" variant="primary" disabled={submitting || word.length !== 5}>
				{submitting ? 'Submitting…' : 'Submit word'}
			</Button>
		</form>
	{/if}
</AppShell>

<script lang="ts">
	import { page } from '$app/state';
	import { ApiError } from '$lib/api';
	import { getLeaderboard } from '$lib/api/leaderboard-api';
	import { AppShell, ListRow } from '$lib/components';
	import type { LeaderboardEntry } from '$lib/types/leaderboard';

	const groupId = $derived(page.params.id ?? '');

	let entries = $state<LeaderboardEntry[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	$effect(() => {
		loadLeaderboard();
	});

	async function loadLeaderboard(): Promise<void> {
		loading = true;
		error = null;

		try {
			const data = await getLeaderboard(groupId);
			entries = data.entries ?? [];
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to load leaderboard';
		} finally {
			loading = false;
		}
	}
</script>

<AppShell title="Leaderboard" showBack={true} backHref="/groups/{groupId}/">
	{#if loading}
		<p class="text-text-muted">Loading…</p>
	{:else if error}
		<p class="text-error">{error}</p>
	{:else if entries.length === 0}
		<p class="text-text-muted">No scores yet.</p>
	{:else}
		<div class="flex flex-col gap-2">
			{#each entries as entry (entry.userId)}
				<ListRow
					title="#{entry.rank} {entry.displayName}"
					subtitle="{entry.totalScore} pts · {entry.roundsWon}/{entry.roundsPlayed} wins · avg {entry.averageAttempts.toFixed(1)}"
				/>
			{/each}
		</div>
	{/if}
</AppShell>

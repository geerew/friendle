<script lang="ts">
	import { ApiError } from '$lib/api';
	import { joinGroup, searchGroups } from '$lib/api/groups-api';
	import { AppShell, ListRow, PlusIcon } from '$lib/components';
	import { Button, Input } from '$lib/components/ui';
	import type { Group } from '$lib/types/group';

	let searchQuery = $state('');
	let searchResults = $state<Group[]>([]);
	let searching = $state(false);
	let message = $state<string | null>(null);
	let error = $state<string | null>(null);

	async function handleSearch(): Promise<void> {
		if (!searchQuery.trim()) return;

		searching = true;
		message = null;
		error = null;

		try {
			searchResults = await searchGroups(searchQuery.trim());
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Search failed';
		} finally {
			searching = false;
		}
	}

	async function handleJoin(group: Group): Promise<void> {
		message = null;

		try {
			await joinGroup(group.id);
			message = `Requested to join ${group.name}`;
		} catch (err) {
			message = err instanceof ApiError ? err.message : 'Failed to join group';
		}
	}
</script>

<AppShell title="Search Groups" backHref="/">
	<div class="flex flex-col gap-2">
		<Input
			placeholder="Search by name"
			bind:value={searchQuery}
			onkeydown={(event) => event.key === 'Enter' && handleSearch()}
		/>
		<Button variant="secondary" onclick={handleSearch} disabled={searching}>
			{searching ? 'Searching…' : 'Search'}
		</Button>
	</div>

	{#if message}
		<p class="text-sm text-text-muted">{message}</p>
	{/if}

	{#if error}
		<p class="text-sm text-error">{error}</p>
	{/if}

	<div class="flex flex-col gap-2">
		{#each searchResults as group (group.id)}
			<ListRow title={group.name} subtitle="{group.memberCount ?? 0} members">
				{#snippet trailing()}
					<button
						type="button"
						class="inline-flex h-9 w-9 items-center justify-center rounded bg-button-primary text-white"
						aria-label="Join {group.name}"
						onclick={() => handleJoin(group)}
					>
						<PlusIcon class="size-5 stroke-2" />
					</button>
				{/snippet}
			</ListRow>
		{/each}
	</div>
</AppShell>

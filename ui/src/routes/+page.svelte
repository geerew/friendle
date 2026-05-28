<script lang="ts">
	import { ApiError } from '$lib/api';
	import { joinGroup, listMyGroups, searchGroups } from '$lib/api/groups-api';
	import { auth } from '$lib/auth.svelte';
	import { AppShell, ListRow } from '$lib/components';
	import { Button, Input } from '$lib/components/ui';
	import type { Group } from '$lib/types/group';

	let groups = $state<Group[]>([]);
	let searchQuery = $state('');
	let searchResults = $state<Group[]>([]);
	let loadingGroups = $state(true);
	let searching = $state(false);
	let error = $state<string | null>(null);
	let joinMessage = $state<string | null>(null);

	$effect(() => {
		loadGroups();
	});

	async function loadGroups(): Promise<void> {
		loadingGroups = true;
		error = null;

		try {
			groups = await listMyGroups();
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to load groups';
		} finally {
			loadingGroups = false;
		}
	}

	async function handleSearch(): Promise<void> {
		if (!searchQuery.trim()) {
			searchResults = [];
			return;
		}

		searching = true;
		joinMessage = null;

		try {
			searchResults = await searchGroups(searchQuery.trim());
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Search failed';
		} finally {
			searching = false;
		}
	}

	async function handleJoin(group: Group): Promise<void> {
		joinMessage = null;

		try {
			await joinGroup(group.id);
			joinMessage = `Requested to join ${group.name}`;
			await loadGroups();
		} catch (err) {
			joinMessage = err instanceof ApiError ? err.message : 'Failed to join group';
		}
	}
</script>

<AppShell title="Friendle" showBack={false} showSettings={true} settingsHref="/settings/">
	<div class="text-center">
		<p class="text-sm text-text-muted">Welcome, {auth.user?.displayName ?? auth.user?.username}</p>
	</div>

	<section class="flex flex-col gap-3">
		<h2 class="section-title">My Groups</h2>

		{#if loadingGroups}
			<p class="text-text-muted">Loading groups…</p>
		{:else if groups.length === 0}
			<p class="rounded border border-border bg-bg-secondary px-4 py-3 text-sm text-text-muted">
				You are not in any groups yet.
			</p>
		{:else}
			<div class="flex flex-col gap-2">
				{#each groups as group (group.id)}
					<ListRow href="/groups/{group.id}/" title={group.name} subtitle="{group.memberCount ?? 0} members" />
				{/each}
			</div>
		{/if}
	</section>

	<section class="flex flex-col gap-3">
		<h2 class="section-title">Search Groups</h2>
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

		{#if joinMessage}
			<p class="text-sm text-text-muted">{joinMessage}</p>
		{/if}

		{#if searchResults.length > 0}
			<div class="flex flex-col gap-2">
				{#each searchResults as group (group.id)}
					<ListRow title={group.name} subtitle="{group.memberCount ?? 0} members">
						{#snippet trailing()}
							<button
								type="button"
								class="inline-flex h-9 w-9 items-center justify-center rounded bg-button-primary text-xl font-bold text-white"
								aria-label="Join {group.name}"
								onclick={() => handleJoin(group)}
							>
								+
							</button>
						{/snippet}
					</ListRow>
				{/each}
			</div>
		{/if}
	</section>

	{#if error}
		<p class="text-sm text-error">{error}</p>
	{/if}

	<div class="mt-auto flex flex-col gap-3 pt-4">
		<Button href="/settings/" variant="secondary">Settings</Button>
		<Button href="/groups/new/" variant="primary">Add Group</Button>
		{#if auth.isAdmin}
			<Button href="/admin/" variant="ghost">Admin</Button>
		{/if}
	</div>
</AppShell>

<script lang="ts">
	import { onMount } from 'svelte';
	import { ApiError } from '$lib/api';
	import { listMyGroups } from '$lib/api/groups-api';
	import { auth } from '$lib/auth.svelte';
	import { AppShell, ListRow } from '$lib/components';
	import { Button } from '$lib/components/ui';
	import type { Group } from '$lib/types/group';
	import { formatMemberCount } from '$lib/utils';

	let groups = $state<Group[]>([]);
	let loadingGroups = $state(true);
	let error = $state<string | null>(null);

	onMount(() => {
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
</script>

<AppShell>
	<section class="flex flex-col gap-3">
		<h2 class="section-title">My Groups</h2>

		{#if loadingGroups}
			<p class="text-sm text-text-muted">Loading groups…</p>
		{:else if error}
			<p class="text-sm text-error">{error}</p>
		{:else if groups.length === 0}
			<p class="text-center text-xs text-text-muted">No groups</p>
		{:else}
			<div class="flex flex-col gap-2">
				{#each groups as group (group.id)}
					<ListRow
						href="/groups/{group.id}/"
						title={group.name}
						subtitle={formatMemberCount(group.memberCount ?? 0)}
					/>
				{/each}
			</div>
		{/if}
	</section>

	<div class="mt-auto flex flex-col gap-3 pt-4">
		<Button href="/groups/new/" variant="primary">Create Group</Button>
		<Button href="/groups/search/" variant="secondary">Search Groups</Button>
	</div>
</AppShell>

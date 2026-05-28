<script lang="ts">
	import { ApiError } from '$lib/api';
	import { deleteAdminGroup, listGroups } from '$lib/api/admin-api';
	import { AppShell, ListRow, Pagination } from '$lib/components';
	import { Button } from '$lib/components/ui';
	import type { AdminGroupModel } from '$lib/models/admin-group-model';

	let groups = $state<AdminGroupModel[]>([]);
	let page = $state(1);
	let perPage = $state(25);
	let totalItems = $state(0);
	let loading = $state(true);
	let error = $state<string | null>(null);

	$effect(() => {
		page;
		perPage;
		void loadGroups();
	});

	async function loadGroups(): Promise<void> {
		loading = true;
		error = null;

		try {
			const data = await listGroups({ page, perPage });
			groups = data.items;
			totalItems = data.totalItems;
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to load groups';
		} finally {
			loading = false;
		}
	}

	async function handleDeleteGroup(group: AdminGroupModel): Promise<void> {
		if (!confirm(`Delete group ${group.name}?`)) return;

		try {
			await deleteAdminGroup(group.id);
			const remainingTotal = totalItems - 1;
			const totalPages = Math.max(1, Math.ceil(remainingTotal / perPage));

			if (page > totalPages) {
				page = totalPages;
			} else {
				await loadGroups();
			}
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to delete group';
		}
	}
</script>

<AppShell
	breadcrumb={[
		{ label: 'Admin', href: '/admin/' },
		{ label: 'Groups' }
	]}
>
	{#if loading}
		<p class="text-text-muted">Loading…</p>
	{:else}
		{#if error}
			<p class="text-sm text-error">{error}</p>
		{/if}

		<div class="flex flex-col gap-2">
			{#if groups.length === 0}
				<p class="text-sm text-text-muted">No groups.</p>
			{:else}
				{#each groups as group (group.id)}
					<ListRow title={group.name} subtitle="{group.memberCount} members">
						{#snippet trailing()}
							<Button variant="destructive" size="inline" onclick={() => handleDeleteGroup(group)}>
								Delete
							</Button>
						{/snippet}
					</ListRow>
				{/each}
			{/if}
		</div>

		<Pagination count={totalItems} bind:page bind:perPage onPageChange={() => {}} onPerPageChange={() => {}} />
	{/if}
</AppShell>

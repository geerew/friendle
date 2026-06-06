<script lang="ts">
	import { ApiError } from '$lib/api';
	import { listGroups } from '$lib/api/groups-api';
	import { AppShell, ListRow, Pagination, Spinner } from '$lib/components';
	import type { GroupModel } from '$lib/models/group-model';
	import { withMinLoadingDelay } from '$lib/utils';

	let groups = $state<GroupModel[]>([]);
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
			const data = await withMinLoadingDelay(listGroups({ page, perPage }));
			groups = data.items;
			totalItems = data.totalItems;
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to load groups';
		} finally {
			loading = false;
		}
	}
</script>

<AppShell
	breadcrumb={[
		{ label: 'Admin', href: '/admin/' },
		{ label: 'Groups' }
	]}
>
	<h2 class="section-title">Groups</h2>

	<div class="flex flex-col gap-6">
		{#if loading}
			<div class="flex min-h-24 items-center justify-center">
				<Spinner class="bg-foreground-alt-2 size-3" />
			</div>
		{:else}
			{#if error}
				<p class="text-foreground-error text-sm">{error}</p>
			{/if}

			<div class="flex flex-col gap-3">
				{#if groups.length === 0}
					<p class="text-foreground-alt-2 text-sm italic">No groups</p>
				{:else}
					{#each groups as group (group.id)}
						<ListRow
							title={group.name}
							subtitle="{group.memberCount} {group.memberCount === 1 ? 'member' : 'members'}"
						/>
					{/each}
				{/if}
			</div>

			<Pagination
				count={totalItems}
				bind:page
				bind:perPage
				selectTriggerClass="h-9 px-2 py-0"
				onPageChange={() => {}}
				onPerPageChange={() => {}}
			/>
		{/if}
	</div>
</AppShell>

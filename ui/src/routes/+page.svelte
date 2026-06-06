<script lang="ts">
	import { ApiError } from '$lib/api';
	import { listSelfGroups } from '$lib/api/groups-api';
	import { AppShell, Pagination, Spinner } from '$lib/components';
	import { GroupList } from '$lib/components/pages';
	import type { GroupModel } from '$lib/models/group-model';
	import { withMinLoadingDelay } from '$lib/utils';

	let groups = $state<GroupModel[]>([]);
	let page = $state(1);
	let perPage = $state(5);
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
			const data = await withMinLoadingDelay(listSelfGroups({ page, perPage }));
			groups = data.items;
			totalItems = data.totalItems;
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to load groups';
			groups = [];
			totalItems = 0;
		} finally {
			loading = false;
		}
	}
</script>

<AppShell>
	<h2 class="section-title">My Groups</h2>

	<div class="flex flex-col gap-6">
		{#if loading}
			<div class="flex min-h-24 items-center justify-center">
				<Spinner class="bg-foreground-alt-2 size-3" />
			</div>
		{:else}
			{#if error}
				<p class="text-foreground-error text-sm">{error}</p>
			{/if}

			{#if groups.length === 0}
				<p class="text-foreground-alt-2 text-sm italic">No groups</p>
			{:else}
				<div class="px-2">
					<GroupList {groups} />
				</div>
			{/if}

			<Pagination
				count={totalItems}
				bind:page
				bind:perPage
				showPerPageSelect={false}
				onPageChange={() => {}}
				onPerPageChange={() => {}}
			/>
		{/if}
	</div>
</AppShell>

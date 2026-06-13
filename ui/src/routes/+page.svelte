<script lang="ts">
	import { listSelfGroups } from '$lib/api/groups-api';
	import { GroupList } from '$lib/components/pages';
	import { Table } from '$lib/components/ui';
	import type { GroupModel } from '$lib/models/group-model';
	import { apiErrorMessage, withMinLoadingDelay } from '$lib/utils';
	import { toast } from 'svelte-sonner';

	let groups = $state<GroupModel[]>([]);
	let page = $state(1);
	let perPage = $state(5);
	let totalItems = $state(0);
	let loading = $state(true);

	$effect(() => {
		page;
		perPage;
		void loadGroups();
	});

	async function loadGroups(): Promise<void> {
		loading = true;

		try {
			const data = await withMinLoadingDelay(listSelfGroups({ page, perPage }));
			groups = data.items;
			totalItems = data.totalItems;
		} catch (err) {
			toast.error(apiErrorMessage(err, 'Failed to load groups'));
			groups = [];
			totalItems = 0;
		} finally {
			loading = false;
		}
	}
</script>

<Table.Root title="My Groups" breadcrumb={[]}>
	<Table.PaginatedBody
		itemCount={groups.length}
		{totalItems}
		bind:page
		bind:perPage
		{loading}
		emptyMessage="No groups"
	>
		{#snippet list()}
			<GroupList {groups} />
		{/snippet}
	</Table.PaginatedBody>
</Table.Root>

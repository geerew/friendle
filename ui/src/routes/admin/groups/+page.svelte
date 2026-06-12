<script lang="ts">
	import { listGroups } from '$lib/api/groups-api';
	import { DeleteGroup } from '$lib/components';
	import { Badge, Button, Separator, Table } from '$lib/components/ui';
	import type { GroupModel } from '$lib/models/group-model';
	import { apiErrorMessage, withMinLoadingDelay } from '$lib/utils';
	import { toast } from 'svelte-sonner';

	let groups = $state<GroupModel[]>([]);
	let page = $state(1);
	let perPage = $state(10);
	let totalItems = $state(0);
	let loading = $state(true);
	let deleteOpen = $state(false);
	let groupToDelete = $state<GroupModel | null>(null);

	$effect(() => {
		page;
		perPage;
		void loadGroups();
	});

	async function loadGroups(): Promise<void> {
		loading = true;

		try {
			const data = await withMinLoadingDelay(listGroups({ page, perPage }));
			groups = data.items;
			totalItems = data.totalItems;
		} catch (err) {
			groups = [];
			totalItems = 0;
			toast.error(apiErrorMessage(err, 'Failed to load groups'));
		} finally {
			loading = false;
		}
	}

	function openDeleteGroup(group: GroupModel): void {
		groupToDelete = group;
		deleteOpen = true;
	}

	async function handleDeleteSuccess(): Promise<void> {
		const remainingTotal = totalItems - 1;
		const totalPages = Math.max(1, Math.ceil(remainingTotal / perPage));

		if (page > totalPages) {
			page = totalPages;
		} else {
			await loadGroups();
		}
	}

	function handleDeleteError(message: string): void {
		toast.error(message);
	}

	function memberLabel(count: number): string {
		return `${count} ${count === 1 ? 'member' : 'members'}`;
	}
</script>

<Table.Root
	title="Groups"
	breadcrumb={[{ label: 'Admin', href: '/admin/' }, { label: 'Groups' }]}
>
	{#snippet header()}
		<Separator />
	{/snippet}

	<Table.PaginatedBody
		itemCount={groups.length}
		{totalItems}
		bind:page
		bind:perPage
		{loading}
		emptyMessage="No groups"
		minimal={false}
		showPerPageSelect
		alwaysShowPagination
		selectTriggerClass="h-9 px-2 py-0"
	>
		{#snippet list()}
			<Table.List>
				{#each groups as group, index (group.id)}
					<Table.Row label={group.name}>
						{#snippet trailing()}
							<Badge>{memberLabel(group.memberCount)}</Badge>
							<Button variant="destructive" size="inline" onclick={() => openDeleteGroup(group)}
								>Delete</Button
							>
						{/snippet}
					</Table.Row>
					{#if index < groups.length - 1}
						<Table.Separator />
					{/if}
				{/each}
			</Table.List>
		{/snippet}
	</Table.PaginatedBody>

	<DeleteGroup
		bind:open={deleteOpen}
		group={groupToDelete}
		onSuccess={handleDeleteSuccess}
		onError={handleDeleteError}
	/>
</Table.Root>

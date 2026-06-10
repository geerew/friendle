<script lang="ts">
	import { ApiError } from '$lib/api';
	import { listGroups } from '$lib/api/groups-api';
	import { DeleteGroup, Pagination, Spinner } from '$lib/components';
	import { Badge, Button, Separator, Table } from '$lib/components/ui';
	import type { GroupModel } from '$lib/models/group-model';
	import { withMinLoadingDelay } from '$lib/utils';

	let groups = $state<GroupModel[]>([]);
	let page = $state(1);
	let perPage = $state(25);
	let totalItems = $state(0);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let deleteOpen = $state(false);
	let groupToDelete = $state<GroupModel | null>(null);

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

	function openDeleteGroup(group: GroupModel): void {
		error = null;
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
		error = message;
	}

	function memberLabel(count: number): string {
		return `${count} ${count === 1 ? 'member' : 'members'}`;
	}
</script>

<Table.Root
	title="Groups"
	breadcrumb={[{ label: 'Admin', href: '/admin/' }, { label: 'Groups' }]}
>
	<div class="flex flex-col gap-6">
		<Separator />

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
			{/if}

			<Pagination
				count={totalItems}
				bind:page
				bind:perPage
				selectTriggerClass="h-9 px-2 py-0"
			/>
		{/if}
	</div>

	<DeleteGroup
		bind:open={deleteOpen}
		group={groupToDelete}
		onSuccess={handleDeleteSuccess}
		onError={handleDeleteError}
	/>
</Table.Root>

<script lang="ts">
	import { ApiError } from '$lib/api';
	import { listGroups } from '$lib/api/admin-api';
	import { AppShell, DeleteGroup, ListRow, Pagination, RightChevronIcon, TrashIcon } from '$lib/components';
	import { Button } from '$lib/components/ui';
	import type { AdminGroupModel } from '$lib/models/admin-group-model';
	import { formatMemberCount } from '$lib/utils';

	let groups = $state<AdminGroupModel[]>([]);
	let page = $state(1);
	let perPage = $state(25);
	let totalItems = $state(0);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let deleteOpen = $state(false);
	let groupToDelete = $state<AdminGroupModel | null>(null);

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

	function openDeleteGroup(group: AdminGroupModel): void {
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
</script>

<AppShell
	breadcrumb={[
		{ label: 'Admin', href: '/admin/' },
		{ label: 'Groups' }
	]}
>
	<div class="flex flex-col gap-6">
		<Button href="/admin/groups/add/" variant="primary" class="w-1/2">+ Add Group</Button>

		<hr class="border-0 border-t border-border" />

		{#if loading}
			<p class="text-text-muted">Loading…</p>
		{:else}
			{#if error}
				<p class="text-sm text-error">{error}</p>
			{/if}

			<div class="flex flex-col gap-3">
				{#if groups.length === 0}
					<p class="text-sm text-text-muted">No groups.</p>
				{:else}
					{#each groups as group (group.id)}
						<ListRow title={group.name} subtitle={formatMemberCount(group.memberCount)}>
							{#snippet trailing()}
								<div class="flex items-center gap-1">
									<Button
										href="/groups/{group.id}/"
										variant="ghost"
										size="icon"
										aria-label="View {group.name}"
									>
										<RightChevronIcon class="h-5 w-5 stroke-2" />
									</Button>
									<Button
										variant="ghost"
										size="icon"
										class="text-error-fg hover:bg-error-bg hover:text-text"
										aria-label="Delete {group.name}"
										onclick={() => openDeleteGroup(group)}
									>
										<TrashIcon class="h-5 w-5 stroke-2" />
									</Button>
								</div>
							{/snippet}
						</ListRow>
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

	<DeleteGroup
		bind:open={deleteOpen}
		group={groupToDelete}
		onSuccess={handleDeleteSuccess}
		onError={handleDeleteError}
	/>
</AppShell>

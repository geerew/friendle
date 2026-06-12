<script lang="ts">
	import { listGroups, searchGroups } from '$lib/api/groups-api';
	import { Spinner } from '$lib/components';
	import { XIcon } from '$lib/components/icons';
	import GroupSearchRow from '$lib/components/pages/groups/search/group-search-row.svelte';
	import { Button, Input, Table } from '$lib/components/ui';
	import type { GroupModel } from '$lib/models/group-model';
	import { apiErrorMessage, withMinLoadingDelay } from '$lib/utils';
	import { toast } from 'svelte-sonner';

	export const SEARCH_DEBOUNCE_MS = 250;

	let query = $state('');
	let searchQuery = $state('');
	let groups = $state<GroupModel[]>([]);

	let page = $state(1);
	let perPage = $state(10);
	let totalItems = $state(0);

	let loading = $state(false);
	let loadRequestId = 0;

	$effect(() => {
		const term = query;
		const timeout = setTimeout(() => {
			const trimmed = term.trim();

			if (trimmed !== searchQuery) {
				page = 1;
				searchQuery = trimmed;
			}
		}, SEARCH_DEBOUNCE_MS);

		return () => clearTimeout(timeout);
	});

	$effect(() => {
		searchQuery;
		page;
		perPage;
		void loadGroups();
	});

	async function loadGroups(): Promise<void> {
		const requestId = ++loadRequestId;
		loading = true;

		try {
			const data = await withMinLoadingDelay(
				searchQuery
					? searchGroups({ name: searchQuery, page, perPage })
					: listGroups({ page, perPage })
			);

			if (requestId !== loadRequestId) {
				return;
			}

			groups = data.items;
			totalItems = data.totalItems;
		} catch (err) {
			if (requestId !== loadRequestId) {
				return;
			}

			groups = [];
			totalItems = 0;
			toast.error(
				apiErrorMessage(err, searchQuery ? 'Failed to search groups' : 'Failed to load groups')
			);
		} finally {
			if (requestId === loadRequestId) {
				loading = false;
			}
		}
	}

	function clearSearch(): void {
		query = '';
	}

	function updateSearchGroup(updated: GroupModel): void {
		groups = groups.map((group) => (group.id === updated.id ? updated : group));
	}
</script>

<Table.Root title="Search Groups">
	<div class="relative w-full">
		<Input
			id="group-search"
			bind:value={query}
			class="pe-11"
			placeholder="Search groups"
			autocomplete="off"
			autofocus
		/>

		<div
			class="pointer-events-none absolute inset-y-0 right-0 flex w-11 items-center justify-center"
		>
			{#if loading && groups.length === 0}
				<Spinner class="bg-foreground-alt-2 size-1.5" />
			{:else if query}
				<Button
					type="button"
					variant="ghost"
					size="icon"
					class="pointer-events-auto size-8 min-w-8 normal-case"
					aria-label="Clear search"
					onclick={clearSearch}
				>
					<XIcon class="size-5 stroke-2" />
				</Button>
			{/if}
		</div>
	</div>

	<Table.PaginatedBody
		itemCount={groups.length}
		{totalItems}
		bind:page
		bind:perPage
		{loading}
		emptyMessage={searchQuery ? 'No groups found' : 'No groups'}
		minimal={false}
		showPerPageSelect
		alwaysShowPagination
		selectTriggerClass="h-9 px-2 py-0"
	>
		{#snippet list()}
			<Table.List>
				{#each groups as group, index (group.id)}
					<GroupSearchRow {group} onGroupChange={updateSearchGroup} />
					{#if index < groups.length - 1}
						<Table.Separator />
					{/if}
				{/each}
			</Table.List>
		{/snippet}
	</Table.PaginatedBody>
</Table.Root>

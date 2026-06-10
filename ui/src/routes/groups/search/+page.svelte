<script lang="ts">
	import { searchGroups } from '$lib/api/groups-api';
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

	const perPage = 7;
	let page = $state(1);
	let totalItems = $state(0);

	let loading = $state(false);
	let searchRequestId = 0;

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
		void runSearch();
	});

	async function runSearch(): Promise<void> {
		if (!searchQuery) {
			groups = [];
			totalItems = 0;
			loading = false;

			return;
		}

		const requestId = ++searchRequestId;
		loading = true;

		try {
			const data = await withMinLoadingDelay(
				searchGroups({
					name: searchQuery,
					page,
					perPage
				})
			);

			if (requestId !== searchRequestId) {
				return;
			}

			groups = data.items;
			totalItems = data.totalItems;
		} catch (err) {
			if (requestId !== searchRequestId) {
				return;
			}

			groups = [];
			totalItems = 0;
			toast.error(apiErrorMessage(err, 'Failed to search groups'));
		} finally {
			if (requestId === searchRequestId) {
				loading = false;
			}
		}
	}

	function clearSearch(): void {
		query = '';
	}

	function markJoinPending(groupId: string): void {
		groups = groups.map((group) =>
			group.id === groupId ? { ...group, joinRequestStatus: 'pending' as const } : group
		);
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

	{#if searchQuery}
		<Table.PaginatedBody
			itemCount={groups.length}
			{totalItems}
			bind:page
			perPage={7}
			{loading}
			emptyMessage="No groups found"
		>
			{#snippet list()}
				<Table.List>
					{#each groups as group, index (group.id)}
						<GroupSearchRow {group} onjoined={markJoinPending} />
						{#if index < groups.length - 1}
							<Table.Separator />
						{/if}
					{/each}
				</Table.List>
			{/snippet}
		</Table.PaginatedBody>
	{/if}
</Table.Root>

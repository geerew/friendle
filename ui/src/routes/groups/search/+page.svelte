<script lang="ts">
	import { searchGroups } from '$lib/api/groups-api';
	import { LoadingOverlay, Pagination, Spinner } from '$lib/components';
	import { XIcon } from '$lib/components/icons';
	import { GroupSearchList } from '$lib/components/pages';
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
		<div class="mt-6 flex flex-col gap-6">
			{#if groups.length === 0 && !loading}
				<div class="flex min-h-16 items-center justify-center">
					<p class="text-foreground-alt-2 text-sm italic">No groups found</p>
				</div>
			{:else if groups.length > 0}
				<LoadingOverlay {loading}>
					<div class="px-2">
						<GroupSearchList bind:groups />
					</div>
				</LoadingOverlay>

				{#if totalItems > perPage}
					<Pagination
						count={totalItems}
						bind:page
						{perPage}
						minimal
						showPerPageSelect={false}
					/>
				{/if}
			{/if}
		</div>
	{/if}
</Table.Root>

<script lang="ts">
	import { untrack } from 'svelte';
	import { ApiError } from '$lib/api';
	import { listSelfGroups } from '$lib/api/groups-api';
	import { AppShell, Pagination, Spinner } from '$lib/components';
	import { RightChevronIcon } from '$lib/components/icons';
	import { Button, Separator } from '$lib/components/ui';
	import type { GroupModel } from '$lib/models/group-model';

	const groupsPerPage = 5;
	const minLoadingMs = 150;

	let groups = $state<GroupModel[]>([]);
	let page = $state(1);
	let perPage = $state(groupsPerPage);
	let totalItems = $state(0);
	let loading = $state(true);
	let hasLoaded = $state(false);
	let error = $state<string | null>(null);
	let listMinHeight = $state<number | undefined>(undefined);

	let listEl = $state<HTMLDivElement | undefined>(undefined);
	let requestId = 0;

	$effect(() => {
		page;
		untrack(() => {
			void loadGroups();
		});
	});

	async function loadGroups(): Promise<void> {
		const id = ++requestId;

		if (hasLoaded && listEl) {
			listMinHeight = listEl.offsetHeight;
		} else {
			listMinHeight = undefined;
		}

		loading = true;
		error = null;

		try {
			const [data] = await Promise.all([
				listSelfGroups({ page, perPage: groupsPerPage }),
				new Promise<void>((resolve) => setTimeout(resolve, minLoadingMs))
			]);

			if (id !== requestId) {
				return;
			}

			groups = data.items;
			totalItems = data.totalItems;
			hasLoaded = true;
		} catch (err) {
			if (id !== requestId) {
				return;
			}

			error = err instanceof ApiError ? err.message : 'Failed to load groups';

			if (!hasLoaded) {
				groups = [];
				totalItems = 0;
			}
		} finally {
			if (id === requestId) {
				loading = false;
			}
		}
	}
</script>

<AppShell>
	<h2 class="section-title">My Groups</h2>

	<div class="flex flex-col gap-7">
		{#if loading && !hasLoaded}
			<div class="flex min-h-24 items-center justify-center">
				<Spinner class="bg-text-muted size-3" />
			</div>
		{:else if error && !hasLoaded}
			<p class="text-error text-sm">{error}</p>
		{:else if groups.length === 0 && !loading}
			<div class="text-text-muted flex min-h-24 items-center justify-center text-sm italic">
				No groups
			</div>
		{:else}
			<div class="flex flex-col gap-4 px-2">
				{#if error}
					<p class="text-error text-sm">{error}</p>
				{/if}

				{#if loading}
					<div
						class="flex items-center justify-center"
						class:min-h-24={listMinHeight === undefined}
						style:min-height={listMinHeight === undefined ? undefined : `${listMinHeight}px`}
					>
						<Spinner class="bg-text-muted size-3" />
					</div>
				{:else}
					<div bind:this={listEl}>
						{#each groups as group, index (group.id)}
							<div class="flex w-full items-center gap-3 py-5 text-left">
								<span class="text-text-secondary min-w-0 flex-1 truncate text-base font-medium">
									{group.name}
								</span>
								<RightChevronIcon class="text-text-muted size-5 shrink-0 stroke-2" />
							</div>
							{#if index < groups.length - 1}
								<Separator dashed />
							{/if}
						{/each}
					</div>
				{/if}

				{#if totalItems > groupsPerPage}
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
		{/if}

		<Separator />

		<div class="flex flex-col gap-3">
			<Button type="button" variant="secondary">Search groups</Button>
			<Button type="button" variant="primary" href="/groups/create/">Create Group</Button>
		</div>
	</div>
</AppShell>

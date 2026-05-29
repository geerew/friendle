<script lang="ts">
	import { ApiError } from '$lib/api';
	import { cancelJoinRequest, joinGroup, searchGroups } from '$lib/api/groups-api';
	import { AppShell, ListRow, Pagination, PlusIcon, XIcon } from '$lib/components';
	import { Badge, Button, Input } from '$lib/components/ui';
	import type { GroupSearchResultModel } from '$lib/models/group-search-model';
	import { formatMemberCount } from '$lib/utils';

	let searchQuery = $state('');
	let activeQuery = $state('');
	let searchResults = $state<GroupSearchResultModel[]>([]);
	let page = $state(1);
	let perPage = $state(25);
	let totalItems = $state(0);
	let searching = $state(false);
	let message = $state<string | null>(null);
	let error = $state<string | null>(null);

	$effect(() => {
		if (!activeQuery) return;

		page;
		perPage;
		void loadResults();
	});

	async function loadResults(): Promise<void> {
		if (!activeQuery) return;

		searching = true;
		error = null;

		try {
			const data = await searchGroups({ q: activeQuery, page, perPage });
			searchResults = data.items;
			totalItems = data.totalItems;
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Search failed';
		} finally {
			searching = false;
		}
	}

	async function handleSearch(): Promise<void> {
		const query = searchQuery.trim();
		if (!query) return;

		message = null;
		error = null;

		if (query !== activeQuery) {
			activeQuery = query;
			page = 1;
			return;
		}

		page = 1;
		await loadResults();
	}

	async function handleJoin(group: GroupSearchResultModel): Promise<void> {
		message = null;
		error = null;

		try {
			await joinGroup(group.id);
			searchResults = searchResults.map((result) =>
				result.id === group.id ? { ...result, joinPending: true } : result
			);
			message = `Requested to join ${group.name}`;
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to join group';
		}
	}

	async function handleCancelJoin(group: GroupSearchResultModel): Promise<void> {
		message = null;
		error = null;

		try {
			await cancelJoinRequest(group.id);
			searchResults = searchResults.map((result) =>
				result.id === group.id ? { ...result, joinPending: false } : result
			);
			message = `Cancelled join request for ${group.name}`;
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to cancel join request';
		}
	}
</script>

<AppShell title="Search Groups">
	<div class="flex flex-col gap-6">
		<div class="flex flex-col gap-2">
			<Input
				placeholder="Search by name"
				bind:value={searchQuery}
				onkeydown={(event) => event.key === 'Enter' && handleSearch()}
			/>
			<Button variant="secondary" onclick={handleSearch} disabled={searching || !searchQuery.trim()}>
				{searching ? 'Searching…' : 'Search'}
			</Button>
		</div>

		{#if message}
			<p class="text-sm text-text-muted">{message}</p>
		{/if}

		{#if error}
			<p class="text-sm text-error">{error}</p>
		{/if}

		{#if activeQuery}
			<div class="flex flex-col gap-6">
				{#if !searching && searchResults.length === 0}
					<p class="text-sm text-text-muted">No groups found.</p>
				{:else}
					<div class="flex flex-col gap-3">
						{#each searchResults as group (group.id)}
							{#if group.isMember}
								<ListRow
									href="/groups/{group.id}/"
									title={group.name}
									subtitle={formatMemberCount(group.memberCount)}
								/>
							{:else}
								<ListRow title={group.name} subtitle={formatMemberCount(group.memberCount)}>
									{#snippet trailing()}
										{#if group.joinPending}
											<div class="flex items-center gap-1">
												<Badge>pending</Badge>
												<Button
													type="button"
													variant="ghost"
													size="icon"
													class="text-text-muted hover:text-text"
													aria-label="Cancel join request for {group.name}"
													onclick={() => handleCancelJoin(group)}
												>
													<XIcon class="size-5 stroke-2" />
												</Button>
											</div>
										{:else}
											<Button
												type="button"
												variant="primary"
												size="icon"
												aria-label="Join {group.name}"
												onclick={() => handleJoin(group)}
											>
												<PlusIcon class="size-5 stroke-2" />
											</Button>
										{/if}
									{/snippet}
								</ListRow>
							{/if}
						{/each}
					</div>
				{/if}

				{#if totalItems > 0}
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
		{/if}
	</div>
</AppShell>

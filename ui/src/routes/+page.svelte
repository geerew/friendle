<script lang="ts">
	import { ApiError } from '$lib/api';
	import { listSelfGroups } from '$lib/api/groups-api';
	import { AppShell, Spinner } from '$lib/components';
	import { GroupList } from '$lib/components/pages';
	import { Button, Separator } from '$lib/components/ui';
	import type { GroupModel } from '$lib/models/group-model';
	import { withMinLoadingDelay } from '$lib/utils';

	const homeGroupLimit = 4;

	let groups = $state<GroupModel[]>([]);
	let totalItems = $state(0);
	let loading = $state(true);
	let error = $state<string | null>(null);

	const showMoreRow = $derived(totalItems > homeGroupLimit);

	$effect(() => {
		void loadGroups();
	});

	async function loadGroups(): Promise<void> {
		loading = true;
		error = null;

		try {
			const data = await withMinLoadingDelay(
				listSelfGroups({ page: 1, perPage: homeGroupLimit }),
				200
			);
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

	<div class="flex flex-col gap-7">
		{#if loading}
			<div class="flex min-h-24 items-center justify-center">
				<Spinner class="bg-foreground-alt-2 size-3" />
			</div>
		{:else if error}
			<p class="text-foreground-error text-sm">{error}</p>
		{:else if groups.length === 0}
			<div class="text-foreground-alt-2 flex min-h-24 items-center justify-center text-sm italic">
				No groups
			</div>
		{:else}
			<div class="flex flex-col gap-2 px-2">
				<GroupList {groups} />

				{#if showMoreRow}
					<div class="flex w-full items-center justify-center">
						<Separator class="bg-foreground-alt-5 w-[95%]" />
					</div>
					<div class="flex w-full items-center justify-center px-8">
						<Button
							href="/groups/"
							variant="ghost"
							class="bg-background-primary/15 text-foreground-alt-1 hover:bg-background-primary/35 enabled:hover:text-foreground-alt-1"
						>
							More
						</Button>
					</div>
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

<script lang="ts">
	import { page } from '$app/state';
	import { ApiError } from '$lib/api';
	import { listGroupMembers } from '$lib/api/groups-api';
	import { AppShell, LoadingOverlay, Pagination, Spinner } from '$lib/components';
	import { GroupMemberList } from '$lib/components/pages';
	import { GROUP_PAGE_KEY, type GroupPageContext } from '$lib/context/group-page';
	import type { GroupMemberModel } from '$lib/models/group-member-model';
	import { withMinLoadingDelay } from '$lib/utils';
	import { groupChildBreadcrumb } from '$lib/utils/group';
	import { getContext } from 'svelte';

	const groupId = $derived(page.params.id ?? '');
	const groupPage = getContext<GroupPageContext>(GROUP_PAGE_KEY);
	const group = $derived(groupPage.group);
	const breadcrumb = $derived(
		group ? groupChildBreadcrumb(group.id, group.name, 'Members') : [{ label: 'Members' }]
	);

	let members = $state<GroupMemberModel[]>([]);
	let pageNum = $state(1);
	let perPage = $state(10);
	let totalItems = $state(0);
	let loading = $state(true);
	let error = $state<string | null>(null);

	$effect(() => {
		groupId;
		pageNum;
		perPage;
		void loadMembers();
	});

	async function loadMembers(): Promise<void> {
		if (!groupId) {
			members = [];
			totalItems = 0;
			error = 'Group not found';
			loading = false;
			return;
		}

		loading = true;
		error = null;

		try {
			const data = await withMinLoadingDelay(
				listGroupMembers(groupId, { page: pageNum, perPage })
			);
			members = data.items;
			totalItems = data.totalItems;
		} catch (err) {
			members = [];
			totalItems = 0;
			error = err instanceof ApiError ? err.message : 'Failed to load members';
		} finally {
			loading = false;
		}
	}
</script>

<AppShell {breadcrumb}>
	<div class="flex flex-col gap-5">
		<h2 class="section-title">Members</h2>

		<div class="flex flex-col gap-6">
			{#if error}
				<p class="text-foreground-error text-sm">{error}</p>
			{/if}

			{#if loading && members.length === 0}
				<div class="flex min-h-24 items-center justify-center">
					<Spinner class="bg-foreground-alt-2 size-3" />
				</div>
			{:else if members.length === 0}
				<div class="flex min-h-16 items-center justify-center">
					<p class="text-foreground-alt-2 text-sm italic">No members</p>
				</div>
			{:else}
				<LoadingOverlay {loading}>
					<div class="px-2">
						<GroupMemberList {members} />
					</div>
				</LoadingOverlay>

				{#if totalItems > perPage}
					<Pagination
						count={totalItems}
						bind:page={pageNum}
						bind:perPage
						minimal
						showPerPageSelect={false}
						onPageChange={() => {}}
						onPerPageChange={() => {}}
					/>
				{/if}
			{/if}
		</div>
	</div>
</AppShell>

<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { listGroupPendingJoinRequests } from '$lib/api/groups-api';
	import { AppShell, LoadingOverlay, Pagination, Spinner } from '$lib/components';
	import { GroupPendingList } from '$lib/components/pages';
	import { GROUP_PAGE_KEY, type GroupPageContext } from '$lib/context/group-page';
	import type { GroupJoinRequestModel } from '$lib/models/group-join-request-model';
	import { groupChildBreadcrumb, isGroupAdmin } from '$lib/utils/group';
	import { apiErrorMessage, withMinLoadingDelay } from '$lib/utils';
	import { getContext } from 'svelte';

	const groupId = $derived(page.params.id ?? '');
	const groupPage = getContext<GroupPageContext>(GROUP_PAGE_KEY);
	const group = $derived(groupPage.group);
	const breadcrumb = $derived(
		group
			? groupChildBreadcrumb(group.id, group.name, 'Pending Requests')
			: [{ label: 'Pending Requests' }]
	);

	let requests = $state<GroupJoinRequestModel[]>([]);
	let pageNum = $state(1);
	let perPage = $state(10);
	let totalItems = $state(0);
	let loading = $state(true);
	let error = $state<string | null>(null);

	$effect(() => {
		groupId;
		group;

		if (group && !isGroupAdmin(group)) {
			void goto(`/groups/${groupId}/`);
			return;
		}

		pageNum;
		perPage;
		void loadPending();
	});

	async function loadPending(): Promise<void> {
		if (!groupId || !group || !isGroupAdmin(group)) {
			return;
		}

		loading = true;
		error = null;

		try {
			const data = await withMinLoadingDelay(
				listGroupPendingJoinRequests(groupId, { page: pageNum, perPage })
			);
			requests = data.items;
			totalItems = data.totalItems;
		} catch (err) {
			requests = [];
			totalItems = 0;
			error = apiErrorMessage(err, 'Failed to load pending requests');
		} finally {
			loading = false;
		}
	}
</script>

<AppShell {breadcrumb}>
	<div class="flex flex-col gap-5">
		<h2 class="section-title">Pending Requests</h2>

		<div class="flex flex-col gap-6">
			{#if error}
				<p class="text-foreground-error text-sm">{error}</p>
			{/if}

			{#if loading && requests.length === 0}
				<div class="flex min-h-24 items-center justify-center">
					<Spinner class="bg-foreground-alt-2 size-3" />
				</div>
			{:else if requests.length === 0}
				<div class="flex min-h-16 items-center justify-center">
					<p class="text-foreground-alt-2 text-sm italic">No pending requests</p>
				</div>
			{:else}
				<LoadingOverlay {loading}>
					<div class="px-2">
						<GroupPendingList {groupId} {requests} onchange={loadPending} />
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

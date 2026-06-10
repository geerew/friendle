<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { listGroupPendingJoinRequests } from '$lib/api/groups-api';
	import { GroupPaginatedListSection, GroupPendingList } from '$lib/components/pages';
	import { Table } from '$lib/components/ui';
	import { GROUP_PAGE_KEY, type GroupPageContext } from '$lib/context/group-page';
	import type { GroupJoinRequestModel } from '$lib/models/group-join-request-model';
	import { apiErrorMessage, withMinLoadingDelay } from '$lib/utils';
	import { groupChildBreadcrumb, isGroupAdmin } from '$lib/utils/group';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

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

		try {
			const data = await withMinLoadingDelay(
				listGroupPendingJoinRequests(groupId, { page: pageNum, perPage })
			);
			requests = data.items;
			totalItems = data.totalItems;
		} catch (err) {
			requests = [];
			totalItems = 0;
			toast.error(apiErrorMessage(err, 'Failed to load pending requests'));
		} finally {
			loading = false;
		}
	}

	// handlePendingChange reloads the list and refreshes group stats
	function handlePendingChange(): void {
		void loadPending();
		void groupPage.reloadGroup({ silent: true });
	}
</script>

<Table.Root title="Pending Requests" {breadcrumb}>
	<GroupPaginatedListSection
		itemCount={requests.length}
		{totalItems}
		bind:page={pageNum}
		bind:perPage
		{loading}
		emptyMessage="No pending requests"
	>
		{#snippet list()}
			<GroupPendingList {groupId} {requests} onchange={handlePendingChange} />
		{/snippet}
	</GroupPaginatedListSection>
</Table.Root>

<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { listGroupRejectedJoinRequests } from '$lib/api/groups-api';
	import { GroupNameSection } from '$lib/components/pages';
	import { Table } from '$lib/components/ui';
	import { GROUP_PAGE_KEY, type GroupPageContext } from '$lib/context/group-page';
	import type { GroupJoinRequestModel } from '$lib/models/group-join-request-model';
	import { apiErrorMessage, withMinLoadingDelay } from '$lib/utils';
	import { groupSettingsChildBreadcrumb, isGroupAdmin } from '$lib/utils/group';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	const groupId = $derived(page.params.id ?? '');
	const groupPage = getContext<GroupPageContext>(GROUP_PAGE_KEY);
	const group = $derived(groupPage.group);
	const breadcrumb = $derived(
		group ? groupSettingsChildBreadcrumb(group.id, 'Rejected') : [{ label: 'Rejected' }]
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
		void loadRejected();
	});

	async function loadRejected(): Promise<void> {
		if (!groupId || !group || !isGroupAdmin(group)) {
			return;
		}

		loading = true;

		try {
			const data = await withMinLoadingDelay(
				listGroupRejectedJoinRequests(groupId, { page: pageNum, perPage })
			);
			requests = data.items;
			totalItems = data.totalItems;
		} catch (err) {
			requests = [];
			totalItems = 0;
			toast.error(apiErrorMessage(err, 'Failed to load rejected requests'));
		} finally {
			loading = false;
		}
	}
</script>

<Table.Root title="Rejected Requests" {breadcrumb}>
	{#snippet header()}
		<GroupNameSection />
	{/snippet}

	<Table.PaginatedBody
		itemCount={requests.length}
		{totalItems}
		bind:page={pageNum}
		bind:perPage
		{loading}
		emptyMessage="No rejected requests"
	>
		{#snippet list()}
			<Table.List>
				{#each requests as request, index (request.userId)}
					<Table.Row label={request.displayName} />
					{#if index < requests.length - 1}
						<Table.Separator />
					{/if}
				{/each}
			</Table.List>
		{/snippet}
	</Table.PaginatedBody>
</Table.Root>

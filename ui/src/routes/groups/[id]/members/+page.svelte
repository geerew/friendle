<script lang="ts">
	import { page } from '$app/state';
	import { listGroupMembers } from '$lib/api/groups-api';
	import { PencilIcon } from '$lib/components/icons';
	import { GroupNameSection, GroupPaginatedListSection } from '$lib/components/pages';
	import GroupRoleBadge from '$lib/components/pages/groups/group-role-badge.svelte';
	import { Table } from '$lib/components/ui';
	import { GROUP_PAGE_KEY, type GroupPageContext } from '$lib/context/group-page';
	import type { GroupMemberModel } from '$lib/models/group-member-model';
	import { apiErrorMessage, withMinLoadingDelay } from '$lib/utils';
	import { groupChildBreadcrumb } from '$lib/utils/group';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	const groupId = $derived(page.params.id ?? '');
	const groupPage = getContext<GroupPageContext>(GROUP_PAGE_KEY);
	const group = $derived(groupPage.group);
	const breadcrumb = $derived(
		group ? groupChildBreadcrumb(group.id, 'Members') : [{ label: 'Members' }]
	);

	let members = $state<GroupMemberModel[]>([]);
	let pageNum = $state(1);
	let perPage = $state(10);
	let totalItems = $state(0);
	let loading = $state(true);

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
			toast.error('Group not found');
			loading = false;
			return;
		}

		loading = true;

		try {
			const data = await withMinLoadingDelay(
				listGroupMembers(groupId, { page: pageNum, perPage })
			);
			members = data.items;
			totalItems = data.totalItems;
		} catch (err) {
			members = [];
			totalItems = 0;
			toast.error(apiErrorMessage(err, 'Failed to load members'));
		} finally {
			loading = false;
		}
	}
</script>

<Table.Root title="Members" {breadcrumb}>
	{#snippet header()}
		<GroupNameSection />
	{/snippet}

	<GroupPaginatedListSection
		itemCount={members.length}
		{totalItems}
		bind:page={pageNum}
		bind:perPage
		{loading}
		emptyMessage="No members"
	>
		{#snippet list()}
			<Table.List>
				{#each members as member, index (member.userId)}
					<Table.Row label={member.displayName}>
						{#snippet trailing()}
							<div class="flex shrink-0 items-center gap-1.5">
								<GroupRoleBadge role={member.groupRole} />
								<button
									type="button"
									class="text-foreground-alt-2 hover:bg-background-alt-1 hover:text-foreground focus-visible:ring-background-primary flex size-8 shrink-0 items-center justify-center rounded-md transition-colors focus-visible:ring-2 focus-visible:outline-none"
									aria-label="Edit {member.displayName}"
								>
									<PencilIcon class="size-4 shrink-0 stroke-2" />
								</button>
							</div>
						{/snippet}
					</Table.Row>
					{#if index < members.length - 1}
						<Table.Separator />
					{/if}
				{/each}
			</Table.List>
		{/snippet}
	</GroupPaginatedListSection>
</Table.Root>

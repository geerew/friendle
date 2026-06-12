<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { GroupNameSection, GroupStatLink } from '$lib/components/pages';
	import { Separator, Table } from '$lib/components/ui';
	import { GROUP_PAGE_KEY, type GroupPageContext } from '$lib/context/group-page';
	import { groupChildBreadcrumb, isGroupAdmin } from '$lib/utils/group';
	import { getContext } from 'svelte';

	const groupId = $derived(page.params.id ?? '');
	const groupPage = getContext<GroupPageContext>(GROUP_PAGE_KEY);
	const group = $derived(groupPage.group);
	const breadcrumb = $derived(
		group ? groupChildBreadcrumb(group.id, 'Settings') : [{ label: 'Settings' }]
	);

	$effect(() => {
		groupId;
		group;

		if (group && !isGroupAdmin(group)) {
			void goto(`/groups/${groupId}/`);
		}
	});
</script>

{#if group && isGroupAdmin(group) && group.adminSummary}
	<Table.Root title="Settings" {breadcrumb} showTitle={false}>
		<GroupNameSection />

		<section class="flex items-stretch justify-between gap-3 py-3">
			<GroupStatLink
				label="Members"
				count={group.memberCount}
				href="/groups/{group.id}/settings/members"
				ariaLabel="View members"
				align="start"
			/>
			<Separator class="h-auto w-px shrink-0 self-stretch" />
			<GroupStatLink
				label="Pending"
				count={group.adminSummary.pendingJoinRequestCount}
				href="/groups/{group.id}/settings/pending"
				ariaLabel="View pending join requests"
			/>
			<Separator class="h-auto w-px shrink-0 self-stretch" />
			<GroupStatLink
				label="Rejected"
				count={group.adminSummary.rejectedJoinRequestCount}
				href="/groups/{group.id}/settings/rejected"
				ariaLabel="View rejected join requests"
				align="end"
			/>
		</section>
	</Table.Root>
{/if}

<script lang="ts">
	import { GroupNameSection, GroupStatLink } from '$lib/components/pages';
	import { Separator, Table } from '$lib/components/ui';
	import { GROUP_PAGE_KEY, type GroupPageContext } from '$lib/context/group-page';
	import { groupHomeBreadcrumb, isGroupAdmin } from '$lib/utils/group';
	import { getContext } from 'svelte';

	const groupPage = getContext<GroupPageContext>(GROUP_PAGE_KEY);
	const group = $derived(groupPage.group);

	const isAdmin = $derived(group != null && isGroupAdmin(group));
	const breadcrumb = $derived(group ? groupHomeBreadcrumb() : [{ label: 'Group' }]);
</script>

{#if group}
	<Table.Root title="Group" {breadcrumb} showTitle={false}>
		<GroupNameSection />

		{#if isAdmin && group.adminSummary}
			<section class="flex items-stretch justify-between gap-3">
				<GroupStatLink
					label="Members"
					count={group.memberCount}
					href="/groups/{group.id}/members"
					ariaLabel="View members"
					align="start"
				/>
				<Separator class="h-auto w-px shrink-0 self-stretch" />
				<GroupStatLink
					label="Pending"
					count={group.adminSummary.pendingJoinRequestCount}
					href="/groups/{group.id}/pending"
					ariaLabel="View pending join requests"
				/>
				<Separator class="h-auto w-px shrink-0 self-stretch" />
				<GroupStatLink
					label="Rejected"
					count={group.adminSummary.rejectedJoinRequestCount}
					href="/groups/{group.id}/rejected"
					ariaLabel="View rejected join requests"
					align="end"
				/>
			</section>
		{:else}
			<section class="flex flex-col gap-3">
				<GroupStatLink
					label="Members"
					count={group.memberCount}
					href="/groups/{group.id}/members"
					ariaLabel="View members"
					align="start"
				/>
			</section>
		{/if}
	</Table.Root>
{/if}

<script lang="ts">
	import { AppShell } from '$lib/components';
	import { GroupNameSection, GroupStatLink } from '$lib/components/pages';
	import { GROUP_PAGE_KEY, type GroupPageContext } from '$lib/context/group-page';
	import { isGroupAdmin } from '$lib/utils/group';
	import { getContext } from 'svelte';

	const groupPage = getContext<GroupPageContext>(GROUP_PAGE_KEY);
	const group = $derived(groupPage.group);

	const isAdmin = $derived(group != null && isGroupAdmin(group));
</script>

{#if group}
	<AppShell breadcrumb={[{ label: 'Group' }]}>
		<div class="flex flex-col gap-5">
			<GroupNameSection />

			{#if isAdmin && group.adminSummary}
				<section class="grid grid-cols-3 gap-3">
					<GroupStatLink
						label="Members"
						count={group.memberCount}
						href="/groups/{group.id}/members"
						ariaLabel="View members"
					/>
					<GroupStatLink
						label="Pending"
						count={group.adminSummary.pendingJoinRequestCount}
						href="/groups/{group.id}/pending"
						ariaLabel="View pending join requests"
					/>
					<GroupStatLink
						label="Rejected"
						count={group.adminSummary.rejectedJoinRequestCount}
						href="/groups/{group.id}/rejected"
						ariaLabel="View rejected join requests"
					/>
				</section>
			{:else}
				<section class="flex flex-col gap-3">
					<GroupStatLink
						label="Members"
						count={group.memberCount}
						href="/groups/{group.id}/members"
						ariaLabel="View members"
					/>
				</section>
			{/if}
		</div>
	</AppShell>
{/if}

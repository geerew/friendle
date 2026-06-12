<script lang="ts">
	import { GroupNameSection, GroupRoundStatus } from '$lib/components/pages';
	import { Table } from '$lib/components/ui';
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
		<GroupNameSection aside={isAdmin ? 'settings' : 'members'} />

		<GroupRoundStatus groupId={group.id} memberThresholdMet={group.memberThresholdMet} />
	</Table.Root>
{/if}

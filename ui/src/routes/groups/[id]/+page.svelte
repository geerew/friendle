<script lang="ts">
	import { GroupNameSection, GroupRoundStatus } from '$lib/components/pages';
	import { Table, Separator } from '$lib/components/ui';
	import { GROUP_PAGE_KEY, type GroupPageContext } from '$lib/context/group-page';
	import { groupHomeBreadcrumb } from '$lib/utils/group';
	import { getContext } from 'svelte';

	const groupPage = getContext<GroupPageContext>(GROUP_PAGE_KEY);
	const group = $derived(groupPage.group);

	const breadcrumb = $derived(group ? groupHomeBreadcrumb() : [{ label: 'Group' }]);
</script>

{#if group}
	<Table.Root title="Group" {breadcrumb} showTitle={false}>
		<div class="flex flex-col gap-5">
			<GroupNameSection showSettings />

			<Separator />

			<GroupRoundStatus groupId={group.id} memberThresholdMet={group.memberThresholdMet} />
		</div>
	</Table.Root>
{/if}

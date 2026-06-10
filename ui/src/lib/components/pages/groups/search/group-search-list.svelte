<script lang="ts">
	import { Table } from '$lib/components/ui';
	import type { GroupModel } from '$lib/models/group-model';
	import GroupSearchRow from './group-search-row.svelte';

	type Props = {
		groups: GroupModel[];
	};

	let { groups = $bindable() }: Props = $props();

	function markJoinPending(groupId: string): void {
		groups = groups.map((group) =>
			group.id === groupId ? { ...group, joinRequestStatus: 'pending' as const } : group
		);
	}
</script>

<Table.List>
	{#each groups as group, index (group.id)}
		<GroupSearchRow {group} onjoined={markJoinPending} />
		{#if index < groups.length - 1}
			<Table.Separator />
		{/if}
	{/each}
</Table.List>

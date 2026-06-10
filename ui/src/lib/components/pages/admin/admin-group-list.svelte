<script lang="ts">
	import { Badge, Button, Table } from '$lib/components/ui';
	import type { GroupModel } from '$lib/models/group-model';

	type Props = {
		groups: GroupModel[];
		onDelete: (group: GroupModel) => void;
	};

	let { groups, onDelete }: Props = $props();

	function memberLabel(count: number): string {
		return `${count} ${count === 1 ? 'member' : 'members'}`;
	}
</script>

<Table.List>
	{#each groups as group, index (group.id)}
		<Table.Row label={group.name}>
			{#snippet trailing()}
				<Badge>{memberLabel(group.memberCount)}</Badge>
				<Button variant="destructive" size="inline" onclick={() => onDelete(group)}>Delete</Button>
			{/snippet}
		</Table.Row>
		{#if index < groups.length - 1}
			<Table.Separator />
		{/if}
	{/each}
</Table.List>

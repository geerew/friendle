<script lang="ts">
	import { RightChevronIcon } from '$lib/components/icons';
	import { Table } from '$lib/components/ui';
	import type { GroupModel } from '$lib/models/group-model';

	type Props = {
		groups: GroupModel[];
	};

	let { groups }: Props = $props();
</script>

<Table.List>
	{#each groups as group, index (group.id)}
		<Table.Row label={group.name} href="/groups/{group.id}/">
			{#snippet trailing()}
				<RightChevronIcon class="text-foreground-alt-2 size-5 shrink-0 stroke-2" />
			{/snippet}
		</Table.Row>
		{#if index < groups.length - 1}
			<Table.Separator />
		{/if}
	{/each}
</Table.List>

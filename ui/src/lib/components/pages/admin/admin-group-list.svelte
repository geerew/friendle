<script lang="ts">
	import { Badge, Button, Separator } from '$lib/components/ui';
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

<div class="flex flex-col gap-2">
	{#each groups as group, index (group.id)}
		<div class="flex w-full items-center gap-3 px-2 py-3">
			<span class="text-foreground-alt-1 min-w-0 flex-1 truncate px-1 text-base font-medium">
				{group.name}
			</span>
			<Badge>{memberLabel(group.memberCount)}</Badge>
			<Button variant="destructive" size="inline" onclick={() => onDelete(group)}>Delete</Button>
		</div>

		{#if index < groups.length - 1}
			<div class="flex w-full items-center justify-center">
				<Separator class="bg-foreground-alt-5 w-full" />
			</div>
		{/if}
	{/each}
</div>

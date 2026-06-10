<script lang="ts">
	import { RightChevronIcon } from '$lib/components/icons';
	import GroupListShell from '../groups/group-list-shell.svelte';
	import GroupListRow from '../groups/group-list-row.svelte';
	import GroupListSeparator from '../groups/group-list-separator.svelte';
	import type { GroupModel } from '$lib/models/group-model';

	type Props = {
		groups: GroupModel[];
	};

	let { groups }: Props = $props();
</script>

<GroupListShell>
	{#each groups as group, index (group.id)}
		<GroupListRow name={group.name} href="/groups/{group.id}/">
			{#snippet trailing()}
				<RightChevronIcon class="text-foreground-alt-2 size-5 shrink-0 stroke-2" />
			{/snippet}
		</GroupListRow>
		{#if index < groups.length - 1}
			<GroupListSeparator />
		{/if}
	{/each}
</GroupListShell>

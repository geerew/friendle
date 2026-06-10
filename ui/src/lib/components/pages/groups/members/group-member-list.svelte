<script lang="ts">
	import { PencilIcon } from '$lib/components/icons';
	import { Table } from '$lib/components/ui';
	import GroupRoleBadge from '../group-role-badge.svelte';
	import type { GroupMemberModel } from '$lib/models/group-member-model';

	type Props = {
		members: GroupMemberModel[];
	};

	let { members }: Props = $props();
</script>

<Table.List>
	{#each members as member, index (member.userId)}
		<Table.Row label={member.displayName}>
			{#snippet trailing()}
				<div class="flex shrink-0 items-center gap-1.5">
					<GroupRoleBadge role={member.groupRole} />
					<button
						type="button"
						class="text-foreground-alt-2 hover:bg-background-alt-1 hover:text-foreground focus-visible:ring-background-primary flex size-8 shrink-0 items-center justify-center rounded-md transition-colors focus-visible:ring-2 focus-visible:outline-none"
						aria-label="Edit {member.displayName}"
					>
						<PencilIcon class="size-4 shrink-0 stroke-2" />
					</button>
				</div>
			{/snippet}
		</Table.Row>
		{#if index < members.length - 1}
			<Table.Separator />
		{/if}
	{/each}
</Table.List>

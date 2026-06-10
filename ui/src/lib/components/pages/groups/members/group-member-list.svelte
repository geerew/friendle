<script lang="ts">
	import { PencilIcon, UserRoundIcon, UserStarIcon } from '$lib/components/icons';
	import GroupStatusBadge from '../group-status-badge.svelte';
	import { Separator } from '$lib/components/ui';
	import type { GroupMemberModel } from '$lib/models/group-member-model';

	type Props = {
		members: GroupMemberModel[];
	};

	let { members }: Props = $props();
</script>

<div class="flex flex-col gap-2">
	{#each members as member, index (member.userId)}
		<div class="flex w-full items-center gap-3 rounded-md px-2 py-3 text-left">
			<div class="flex w-full items-center gap-3 px-1">
				<span
					class="text-foreground-alt-1 min-w-0 flex-1 truncate text-base leading-5 font-medium"
				>
					{member.displayName}
				</span>
				<div class="flex shrink-0 items-center gap-1.5">
					{#if member.groupRole === 'group_admin'}
						<GroupStatusBadge variant="admin" label="Admin">
							{#snippet icon()}
								<UserStarIcon class="size-3.5 shrink-0 stroke-2" />
							{/snippet}
						</GroupStatusBadge>
					{:else}
						<GroupStatusBadge variant="member" label="Member">
							{#snippet icon()}
								<UserRoundIcon class="size-3.5 shrink-0 stroke-2" />
							{/snippet}
						</GroupStatusBadge>
					{/if}
					<button
						type="button"
						class="text-foreground-alt-2 hover:bg-background-alt-1 hover:text-foreground focus-visible:ring-background-primary flex size-8 shrink-0 items-center justify-center rounded-md transition-colors focus-visible:ring-2 focus-visible:outline-none"
						aria-label="Edit {member.displayName}"
					>
						<PencilIcon class="size-4 shrink-0 stroke-2" />
					</button>
				</div>
			</div>
		</div>
		{#if index < members.length - 1}
			<div class="flex w-full items-center justify-center">
				<Separator class="bg-foreground-alt-5 w-[95%]" />
			</div>
		{/if}
	{/each}
</div>

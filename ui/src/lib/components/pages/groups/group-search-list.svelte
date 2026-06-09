<script lang="ts">
	import {
		PlusIcon,
		RightChevronIcon,
		UserRoundIcon,
		UserRoundPlusIcon,
		UserRoundXIcon,
		UserStarIcon
	} from '$lib/components/icons';
	import GroupStatusBadge from './group-status-badge.svelte';
	import { Separator } from '$lib/components/ui';
	import type { GroupModel } from '$lib/models/group-model';

	type Props = {
		groups: GroupModel[];
	};

	let { groups }: Props = $props();

	function isMember(group: GroupModel): boolean {
		return group.groupRole === 'group_admin' || group.groupRole === 'group_user';
	}
</script>

<div class="flex flex-col gap-2">
	{#each groups as group, index (group.id)}
		{#if isMember(group)}
			<a
				href="/groups/{group.id}/"
				class="hover:bg-background-alt-1 flex w-full items-center gap-3 rounded-md px-2 py-3 text-left transition-all"
			>
				<div class="flex w-full items-start gap-3 px-1">
					<span class="text-foreground-alt-1 min-w-0 flex-1 break-words text-base font-medium">
						{group.name}
					</span>
					{#if group.groupRole === 'group_admin'}
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
					<RightChevronIcon class="text-foreground-alt-2 size-5 shrink-0 stroke-2" />
				</div>
			</a>
		{:else if group.joinRequestStatus === 'pending'}
			<div
				class="flex w-full cursor-default items-center gap-3 rounded-md px-2 py-3 text-left transition-all"
			>
				<div class="flex w-full items-start gap-3 px-1">
					<span class="text-foreground-alt-1 min-w-0 flex-1 break-words text-base font-medium">
						{group.name}
					</span>
					<GroupStatusBadge variant="pending" label="Pending accept">
						{#snippet icon()}
							<UserRoundPlusIcon class="size-3.5 shrink-0 stroke-2" />
						{/snippet}
					</GroupStatusBadge>
				</div>
			</div>
		{:else if group.joinRequestStatus === 'rejected'}
			<div
				class="flex w-full cursor-default items-center gap-3 rounded-md px-2 py-3 text-left transition-all"
			>
				<div class="flex w-full items-start gap-3 px-1">
					<span class="text-foreground-alt-1 min-w-0 flex-1 break-words text-base font-medium">
						{group.name}
					</span>
					<GroupStatusBadge variant="rejected" label="Rejected">
						{#snippet icon()}
							<UserRoundXIcon class="size-3.5 shrink-0 stroke-2" />
						{/snippet}
					</GroupStatusBadge>
				</div>
			</div>
		{:else}
			<div
				class="flex w-full cursor-default items-center gap-3 rounded-md px-2 py-3 text-left transition-all"
			>
				<div class="flex w-full items-start gap-3 px-1">
					<span class="text-foreground-alt-1 min-w-0 flex-1 break-words text-base font-medium">
						{group.name}
					</span>
					<PlusIcon class="text-foreground-alt-2 size-5 shrink-0 stroke-2" />
				</div>
			</div>
		{/if}
		{#if index < groups.length - 1}
			<div class="flex w-full items-center justify-center">
				<Separator class="bg-foreground-alt-5 w-[95%]" />
			</div>
		{/if}
	{/each}
</div>

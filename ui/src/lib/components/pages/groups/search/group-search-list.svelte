<script lang="ts">
	import {
		RightChevronIcon,
		UserRoundIcon,
		UserRoundPlusIcon,
		UserRoundXIcon,
		UserStarIcon
	} from '$lib/components/icons';
	import GroupJoinButton from '../group-join-button.svelte';
	import GroupStatusBadge from '../group-status-badge.svelte';
	import { Separator } from '$lib/components/ui';
	import type { GroupModel } from '$lib/models/group-model';
	import { isGroupMember } from '$lib/utils/group';

	type Props = {
		groups: GroupModel[];
	};

	let { groups = $bindable() }: Props = $props();

	let joinErrors = $state<Record<string, string>>({});

	function markJoinPending(groupId: string): void {
		groups = groups.map((group) =>
			group.id === groupId ? { ...group, joinRequestStatus: 'pending' as const } : group
		);
		joinErrors = { ...joinErrors, [groupId]: '' };
	}

	function setJoinError(groupId: string, message: string): void {
		joinErrors = { ...joinErrors, [groupId]: message };
	}
</script>

<div class="flex flex-col gap-2">
	{#each groups as group, index (group.id)}
		<div class="flex flex-col gap-1">
			{#if isGroupMember(group)}
				<a
					href="/groups/{group.id}/"
					class="hover:bg-background-alt-1 flex w-full items-start gap-3 rounded-md px-2 py-3 text-left transition-all"
				>
					<div class="flex w-full items-start gap-3 px-1">
						<span
							class="text-foreground-alt-1 min-w-0 flex-1 wrap-break-word text-base leading-5 font-medium"
						>
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
						<div class="flex h-5 w-5 shrink-0 items-center justify-center self-start">
							<RightChevronIcon class="text-foreground-alt-2 size-5 shrink-0 stroke-2" />
						</div>
					</div>
				</a>
			{:else if group.joinRequestStatus === 'pending'}
				<div
					class="flex w-full cursor-default items-start gap-3 rounded-md px-2 py-3 text-left transition-all"
				>
					<div class="flex w-full items-start gap-3 px-1">
						<span
							class="text-foreground-alt-1 min-w-0 flex-1 wrap-break-word text-base leading-5 font-medium"
						>
							{group.name}
						</span>
						<GroupStatusBadge variant="pending" label="Pending">
							{#snippet icon()}
								<UserRoundPlusIcon class="size-3.5 shrink-0 stroke-2" />
							{/snippet}
						</GroupStatusBadge>
						<div
							class="flex h-5 w-5 shrink-0 items-center justify-center self-start"
							aria-hidden="true"
						></div>
					</div>
				</div>
			{:else if group.joinRequestStatus === 'rejected'}
				<div
					class="flex w-full cursor-default items-start gap-3 rounded-md px-2 py-3 text-left transition-all"
				>
					<div class="flex w-full items-start gap-3 px-1">
						<span
							class="text-foreground-alt-1 min-w-0 flex-1 wrap-break-word text-base leading-5 font-medium"
						>
							{group.name}
						</span>
						<GroupStatusBadge variant="rejected" label="Rejected">
							{#snippet icon()}
								<UserRoundXIcon class="size-3.5 shrink-0 stroke-2" />
							{/snippet}
						</GroupStatusBadge>
						<div
							class="flex h-5 w-5 shrink-0 items-center justify-center self-start"
							aria-hidden="true"
						></div>
					</div>
				</div>
			{:else}
				<div
					class="flex w-full cursor-default items-start gap-3 rounded-md px-2 py-3 text-left transition-all"
				>
					<div class="flex w-full items-start gap-3 px-1">
						<span
							class="text-foreground-alt-1 min-w-0 flex-1 wrap-break-word text-base leading-5 font-medium"
						>
							{group.name}
						</span>
						<div class="flex h-5 w-5 shrink-0 items-center justify-center self-start">
							<GroupJoinButton
								groupId={group.id}
								onjoined={() => markJoinPending(group.id)}
								onerror={(message) => setJoinError(group.id, message)}
							/>
						</div>
					</div>
				</div>
			{/if}
			{#if joinErrors[group.id]}
				<p class="text-foreground-error px-3 text-xs">{joinErrors[group.id]}</p>
			{/if}
		</div>
		{#if index < groups.length - 1}
			<div class="flex w-full items-center justify-center">
				<Separator class="bg-foreground-alt-5 w-[95%]" />
			</div>
		{/if}
	{/each}
</div>

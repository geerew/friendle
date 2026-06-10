<script lang="ts">
	import { RightChevronIcon, UserRoundPlusIcon, UserRoundXIcon } from '$lib/components/icons';
	import GroupJoinButton from '../group-join-button.svelte';
	import GroupListShell from '../group-list-shell.svelte';
	import GroupListRow from '../group-list-row.svelte';
	import GroupListSeparator from '../group-list-separator.svelte';
	import GroupRoleBadge from '../group-role-badge.svelte';
	import GroupStatusBadge from '../group-status-badge.svelte';
	import type { GroupModel } from '$lib/models/group-model';
	import { isGroupMember } from '$lib/utils/group';

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

<GroupListShell>
	{#each groups as group, index (group.id)}
		{#if isGroupMember(group)}
			<GroupListRow name={group.name} href="/groups/{group.id}/" align="start" wrapName>
				{#snippet trailing()}
					<GroupRoleBadge role={group.groupRole ?? 'group_user'} />
					<div class="flex h-5 w-5 shrink-0 items-center justify-center self-start">
						<RightChevronIcon class="text-foreground-alt-2 size-5 shrink-0 stroke-2" />
					</div>
				{/snippet}
			</GroupListRow>
		{:else if group.joinRequestStatus === 'pending'}
			<GroupListRow name={group.name} align="start" wrapName>
				{#snippet trailing()}
					<GroupStatusBadge variant="pending" label="Pending">
						{#snippet icon()}
							<UserRoundPlusIcon class="size-3.5 shrink-0 stroke-2" />
						{/snippet}
					</GroupStatusBadge>
					<div class="flex h-5 w-5 shrink-0 items-center justify-center self-start" aria-hidden="true">
					</div>
				{/snippet}
			</GroupListRow>
		{:else if group.joinRequestStatus === 'rejected'}
			<GroupListRow name={group.name} align="start" wrapName>
				{#snippet trailing()}
					<GroupStatusBadge variant="rejected" label="Rejected">
						{#snippet icon()}
							<UserRoundXIcon class="size-3.5 shrink-0 stroke-2" />
						{/snippet}
					</GroupStatusBadge>
					<div class="flex h-5 w-5 shrink-0 items-center justify-center self-start" aria-hidden="true">
					</div>
				{/snippet}
			</GroupListRow>
		{:else}
			<GroupListRow name={group.name} align="start" wrapName>
				{#snippet trailing()}
					<div class="flex h-5 w-5 shrink-0 items-center justify-center self-start">
						<GroupJoinButton
							groupId={group.id}
							onjoined={() => markJoinPending(group.id)}
						/>
					</div>
				{/snippet}
			</GroupListRow>
		{/if}
		{#if index < groups.length - 1}
			<GroupListSeparator />
		{/if}
	{/each}
</GroupListShell>

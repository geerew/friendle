<script lang="ts">
	import { RightChevronIcon, UserRoundPlusIcon, UserRoundXIcon } from '$lib/components/icons';
	import { Table } from '$lib/components/ui';
	import GroupJoinButton from '../group-join-button.svelte';
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

<Table.List>
	{#each groups as group, index (group.id)}
		{#if isGroupMember(group)}
			<Table.Row label={group.name} href="/groups/{group.id}/" align="start" wrapLabel>
				{#snippet trailing()}
					<GroupRoleBadge role={group.groupRole ?? 'group_user'} />
					<div class="flex h-5 w-5 shrink-0 items-center justify-center self-start">
						<RightChevronIcon class="text-foreground-alt-2 size-5 shrink-0 stroke-2" />
					</div>
				{/snippet}
			</Table.Row>
		{:else if group.joinRequestStatus === 'pending'}
			<Table.Row label={group.name} align="start" wrapLabel>
				{#snippet trailing()}
					<GroupStatusBadge variant="pending" label="Pending">
						{#snippet icon()}
							<UserRoundPlusIcon class="size-3.5 shrink-0 stroke-2" />
						{/snippet}
					</GroupStatusBadge>
					<div
						class="flex h-5 w-5 shrink-0 items-center justify-center self-start"
						aria-hidden="true"
					></div>
				{/snippet}
			</Table.Row>
		{:else if group.joinRequestStatus === 'rejected'}
			<Table.Row label={group.name} align="start" wrapLabel>
				{#snippet trailing()}
					<GroupStatusBadge variant="rejected" label="Rejected">
						{#snippet icon()}
							<UserRoundXIcon class="size-3.5 shrink-0 stroke-2" />
						{/snippet}
					</GroupStatusBadge>
					<div
						class="flex h-5 w-5 shrink-0 items-center justify-center self-start"
						aria-hidden="true"
					></div>
				{/snippet}
			</Table.Row>
		{:else}
			<Table.Row label={group.name} align="start" wrapLabel>
				{#snippet trailing()}
					<div class="flex h-5 w-5 shrink-0 items-center justify-center self-start">
						<GroupJoinButton groupId={group.id} onjoined={() => markJoinPending(group.id)} />
					</div>
				{/snippet}
			</Table.Row>
		{/if}
		{#if index < groups.length - 1}
			<Table.Separator />
		{/if}
	{/each}
</Table.List>

<script lang="ts">
	import { RightChevronIcon, UserRoundPlusIcon, UserRoundXIcon } from '$lib/components/icons';
	import { Table } from '$lib/components/ui';
	import GroupJoinButton from '../group-join-button.svelte';
	import GroupRoleBadge from '../group-role-badge.svelte';
	import GroupStatusBadge from '../group-status-badge.svelte';
	import type { GroupModel } from '$lib/models/group-model';
	import { isGroupMember } from '$lib/utils/group';

	type Props = {
		group: GroupModel;
		onjoined?: (groupId: string) => void;
	};

	let { group, onjoined }: Props = $props();

	const isMember = $derived(isGroupMember(group));
	const isPending = $derived(group.joinRequestStatus === 'pending');
	const isRejected = $derived(group.joinRequestStatus === 'rejected');
</script>

{#if isMember}
	<Table.Row label={group.name} href="/groups/{group.id}/" wrapLabel>
		{#snippet trailing()}
			<div class="flex shrink-0 items-center gap-1.5">
				<GroupRoleBadge role={group.groupRole ?? 'group_user'} />
				<RightChevronIcon class="text-foreground-alt-2 size-5 shrink-0 stroke-2" />
			</div>
		{/snippet}
	</Table.Row>
{:else if isPending}
	<Table.Row label={group.name} wrapLabel>
		{#snippet trailing()}
			<div class="flex shrink-0 items-center gap-1.5">
				<GroupStatusBadge variant="pending" label="Pending">
					{#snippet icon()}
						<UserRoundPlusIcon class="size-3.5 shrink-0 stroke-2" />
					{/snippet}
				</GroupStatusBadge>
				<div class="size-5 shrink-0" aria-hidden="true"></div>
			</div>
		{/snippet}
	</Table.Row>
{:else if isRejected}
	<Table.Row label={group.name} wrapLabel>
		{#snippet trailing()}
			<div class="flex shrink-0 items-center gap-1.5">
				<GroupStatusBadge variant="rejected" label="Rejected">
					{#snippet icon()}
						<UserRoundXIcon class="size-3.5 shrink-0 stroke-2" />
					{/snippet}
				</GroupStatusBadge>
				<div class="size-5 shrink-0" aria-hidden="true"></div>
			</div>
		{/snippet}
	</Table.Row>
{:else}
	<Table.Row label={group.name} wrapLabel>
		{#snippet trailing()}
			<div class="flex shrink-0 items-center">
				<GroupJoinButton groupId={group.id} onjoined={() => onjoined?.(group.id)} />
			</div>
		{/snippet}
	</Table.Row>
{/if}

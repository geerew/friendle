<script lang="ts">
	import { requestGroupJoin } from '$lib/api/groups-api';
	import {
		PlusIcon,
		RightChevronIcon,
		UserRoundPlusIcon,
		UserRoundXIcon
	} from '$lib/components/icons';
	import { Button, StatusBadge, Table } from '$lib/components/ui';
	import GroupRoleBadge from '../group-role-badge.svelte';
	import type { GroupModel } from '$lib/models/group-model';
	import { apiErrorMessage, withMinLoadingDelay } from '$lib/utils';
	import { isGroupMember } from '$lib/utils/group';
	import { toast } from 'svelte-sonner';

	type Props = {
		group: GroupModel;
		onjoined?: (groupId: string) => void;
	};

	let { group, onjoined }: Props = $props();

	let joining = $state(false);

	const isMember = $derived(isGroupMember(group));
	const isPending = $derived(group.joinRequestStatus === 'pending');
	const isRejected = $derived(group.joinRequestStatus === 'rejected');

	// handleJoinClick submits a join request for the group
	async function handleJoinClick(): Promise<void> {
		if (joining) {
			return;
		}

		joining = true;

		try {
			await withMinLoadingDelay(requestGroupJoin(group.id));
			onjoined?.(group.id);
		} catch (err) {
			toast.error(apiErrorMessage(err, 'Failed to request join'));
		} finally {
			joining = false;
		}
	}
</script>

{#if isMember}
	<Table.Row label={group.name} href="/groups/{group.id}/" wrapLabel>
		{#snippet trailing()}
			<GroupRoleBadge role={group.groupRole ?? 'group_user'} />
			<RightChevronIcon class="text-foreground-alt-2 size-5 shrink-0 stroke-2" />
		{/snippet}
	</Table.Row>
{:else if isPending}
	<Table.Row label={group.name} wrapLabel>
		{#snippet trailing()}
			<StatusBadge variant="pending" label="Pending">
				{#snippet icon()}
					<UserRoundPlusIcon class="size-3.5 shrink-0 stroke-2" />
				{/snippet}
			</StatusBadge>
			<div class="size-5 shrink-0" aria-hidden="true"></div>
		{/snippet}
	</Table.Row>
{:else if isRejected}
	<Table.Row label={group.name} wrapLabel>
		{#snippet trailing()}
			<StatusBadge variant="rejected" label="Rejected">
				{#snippet icon()}
					<UserRoundXIcon class="size-3.5 shrink-0 stroke-2" />
				{/snippet}
			</StatusBadge>
			<div class="size-5 shrink-0" aria-hidden="true"></div>
		{/snippet}
	</Table.Row>
{:else}
	<Table.Row label={group.name} wrapLabel>
		{#snippet trailing()}
			<Button
				type="button"
				variant="ghost"
				size="inline"
				class="text-foreground-alt-2 h-5 w-5 min-h-5 min-w-5 shrink-0 p-0 normal-case hover:bg-transparent"
				aria-label="Request to join group"
				loading={joining}
				onclick={handleJoinClick}
			>
				<PlusIcon class="size-5 shrink-0 stroke-2" />
			</Button>
		{/snippet}
	</Table.Row>
{/if}

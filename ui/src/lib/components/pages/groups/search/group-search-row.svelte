<script lang="ts">
	import { cancelGroupJoinRequest, getGroup, requestGroupJoin } from '$lib/api/groups-api';
	import {
		PlusIcon,
		RightChevronIcon,
		UserRoundPlusIcon,
		UserRoundXIcon,
		XIcon
	} from '$lib/components/icons';
	import { Button, StatusBadge, Table } from '$lib/components/ui';
	import GroupRoleBadge from '../group-role-badge.svelte';
	import type { GroupModel } from '$lib/models/group-model';
	import { apiErrorMessage, isJoinRequestNotFound, withMinLoadingDelay } from '$lib/utils';
	import { isGroupMember, joinRequestResolvedNotice } from '$lib/utils/group';
	import { toast } from 'svelte-sonner';

	type Props = {
		group: GroupModel;
		onGroupChange?: (group: GroupModel) => void;
	};

	let { group, onGroupChange }: Props = $props();

	let joining = $state(false);
	let cancelling = $state(false);

	const isMember = $derived(isGroupMember(group));
	const isPending = $derived(group.joinRequestStatus === 'pending');
	const isRejected = $derived(group.joinRequestStatus === 'rejected');

	// showJoinRequestResolvedNotice toasts and applies the viewer's current group status
	function showJoinRequestResolvedNotice(updated: GroupModel): void {
		onGroupChange?.(updated);

		const notice = joinRequestResolvedNotice(updated);

		if (notice.variant === 'success') {
			toast.success(notice.message);
		} else {
			toast.error(notice.message);
		}
	}

	// handleJoinClick submits a join request for the group
	async function handleJoinClick(): Promise<void> {
		if (joining) {
			return;
		}

		joining = true;

		try {
			const updated = await withMinLoadingDelay(requestGroupJoin(group.id));
			onGroupChange?.(updated);
		} catch (err) {
			toast.error(apiErrorMessage(err, 'Failed to request join'));
		} finally {
			joining = false;
		}
	}

	// handleCancelClick withdraws the caller's pending join request
	async function handleCancelClick(): Promise<void> {
		if (cancelling) {
			return;
		}

		cancelling = true;

		try {
			const updated = await withMinLoadingDelay(cancelGroupJoinRequest(group.id));
			onGroupChange?.(updated);
		} catch (err) {
			if (isJoinRequestNotFound(err)) {
				try {
					const updated = await getGroup(group.id);
					showJoinRequestResolvedNotice(updated);
				} catch {
					toast.error('Join request is no longer pending');
				}
			} else {
				toast.error(apiErrorMessage(err, 'Failed to cancel join request'));
			}
		} finally {
			cancelling = false;
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
			<Button
				type="button"
				variant="ghost"
				size="inline"
				class="text-foreground-alt-2 hover:bg-background-error hover:text-foreground h-5 min-h-5 w-5 min-w-5 shrink-0 p-0 normal-case"
				aria-label="Cancel join request"
				loading={cancelling}
				onclick={handleCancelClick}
			>
				<XIcon class="size-5 shrink-0 stroke-2" />
			</Button>
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
				class="text-foreground-alt-2 h-5 min-h-5 w-5 min-w-5 shrink-0 p-0 normal-case hover:bg-transparent"
				aria-label="Request to join group"
				loading={joining}
				onclick={handleJoinClick}
			>
				<PlusIcon class="size-5 shrink-0 stroke-2" />
			</Button>
		{/snippet}
	</Table.Row>
{/if}

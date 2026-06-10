<script lang="ts">
	import {
		approveGroupJoinRequest,
		declineGroupJoinRequest
	} from '$lib/api/groups-api';
	import { TickIcon, XIcon } from '$lib/components/icons';
	import { Button, Table } from '$lib/components/ui';
	import type { GroupJoinRequestModel } from '$lib/models/group-join-request-model';
	import { apiErrorMessage, withMinLoadingDelay } from '$lib/utils';
	import { toast } from 'svelte-sonner';

	type Props = {
		groupId: string;
		requests: GroupJoinRequestModel[];
		onchange?: () => void;
	};

	let { groupId, requests, onchange }: Props = $props();

	let acting = $state<{ userId: string; action: 'approve' | 'decline' } | null>(null);

	// handleApprove approves a pending join request
	async function handleApprove(userId: string): Promise<void> {
		if (acting) {
			return;
		}

		acting = { userId, action: 'approve' };

		try {
			await withMinLoadingDelay(approveGroupJoinRequest(groupId, userId));
			onchange?.();
		} catch (err) {
			toast.error(apiErrorMessage(err, 'Failed to approve request'));
		} finally {
			acting = null;
		}
	}

	// handleDecline rejects a pending join request
	async function handleDecline(userId: string): Promise<void> {
		if (acting) {
			return;
		}

		acting = { userId, action: 'decline' };

		try {
			await withMinLoadingDelay(declineGroupJoinRequest(groupId, userId));
			onchange?.();
		} catch (err) {
			toast.error(apiErrorMessage(err, 'Failed to decline request'));
		} finally {
			acting = null;
		}
	}
</script>

<Table.List>
	{#each requests as request, index (request.userId)}
		<Table.Row label={request.displayName}>
			{#snippet trailing()}
				<div class="flex shrink-0 items-center gap-1.5">
					<Button
						type="button"
						variant="ghost"
						size="icon"
						class="text-foreground-alt-2 hover:bg-background-primary/25 hover:text-background-primary size-8 min-h-8 min-w-8 shrink-0 p-0 normal-case"
						loading={acting?.userId === request.userId && acting.action === 'approve'}
						disabled={acting != null}
						aria-label="Approve {request.displayName}"
						onclick={() => handleApprove(request.userId)}
					>
						<TickIcon class="size-4 shrink-0 stroke-2" />
					</Button>
					<Button
						type="button"
						variant="ghost"
						size="icon"
						class="text-foreground-alt-2 hover:bg-background-error hover:text-foreground size-8 min-h-8 min-w-8 shrink-0 p-0 normal-case"
						loading={acting?.userId === request.userId && acting.action === 'decline'}
						disabled={acting != null}
						aria-label="Decline {request.displayName}"
						onclick={() => handleDecline(request.userId)}
					>
						<XIcon class="size-4 shrink-0 stroke-2" />
					</Button>
				</div>
			{/snippet}
		</Table.Row>
		{#if index < requests.length - 1}
			<Table.Separator />
		{/if}
	{/each}
</Table.List>

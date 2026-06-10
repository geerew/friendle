<script lang="ts">
	import {
		approveGroupJoinRequest,
		declineGroupJoinRequest
	} from '$lib/api/groups-api';
	import { TickIcon, XIcon } from '$lib/components/icons';
	import { Button, Separator } from '$lib/components/ui';
	import type { GroupJoinRequestModel } from '$lib/models/group-join-request-model';
	import { apiErrorMessage, withMinLoadingDelay } from '$lib/utils';

	type Props = {
		groupId: string;
		requests: GroupJoinRequestModel[];
		onchange?: () => void;
	};

	let { groupId, requests, onchange }: Props = $props();

	let acting = $state<{ userId: string; action: 'approve' | 'decline' } | null>(null);
	let rowErrors = $state<Record<string, string>>({});

	// handleApprove approves a pending join request
	async function handleApprove(userId: string): Promise<void> {
		if (acting) {
			return;
		}

		acting = { userId, action: 'approve' };
		rowErrors = { ...rowErrors, [userId]: '' };

		try {
			await withMinLoadingDelay(approveGroupJoinRequest(groupId, userId));
			onchange?.();
		} catch (err) {
			rowErrors = {
				...rowErrors,
				[userId]: apiErrorMessage(err, 'Failed to approve request')
			};
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
		rowErrors = { ...rowErrors, [userId]: '' };

		try {
			await withMinLoadingDelay(declineGroupJoinRequest(groupId, userId));
			onchange?.();
		} catch (err) {
			rowErrors = {
				...rowErrors,
				[userId]: apiErrorMessage(err, 'Failed to decline request')
			};
		} finally {
			acting = null;
		}
	}
</script>

<div class="flex flex-col gap-2">
	{#each requests as request, index (request.userId)}
		<div class="flex flex-col gap-1">
			<div class="flex w-full items-center gap-3 rounded-md px-2 py-3 text-left">
				<div class="flex w-full items-center gap-3 px-1">
					<span
						class="text-foreground-alt-1 min-w-0 flex-1 truncate text-base leading-5 font-medium"
					>
						{request.displayName}
					</span>
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
				</div>
			</div>
			{#if rowErrors[request.userId]}
				<p class="text-foreground-error px-3 text-xs">{rowErrors[request.userId]}</p>
			{/if}
		</div>
		{#if index < requests.length - 1}
			<div class="flex w-full items-center justify-center">
				<Separator class="bg-foreground-alt-5 w-[95%]" />
			</div>
		{/if}
	{/each}
</div>

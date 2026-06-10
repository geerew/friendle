<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import {
		approveGroupJoinRequest,
		declineGroupJoinRequest,
		listGroupPendingJoinRequests
	} from '$lib/api/groups-api';
	import { TickIcon, XIcon } from '$lib/components/icons';
	import { GroupNameSection } from '$lib/components/pages';
	import { Button, Table } from '$lib/components/ui';
	import { GROUP_PAGE_KEY, type GroupPageContext } from '$lib/context/group-page';
	import type { GroupJoinRequestModel } from '$lib/models/group-join-request-model';
	import { apiErrorMessage, withMinLoadingDelay } from '$lib/utils';
	import { groupChildBreadcrumb, isGroupAdmin } from '$lib/utils/group';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	const groupId = $derived(page.params.id ?? '');
	const groupPage = getContext<GroupPageContext>(GROUP_PAGE_KEY);
	const group = $derived(groupPage.group);
	const breadcrumb = $derived(
		group ? groupChildBreadcrumb(group.id, 'Pending') : [{ label: 'Pending' }]
	);

	let requests = $state<GroupJoinRequestModel[]>([]);
	let pageNum = $state(1);
	let perPage = $state(10);
	let totalItems = $state(0);
	let loading = $state(true);
	let acting = $state<{ userId: string; action: 'approve' | 'decline' } | null>(null);

	$effect(() => {
		groupId;
		group;

		if (group && !isGroupAdmin(group)) {
			void goto(`/groups/${groupId}/`);
			return;
		}

		pageNum;
		perPage;
		void loadPending();
	});

	async function loadPending(): Promise<void> {
		if (!groupId || !group || !isGroupAdmin(group)) {
			return;
		}

		loading = true;

		try {
			const data = await withMinLoadingDelay(
				listGroupPendingJoinRequests(groupId, { page: pageNum, perPage })
			);
			requests = data.items;
			totalItems = data.totalItems;
		} catch (err) {
			requests = [];
			totalItems = 0;
			toast.error(apiErrorMessage(err, 'Failed to load pending requests'));
		} finally {
			loading = false;
		}
	}

	// handlePendingChange reloads the list and refreshes group stats
	function handlePendingChange(): void {
		void loadPending();
		void groupPage.reloadGroup({ silent: true });
	}

	// handleApprove approves a pending join request
	async function handleApprove(userId: string): Promise<void> {
		if (acting) {
			return;
		}

		acting = { userId, action: 'approve' };

		try {
			await withMinLoadingDelay(approveGroupJoinRequest(groupId, userId));
			handlePendingChange();
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
			handlePendingChange();
		} catch (err) {
			toast.error(apiErrorMessage(err, 'Failed to decline request'));
		} finally {
			acting = null;
		}
	}
</script>

<Table.Root title="Pending Requests" {breadcrumb}>
	{#snippet header()}
		<GroupNameSection />
	{/snippet}

	<Table.PaginatedBody
		itemCount={requests.length}
		{totalItems}
		bind:page={pageNum}
		bind:perPage
		{loading}
		emptyMessage="No pending requests"
	>
		{#snippet list()}
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
		{/snippet}
	</Table.PaginatedBody>
</Table.Root>

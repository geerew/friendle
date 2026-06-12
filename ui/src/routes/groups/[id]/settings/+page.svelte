<script lang="ts">
	import { goto } from '$app/navigation';
	import { leaveGroup } from '$lib/api/groups-api';
	import { DeleteGroup } from '$lib/components';
	import { GroupNameSection, GroupStatLink } from '$lib/components/pages';
	import { Button, DestroyDialog, Separator, Table } from '$lib/components/ui';
	import { GROUP_PAGE_KEY, type GroupPageContext } from '$lib/context/group-page';
	import { apiErrorMessage, withMinLoadingDelay } from '$lib/utils';
	import { groupChildBreadcrumb, isGroupAdmin } from '$lib/utils/group';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	const groupPage = getContext<GroupPageContext>(GROUP_PAGE_KEY);
	const group = $derived(groupPage.group);
	const isAdmin = $derived(group != null && isGroupAdmin(group));
	const breadcrumb = $derived(
		group ? groupChildBreadcrumb(group.id, 'Settings') : [{ label: 'Settings' }]
	);

	let deleteGroupOpen = $state(false);
	let leaveConfirmOpen = $state(false);
	let leaving = $state(false);

	async function handleDeleteSuccess(): Promise<void> {
		deleteGroupOpen = false;
		await goto('/');
	}

	function handleDeleteError(message: string): void {
		toast.error(message);
	}

	async function confirmLeaveGroup(): Promise<void> {
		if (!group) {
			return;
		}

		leaving = true;

		try {
			await withMinLoadingDelay(leaveGroup(group.id));
			leaveConfirmOpen = false;
			await goto('/');
		} catch (err) {
			toast.error(apiErrorMessage(err, 'Failed to leave group'));
		} finally {
			leaving = false;
		}
	}
</script>

{#if group}
	<Table.Root title="Settings" {breadcrumb} showTitle={false}>
		<div class="flex flex-col gap-5">
			<GroupNameSection />

			{#if isAdmin && group.adminSummary}
				<section class="flex items-stretch justify-between gap-3">
					<GroupStatLink
						label="Members"
						count={group.memberCount}
						href="/groups/{group.id}/settings/members"
						ariaLabel="View members"
						align="start"
					/>
					<Separator class="h-auto w-px shrink-0 self-stretch" />
					<GroupStatLink
						label="Pending"
						count={group.adminSummary.pendingJoinRequestCount}
						href="/groups/{group.id}/settings/pending"
						ariaLabel="View pending join requests"
					/>
					<Separator class="h-auto w-px shrink-0 self-stretch" />
					<GroupStatLink
						label="Rejected"
						count={group.adminSummary.rejectedJoinRequestCount}
						href="/groups/{group.id}/settings/rejected"
						ariaLabel="View rejected join requests"
						align="end"
					/>
				</section>
			{:else}
				<section class="flex flex-col gap-3">
					<GroupStatLink
						label="Members"
						count={group.memberCount}
						href="/groups/{group.id}/settings/members"
						ariaLabel="View members"
						align="start"
					/>
				</section>
			{/if}

			<Separator />

			<section class="flex flex-col gap-3">
				<div class="flex items-center justify-between gap-3">
					<h2 class="section-title">Leave group</h2>
					<Button
						type="button"
						variant="destructive"
						size="inline"
						onclick={() => (leaveConfirmOpen = true)}
					>
						Leave
					</Button>
				</div>
				<p class="text-foreground-alt-2 text-sm">
					Leaving this group will permanently delete all your group member data
				</p>
			</section>

			{#if isAdmin}
				<Separator />

				<section class="flex flex-col gap-3">
					<div class="flex items-center justify-between gap-3">
						<h2 class="section-title section-title-destructive">Delete group</h2>
						<Button
							type="button"
							variant="destructive"
							size="inline"
							onclick={() => (deleteGroupOpen = true)}
						>
							Delete
						</Button>
					</div>
					<p class="text-foreground-alt-2 text-sm">Permanently delete this Friendle group</p>
				</section>
			{/if}
		</div>
	</Table.Root>

	<DeleteGroup
		bind:open={deleteGroupOpen}
		{group}
		onSuccess={handleDeleteSuccess}
		onError={handleDeleteError}
	/>

	<DestroyDialog
		bind:open={leaveConfirmOpen}
		title="Are you sure you want to leave this group?"
		detail={group.name}
		description="All your group member data will be permanently deleted"
		confirmLabel="Leave"
		loading={leaving}
		onConfirm={confirmLeaveGroup}
	/>
{/if}

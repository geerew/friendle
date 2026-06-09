<script lang="ts">
	import { page } from '$app/state';
	import { ApiError } from '$lib/api';
	import { getGroup } from '$lib/api/groups-api';
	import { AppShell, Spinner } from '$lib/components';
	import { GroupJoinButton } from '$lib/components/pages';
	import { Separator } from '$lib/components/ui';
	import type { GroupModel } from '$lib/models/group-model';
	import { toast } from 'svelte-sonner';

	const groupId = $derived(page.params.id ?? '');

	let group = $state<GroupModel | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);

	const isGroupAdmin = $derived(group?.groupRole === 'group_admin');
	const isGroupMember = $derived(group?.groupRole === 'group_user');
	const canRequestJoin = $derived(
		group != null && !isGroupAdmin && !isGroupMember && group.joinRequestStatus == null
	);
	const joinRequestMessage = $derived.by(() => {
		if (group?.joinRequestStatus === 'pending') {
			return 'Your join request is pending';
		}

		if (group?.joinRequestStatus === 'rejected') {
			return 'Your join request was rejected';
		}

		return null;
	});

	function markJoinPending(): void {
		if (!group) {
			return;
		}

		group = { ...group, joinRequestStatus: 'pending' };
	}

	$effect(() => {
		groupId;
		void loadGroup();
	});

	async function loadGroup(): Promise<void> {
		if (!groupId) {
			group = null;
			error = 'Group not found';
			loading = false;
			return;
		}

		loading = true;
		error = null;

		try {
			group = await getGroup(groupId);
		} catch (err) {
			group = null;
			error = err instanceof ApiError ? err.message : 'Failed to load group';
		} finally {
			loading = false;
		}
	}
</script>

<AppShell breadcrumb={[{ label: 'Group' }]}>
	{#if loading}
		<div class="flex min-h-24 items-center justify-center">
			<Spinner class="bg-foreground-alt-2 size-3" />
		</div>
	{:else if error}
		<p class="text-foreground-error-alt-1 text-sm">{error}</p>
	{:else if group}
		<div class="flex flex-col gap-5">
			<section class="flex flex-col gap-3">
				<h2 class="section-title">Group name</h2>
				<p class="text-background-primary text-2xl">{group.name}</p>
			</section>

			{#if isGroupAdmin && group.adminSummary}
				<Separator />

				<section class="grid grid-cols-3 gap-3">
					<div class="flex flex-col items-center gap-3 text-center">
						<h2 class="section-title">Members</h2>
						<p class="text-background-primary text-2xl tabular-nums">{group.memberCount}</p>
					</div>
					<div class="flex flex-col items-center gap-3 text-center">
						<h2 class="section-title">Pending</h2>
						<p class="text-background-primary text-2xl tabular-nums">
							{group.adminSummary.pendingJoinRequestCount}
						</p>
					</div>
					<div class="flex flex-col items-center gap-3 text-center">
						<h2 class="section-title">Rejected</h2>
						<p class="text-background-primary text-2xl tabular-nums">
							{group.adminSummary.rejectedJoinRequestCount}
						</p>
					</div>
				</section>
			{:else if isGroupMember}
				<Separator />

				<section class="flex flex-col gap-3">
					<h2 class="section-title">Members</h2>
					<p class="text-background-primary text-2xl tabular-nums">{group.memberCount}</p>
				</section>
			{:else}
				<Separator />

				<section class="flex flex-col gap-3">
					<p class="text-foreground text-sm">You are not a member of this group</p>
					{#if joinRequestMessage}
						<p class="text-foreground-alt-2 text-sm">{joinRequestMessage}</p>
					{/if}
					{#if canRequestJoin}
						<GroupJoinButton
							groupId={group.id}
							appearance="button"
							onjoined={markJoinPending}
							onerror={(message) => toast.error(message)}
						/>
					{/if}
				</section>
			{/if}
		</div>
	{/if}
</AppShell>

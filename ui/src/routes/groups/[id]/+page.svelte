<script lang="ts">
	import { page } from '$app/state';
	import { ApiError } from '$lib/api';
	import { getGroup } from '$lib/api/groups-api';
	import { AppShell, ListRow } from '$lib/components';
	import { Button } from '$lib/components/ui';
	import type { GroupDetail } from '$lib/types/group';
	import { formatMemberCount } from '$lib/utils';

	const groupId = $derived(page.params.id ?? '');

	let group = $state<GroupDetail | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);

	$effect(() => {
		loadGroup();
	});

	async function loadGroup(): Promise<void> {
		loading = true;
		error = null;

		try {
			group = await getGroup(groupId);
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to load group';
		} finally {
			loading = false;
		}
	}

	const round = $derived(group?.currentRound ?? null);
</script>

<AppShell title={group?.name ?? 'Group'}>
	{#if loading}
		<p class="text-text-muted">Loading…</p>
	{:else if error}
		<p class="text-error">{error}</p>
	{:else if group}
		<div class="rounded border border-border bg-bg-secondary px-4 py-3 text-sm text-text-muted">
			<p>{formatMemberCount(group.memberCount ?? group.members?.length ?? 0)}</p>
			<p>Round every {group.intervalHours}h ({group.timezone})</p>
			{#if round}
				<p class="mt-2 capitalize">Today: {round.status.replace('_', ' ')}</p>
				{#if round.pickerUsername}
					<p>Picker: {round.pickerUsername}</p>
				{/if}
			{/if}
		</div>

		<div class="flex flex-col gap-3">
			{#if round?.status === 'awaiting_word' && round.isPicker}
				<Button href="/groups/{groupId}/pick/" variant="primary">Pick today's word</Button>
			{/if}

			{#if round?.status === 'active'}
				<Button href="/groups/{groupId}/play/" variant="primary">Play today's word</Button>
			{/if}

			<Button href="/groups/{groupId}/leaderboard/" variant="secondary">Leaderboard</Button>
			<Button href="/groups/{groupId}/settings/" variant="secondary">Group settings</Button>
		</div>

		{#if group.members && group.members.length > 0}
			<section class="flex flex-col gap-2">
				<h2 class="section-title">Members</h2>
				{#each group.members as member (member.id)}
					<ListRow
						title={member.displayName}
						subtitle="{member.username} · {member.groupRole.replace('_', ' ')}"
					/>
				{/each}
			</section>
		{/if}
	{/if}
</AppShell>

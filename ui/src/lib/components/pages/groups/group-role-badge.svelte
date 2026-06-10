<script lang="ts">
	import { UserRoundIcon, UserStarIcon } from '$lib/components/icons';
	import GroupStatusBadge from './group-status-badge.svelte';

	type Props = {
		role: 'group_admin' | 'group_user';
	};

	let { role }: Props = $props();
</script>

{#if role === 'group_admin'}
	<GroupStatusBadge variant="admin" label="Admin">
		{#snippet icon()}
			<UserStarIcon class="size-3.5 shrink-0 stroke-2" />
		{/snippet}
	</GroupStatusBadge>
{:else}
	<GroupStatusBadge variant="member" label="Member">
		{#snippet icon()}
			<UserRoundIcon class="size-3.5 shrink-0 stroke-2" />
		{/snippet}
	</GroupStatusBadge>
{/if}

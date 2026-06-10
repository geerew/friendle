<script lang="ts">
	import { UserRoundIcon, UserStarIcon } from '$lib/components/icons';
	import { StatusBadge } from '$lib/components/ui';

	type Props = {
		role: 'group_admin' | 'group_user';
	};

	let { role }: Props = $props();
</script>

{#if role === 'group_admin'}
	<StatusBadge variant="admin" label="Admin">
		{#snippet icon()}
			<UserStarIcon class="size-3.5 shrink-0 stroke-2" />
		{/snippet}
	</StatusBadge>
{:else}
	<StatusBadge variant="member" label="Member">
		{#snippet icon()}
			<UserRoundIcon class="size-3.5 shrink-0 stroke-2" />
		{/snippet}
	</StatusBadge>
{/if}

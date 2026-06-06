<script lang="ts">
	import { ShieldUserIcon } from '$lib/components/icons';
	import { Button, Separator } from '$lib/components/ui';
	import type { AdminUserModel } from '$lib/models/admin-user-model';

	type Props = {
		users: AdminUserModel[];
		onDelete: (user: AdminUserModel) => void;
	};

	let { users, onDelete }: Props = $props();

	const rowClass = 'flex w-full items-center gap-3 px-2 py-3';
</script>

<div class="flex flex-col gap-2">
	{#each users as user, index (user.id)}
		<div class={rowClass}>
			<span class="text-foreground-alt-1 min-w-0 flex-1 truncate px-1 text-base font-medium">
				{user.username}
			</span>
			{#if user.siteRole === 'site_admin'}
				<ShieldUserIcon
					class="text-background-primary size-5 shrink-0 stroke-2"
					aria-label="Admin"
				/>
			{/if}
			<Button variant="destructive" size="inline" onclick={() => onDelete(user)}>Delete</Button>
		</div>

		{#if index < users.length - 1}
			<div class="flex w-full items-center justify-center">
				<Separator class="bg-foreground-alt-5 w-[95%]" />
			</div>
		{/if}
	{/each}
</div>

<script lang="ts">
	import { ShieldUserIcon } from '$lib/components/icons';
	import { Button, Table } from '$lib/components/ui';
	import type { AdminUserModel } from '$lib/models/admin-user-model';

	type Props = {
		users: AdminUserModel[];
		onDelete: (user: AdminUserModel) => void;
	};

	let { users, onDelete }: Props = $props();
</script>

<Table.List>
	{#each users as user, index (user.id)}
		<Table.Row label={user.username}>
			{#snippet trailing()}
				{#if user.siteRole === 'site_admin'}
					<ShieldUserIcon
						class="text-background-primary size-5 shrink-0 stroke-2"
						aria-label="Admin"
					/>
				{/if}
				<Button variant="destructive" size="inline" onclick={() => onDelete(user)}>Delete</Button>
			{/snippet}
		</Table.Row>
		{#if index < users.length - 1}
			<Table.Separator />
		{/if}
	{/each}
</Table.List>

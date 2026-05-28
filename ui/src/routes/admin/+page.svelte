<script lang="ts">
	import { goto } from '$app/navigation';
	import { ApiError } from '$lib/api';
	import { deleteAdminGroup, deleteUser, listGroups, listUsers } from '$lib/api/admin-api';
	import { auth } from '$lib/auth.svelte';
	import { AppShell, ListRow } from '$lib/components';
	import { Button } from '$lib/components/ui';
	import type { AdminGroup, AdminUser } from '$lib/types/admin';

	let users = $state<AdminUser[]>([]);
	let groups = $state<AdminGroup[]>([]);
	let loading = $state(true);
	let error = $state<string | null>(null);

	$effect(() => {
		if (!auth.isAdmin) {
			void goto('/');
			return;
		}

		loadAdmin();
	});

	async function loadAdmin(): Promise<void> {
		loading = true;
		error = null;

		try {
			const [userData, groupData] = await Promise.all([listUsers(), listGroups()]);
			users = userData.items ?? [];
			groups = groupData.items ?? [];
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to load admin data';
		} finally {
			loading = false;
		}
	}

	async function handleDeleteUser(user: AdminUser): Promise<void> {
		if (!confirm(`Delete user ${user.username}?`)) return;

		try {
			await deleteUser(user.id);
			await loadAdmin();
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to delete user';
		}
	}

	async function handleDeleteGroup(group: AdminGroup): Promise<void> {
		if (!confirm(`Delete group ${group.name}?`)) return;

		try {
			await deleteAdminGroup(group.id);
			await loadAdmin();
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to delete group';
		}
	}
</script>

<AppShell title="Admin" showHome={true}>
	{#if loading}
		<p class="text-text-muted">Loading…</p>
	{:else}
		<section class="flex flex-col gap-2">
			<h2 class="section-title">Users ({users.length})</h2>
			{#each users as user (user.id)}
				<ListRow title={user.displayName} subtitle="{user.username} · {user.role}">
					{#snippet trailing()}
						<Button variant="destructive" size="inline" onclick={() => handleDeleteUser(user)}>
							Delete
						</Button>
					{/snippet}
				</ListRow>
			{/each}
		</section>

		<section class="flex flex-col gap-2">
			<h2 class="section-title">Groups ({groups.length})</h2>
			{#each groups as group (group.id)}
				<ListRow title={group.name} subtitle="{group.memberCount} members">
					{#snippet trailing()}
						<Button variant="destructive" size="inline" onclick={() => handleDeleteGroup(group)}>
							Delete
						</Button>
					{/snippet}
				</ListRow>
			{/each}
		</section>

		{#if error}
			<p class="text-sm text-error">{error}</p>
		{/if}

		<Button variant="secondary" onclick={loadAdmin}>Refresh</Button>
	{/if}
</AppShell>

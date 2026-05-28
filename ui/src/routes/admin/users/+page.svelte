<script lang="ts">
	import { ApiError } from '$lib/api';
	import { deleteUser, listUsers } from '$lib/api/admin-api';
	import { AppShell, ListRow, Pagination } from '$lib/components';
	import { Button } from '$lib/components/ui';
	import type { AdminUserModel } from '$lib/models/admin-user-model';

	let users = $state<AdminUserModel[]>([]);
	let page = $state(1);
	let perPage = $state(25);
	let totalItems = $state(0);
	let loading = $state(true);
	let error = $state<string | null>(null);

	$effect(() => {
		page;
		perPage;
		void loadUsers();
	});

	async function loadUsers(): Promise<void> {
		loading = true;
		error = null;

		try {
			const data = await listUsers({ page, perPage });
			users = data.items;
			totalItems = data.totalItems;
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to load users';
		} finally {
			loading = false;
		}
	}

	async function handleDeleteUser(user: AdminUserModel): Promise<void> {
		if (!confirm(`Delete user ${user.username}?`)) return;

		try {
			await deleteUser(user.id);
			const remainingTotal = totalItems - 1;
			const totalPages = Math.max(1, Math.ceil(remainingTotal / perPage));

			if (page > totalPages) {
				page = totalPages;
			} else {
				await loadUsers();
			}
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to delete user';
		}
	}
</script>

<AppShell
	breadcrumb={[
		{ label: 'Admin', href: '/admin/' },
		{ label: 'Users' }
	]}
>
	<div class="flex flex-col gap-6">
		<Button href="/admin/users/add/" variant="primary" class="w-1/2">+ Add User</Button>

		<hr class="border-0 border-t border-border" />

		{#if loading}
			<p class="text-text-muted">Loading…</p>
		{:else}
			{#if error}
				<p class="text-sm text-error">{error}</p>
			{/if}

			<div class="flex flex-col gap-2">
				{#if users.length === 0}
					<p class="text-sm text-text-muted">No users.</p>
				{:else}
					{#each users as user (user.id)}
						<ListRow title={user.displayName} subtitle="{user.username} · {user.siteRole}">
							{#snippet trailing()}
								<Button variant="destructive" size="inline" onclick={() => handleDeleteUser(user)}>
									Delete
								</Button>
							{/snippet}
						</ListRow>
					{/each}
				{/if}
			</div>

			<Pagination
				count={totalItems}
				bind:page
				bind:perPage
				selectTriggerClass="h-9 px-2 py-0"
				onPageChange={() => {}}
				onPerPageChange={() => {}}
			/>
		{/if}
	</div>
</AppShell>

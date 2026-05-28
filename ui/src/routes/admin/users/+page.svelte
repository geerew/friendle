<script lang="ts">
	import { ApiError } from '$lib/api';
	import { listUsers } from '$lib/api/admin-api';
	import { AppShell, DeleteUser, ListRow, Pagination } from '$lib/components';
	import { Button } from '$lib/components/ui';
	import { formatSiteRole, type AdminUserModel } from '$lib/models/admin-user-model';

	let users = $state<AdminUserModel[]>([]);
	let page = $state(1);
	let perPage = $state(25);
	let totalItems = $state(0);
	let loading = $state(true);
	let error = $state<string | null>(null);
	let deleteOpen = $state(false);
	let userToDelete = $state<AdminUserModel | null>(null);

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

	function openDeleteUser(user: AdminUserModel): void {
		error = null;
		userToDelete = user;
		deleteOpen = true;
	}

	async function handleDeleteSuccess(): Promise<void> {
		const remainingTotal = totalItems - 1;
		const totalPages = Math.max(1, Math.ceil(remainingTotal / perPage));

		if (page > totalPages) {
			page = totalPages;
		} else {
			await loadUsers();
		}
	}

	function handleDeleteError(message: string): void {
		error = message;
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

			<div class="flex flex-col gap-3">
				{#if users.length === 0}
					<p class="text-sm text-text-muted">No users.</p>
				{:else}
					{#each users as user (user.id)}
						<ListRow
							title={user.displayName}
							subtitle="{user.username} · {formatSiteRole(user.siteRole)}"
						>
							{#snippet trailing()}
								<Button variant="destructive" size="inline" onclick={() => openDeleteUser(user)}>
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

	<DeleteUser
		bind:open={deleteOpen}
		user={userToDelete}
		onSuccess={handleDeleteSuccess}
		onError={handleDeleteError}
	/>
</AppShell>

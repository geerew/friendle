<script lang="ts">
	import { ApiError } from '$lib/api';
	import { listUsers } from '$lib/api/admin-api';
	import { DeleteUser, Pagination, Spinner } from '$lib/components';
	import { AdminUserList } from '$lib/components/pages';
	import { Button, Separator, Table } from '$lib/components/ui';
	import type { AdminUserModel } from '$lib/models/admin-user-model';
	import { withMinLoadingDelay } from '$lib/utils';

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
			const data = await withMinLoadingDelay(listUsers({ page, perPage }));
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

<Table.Root
	title="Users"
	breadcrumb={[{ label: 'Admin', href: '/admin/' }, { label: 'Users' }]}
>
	<div class="flex flex-col gap-6">
		<Button href="/admin/users/add/" variant="primary" class="w-auto self-start px-6"
			>Add User</Button
		>

		<Separator />

		{#if loading}
			<div class="flex min-h-24 items-center justify-center">
				<Spinner class="bg-foreground-alt-2 size-3" />
			</div>
		{:else}
			{#if error}
				<p class="text-foreground-error text-sm">{error}</p>
			{/if}

			{#if users.length === 0}
				<p class="text-foreground-alt-2 text-sm italic">No users</p>
			{:else}
				<AdminUserList {users} onDelete={openDeleteUser} />
			{/if}

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
</Table.Root>

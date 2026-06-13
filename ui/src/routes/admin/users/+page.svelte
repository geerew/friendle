<script lang="ts">
	import { listUsers } from '$lib/api/admin-api';
	import { DeleteUser } from '$lib/components';
	import { ShieldUserIcon } from '$lib/components/icons';
	import { Button, Separator, Table } from '$lib/components/ui';
	import type { AdminUserModel } from '$lib/models/admin-user-model';
	import { apiErrorMessage, withMinLoadingDelay } from '$lib/utils';
	import { toast } from 'svelte-sonner';

	let users = $state<AdminUserModel[]>([]);
	let page = $state(1);
	let perPage = $state(25);
	let totalItems = $state(0);
	let loading = $state(true);
	let deleteOpen = $state(false);
	let userToDelete = $state<AdminUserModel | null>(null);

	$effect(() => {
		page;
		perPage;
		void loadUsers();
	});

	async function loadUsers(): Promise<void> {
		loading = true;

		try {
			const data = await withMinLoadingDelay(listUsers({ page, perPage }));
			users = data.items;
			totalItems = data.totalItems;
		} catch (err) {
			users = [];
			totalItems = 0;
			toast.error(apiErrorMessage(err, 'Failed to load users'));
		} finally {
			loading = false;
		}
	}

	function openDeleteUser(user: AdminUserModel): void {
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
		toast.error(message);
	}
</script>

<Table.Root title="Users" breadcrumb={[{ label: 'Admin', href: '/admin/' }, { label: 'Users' }]}>
	{#snippet header()}
		<div class="flex flex-col gap-6">
			<Button href="/admin/users/add/" variant="primary" class="w-auto self-start px-6"
				>Add User</Button
			>
			<Separator />
		</div>
	{/snippet}

	<Table.PaginatedBody
		itemCount={users.length}
		{totalItems}
		bind:page
		bind:perPage
		{loading}
		emptyMessage="No users"
		minimal={false}
		showPerPageSelect
		alwaysShowPagination
		selectTriggerClass="h-9 px-2 py-0"
	>
		{#snippet list()}
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
							<Button variant="destructive" size="inline" onclick={() => openDeleteUser(user)}
								>Delete</Button
							>
						{/snippet}
					</Table.Row>
					{#if index < users.length - 1}
						<Table.Separator />
					{/if}
				{/each}
			</Table.List>
		{/snippet}
	</Table.PaginatedBody>

	<DeleteUser
		bind:open={deleteOpen}
		user={userToDelete}
		onSuccess={handleDeleteSuccess}
		onError={handleDeleteError}
	/>
</Table.Root>

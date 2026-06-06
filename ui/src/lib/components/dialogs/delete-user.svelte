<script lang="ts">
	import { ApiError } from '$lib/api';
	import { deleteUser } from '$lib/api/admin-api';
	import { Button, Drawer } from '$lib/components/ui';
	import type { AdminUserModel } from '$lib/models/admin-user-model';

	type Props = {
		open?: boolean;
		user: AdminUserModel | null;
		onSuccess?: () => void;
		onError?: (message: string) => void;
	};

	let { open = $bindable(false), user, onSuccess, onError }: Props = $props();

	let isPosting = $state(false);

	$effect(() => {
		if (open) {
			isPosting = false;
		}
	});

	async function doDelete(): Promise<void> {
		if (!user) return;

		isPosting = true;

		try {
			await deleteUser(user.id);
			open = false;
			onSuccess?.();
		} catch (err) {
			onError?.(err instanceof ApiError ? err.message : 'Failed to delete user');
		} finally {
			isPosting = false;
		}
	}
</script>

{#snippet alertContents()}
	<Drawer.Alert>
		<div class="flex flex-col gap-2 text-center text-foreground">
			<span class="text-lg">Are you sure you want to delete this user?</span>
			{#if user}
				<span class="font-semibold text-background-primary">{user.displayName}</span>
			{/if}
			<span class="text-sm text-foreground-alt-2">All associated data will be deleted</span>
		</div>
	</Drawer.Alert>
{/snippet}

{#snippet deleteButton()}
	<Button variant="destructive" class="w-full" loading={isPosting} onclick={doDelete}>
		Delete
	</Button>
{/snippet}

<Drawer.Root bind:open>
	<Drawer.Content handleClass="bg-foreground-alt-2">
		<div class="overflow-hidden rounded-lg">
			{@render alertContents()}

			<Drawer.Footer class="grid h-auto grid-cols-2 gap-2 px-4 py-4">
				<Drawer.CloseButton class="w-full">Cancel</Drawer.CloseButton>
				{@render deleteButton()}
			</Drawer.Footer>
		</div>
	</Drawer.Content>
</Drawer.Root>

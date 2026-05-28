<script lang="ts">
	import { ApiError } from '$lib/api';
	import { deleteAdminGroup } from '$lib/api/admin-api';
	import { Button, Drawer } from '$lib/components/ui';
	import type { AdminGroupModel } from '$lib/models/admin-group-model';

	type Props = {
		open?: boolean;
		group: AdminGroupModel | null;
		onSuccess?: () => void;
		onError?: (message: string) => void;
	};

	let { open = $bindable(false), group, onSuccess, onError }: Props = $props();

	let isPosting = $state(false);

	$effect(() => {
		if (open) {
			isPosting = false;
		}
	});

	async function doDelete(): Promise<void> {
		if (!group) return;

		isPosting = true;

		try {
			await deleteAdminGroup(group.id);
			open = false;
			onSuccess?.();
		} catch (err) {
			onError?.(err instanceof ApiError ? err.message : 'Failed to delete group');
		} finally {
			isPosting = false;
		}
	}
</script>

{#snippet alertContents()}
	<Drawer.Alert>
		<div class="flex flex-col gap-2 text-center text-text">
			<span class="text-lg">Are you sure you want to delete this group?</span>
			{#if group}
				<span class="font-semibold text-button-primary">{group.name}</span>
			{/if}
			<span class="text-sm text-text-muted">All associated data will be deleted</span>
		</div>
	</Drawer.Alert>
{/snippet}

{#snippet deleteButton()}
	<Button variant="destructive" class="w-full" disabled={isPosting} onclick={doDelete}>
		{isPosting ? 'Deleting…' : 'Delete'}
	</Button>
{/snippet}

<Drawer.Root bind:open>
	<Drawer.Content handleClass="bg-text-muted">
		<div class="overflow-hidden rounded-lg">
			{@render alertContents()}

			<Drawer.Footer class="grid h-auto grid-cols-2 gap-2 px-4 py-4">
				<Drawer.CloseButton class="w-full">Cancel</Drawer.CloseButton>
				{@render deleteButton()}
			</Drawer.Footer>
		</div>
	</Drawer.Content>
</Drawer.Root>

<script lang="ts">
	import { ApiError } from '$lib/api';
	import { deleteUser } from '$lib/api/admin-api';
	import { DestroyDialog } from '$lib/components/ui';
	import type { AdminUserModel } from '$lib/models/admin-user-model';
	import { withMinLoadingDelay } from '$lib/utils';

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
			await withMinLoadingDelay(deleteUser(user.id));
			open = false;
			onSuccess?.();
		} catch (err) {
			onError?.(err instanceof ApiError ? err.message : 'Failed to delete user');
		} finally {
			isPosting = false;
		}
	}
</script>

<DestroyDialog
	bind:open
	title="Are you sure you want to delete this user?"
	detail={user ? user.displayName || user.username : undefined}
	description="All associated data will be deleted"
	confirmLabel="Delete"
	loading={isPosting}
	onConfirm={doDelete}
/>

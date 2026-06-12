<script lang="ts">
	import { deleteGroup } from '$lib/api/groups-api';
	import { DestroyDialog } from '$lib/components/ui';
	import type { GroupModel } from '$lib/models/group-model';
	import { apiErrorMessage, withMinLoadingDelay } from '$lib/utils';

	type Props = {
		open?: boolean;
		group: GroupModel | null;
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
			await withMinLoadingDelay(deleteGroup(group.id));
			open = false;
			onSuccess?.();
		} catch (err) {
			onError?.(apiErrorMessage(err, 'Failed to delete group'));
		} finally {
			isPosting = false;
		}
	}
</script>

<DestroyDialog
	bind:open
	title="Are you sure you want to delete this group?"
	detail={group?.name}
	description="All associated group data will be permanently deleted"
	confirmLabel="Delete"
	loading={isPosting}
	onConfirm={doDelete}
/>

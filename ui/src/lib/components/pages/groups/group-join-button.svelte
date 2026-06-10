<script lang="ts">
	import { requestGroupJoin } from '$lib/api/groups-api';
	import { PlusIcon } from '$lib/components/icons';
	import { Button } from '$lib/components/ui';
	import { apiErrorMessage, withMinLoadingDelay } from '$lib/utils';

	type Props = {
		groupId: string;
		onjoined?: () => void;
		onerror?: (message: string) => void;
	};

	let { groupId, onjoined, onerror }: Props = $props();

	let loading = $state(false);

	// handleClick submits a join request for the group
	async function handleClick(): Promise<void> {
		if (loading) {
			return;
		}

		loading = true;

		try {
			await withMinLoadingDelay(requestGroupJoin(groupId));
			onjoined?.();
		} catch (err) {
			onerror?.(apiErrorMessage(err, 'Failed to request join'));
		} finally {
			loading = false;
		}
	}
</script>

<Button
	type="button"
	variant="ghost"
	size="inline"
	class="text-foreground-alt-2 h-5 w-5 min-h-5 min-w-5 shrink-0 p-0 normal-case hover:bg-transparent"
	aria-label="Request to join group"
	{loading}
	onclick={handleClick}
>
	<PlusIcon class="size-5 shrink-0 stroke-2" />
</Button>

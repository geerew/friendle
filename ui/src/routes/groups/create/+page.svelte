<script lang="ts">
	import { goto } from '$app/navigation';
	import { createGroup } from '$lib/api/groups-api';
	import { Button, Input, Table } from '$lib/components/ui';
	import { MAX_GROUP_NAME_LENGTH } from '$lib/models/group-model';
	import { apiErrorMessage, withMinLoadingDelay } from '$lib/utils';
	import { toast } from 'svelte-sonner';

	let name = $state('');
	let submitting = $state(false);

	const canSubmit = $derived(name.trim() !== '' && !submitting);

	async function handleCreate(event: SubmitEvent): Promise<void> {
		event.preventDefault();

		const trimmed = name.trim();
		if (!trimmed) {
			toast.error('Group name is required');
			return;
		}

		if ([...trimmed].length > MAX_GROUP_NAME_LENGTH) {
			toast.error(`Group name must be no more than ${MAX_GROUP_NAME_LENGTH} characters`);
			return;
		}

		submitting = true;

		try {
			await withMinLoadingDelay(createGroup({ name: trimmed }));
			toast.success('Group created');
			await goto('/');
		} catch (err) {
			toast.error(apiErrorMessage(err, 'Failed to create group'));
		} finally {
			submitting = false;
		}
	}
</script>

<Table.Root title="Create Group">
	<form class="flex flex-col gap-4" onsubmit={handleCreate}>
		<Input bind:value={name} placeholder="Group Name" autofocus required />

		<Button type="submit" variant="primary" disabled={!canSubmit} loading={submitting}>
			Create Group
		</Button>
	</form>
</Table.Root>

<script lang="ts">
	import { goto } from '$app/navigation';
	import { createGroup } from '$lib/api/groups-api';
	import { AppShell } from '$lib/components';
	import { Button, Field, Input } from '$lib/components/ui';
	import { apiErrorMessage, withMinLoadingDelay } from '$lib/utils';
	import { toast } from 'svelte-sonner';

	const maxGroupNameLength = 64;

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

		if (trimmed.length > maxGroupNameLength) {
			toast.error(`Group name must be no more than ${maxGroupNameLength} characters`);
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

<AppShell title="Create Group">
	<form class="flex flex-col gap-4" onsubmit={handleCreate}>
		<Input bind:value={name} placeholder="Group Name" autofocus required />

		<Button type="submit" variant="primary" disabled={!canSubmit} loading={submitting}>
			Create Group
		</Button>
	</form>
</AppShell>

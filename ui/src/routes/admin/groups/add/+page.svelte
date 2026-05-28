<script lang="ts">
	import { goto } from '$app/navigation';
	import { ApiError } from '$lib/api';
	import { createGroup } from '$lib/api/groups-api';
	import { AppShell } from '$lib/components';
	import { Button, Field, Input } from '$lib/components/ui';

	let name = $state('');
	let submitting = $state(false);
	let error = $state<string | null>(null);

	const submitDisabled = $derived(name.trim() === '');

	async function handleAdd(event: SubmitEvent): Promise<void> {
		event.preventDefault();
		submitting = true;
		error = null;

		try {
			await createGroup({ name: name.trim() });
			await goto('/admin/groups/');
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to create group';
		} finally {
			submitting = false;
		}
	}
</script>

<AppShell
	breadcrumb={[
		{ label: 'Admin', href: '/admin/' },
		{ label: 'Groups', href: '/admin/groups/' },
		{ label: 'Add Group' }
	]}
>
	<form class="flex flex-col gap-4" onsubmit={handleAdd}>
		<Field label="Group name">
			<Input bind:value={name} required maxlength={64} />
		</Field>

		{#if error}
			<p class="text-sm text-error">{error}</p>
		{/if}

		<div class="grid grid-cols-2 gap-2">
			<Button href="/admin/groups/" variant="secondary">Cancel</Button>
			<Button type="submit" variant="primary" disabled={submitDisabled || submitting}>
				{submitting ? 'Adding…' : 'Add'}
			</Button>
		</div>
	</form>
</AppShell>

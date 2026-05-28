<script lang="ts">
	import { goto } from '$app/navigation';
	import { ApiError } from '$lib/api';
	import { createGroup } from '$lib/api/groups-api';
	import { AppShell } from '$lib/components';
	import { Button, Field, Input } from '$lib/components/ui';

	let name = $state('');
	let intervalHours = $state(24);
	let timezone = $state('UTC');
	let submitting = $state(false);
	let error = $state<string | null>(null);

	async function handleSubmit(event: SubmitEvent): Promise<void> {
		event.preventDefault();
		submitting = true;
		error = null;

		try {
			const group = await createGroup({ name, intervalHours, timezone });
			await goto(`/groups/${group.id}/`);
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to create group';
		} finally {
			submitting = false;
		}
	}
</script>

<AppShell title="Create Group">
	<form class="flex flex-col gap-4" onsubmit={handleSubmit}>
		<Field label="Group name">
			<Input bind:value={name} required maxlength={64} />
		</Field>

		<Field label="Round interval (hours)">
			<Input type="number" min={1} max={168} bind:value={intervalHours} required />
		</Field>

		<Field label="Timezone">
			<Input bind:value={timezone} required />
		</Field>

		{#if error}
			<p class="text-sm text-error">{error}</p>
		{/if}

		<Button type="submit" variant="primary" disabled={submitting}>
			{submitting ? 'Creating…' : 'Create group'}
		</Button>
	</form>
</AppShell>

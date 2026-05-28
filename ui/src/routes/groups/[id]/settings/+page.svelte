<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { ApiError } from '$lib/api';
	import { deleteGroup, getGroup, leaveGroup, updateGroup } from '$lib/api/groups-api';
	import { AppShell } from '$lib/components';
	import { Button, Field, Input } from '$lib/components/ui';

	const groupId = $derived(page.params.id ?? '');

	let name = $state('');
	let intervalHours = $state(24);
	let timezone = $state('UTC');
	let myRole = $state<string | null>(null);
	let loading = $state(true);
	let saving = $state(false);
	let error = $state<string | null>(null);
	let message = $state<string | null>(null);

	$effect(() => {
		loadGroup();
	});

	async function loadGroup(): Promise<void> {
		loading = true;
		error = null;

		try {
			const group = await getGroup(groupId);
			name = group.name;
			intervalHours = group.intervalHours;
			timezone = group.timezone;
			myRole = group.myRole ?? null;
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to load group';
		} finally {
			loading = false;
		}
	}

	async function handleSave(event: SubmitEvent): Promise<void> {
		event.preventDefault();
		saving = true;
		error = null;
		message = null;

		try {
			await updateGroup(groupId, { name, intervalHours, timezone });
			message = 'Settings saved';
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to save settings';
		} finally {
			saving = false;
		}
	}

	async function handleLeave(): Promise<void> {
		if (!confirm('Leave this group?')) return;

		try {
			await leaveGroup(groupId);
			await goto('/');
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to leave group';
		}
	}

	async function handleDelete(): Promise<void> {
		if (!confirm('Delete this group permanently?')) return;

		try {
			await deleteGroup(groupId);
			await goto('/');
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to delete group';
		}
	}
</script>

<AppShell
	breadcrumb={[
		{ label: name || 'Group', href: `/groups/${groupId}/` },
		{ label: 'Group Settings' }
	]}
>
	{#if loading}
		<p class="text-text-muted">Loading…</p>
	{:else}
		<form class="flex flex-col gap-4" onsubmit={handleSave}>
			<Field label="Group name">
				<Input
					bind:value={name}
					required
					maxlength={64}
					disabled={myRole !== 'group_admin'}
				/>
			</Field>

			<Field label="Round interval (hours)">
				<Input
					type="number"
					min={1}
					max={168}
					bind:value={intervalHours}
					required
					disabled={myRole !== 'group_admin'}
				/>
			</Field>

			<Field label="Timezone">
				<Input bind:value={timezone} required disabled={myRole !== 'group_admin'} />
			</Field>

			{#if message}
				<p class="text-sm text-tile-correct">{message}</p>
			{/if}

			{#if error}
				<p class="text-sm text-error">{error}</p>
			{/if}

			{#if myRole === 'group_admin'}
				<Button type="submit" variant="primary" disabled={saving}>
					{saving ? 'Saving…' : 'Save changes'}
				</Button>
				<Button type="button" variant="destructive" onclick={handleDelete}>Delete group</Button>
			{/if}

			<Button type="button" variant="ghost" onclick={handleLeave}>Leave group</Button>
		</form>
	{/if}
</AppShell>

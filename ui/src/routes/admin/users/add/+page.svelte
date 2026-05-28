<script lang="ts">
	import { goto } from '$app/navigation';
	import { ApiError } from '$lib/api';
	import { createUser } from '$lib/api/user-api';
	import { AppShell } from '$lib/components';
	import { Button, Field, Input, RadioGroup } from '$lib/components/ui';
	import { SelectSiteRoles, type SiteRole } from '$lib/models/admin-user-model';

	let username = $state('');
	let displayName = $state('');
	let siteRole = $state<SiteRole>('site_user');
	let password = $state('');
	let confirmPassword = $state('');
	let error = $state<string | null>(null);
	let submitting = $state(false);
	let previousUsername = $state('');

	const submitDisabled = $derived(
		username === '' || password === '' || confirmPassword === ''
	);

	$effect(() => {
		if (displayName === '' || displayName === previousUsername) {
			displayName = username;
		}

		previousUsername = username;
	});

	async function handleCreate(event: SubmitEvent): Promise<void> {
		event.preventDefault();

		if (password !== confirmPassword) {
			error = 'Passwords do not match';
			return;
		}

		submitting = true;
		error = null;

		try {
			await createUser({
				username,
				displayName,
				password,
				siteRole
			});
			await goto('/admin/users/');
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to create user';
		} finally {
			submitting = false;
		}
	}
</script>

<AppShell
	breadcrumb={[
		{ label: 'Admin', href: '/admin/' },
		{ label: 'Users', href: '/admin/users/' },
		{ label: 'Add User' }
	]}
>
	<form class="flex flex-col gap-4" onsubmit={handleCreate}>
		<Field label="Username">
			<Input bind:value={username} autocomplete="username" required />
		</Field>

		<Field label="Display name">
			<Input bind:value={displayName} autocomplete="name" />
		</Field>

		<Field label="Role">
			<RadioGroup name="siteRole" items={SelectSiteRoles} bind:value={siteRole} required />
		</Field>

		<Field label="Password">
			<Input password bind:value={password} autocomplete="new-password" required />
		</Field>

		<Field label="Confirm password">
			<Input password bind:value={confirmPassword} autocomplete="new-password" required />
		</Field>

		{#if error}
			<p class="text-sm text-error">{error}</p>
		{/if}

		<div class="grid grid-cols-2 gap-2">
			<Button href="/admin/users/" variant="secondary">Cancel</Button>
			<Button type="submit" variant="primary" disabled={submitDisabled || submitting}>
				{submitting ? 'Adding…' : 'Add'}
			</Button>
		</div>
	</form>
</AppShell>

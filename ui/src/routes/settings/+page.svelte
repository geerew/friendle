<script lang="ts">
	import { goto } from '$app/navigation';
	import { ApiError } from '$lib/api';
	import { deleteMe, logout, updateMe } from '$lib/api/auth-api';
	import { auth } from '$lib/auth.svelte';
	import { AppShell } from '$lib/components';
	import { Button, Field, Input } from '$lib/components/ui';

	let displayName = $state(auth.user?.displayName ?? '');
	let currentPassword = $state('');
	let newPassword = $state('');
	let saving = $state(false);
	let loggingOut = $state(false);
	let error = $state<string | null>(null);
	let message = $state<string | null>(null);

	$effect(() => {
		displayName = auth.user?.displayName ?? '';
	});

	async function handleSave(event: SubmitEvent): Promise<void> {
		event.preventDefault();
		saving = true;
		error = null;
		message = null;

		try {
			const user = await updateMe({
				displayName,
				password: newPassword || undefined,
				currentPassword: currentPassword || undefined
			});
			auth.setUser(user);
			currentPassword = '';
			newPassword = '';
			message = 'Profile updated';
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to update profile';
		} finally {
			saving = false;
		}
	}

	async function handleLogout(): Promise<void> {
		loggingOut = true;

		try {
			await logout();
			auth.clear();
			await goto('/auth/login/');
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Logout failed';
		} finally {
			loggingOut = false;
		}
	}

	async function handleDeleteAccount(): Promise<void> {
		if (!currentPassword) {
			error = 'Enter your current password to delete your account';
			return;
		}

		if (!confirm('Delete your account permanently?')) return;

		try {
			await deleteMe({ currentPassword });
			auth.clear();
			await goto('/auth/login/');
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to delete account';
		}
	}
</script>

<AppShell title="Settings" backHref="/">
	<form class="flex flex-col gap-4" onsubmit={handleSave}>
		<Field label="Username">
			<Input value={auth.user?.username ?? ''} disabled />
		</Field>

		<Field label="Display name">
			<Input bind:value={displayName} required />
		</Field>

		<Field label="Current password">
			<Input password bind:value={currentPassword} autocomplete="current-password" />
		</Field>

		<Field label="New password">
			<Input password bind:value={newPassword} autocomplete="new-password" />
		</Field>

		{#if message}
			<p class="text-sm text-tile-correct">{message}</p>
		{/if}

		{#if error}
			<p class="text-sm text-error">{error}</p>
		{/if}

		<Button type="submit" variant="primary" disabled={saving}>
			{saving ? 'Saving…' : 'Save profile'}
		</Button>
	</form>

	<div class="mt-auto flex flex-col gap-3 pt-6">
		<Button variant="secondary" onclick={handleLogout} disabled={loggingOut}>
			{loggingOut ? 'Signing out…' : 'Sign out'}
		</Button>
		<Button variant="danger" onclick={handleDeleteAccount}>Delete account</Button>
	</div>
</AppShell>

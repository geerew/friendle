<script lang="ts">
	import { goto } from '$app/navigation';
	import { ApiError } from '$lib/api';
	import { deleteMe, updateMe } from '$lib/api/auth-api';
	import { auth } from '$lib/auth.svelte';
	import { AppShell } from '$lib/components';
	import { Button, Input } from '$lib/components/ui';
	import { isPasswordFieldError } from '$lib/utils';

	const minPasswordLength = 8;

	let isEditingDisplayName = $state(false);
	let displayName = $state('');
	let originalDisplayName = $state('');
	let savingDisplayName = $state(false);

	let isEditingPassword = $state(false);
	let currentPassword = $state('');
	let newPassword = $state('');
	let confirmPassword = $state('');
	let savingPassword = $state(false);
	let previousCurrentPassword = $state('');
	let previousNewPassword = $state('');
	let previousConfirmPassword = $state('');

	let isDeletingAccount = $state(false);
	let deletePassword = $state('');
	let deletingAccount = $state(false);

	let error = $state<string | null>(null);
	let message = $state<string | null>(null);

	$effect(() => {
		if (!isEditingDisplayName && auth.user?.displayName) {
			displayName = auth.user.displayName;
		}
	});

	$effect(() => {
		if (
			isEditingPassword &&
			isPasswordFieldError(error) &&
			(currentPassword !== previousCurrentPassword ||
				newPassword !== previousNewPassword ||
				confirmPassword !== previousConfirmPassword)
		) {
			error = null;
		}

		previousCurrentPassword = currentPassword;
		previousNewPassword = newPassword;
		previousConfirmPassword = confirmPassword;
	});

	function clearFeedback(): void {
		error = null;
		message = null;
	}

	function startEditingDisplayName(): void {
		clearFeedback();
		originalDisplayName = auth.user?.displayName ?? '';
		displayName = originalDisplayName;
		isEditingDisplayName = true;
	}

	function cancelEditingDisplayName(): void {
		displayName = originalDisplayName;
		isEditingDisplayName = false;
	}

	async function saveDisplayName(): Promise<void> {
		const trimmed = displayName.trim();
		if (!trimmed || trimmed === auth.user?.displayName) {
			cancelEditingDisplayName();
			return;
		}

		savingDisplayName = true;
		clearFeedback();

		try {
			const user = await updateMe({ displayName: trimmed });
			auth.setUser(user);
			isEditingDisplayName = false;
			message = 'Display name updated';
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to update display name';
			displayName = originalDisplayName;
		} finally {
			savingDisplayName = false;
		}
	}

	function startEditingPassword(): void {
		clearFeedback();
		currentPassword = '';
		newPassword = '';
		confirmPassword = '';
		isEditingPassword = true;
	}

	function cancelEditingPassword(): void {
		currentPassword = '';
		newPassword = '';
		confirmPassword = '';
		isEditingPassword = false;
	}

	async function savePassword(): Promise<void> {
		if (!currentPassword || !newPassword || !confirmPassword) {
			error = 'All password fields are required';
			return;
		}

		if (newPassword.length < minPasswordLength) {
			error = 'Password must be at least 8 characters';
			return;
		}

		if (newPassword !== confirmPassword) {
			error = 'Passwords do not match';
			return;
		}

		savingPassword = true;
		clearFeedback();

		try {
			await updateMe({ currentPassword, password: newPassword });
			cancelEditingPassword();
			message = 'Password updated';
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to update password';
		} finally {
			savingPassword = false;
		}
	}

	function startDeletingAccount(): void {
		clearFeedback();
		deletePassword = '';
		isDeletingAccount = true;
	}

	function cancelDeletingAccount(): void {
		deletePassword = '';
		isDeletingAccount = false;
	}

	async function handleDeleteAccount(): Promise<void> {
		if (!deletePassword) {
			error = 'Enter your current password to delete your account';
			return;
		}

		if (!confirm('Delete your account permanently?')) return;

		deletingAccount = true;
		clearFeedback();

		try {
			await deleteMe({ currentPassword: deletePassword });
			auth.clear();
			await goto('/auth/login/');
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to delete account';
		} finally {
			deletingAccount = false;
		}
	}
</script>

{#if auth.user}
	<AppShell title="Settings">
		<div class="flex flex-col gap-5">
			{#if message}
				<p class="text-sm text-tile-correct">{message}</p>
			{/if}

			{#if error}
				<p class="text-sm text-error">{error}</p>
			{/if}

			<section class="flex flex-col gap-3">
				<h2 class="section-title">Username</h2>
				<p class="text-2xl text-button-primary">{auth.user.username}</p>
			</section>

			<div class="h-px shrink-0 bg-border"></div>

			<section class="flex flex-col gap-3">
				<div class="flex items-center justify-between gap-3">
					<h2 class="section-title">Display name</h2>
					{#if isEditingDisplayName}
						<Button
							type="button"
							variant="ghost"
							size="inline"
							onclick={cancelEditingDisplayName}
						>
							Cancel
						</Button>
					{:else}
						<Button type="button" variant="ghost" size="inline" onclick={startEditingDisplayName}>
							Edit
						</Button>
					{/if}
				</div>

				{#if isEditingDisplayName}
					<div class="flex flex-col gap-2">
						<Input bind:value={displayName} required />
						<Button
							type="button"
							variant="primary"
							size="inline"
							class="w-36 uppercase tracking-wide"
							disabled={!displayName.trim() ||
								displayName.trim() === auth.user.displayName ||
								savingDisplayName}
							onclick={saveDisplayName}
						>
							{savingDisplayName ? 'Saving…' : 'Save'}
						</Button>
					</div>
				{:else}
					<p class="text-2xl text-button-primary">{auth.user.displayName}</p>
				{/if}
			</section>

			<div class="h-px shrink-0 bg-border"></div>

			<section class="flex flex-col gap-3">
				<div class="flex items-center justify-between gap-3">
					<h2 class="section-title">Password</h2>
					{#if isEditingPassword}
						<Button
							type="button"
							variant="ghost"
							size="inline"
							onclick={cancelEditingPassword}
						>
							Cancel
						</Button>
					{:else}
						<Button type="button" variant="ghost" size="inline" onclick={startEditingPassword}>
							Edit
						</Button>
					{/if}
				</div>

				{#if isEditingPassword}
					<div class="flex flex-col gap-3">
						<div class="flex flex-col gap-2">
							<span class="text-sm text-text-muted">Current password</span>
							<Input password bind:value={currentPassword} autocomplete="current-password" />
						</div>
						<div class="flex flex-col gap-2">
							<span class="text-sm text-text-muted">New password</span>
							<Input password bind:value={newPassword} autocomplete="new-password" />
						</div>
						<div class="flex flex-col gap-2">
							<span class="text-sm text-text-muted">Confirm password</span>
							<Input password bind:value={confirmPassword} autocomplete="new-password" />
						</div>
						<Button
							type="button"
							variant="primary"
							size="inline"
							class="w-36 uppercase tracking-wide"
							disabled={!currentPassword ||
								!newPassword ||
								!confirmPassword ||
								savingPassword}
							onclick={savePassword}
						>
							{savingPassword ? 'Saving…' : 'Save'}
						</Button>
					</div>
				{/if}
			</section>

			<div class="h-px shrink-0 bg-border"></div>

			<section class="flex flex-col gap-1">
				<div class="flex items-center justify-between gap-3">
					<h2 class="section-title">Delete account</h2>
					{#if isDeletingAccount}
						<Button
							type="button"
							variant="ghost"
							size="inline"
							onclick={cancelDeletingAccount}
						>
							Cancel
						</Button>
					{:else}
						<Button
							type="button"
							variant="destructive"
							size="inline"
							onclick={startDeletingAccount}
						>
							Delete
						</Button>
					{/if}
				</div>
				<p class="text-sm text-text-muted">Permanently delete your Friendle account</p>

				{#if isDeletingAccount}
					<div class="mt-2 flex flex-col gap-3">
						<div class="flex flex-col gap-2">
							<span class="text-sm text-text-muted">Current password</span>
							<Input password bind:value={deletePassword} autocomplete="current-password" />
						</div>
						<Button
							type="button"
							variant="destructive"
							size="inline"
							class="w-36 uppercase tracking-wide"
							disabled={!deletePassword || deletingAccount}
							onclick={handleDeleteAccount}
						>
							{deletingAccount ? 'Deleting…' : 'Delete account'}
						</Button>
					</div>
				{/if}
			</section>
		</div>
	</AppShell>
{/if}

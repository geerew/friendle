<script lang="ts">
	import { goto } from '$app/navigation';
	import { deleteMe, updateMe } from '$lib/api/auth-api';
	import { auth } from '$lib/auth.svelte';
	import { AppShell, EditableSection } from '$lib/components';
	import { Button, DestroyDialog, Input, Separator } from '$lib/components/ui';
	import { apiErrorMessage, withMinLoadingDelay } from '$lib/utils';
	import { toast } from 'svelte-sonner';

	const minPasswordLength = 8;

	let editingDisplayName = $state(false);
	let displayNameDraft = $state('');
	let savingDisplayName = $state(false);

	let editingPassword = $state(false);
	let currentPassword = $state('');
	let newPassword = $state('');
	let confirmPassword = $state('');
	let savingPassword = $state(false);

	let deletingAccount = $state(false);
	let deletePassword = $state('');
	let deleteConfirmOpen = $state(false);
	let deleting = $state(false);

	let passwordSubmitBlocked = $state(false);

	const canSaveDisplayName = $derived(
		displayNameDraft.trim() !== '' &&
			displayNameDraft.trim() !== auth.user?.displayName &&
			!savingDisplayName
	);

	const canSavePassword = $derived(
		currentPassword !== '' &&
			newPassword !== '' &&
			confirmPassword !== '' &&
			!savingPassword &&
			!passwordSubmitBlocked
	);

	// Re-enable password save after a failed attempt once the user edits a field
	function onPasswordFieldInput(): void {
		passwordSubmitBlocked = false;
	}

	// Enter display name edit mode with the current value
	function openDisplayNameEdit(): void {
		displayNameDraft = auth.user?.displayName ?? '';
		editingDisplayName = true;
	}

	// Close display name edit mode without saving
	function closeDisplayNameEdit(): void {
		editingDisplayName = false;
	}

	// Persist a new display name via the API
	async function saveDisplayName(): Promise<void> {
		const trimmed = displayNameDraft.trim();
		if (!trimmed || trimmed === auth.user?.displayName) {
			closeDisplayNameEdit();
			return;
		}

		savingDisplayName = true;

		try {
			const user = await withMinLoadingDelay(updateMe({ displayName: trimmed }));
			auth.setUser(user);
			closeDisplayNameEdit();
			toast.success('Display name updated');
		} catch (err) {
			toast.error(apiErrorMessage(err, 'Failed to update display name'));
			displayNameDraft = auth.user?.displayName ?? '';
		} finally {
			savingDisplayName = false;
		}
	}

	// Enter password edit mode with empty fields
	function openPasswordEdit(): void {
		passwordSubmitBlocked = false;
		currentPassword = '';
		newPassword = '';
		confirmPassword = '';
		editingPassword = true;
	}

	// Discard password edits and close edit mode
	function closePasswordEdit(): void {
		currentPassword = '';
		newPassword = '';
		confirmPassword = '';
		passwordSubmitBlocked = false;
		editingPassword = false;
	}

	// Validate and persist a new password via the API
	async function savePassword(): Promise<void> {
		if (!currentPassword || !newPassword || !confirmPassword) {
			toast.error('All password fields are required');
			passwordSubmitBlocked = true;
			return;
		}

		if (newPassword.length < minPasswordLength) {
			toast.error(`New password must be at least ${minPasswordLength} characters`);
			passwordSubmitBlocked = true;
			return;
		}

		if (newPassword !== confirmPassword) {
			toast.error('Passwords do not match');
			passwordSubmitBlocked = true;
			return;
		}

		savingPassword = true;

		try {
			await withMinLoadingDelay(updateMe({ currentPassword, password: newPassword }));
			closePasswordEdit();
			toast.success('Password updated');
		} catch (err) {
			toast.error(apiErrorMessage(err, 'Failed to update password'));
			passwordSubmitBlocked = true;
		} finally {
			savingPassword = false;
		}
	}

	// Show the delete account form
	function openDeleteAccount(): void {
		deletePassword = '';
		deletingAccount = true;
	}

	// Hide the delete account form
	function closeDeleteAccount(): void {
		deletePassword = '';
		deletingAccount = false;
		deleteConfirmOpen = false;
	}

	// Open the destroy confirmation dialog
	function requestDeleteAccount(): void {
		if (!deletePassword) {
			toast.error('Enter your current password to delete your account');
			return;
		}

		deleteConfirmOpen = true;
	}

	// Delete the signed-in account and redirect to login
	async function confirmDeleteAccount(): Promise<void> {
		deleting = true;

		try {
			await withMinLoadingDelay(deleteMe({ currentPassword: deletePassword }));
			deleteConfirmOpen = false;
			closeDeleteAccount();
			auth.clear();
			await goto('/auth/login/');
		} catch (err) {
			toast.error(apiErrorMessage(err, 'Failed to delete account'));
		} finally {
			deleting = false;
		}
	}
</script>

{#snippet saveButton(loading: boolean, canSave: boolean)}
	<Button
		type="submit"
		variant="primary"
		class="w-auto min-w-16 self-start px-5 text-xs"
		disabled={!canSave}
		{loading}
	>
		Save
	</Button>
{/snippet}

{#if auth.user}
	<AppShell title="Settings">
		<div class="flex flex-col gap-5">
			<section class="flex flex-col gap-3">
				<h2 class="section-title">Username</h2>
				<p class="text-background-primary text-2xl">{auth.user.username}</p>
			</section>

			<Separator />

			<EditableSection
				title="Display name"
				editing={editingDisplayName}
				cancelDisabled={savingDisplayName}
				onEdit={openDisplayNameEdit}
				onCancel={closeDisplayNameEdit}
			>
				{#if editingDisplayName}
					<form
						class="flex flex-col gap-4"
						onsubmit={(event) => {
							event.preventDefault();
							void saveDisplayName();
						}}
					>
						<Input bind:value={displayNameDraft} autofocus selectOnMount required />
						{@render saveButton(savingDisplayName, canSaveDisplayName)}
					</form>
				{:else}
					<p class="text-background-primary text-2xl">{auth.user.displayName}</p>
				{/if}
			</EditableSection>

			<Separator />

			<EditableSection
				title="Password"
				editing={editingPassword}
				cancelDisabled={savingPassword}
				onEdit={openPasswordEdit}
				onCancel={closePasswordEdit}
			>
				{#if editingPassword}
					<form
						class="flex flex-col gap-3"
						onsubmit={(event) => {
							event.preventDefault();
							void savePassword();
						}}
					>
						<label class="flex flex-col gap-2">
							<span class="text-foreground-alt-2 text-sm">Current password</span>
							<Input
								password
								bind:value={currentPassword}
								autocomplete="current-password"
								autofocus
								oninput={onPasswordFieldInput}
							/>
						</label>
						<label class="flex flex-col gap-2">
							<span class="text-foreground-alt-2 text-sm">New password</span>
							<Input
								password
								bind:value={newPassword}
								autocomplete="new-password"
								oninput={onPasswordFieldInput}
							/>
						</label>
						<label class="flex flex-col gap-2">
							<span class="text-foreground-alt-2 text-sm">Confirm password</span>
							<Input
								password
								bind:value={confirmPassword}
								autocomplete="new-password"
								oninput={onPasswordFieldInput}
							/>
						</label>
						{@render saveButton(savingPassword, canSavePassword)}
					</form>
				{/if}
			</EditableSection>

			<Separator />

			<EditableSection
				title="Delete account"
				editing={deletingAccount}
				cancelDisabled={deleting}
				editLabel="Delete"
				editVariant="destructive"
				hint="Permanently delete your Friendle account"
				onEdit={openDeleteAccount}
				onCancel={closeDeleteAccount}
			>
				{#if deletingAccount}
					<div class="flex flex-col gap-3">
						<label class="flex flex-col gap-2">
							<span class="text-foreground-alt-2 text-sm">Current password</span>
							<Input
								password
								bind:value={deletePassword}
								autocomplete="current-password"
								autofocus
							/>
						</label>
						<Button
							type="button"
							variant="destructive"
							class="w-36 px-0 text-xs"
							disabled={!deletePassword || deleting}
							onclick={requestDeleteAccount}
						>
							Delete account
						</Button>
					</div>
				{/if}
			</EditableSection>
		</div>
	</AppShell>

	<DestroyDialog
		bind:open={deleteConfirmOpen}
		title="Are you sure you want to delete your account?"
		description="All associated data will be permanently deleted"
		confirmLabel="Delete account"
		loading={deleting}
		onConfirm={confirmDeleteAccount}
	/>
{/if}

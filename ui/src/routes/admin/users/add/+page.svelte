<script lang="ts">
	import { goto } from '$app/navigation';
	import { ApiError } from '$lib/api';
	import { createUser } from '$lib/api/admin-api';
	import { Button, Field, Input, RadioGroup, Table } from '$lib/components/ui';
	import { SelectSiteRoles, type SiteRole } from '$lib/models/admin-user-model';
	import { isPasswordFieldError, withMinLoadingDelay } from '$lib/utils';

	const minPasswordLength = 8;

	let username = $state('');
	let displayName = $state('');
	let siteRole = $state<SiteRole>('site_user');
	let password = $state('');
	let confirmPassword = $state('');
	let passwordTooShortError = $state(false);
	let passwordMismatchError = $state(false);
	let error = $state<string | null>(null);
	let submitting = $state(false);
	let previousUsername = $state('');
	let previousPassword = $state('');
	let previousConfirmPassword = $state('');

	const submitDisabled = $derived(
		username === '' || password === '' || confirmPassword === ''
	);

	$effect(() => {
		if (displayName === '' || displayName === previousUsername) {
			displayName = username;
		}

		previousUsername = username;
	});

	$effect(() => {
		if (
			passwordMismatchError &&
			(password !== previousPassword || confirmPassword !== previousConfirmPassword)
		) {
			passwordMismatchError = false;
		}

		if (passwordTooShortError && (password !== previousPassword || confirmPassword !== previousConfirmPassword)) {
			passwordTooShortError = false;
		}

		if (
			isPasswordFieldError(error) &&
			(password !== previousPassword || confirmPassword !== previousConfirmPassword)
		) {
			error = null;
		}

		previousPassword = password;
		previousConfirmPassword = confirmPassword;
	});

	async function handleCreate(event: SubmitEvent): Promise<void> {
		event.preventDefault();

		passwordTooShortError = false;
		passwordMismatchError = false;
		error = null;

		if (password.length < minPasswordLength) {
			passwordTooShortError = true;
			return;
		}

		if (password !== confirmPassword) {
			passwordMismatchError = true;
			return;
		}

		submitting = true;

		try {
			await withMinLoadingDelay(
				createUser({
					username,
					displayName,
					password,
					siteRole
				})
			);
			await goto('/admin/users/');
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Failed to create user';
		} finally {
			submitting = false;
		}
	}
</script>

<Table.Root
	title="Add User"
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

		{#if passwordTooShortError}
			<p class="text-sm text-foreground-error">Password must be at least 8 characters</p>
		{/if}

		{#if passwordMismatchError}
			<p class="text-sm text-foreground-error">Passwords do not match</p>
		{/if}

		{#if error}
			<p class="text-sm text-foreground-error">{error}</p>
		{/if}

		<div class="grid grid-cols-2 gap-2">
			<Button href="/admin/users/" variant="secondary">Cancel</Button>
			<Button type="submit" variant="primary" disabled={submitDisabled} loading={submitting}>
				Add
			</Button>
		</div>
	</form>
</Table.Root>

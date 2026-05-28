<script lang="ts">
	import { Button, Field, Input } from '$lib/components/ui';

	type Props = {
		submitLabel: string;
		submittingLabel: string;
		error?: string | null;
		submitting?: boolean;
		onsubmit: (data: { username: string; password: string }) => void | Promise<void>;
	};

	let {
		submitLabel,
		submittingLabel,
		error = null,
		submitting = false,
		onsubmit
	}: Props = $props();

	let username = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let passwordMismatchError = $state(false);
	let passwordTooShortError = $state(false);
	let previousPassword = $state('');
	let previousConfirmPassword = $state('');

	const minPasswordLength = 8;

	$effect(() => {
		if (
			passwordMismatchError &&
			(password !== previousPassword || confirmPassword !== previousConfirmPassword)
		) {
			passwordMismatchError = false;
		}

		if (passwordTooShortError && password !== previousPassword) {
			passwordTooShortError = false;
		}

		previousPassword = password;
		previousConfirmPassword = confirmPassword;
	});

	async function handleSubmit(event: SubmitEvent): Promise<void> {
		event.preventDefault();

		if (password.length < minPasswordLength) {
			passwordTooShortError = true;
			return;
		}

		if (password !== confirmPassword) {
			passwordMismatchError = true;
			return;
		}

		await onsubmit({ username, password });
	}
</script>

<form class="flex flex-col gap-4" onsubmit={handleSubmit}>
	<Field label="Username">
		<Input bind:value={username} autocomplete="username" required />
	</Field>

	<Field label="Password">
		<Input password bind:value={password} autocomplete="new-password" required />
	</Field>

	<Field label="Confirm password">
		<Input password bind:value={confirmPassword} autocomplete="new-password" required />
	</Field>

	{#if passwordTooShortError}
		<p class="text-sm text-error">Password must be at least 8 characters</p>
	{/if}

	{#if passwordMismatchError}
		<p class="text-sm text-error">Passwords do not match</p>
	{/if}

	{#if error}
		<p class="text-sm text-error">{error}</p>
	{/if}

	<Button type="submit" variant="primary" disabled={submitting}>
		{submitting ? submittingLabel : submitLabel}
	</Button>
</form>

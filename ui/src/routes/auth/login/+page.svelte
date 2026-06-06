<script lang="ts">
	import { goto } from '$app/navigation';
	import { ApiError } from '$lib/api';
	import { getSignupStatus, login } from '$lib/api/auth-api';
	import { auth } from '$lib/auth.svelte';
	import { AuthHeader } from '$lib/components';
	import { Button, Field, Input } from '$lib/components/ui';
	import { withMinLoadingDelay } from '$lib/utils';

	let username = $state('');
	let password = $state('');
	let signupEnabled = $state<boolean | undefined>(undefined);
	let submitting = $state(false);
	let error = $state<string | null>(null);

	$effect(() => {
		void getSignupStatus()
			.then((status) => {
				signupEnabled = status.enabled;
			})
			.catch(() => {
				signupEnabled = false;
			});
	});

	async function handleSubmit(event: SubmitEvent): Promise<void> {
		event.preventDefault();
		submitting = true;
		error = null;

		try {
			const user = await withMinLoadingDelay(login({ username, password }));

			auth.setUser(user);

			await goto('/');
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Login failed';
		} finally {
			submitting = false;
		}
	}
</script>

<div class="app-shell page-content justify-center gap-10 py-10">
	<AuthHeader subtitle="Sign in to your account" />

	<form class="flex flex-col gap-8" onsubmit={handleSubmit}>
		<div class="flex flex-col gap-4">
			<Field label="Username">
				<Input bind:value={username} autocomplete="username" required />
			</Field>

			<Field label="Password">
				<Input password bind:value={password} autocomplete="current-password" required />
			</Field>

			{#if error}
				<p class="text-foreground-error text-sm">{error}</p>
			{/if}
		</div>

		<Button type="submit" variant="primary" loading={submitting}>
			Sign in
		</Button>
	</form>

	{#if signupEnabled}
		<p class="text-foreground-alt-2 text-center text-sm">
			No account?
			<a href="/auth/register/" class="text-background-primary">Register</a>
		</p>
	{/if}
</div>

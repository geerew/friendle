<script lang="ts">
	import { goto } from '$app/navigation';
	import { ApiError } from '$lib/api';
	import { getSignupStatus, login } from '$lib/api/auth-api';
	import { auth } from '$lib/auth.svelte';
	import { Logo } from '$lib/components';
	import { Button, Field, Input } from '$lib/components/ui';

	let username = $state('');
	let password = $state('');
	let signupEnabled = $state(false);
	let loading = $state(true);
	let submitting = $state(false);
	let error = $state<string | null>(null);

	$effect(() => {
		getSignupStatus()
			.then((status) => {
				signupEnabled = status.enabled;
			})
			.catch(() => {
				signupEnabled = false;
			})
			.finally(() => {
				loading = false;
			});
	});

	async function handleSubmit(event: SubmitEvent): Promise<void> {
		event.preventDefault();
		submitting = true;
		error = null;

		try {
			const user = await login({ username, password });
			auth.setUser(user);
			await goto('/');
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Login failed';
		} finally {
			submitting = false;
		}
	}
</script>

<div class="app-shell page-content justify-center gap-6 py-10">
	<div class="text-center">
		<Logo />
	</div>

	{#if loading}
		<p class="text-center text-text-muted">Loading…</p>
	{:else}
		<form class="flex flex-col gap-4" onsubmit={handleSubmit}>
			<Field label="Username">
				<Input bind:value={username} autocomplete="username" required />
			</Field>

			<Field label="Password">
				<Input password bind:value={password} autocomplete="current-password" required />
			</Field>

			{#if error}
				<p class="text-sm text-error">{error}</p>
			{/if}

			<Button type="submit" variant="primary" disabled={submitting}>
				{submitting ? 'Signing in…' : 'Sign in'}
			</Button>
		</form>

		{#if signupEnabled}
			<p class="text-center text-sm text-text-muted">
				No account?
				<a href="/auth/register/" class="text-tile-correct">Register</a>
			</p>
		{/if}
	{/if}
</div>

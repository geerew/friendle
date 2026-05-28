<script lang="ts">
	import { goto } from '$app/navigation';
	import { ApiError } from '$lib/api';
	import { register } from '$lib/api/auth-api';
	import { auth } from '$lib/auth.svelte';
	import { AuthRegisterForm } from '$lib/components';

	let submitting = $state(false);
	let error = $state<string | null>(null);

	async function handleSubmit(data: { username: string; password: string }): Promise<void> {
		submitting = true;
		error = null;

		try {
			const user = await register(data);
			auth.setUser(user);
			await goto('/');
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Registration failed';
		} finally {
			submitting = false;
		}
	}
</script>

<div class="app-shell page-content justify-center gap-6 py-10">
	<div class="text-center">
		<h1 class="text-2xl font-bold">Create account</h1>
		<p class="text-text-muted">Join Friendle</p>
	</div>

	<AuthRegisterForm
		submitLabel="Register"
		submittingLabel="Creating…"
		{error}
		{submitting}
		onsubmit={handleSubmit}
	/>

	<p class="text-center text-sm text-text-muted">
		Already have an account?
		<a href="/auth/login/" class="text-tile-correct">Sign in</a>
	</p>
</div>

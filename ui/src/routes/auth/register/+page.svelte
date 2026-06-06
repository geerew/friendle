<script lang="ts">
	import { goto } from '$app/navigation';
	import { ApiError } from '$lib/api';
	import { register } from '$lib/api/auth-api';
	import { auth } from '$lib/auth.svelte';
	import { AuthHeader, AuthRegisterForm } from '$lib/components';

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
	<AuthHeader subtitle="Create your account" />

	<AuthRegisterForm
		submitLabel="Register"
		bind:error
		{submitting}
		onsubmit={handleSubmit}
	/>

	<p class="text-center text-sm text-foreground-alt-2">
		Already have an account?
		<a href="/auth/login/" class="text-background-primary">Sign in</a>
	</p>
</div>

<script lang="ts">
	import { goto } from '$app/navigation';
	import { ApiError } from '$lib/api';
	import { register } from '$lib/api/auth-api';
	import { auth } from '$lib/auth.svelte';
	import { AuthHeader, AuthRegisterForm } from '$lib/components';
	import { withMinLoadingDelay } from '$lib/utils';

	let submitting = $state(false);
	let error = $state<string | null>(null);

	async function handleSubmit(data: { username: string; password: string }): Promise<void> {
		submitting = true;
		error = null;

		try {
			const user = await withMinLoadingDelay(register(data));

			auth.setUser(user);
			await goto('/');
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Registration failed';
		} finally {
			submitting = false;
		}
	}
</script>

<div class="app-shell page-content justify-center gap-10 py-10">
	<AuthHeader subtitle="Create your account" />

	<AuthRegisterForm
		submitLabel="Register"
		bind:error
		{submitting}
		onsubmit={handleSubmit}
	/>

	<p class="text-foreground-alt-2 text-center text-sm">
		Already have an account?
		<a href="/auth/login/" class="text-background-primary">Sign in</a>
	</p>
</div>

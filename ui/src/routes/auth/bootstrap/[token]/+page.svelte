<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { ApiError } from '$lib/api';
	import { bootstrap } from '$lib/api/auth-api';
	import { auth } from '$lib/auth.svelte';
	import { AuthHeader, AuthRegisterForm } from '$lib/components';

	const token = $derived(page.params.token ?? '');

	let submitting = $state(false);
	let error = $state<string | null>(null);

	async function handleSubmit(data: { username: string; password: string }): Promise<void> {
		submitting = true;
		error = null;

		try {
			const user = await bootstrap(token, data);
			auth.setUser(user);
			await goto('/');
		} catch (err) {
			error = err instanceof ApiError ? err.message : 'Bootstrap failed';
		} finally {
			submitting = false;
		}
	}
</script>

<div class="app-shell page-content justify-center gap-6 py-10">
	<AuthHeader subtitle="Create the first administrator account" />

	<AuthRegisterForm
		submitLabel="Create admin"
		bind:error
		{submitting}
		onsubmit={handleSubmit}
	/>
</div>

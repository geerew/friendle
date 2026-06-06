<script lang="ts">
	import '../app.css';

	import { page } from '$app/state';
	import { auth } from '$lib/auth.svelte';
	import { Toaster } from 'svelte-sonner';

	let { children } = $props();

	const isAuthPath = $derived(page.url.pathname.startsWith('/auth'));

	$effect(() => {
		if (isAuthPath) return;
		void auth.load();
	});
</script>

<Toaster theme="dark" richColors />

{#if !isAuthPath}
	{#if auth.initialized && auth.error && !auth.user}
		<div class="app-shell page-content">
			<p class="text-foreground-error">{auth.error}</p>
		</div>
	{:else if auth.initialized && !auth.user}
		<!-- session expired; redirect to login is in flight -->
	{:else}
		{@render children()}
	{/if}
{:else}
	{@render children()}
{/if}

<script lang="ts">
	import '@fontsource/anton';
	import '../app.css';

	import { page } from '$app/state';
	import { auth } from '$lib/auth.svelte';

	let { children } = $props();

	const isAuthPath = $derived(page.url.pathname.startsWith('/auth'));

	$effect(() => {
		page.url.pathname;
		if (isAuthPath) return;
		auth.load();
	});
</script>

{#if !isAuthPath}
	{#if auth.loading}
		<div class="app-shell items-center justify-center py-20 text-text-muted">Loading…</div>
	{:else if auth.error && !auth.user}
		<div class="app-shell page-content">
			<p class="text-error">{auth.error}</p>
		</div>
	{:else if auth.user}
		{@render children()}
	{/if}
{:else}
	{@render children()}
{/if}

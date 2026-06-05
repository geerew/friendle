<script lang="ts">
	import '@fontsource/anton';
	import '../app.css';

	import { page } from '$app/state';
	import { Spinner } from '$lib/components';
	import { auth } from '$lib/auth.svelte';
	import { Toaster } from 'svelte-sonner';

	let { children } = $props();

	const isAuthPath = $derived(page.url.pathname.startsWith('/auth'));

	$effect(() => {
		page.url.pathname;
		if (isAuthPath) return;
		auth.load();
	});
</script>

<Toaster theme="dark" richColors />

{#if !isAuthPath}
	{#if auth.loading}
		<div class="app-shell items-center justify-center py-20">
			<Spinner class="size-6 bg-text-muted" />
		</div>
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

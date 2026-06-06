<script lang="ts">
	import '../app.css';

	import { page } from '$app/state';
	import { AppShell, Spinner } from '$lib/components';
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
	{#if !auth.initialized}
		<AppShell showMenu={false}>
			<div class="flex justify-center pt-8">
				<Spinner class="bg-foreground-alt-2 size-3" />
			</div>
		</AppShell>
	{:else if auth.error && !auth.user}
		<div class="app-shell page-content">
			<p class="text-foreground-error">{auth.error}</p>
		</div>
	{:else if auth.user}
		{@render children()}
	{/if}
{:else}
	{@render children()}
{/if}

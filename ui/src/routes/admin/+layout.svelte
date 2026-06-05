<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/auth.svelte';
	import type { Snippet } from 'svelte';

	type Props = {
		children: Snippet;
	};

	let { children }: Props = $props();

	$effect(() => {
		if (!auth.initialized) return;

		if (!auth.isAdmin) {
			void goto('/');
		}
	});
</script>

{@render children()}

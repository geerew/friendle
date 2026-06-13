<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		label: string;
		error?: string | null;
		children: Snippet;
	};

	let { label, error = null, children }: Props = $props();
</script>

<label class="flex flex-col gap-2">
	<span class="section-title">{label}</span>
	{@render children()}
	{#if error}
		<p class="text-foreground-error text-sm">{error}</p>
	{/if}
</label>

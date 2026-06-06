<script lang="ts">
	import Spinner from './spinner.svelte';
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';

	type Props = {
		loading?: boolean;
		class?: string;
		children: Snippet;
	};

	let { loading = false, class: className = '', children }: Props = $props();
</script>

<div class={cn('relative', className)}>
	<div class:select-none={loading} class:pointer-events-none={loading}>
		{@render children()}
	</div>

	{#if loading}
		<div
			class="absolute inset-0 z-10 flex items-center justify-center rounded-md bg-background/70"
			aria-busy="true"
			aria-live="polite"
		>
			<Spinner class="bg-foreground-alt-2 size-3" />
		</div>
	{/if}
</div>

<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';

	type Props = {
		class?: string;
		children: Snippet;
	};

	let { class: className = '', children }: Props = $props();
</script>

<footer
	class={cn(
		'border-foreground-alt-4 bg-background flex h-16 w-full shrink-0 items-center justify-end gap-2 border-t px-5 py-2.5',
		className
	)}
>
	{@render children()}
</footer>

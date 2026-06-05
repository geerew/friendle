<script lang="ts">
	import { cn } from '$lib/utils';
	import { Dialog, type WithoutChild } from 'bits-ui';
	import type { Snippet } from 'svelte';

	type Props = WithoutChild<Dialog.ContentProps> & {
		class?: string;
		children: Snippet;
	};

	let { class: className = '', children, ...restProps }: Props = $props();
</script>

<Dialog.Content
	class={cn(
		'fixed top-20 left-1/2 z-50 w-[calc(100%-2rem)] max-w-md -translate-x-1/2 overflow-hidden rounded-lg border border-border bg-bg-secondary',
		className
	)}
	{...restProps}
>
	{@render children()}
</Dialog.Content>

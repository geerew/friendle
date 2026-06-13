<script lang="ts">
	import { cn } from '$lib/utils';
	import { DropdownMenu, type DropdownMenuTriggerProps, type WithoutChild } from 'bits-ui';
	import type { Snippet } from 'svelte';

	type Props = WithoutChild<DropdownMenuTriggerProps> & {
		ref?: HTMLButtonElement | null;
		class?: string;
		children: Snippet;
	};

	let { ref = $bindable(null), class: containerClass, children, ...restProps }: Props = $props();
</script>

<DropdownMenu.Trigger
	bind:ref
	class={cn(
		'text-foreground-alt-2 hover:text-foreground inline-flex h-10 w-10 cursor-pointer items-center justify-center rounded transition-colors disabled:cursor-not-allowed disabled:opacity-50',
		containerClass
	)}
	{...restProps}
>
	{@render children()}
</DropdownMenu.Trigger>

<script lang="ts">
	import { cn } from '$lib/utils';
	import { DropdownMenu, type DropdownMenuItemProps, type WithoutChild } from 'bits-ui';
	import type { Snippet } from 'svelte';

	type Props = WithoutChild<DropdownMenuItemProps> & {
		ref?: HTMLButtonElement | null;
		class?: string;
		children: Snippet;
	};

	let { ref = $bindable(null), class: containerClass, children, ...restProps }: Props = $props();
</script>

<DropdownMenu.Item
	bind:ref
	class={cn(
		'inline-flex w-full cursor-pointer items-center rounded-md px-2 py-1.5 text-text-muted transition-colors outline-none hover:bg-bg hover:text-text data-disabled:cursor-default data-disabled:opacity-50 data-disabled:hover:bg-transparent',
		containerClass
	)}
	{...restProps}
>
	{@render children()}
</DropdownMenu.Item>

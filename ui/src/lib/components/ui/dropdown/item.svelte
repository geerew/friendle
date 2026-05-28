<script lang="ts">
	import { cn } from '$lib/utils';
	import { DropdownMenu, type DropdownMenuItemProps, type WithoutChild } from 'bits-ui';
	import type { Snippet } from 'svelte';

	type Variant = 'default' | 'destructive';

	type Props = WithoutChild<DropdownMenuItemProps> & {
		ref?: HTMLButtonElement | null;
		class?: string;
		variant?: Variant;
		children: Snippet;
	};

	let {
		ref = $bindable(null),
		class: containerClass,
		variant = 'default',
		children,
		...restProps
	}: Props = $props();

	const variantClasses: Record<Variant, string> = {
		default:
			'text-text-muted hover:bg-button-primary/25 hover:text-text data-disabled:hover:bg-transparent',
		destructive:
			'text-error-fg hover:bg-error-bg hover:text-text data-disabled:hover:bg-transparent'
	};
</script>

<DropdownMenu.Item
	bind:ref
	class={cn(
		'inline-flex w-full cursor-pointer items-center rounded-md px-2 py-1.5 transition-colors outline-none data-disabled:cursor-default data-disabled:opacity-50',
		variantClasses[variant],
		containerClass
	)}
	{...restProps}
>
	{@render children()}
</DropdownMenu.Item>

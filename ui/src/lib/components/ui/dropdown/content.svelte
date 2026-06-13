<script lang="ts">
	import { cn } from '$lib/utils';
	import { DropdownMenu, type DropdownMenuContentProps, type WithoutChild } from 'bits-ui';
	import type { Snippet } from 'svelte';

	type Props = WithoutChild<DropdownMenuContentProps> & {
		ref?: HTMLDivElement | null;
		class?: string;
		portalProps?: DropdownMenu.PortalProps;
		children: Snippet;
	};

	let {
		ref = $bindable(null),
		class: containerClass,
		portalProps,
		children,
		...restProps
	}: Props = $props();
</script>

<DropdownMenu.Portal {...portalProps}>
	<DropdownMenu.Content
		bind:ref
		align="end"
		sideOffset={4}
		class={cn(
			'border-foreground-alt-4 bg-background-alt-1 z-50 flex min-w-36 flex-col gap-1 rounded-md border p-1.5 text-sm shadow-lg outline-none select-none',
			containerClass
		)}
		{...restProps}
	>
		{@render children?.()}
	</DropdownMenu.Content>
</DropdownMenu.Portal>

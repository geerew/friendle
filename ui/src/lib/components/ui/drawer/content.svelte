<script lang="ts">
	import { cn } from '$lib/utils';
	import type { WithoutChild } from 'bits-ui';
	import type { Snippet } from 'svelte';
	import { Drawer as DrawerPrimitive } from 'vaul-svelte';
	import DrawerOverlay from './overlay.svelte';

	type Props = WithoutChild<DrawerPrimitive.ContentProps> & {
		ref?: HTMLDivElement | null;
		class?: string;
		handleClass?: string;
		portalProps?: DrawerPrimitive.PortalProps;
		children: Snippet;
	};

	let {
		ref = $bindable(null),
		class: containerClass,
		handleClass = '',
		portalProps,
		children,
		...restProps
	}: Props = $props();
</script>

<DrawerPrimitive.Portal {...portalProps}>
	<DrawerOverlay />
	<DrawerPrimitive.Content
		bind:ref
		data-slot="drawer-content"
		class={cn(
			'group/drawer-content fixed z-50 flex w-full max-w-md flex-col border-foreground-alt-4 bg-background',
			'data-[vaul-drawer-direction=bottom]:bottom-0 data-[vaul-drawer-direction=bottom]:left-1/2 data-[vaul-drawer-direction=bottom]:mt-24 data-[vaul-drawer-direction=bottom]:max-h-[80vh] data-[vaul-drawer-direction=bottom]:-translate-x-1/2 data-[vaul-drawer-direction=bottom]:rounded-t-lg data-[vaul-drawer-direction=bottom]:border-t',
			containerClass
		)}
		{...restProps}
	>
		<div
			class={cn(
				'mx-auto mt-2 mb-2 h-1.5 w-[100px] shrink-0 rounded-full bg-foreground-alt-4 group-data-[vaul-drawer-direction=bottom]/drawer-content:block',
				handleClass
			)}
		></div>
		{@render children?.()}
	</DrawerPrimitive.Content>
</DrawerPrimitive.Portal>

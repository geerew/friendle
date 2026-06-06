<script lang="ts">
	import { cn } from '$lib/utils';
	import type { WithElementRef } from 'bits-ui';
	import type { HTMLAttributes } from 'svelte/elements';

	let {
		ref = $bindable(null),
		class: className = '',
		children,
		...restProps
	}: WithElementRef<HTMLAttributes<HTMLDivElement>> = $props();
</script>

<div
	bind:this={ref}
	data-slot="drawer-footer"
		class={cn(
		'flex h-16 w-full shrink-0 items-center justify-end gap-2 border-t border-foreground-alt-4 bg-background px-5 py-2.5',
		className as string
	)}
	{...restProps}
>
	{@render children?.()}
</div>

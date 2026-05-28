<script lang="ts">
	import { cn } from '$lib/utils';
	import type { WithoutChildrenOrChild } from 'bits-ui';
	import type { Snippet } from 'svelte';
	import { Drawer as DrawerPrimitive } from 'vaul-svelte';
	import Button from '../button.svelte';

	type Props = WithoutChildrenOrChild<DrawerPrimitive.CloseProps> & {
		class?: string;
		children: Snippet;
	};

	let { class: containerClass, children, ...restProps }: Props = $props();
</script>

<DrawerPrimitive.Close {...restProps}>
	{#snippet child({ props })}
		<Button {...props} variant="secondary" class={cn('w-full', containerClass)}>
			{@render children()}
		</Button>
	{/snippet}
</DrawerPrimitive.Close>

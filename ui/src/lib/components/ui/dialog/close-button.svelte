<script lang="ts">
	import { cn } from '$lib/utils';
	import type { WithoutChildrenOrChild } from 'bits-ui';
	import { Dialog } from 'bits-ui';
	import type { Snippet } from 'svelte';
	import Button from '../button.svelte';

	type Props = WithoutChildrenOrChild<Dialog.CloseProps> & {
		class?: string;
		disabled?: boolean;
		children: Snippet;
	};

	let { class: className = '', disabled = false, children, ...restProps }: Props = $props();
</script>

<Dialog.Close {...restProps}>
	{#snippet child({ props })}
		<Button {...props} variant="secondary" class={cn('w-24', className)} {disabled}>
			{@render children()}
		</Button>
	{/snippet}
</Dialog.Close>

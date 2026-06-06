<script lang="ts">
	import { TickIcon } from '$lib/components/icons';
	import { cn } from '$lib/utils';
	import { DropdownMenu, type WithoutChild } from 'bits-ui';

	let {
		ref = $bindable(null),
		class: className,
		children: childrenProp,
		...restProps
	}: WithoutChild<DropdownMenu.RadioItemProps> = $props();
</script>

<DropdownMenu.RadioItem
	bind:ref
	class={cn(
		'inline-flex w-full cursor-pointer items-center gap-2.5 rounded-md px-1 py-1 text-sm text-foreground-alt-2 transition-colors select-none hover:bg-background-primary/25 hover:text-foreground',
		className
	)}
	{...restProps}
>
	{#snippet children({ checked })}
		{#if checked}
			<TickIcon class="size-3.5 shrink-0 stroke-2 text-foreground-alt-2" />
		{:else}
			<span class="size-3.5 shrink-0"></span>
		{/if}
		{@render childrenProp?.({ checked })}
	{/snippet}
</DropdownMenu.RadioItem>

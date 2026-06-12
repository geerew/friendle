<script lang="ts">
	import { cn } from '$lib/utils';
	import { Label, RadioGroup, useId, type WithoutChildrenOrChild } from 'bits-ui';

	type Props = WithoutChildrenOrChild<RadioGroup.RootProps> & {
		items: { value: string; label: string; disabled?: boolean }[];
		rootClass?: string;
		itemRowClass?: string;
		itemClass?: string;
		labelClass?: string;
	};

	let {
		value = $bindable(),
		items,
		rootClass,
		itemRowClass,
		itemClass,
		labelClass,
		...restProps
	}: Props = $props();
</script>

<RadioGroup.Root
	bind:value={value as never}
	class={cn('flex flex-col gap-3', rootClass)}
	{...restProps}
>
	{#each items as item (item.value)}
		{@const id = useId()}
		<div class={cn('flex items-center select-none', itemRowClass)}>
			<RadioGroup.Item
				{id}
				value={item.value}
				disabled={item.disabled}
				class={cn(
					'border-foreground-alt-4 bg-background-alt-1 hover:border-foreground-alt-3 data-[state=checked]:border-background-primary size-5 shrink-0 cursor-pointer rounded-full border transition-colors data-disabled:cursor-not-allowed data-disabled:opacity-50 data-[state=checked]:border-[6px]',
					itemClass
				)}
			/>
			<Label.Root for={id} class={cn('text-foreground cursor-pointer ps-3 text-sm', labelClass)}>
				{item.label}
			</Label.Root>
		</div>
	{/each}
</RadioGroup.Root>

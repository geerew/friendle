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

<RadioGroup.Root bind:value={value as never} class={cn('flex flex-col gap-3', rootClass)} {...restProps}>
	{#each items as item (item.value)}
		{@const id = useId()}
		<div class={cn('flex select-none items-center', itemRowClass)}>
			<RadioGroup.Item
				{id}
				value={item.value}
				disabled={item.disabled}
				class={cn(
					'size-5 shrink-0 cursor-pointer rounded-full border border-border bg-bg-secondary transition-colors hover:border-text-muted data-[state=checked]:border-button-primary data-[state=checked]:border-[6px] data-disabled:cursor-not-allowed data-disabled:opacity-50',
					itemClass
				)}
			/>
			<Label.Root
				for={id}
				class={cn('cursor-pointer ps-3 text-sm text-text', labelClass)}
			>
				{item.label}
			</Label.Root>
		</div>
	{/each}
</RadioGroup.Root>

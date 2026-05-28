<script lang="ts">
	import { cn } from '$lib/utils';
	import { RightChevronIcon } from '$lib/components/icons';
	import { Select, type WithoutChildren } from 'bits-ui';

	type Props = WithoutChildren<
		Select.RootProps & {
			placeholder?: string;
			items: { value: string; label: string; disabled?: boolean }[];
			triggerClass?: string;
			contentProps?: Omit<WithoutChildren<Select.ContentProps>, 'class'>;
			contentClass?: string;
			itemClass?: string;
		}
	>;

	let {
		value = $bindable(),
		placeholder,
		items,
		triggerClass,
		contentProps,
		contentClass,
		itemClass,
		...restProps
	}: Props = $props();

	const selectedLabel = $derived(items.find((item) => item.value === value)?.label);
</script>

<Select.Root bind:value={value as never} {...restProps}>
	<Select.Trigger
		class={cn(
			'field inline-flex shrink-0 items-center justify-between gap-1 hover:cursor-pointer hover:border-text-muted data-[state=open]:border-tile-correct data-[state=open]:ring-2 data-[state=open]:ring-tile-correct [&[data-state=open]>svg]:rotate-90',
			triggerClass
		)}
	>
		<span class={selectedLabel ? '' : 'text-text-muted'}>
			{selectedLabel ?? placeholder}
		</span>
		<RightChevronIcon class="size-3.5 shrink-0 stroke-2 text-text-muted transition-transform duration-200" />
	</Select.Trigger>

	<Select.Portal>
		<Select.Content
			class={cn(
				'z-50 w-(--bits-select-anchor-width) min-w-(--bits-select-anchor-width) rounded-md border border-border bg-bg-secondary py-1.5 shadow-lg outline-none select-none',
				contentClass
			)}
			{...contentProps}
		>
			<Select.Viewport>
				{#each items as { value: itemValue, label, disabled } (itemValue)}
					<Select.Item
						value={itemValue}
						{label}
						{disabled}
						class={cn(
							'flex h-9 w-full cursor-pointer items-center px-2 text-sm text-text-muted transition-colors outline-none select-none hover:bg-button-primary/25 hover:text-text data-disabled:cursor-not-allowed data-disabled:opacity-50 data-highlighted:bg-button-primary/25 data-highlighted:text-text',
							itemClass
						)}
					>
						{#snippet children({ selected })}
							{label}
							{#if selected}
								<span class="ml-auto text-button-primary">✓</span>
							{/if}
						{/snippet}
					</Select.Item>
				{/each}
			</Select.Viewport>
		</Select.Content>
	</Select.Portal>
</Select.Root>

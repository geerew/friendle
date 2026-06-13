<script lang="ts">
	import { RightChevronIcon, SortAscendingIcon, SortDescendingIcon } from '$lib/components/icons';
	import { Dropdown } from '$lib/components/ui';
	import type { SortColumns, SortDirection } from '$lib/types/sort';

	type Props = {
		columns: SortColumns;
		selectedColumn: string;
		selectedDirection: SortDirection;
		disabled?: boolean;
		onUpdate?: () => void;
	};

	let {
		columns,
		selectedColumn = $bindable(),
		selectedDirection = $bindable(),
		disabled = false,
		onUpdate
	}: Props = $props();

	selectedDirection = selectedDirection || 'desc';

	$effect(() => {
		if (columns.length === 0) return;

		if (!selectedColumn || !columns.find((column) => column.column === selectedColumn)) {
			selectedColumn = columns[0].column;
		}
	});

	const selectedColumnMeta = $derived(columns.find((column) => column.column === selectedColumn));
</script>

<Dropdown.Root>
	<Dropdown.Trigger
		class="field hover:border-foreground-alt-3 data-[state=open]:border-background-primary data-[state=open]:ring-background-primary inline-flex h-9 w-36 shrink-0 items-center justify-between gap-1 px-2 hover:cursor-pointer data-[state=open]:ring-2 [&[data-state=open]>svg:last-child]:rotate-90"
		{disabled}
	>
		<div class="flex min-w-0 items-center gap-1.5">
			{#if selectedDirection === 'asc'}
				<SortAscendingIcon class="size-4 shrink-0" />
			{:else}
				<SortDescendingIcon class="size-4 shrink-0" />
			{/if}
			<span class="truncate">{selectedColumnMeta?.label ?? 'Sort'}</span>
		</div>
		<RightChevronIcon
			class="text-foreground-alt-2 size-3.5 shrink-0 stroke-2 transition-transform duration-200"
		/>
	</Dropdown.Trigger>

	<Dropdown.Content class="w-36 max-w-36 min-w-0" sideOffset={8}>
		<div class="flex flex-col gap-1.5">
			<Dropdown.RadioGroup bind:value={selectedColumn}>
				{#each columns as column (column.column)}
					<Dropdown.RadioItem
						value={column.column}
						onclick={() => {
							if (selectedColumn === column.column) return;
							onUpdate?.();
						}}
					>
						{column.label}
					</Dropdown.RadioItem>
				{/each}
			</Dropdown.RadioGroup>
		</div>

		<Dropdown.Separator />

		<div class="flex flex-col gap-1.5">
			<Dropdown.RadioGroup bind:value={selectedDirection}>
				<Dropdown.RadioItem
					value="asc"
					onclick={() => {
						if (selectedDirection === 'asc') return;
						onUpdate?.();
					}}
				>
					<SortAscendingIcon class="text-foreground-alt-2 size-4" />
					{selectedColumnMeta?.asc ?? 'Ascending'}
				</Dropdown.RadioItem>

				<Dropdown.RadioItem
					value="desc"
					onclick={() => {
						if (selectedDirection === 'desc') return;
						onUpdate?.();
					}}
				>
					<SortDescendingIcon class="text-foreground-alt-2 size-4" />
					{selectedColumnMeta?.desc ?? 'Descending'}
				</Dropdown.RadioItem>
			</Dropdown.RadioGroup>
		</div>
	</Dropdown.Content>
</Dropdown.Root>

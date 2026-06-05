<script lang="ts">
	import { SelectPaginationPerPage } from '$lib/models/pagination-model';
	import { LeftChevronIcon, RightChevronIcon } from '$lib/components/icons';
	import { Select } from '$lib/components/ui';
	import { cn } from '$lib/utils';
	import { Pagination } from 'bits-ui';

	type Props = {
		count: number;
		page: number;
		onPageChange: () => void;
		perPage: number;
		onPerPageChange: () => void;
		selectTriggerClass?: string;
		showPerPageSelect?: boolean;
	};

	let {
		count,
		page = $bindable(),
		onPageChange,
		perPage = $bindable(),
		onPerPageChange,
		selectTriggerClass,
		showPerPageSelect = true
	}: Props = $props();

	let perPageValue = $state(`${perPage}`);

	$effect(() => {
		perPageValue = `${perPage}`;
	});

	const pageButtonClass =
		'inline-flex size-10 items-center justify-center rounded-lg text-sm font-medium text-text-muted transition-colors select-none hover:cursor-pointer hover:bg-button-primary/25 hover:text-text data-selected:bg-button-primary data-selected:text-white';

	const navButtonClass =
		'inline-flex h-10 flex-row items-center justify-center gap-1 rounded-lg px-2 text-sm font-medium text-text-muted transition-colors select-none hover:cursor-pointer hover:bg-button-primary/25 hover:text-text disabled:cursor-not-allowed disabled:opacity-50 hover:disabled:bg-transparent';
</script>

<Pagination.Root {count} {perPage} bind:page class="flex w-full justify-center" {onPageChange}>
	{#snippet children({ pages, range })}
		<div class="flex w-full flex-col gap-4">
			{#if count > perPage}
				<div class="flex items-center justify-center gap-3">
					<Pagination.PrevButton class={navButtonClass}>
						<LeftChevronIcon class="size-5 stroke-2" />
						<span class="text-xs tracking-wide uppercase">Previous</span>
					</Pagination.PrevButton>

					<div class="flex items-center gap-2">
						{#each pages as pageItem (pageItem.key)}
							{#if pageItem.type === 'ellipsis'}
								<div class="px-1 text-sm font-medium text-text-muted select-none">…</div>
							{:else}
								<Pagination.Page page={pageItem} class={pageButtonClass}>
									{pageItem.value}
								</Pagination.Page>
							{/if}
						{/each}
					</div>

					<Pagination.NextButton class={navButtonClass}>
						<span class="text-xs tracking-wide uppercase">Next</span>
						<RightChevronIcon class="size-5 stroke-2" />
					</Pagination.NextButton>
				</div>
			{/if}

			<div
				class={cn(
					'grid w-full gap-3',
					showPerPageSelect ? 'grid-cols-2' : 'grid-cols-1'
				)}
			>
				{#if showPerPageSelect}
					<div class="flex items-center">
						<Select
							type="single"
							items={SelectPaginationPerPage}
							bind:value={perPageValue}
							contentProps={{ sideOffset: 8, loop: true }}
							triggerClass={cn('w-24', selectTriggerClass)}
							onValueChange={(v) => {
								perPage = +v;

								if (page > Math.ceil(count / perPage)) {
									page = Math.max(1, Math.ceil(count / perPage));
								}

								onPerPageChange();
							}}
						/>
					</div>
				{/if}

				<p
					class={cn(
						'flex min-w-0 items-center text-sm whitespace-nowrap text-text-muted',
						showPerPageSelect ? 'justify-end text-end' : 'justify-center text-center'
					)}
				>
					{range.start} - {range.end} / {count}
				</p>
			</div>
		</div>
	{/snippet}
</Pagination.Root>

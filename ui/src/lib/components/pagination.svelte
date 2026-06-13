<script lang="ts">
	import { SelectPaginationPerPage } from '$lib/models/pagination-model';
	import { LeftChevronIcon, RightChevronIcon } from '$lib/components/icons';
	import { Select } from '$lib/components/ui';
	import { cn } from '$lib/utils';
	import { Pagination } from 'bits-ui';

	type Props = {
		count: number;
		page: number;
		onPageChange?: () => void;
		perPage: number;
		onPerPageChange?: () => void;
		selectTriggerClass?: string;
		showPerPageSelect?: boolean;
		minimal?: boolean;
	};

	let {
		count,
		page = $bindable(),
		onPageChange = () => {},
		perPage = $bindable(),
		onPerPageChange = () => {},
		selectTriggerClass,
		showPerPageSelect = true,
		minimal = false
	}: Props = $props();

	let perPageValue = $state(`${perPage}`);

	$effect(() => {
		perPageValue = `${perPage}`;
	});
</script>

<Pagination.Root {count} {perPage} bind:page class="flex w-full justify-center" {onPageChange}>
	{#snippet children({ pages, range })}
		<div class="flex w-full flex-col gap-4">
			{#if count > perPage}
				<div class="grid w-full grid-cols-[1fr_auto_1fr] items-center gap-3">
					<div class="flex justify-end">
						<Pagination.PrevButton
							class="text-foreground-alt-2 hover:bg-background-primary/25 hover:text-foreground inline-flex h-10 flex-row items-center justify-center gap-1 rounded-lg px-2 text-sm font-medium transition-colors select-none hover:cursor-pointer disabled:cursor-not-allowed disabled:opacity-50 hover:disabled:bg-transparent"
						>
							<LeftChevronIcon class="size-5 stroke-2" />
							<span class="text-xs tracking-wide uppercase">Prev</span>
						</Pagination.PrevButton>
					</div>

					<div class="flex items-center justify-center gap-2">
						{#each pages as pageItem (pageItem.key)}
							{#if pageItem.type === 'ellipsis'}
								<div class="text-foreground-alt-2 px-1 text-sm font-medium select-none">…</div>
							{:else}
								<Pagination.Page
									page={pageItem}
									class="text-foreground-alt-2 hover:bg-background-primary/25 hover:text-foreground data-selected:bg-background-primary inline-flex size-10 items-center justify-center rounded-lg text-sm font-medium transition-colors select-none hover:cursor-pointer data-selected:text-white"
								>
									{pageItem.value}
								</Pagination.Page>
							{/if}
						{/each}
					</div>

					<div class="flex justify-start">
						<Pagination.NextButton
							class="text-foreground-alt-2 hover:bg-background-primary/25 hover:text-foreground inline-flex h-10 flex-row items-center justify-center gap-1 rounded-lg px-2 text-sm font-medium transition-colors select-none hover:cursor-pointer disabled:cursor-not-allowed disabled:opacity-50 hover:disabled:bg-transparent"
						>
							<span class="text-xs tracking-wide uppercase">Next</span>
							<RightChevronIcon class="size-5 stroke-2" />
						</Pagination.NextButton>
					</div>
				</div>
			{/if}

			{#if !minimal}
				<div class={cn('grid w-full gap-3', showPerPageSelect ? 'grid-cols-2' : 'grid-cols-1')}>
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
							'text-foreground-alt-2 flex min-w-0 items-center text-sm whitespace-nowrap',
							showPerPageSelect ? 'justify-end text-end' : 'justify-center text-center'
						)}
					>
						{range.start} - {range.end} / {count}
					</p>
				</div>
			{:else if count > 0}
				<p
					class="text-foreground-alt-2 flex min-w-0 items-center justify-center text-center text-sm whitespace-nowrap"
				>
					{range.start} - {range.end} / {count}
				</p>
			{/if}
		</div>
	{/snippet}
</Pagination.Root>

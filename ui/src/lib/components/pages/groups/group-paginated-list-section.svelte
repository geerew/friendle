<script lang="ts">
	import { LoadingOverlay, Pagination, Spinner } from '$lib/components';
	import type { Snippet } from 'svelte';

	type Props = {
		itemCount: number;
		totalItems: number;
		page?: number;
		perPage?: number;
		loading: boolean;
		emptyMessage: string;
		list: Snippet;
	};

	let {
		itemCount,
		totalItems,
		perPage = $bindable(10),
		page = $bindable(1),
		loading,
		emptyMessage,
		list
	}: Props = $props();
</script>

<div class="flex flex-col gap-6">
	{#if loading && itemCount === 0}
		<div class="flex min-h-24 items-center justify-center">
			<Spinner class="bg-foreground-alt-2 size-3" />
		</div>
	{:else if itemCount === 0}
		<div class="flex min-h-16 items-center justify-center">
			<p class="text-foreground-alt-2 text-sm italic">{emptyMessage}</p>
		</div>
	{:else}
		<LoadingOverlay {loading}>
			<div class="px-2">
				{@render list()}
			</div>
		</LoadingOverlay>

		{#if totalItems > perPage}
			<Pagination
				count={totalItems}
				bind:page
				{perPage}
				minimal
				showPerPageSelect={false}
			/>
		{/if}
	{/if}
</div>

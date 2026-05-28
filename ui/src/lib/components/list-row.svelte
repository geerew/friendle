<script lang="ts">
	import { RightChevronIcon } from '$lib/components/icons';
	import type { Snippet } from 'svelte';

	type Props = {
		href?: string;
		onclick?: () => void;
		title: string;
		subtitle?: string;
		trailing?: Snippet;
	};

	let { href, onclick, title, subtitle, trailing }: Props = $props();

	const interactive =
		'flex w-full cursor-pointer items-center gap-3 rounded border border-border bg-bg-secondary px-4 py-3 text-left transition-colors hover:border-text-muted';

	const staticRow =
		'flex w-full items-center gap-3 rounded border border-border bg-bg-secondary px-4 py-3 text-left';
</script>

{#snippet content()}
	<div class="flex min-w-0 flex-1 flex-col gap-1">
		<div class="truncate font-semibold">{title}</div>
		{#if subtitle}
			<div class="truncate text-sm text-text-muted">{subtitle}</div>
		{/if}
	</div>
	{#if trailing}
		<div class="shrink-0">{@render trailing()}</div>
	{:else if href}
		<RightChevronIcon class="h-5 w-5 shrink-0 stroke-2 text-text-muted" />
	{/if}
{/snippet}

{#if href}
	<a {href} class={interactive}>
		{@render content()}
	</a>
{:else if onclick}
	<button type="button" class={interactive} {onclick}>
		{@render content()}
	</button>
{:else}
	<div class={staticRow}>
		{@render content()}
	</div>
{/if}

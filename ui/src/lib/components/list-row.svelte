<script lang="ts">
	import { RightChevronIcon } from '$lib/components/icons';
	import type { Snippet } from 'svelte';

	type Props = {
		href?: string;
		onclick?: () => void;
		title: string;
		subtitle?: string;
		meta?: Snippet;
		trailing?: Snippet;
	};

	let { href, onclick, title, subtitle, meta, trailing }: Props = $props();

	const interactive =
		'flex w-full cursor-pointer items-center gap-3 rounded border border-foreground-alt-4 bg-background-alt-1 px-4 py-3 text-left transition-colors hover:border-foreground-alt-3';

	const staticRow =
		'flex w-full items-center gap-3 rounded border border-foreground-alt-4 bg-background-alt-1 px-4 py-3 text-left';
</script>

{#snippet content()}
	<div class="flex min-w-0 flex-1 flex-col gap-1">
		<div class="truncate font-semibold">{title}</div>
		{#if meta}
			<div class="flex min-w-0 flex-wrap items-center gap-2">{@render meta()}</div>
		{:else if subtitle}
			<div class="truncate text-sm text-foreground-alt-2">{subtitle}</div>
		{/if}
	</div>
	{#if trailing}
		<div class="shrink-0">{@render trailing()}</div>
	{:else if href}
		<RightChevronIcon class="h-5 w-5 shrink-0 stroke-2 text-foreground-alt-2" />
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

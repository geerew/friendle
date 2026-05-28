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

	const base =
		'flex w-full items-center gap-3 rounded border border-border bg-bg-secondary px-4 py-3 text-left transition-colors hover:border-text-muted';
</script>

{#if href}
	<a {href} class={base}>
		<div class="min-w-0 flex-1">
			<div class="truncate font-semibold">{title}</div>
			{#if subtitle}
				<div class="truncate text-sm text-text-muted">{subtitle}</div>
			{/if}
		</div>
		{#if trailing}
			<div class="shrink-0">{@render trailing()}</div>
		{:else}
			<RightChevronIcon class="h-5 w-5 shrink-0 stroke-2 text-text-muted" />
		{/if}
	</a>
{:else}
	<button type="button" class={base} {onclick}>
		<div class="min-w-0 flex-1">
			<div class="truncate font-semibold">{title}</div>
			{#if subtitle}
				<div class="truncate text-sm text-text-muted">{subtitle}</div>
			{/if}
		</div>
		{#if trailing}
			<div class="shrink-0">{@render trailing()}</div>
		{/if}
	</button>
{/if}

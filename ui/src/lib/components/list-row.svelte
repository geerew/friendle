<script lang="ts">
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
			<svg viewBox="0 0 24 24" class="h-5 w-5 shrink-0 text-text-muted" fill="none" stroke="currentColor" stroke-width="2">
				<path d="M9 18l6-6-6-6" />
			</svg>
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

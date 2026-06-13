<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';

	type Props = {
		label: string;
		leading?: Snippet;
		trailing?: Snippet;
		href?: string;
		class?: string;
		wrapLabel?: boolean;
		align?: 'center' | 'start';
	};

	let {
		label,
		leading,
		trailing,
		href,
		class: className = '',
		wrapLabel = false,
		align = 'center'
	}: Props = $props();

	const rowClass = $derived(
		cn(
			'flex w-full gap-3 rounded-md px-3 text-left',
			wrapLabel ? 'min-h-11 py-3' : 'h-11 items-center',
			href ? 'hover:bg-background-alt-1 transition-all' : 'cursor-default',
			className
		)
	);
	const innerClass = $derived(
		cn('flex w-full min-w-0 gap-3', align === 'start' ? 'items-start' : 'items-center')
	);
	const labelClass = $derived(
		cn(
			'text-foreground-alt-1 min-w-0 flex-1 text-base leading-5 font-medium',
			wrapLabel ? 'wrap-break-word' : 'truncate'
		)
	);
</script>

{#if href}
	<a {href} class={rowClass}>
		<div class={innerClass}>
			{#if leading}
				<div class="flex w-2 shrink-0 items-center justify-center">
					{@render leading()}
				</div>
			{/if}
			<span class={labelClass}>
				{label}
			</span>
			{#if trailing}
				<div class="flex shrink-0 items-center gap-1.5">
					{@render trailing()}
				</div>
			{/if}
		</div>
	</a>
{:else}
	<div class={rowClass}>
		<div class={innerClass}>
			{#if leading}
				<div class="flex w-2 shrink-0 items-center justify-center">
					{@render leading()}
				</div>
			{/if}
			<span class={labelClass}>
				{label}
			</span>
			{#if trailing}
				<div class="flex shrink-0 items-center gap-1.5">
					{@render trailing()}
				</div>
			{/if}
		</div>
	</div>
{/if}

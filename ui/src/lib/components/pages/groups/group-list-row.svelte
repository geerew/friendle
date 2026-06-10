<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';

	type Props = {
		name: string;
		trailing: Snippet;
		href?: string;
		class?: string;
		wrapName?: boolean;
		align?: 'center' | 'start';
	};

	let {
		name,
		trailing,
		href,
		class: className = '',
		wrapName = false,
		align = 'center'
	}: Props = $props();

	const rowClass = $derived(
		cn(
			'flex w-full gap-3 rounded-md px-2 py-3 text-left',
			href ? 'hover:bg-background-alt-1 transition-all' : 'cursor-default',
			className
		)
	);
	const innerClass = $derived(
		cn('flex w-full gap-3 px-1', align === 'start' ? 'items-start' : 'items-center')
	);
	const nameClass = $derived(
		cn(
			'text-foreground-alt-1 min-w-0 flex-1 text-base leading-5 font-medium',
			wrapName ? 'wrap-break-word' : 'truncate'
		)
	);
</script>

{#if href}
	<a {href} class={rowClass}>
		<div class={innerClass}>
			<span class={nameClass}>
				{name}
			</span>
			{@render trailing()}
		</div>
	</a>
{:else}
	<div class={rowClass}>
		<div class={innerClass}>
			<span class={nameClass}>
				{name}
			</span>
			{@render trailing()}
		</div>
	</div>
{/if}

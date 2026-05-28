<script lang="ts">
	import { cn } from '$lib/utils';

	type Props = {
		text?: string;
		class?: string;
		href?: string;
		variant?: 'default' | 'header';
	};

	let { text = 'Friendle', class: className = '', href, variant = 'default' }: Props = $props();

	const classes = $derived(
		cn(
			variant === 'header'
				? href
					? 'truncate text-sm font-semibold tracking-[0.18em] text-text-muted uppercase no-underline transition-colors hover:text-text'
					: 'truncate cursor-default select-none text-sm font-semibold tracking-[0.18em] text-text-muted uppercase'
				: 'font-[Anton,sans-serif] text-4xl font-normal tracking-[0.2em] uppercase text-text no-underline',
			className
		)
	);
</script>

{#if href}
	<a {href} class={classes}>{text}</a>
{:else if variant === 'default'}
	<h1 class={classes}>{text}</h1>
{:else}
	<span class={classes}>{text}</span>
{/if}

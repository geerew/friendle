<script lang="ts">
	import { cn } from '$lib/utils';

	type Props = {
		class?: string;
		href?: string;
		variant?: 'default' | 'header';
	};

	let { class: className = '', href, variant = 'default' }: Props = $props();

	const classes = $derived(
		cn(
			variant === 'header'
				? 'text-sm font-semibold tracking-[0.18em] text-text-muted uppercase no-underline transition-colors hover:text-text'
				: 'font-[Anton,sans-serif] text-3xl font-normal tracking-[0.2em] uppercase text-text no-underline',
			className
		)
	);
</script>

{#if href}
	<a {href} class={classes}>Friendle</a>
{:else}
	<h1 class={classes}>Friendle</h1>
{/if}
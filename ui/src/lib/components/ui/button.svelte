<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';

	type Variant = 'primary' | 'secondary' | 'ghost' | 'destructive';
	type Size = 'default' | 'inline' | 'icon';

	type Props = {
		href?: string;
		variant?: Variant;
		size?: Size;
		class?: string;
		disabled?: boolean;
		type?: 'button' | 'submit' | 'reset';
		onclick?: (event: MouseEvent) => void;
		'aria-label'?: string;
		children: Snippet;
	};

	let {
		href,
		variant = 'secondary',
		size = 'default',
		class: className = '',
		disabled = false,
		type = 'button',
		onclick,
		'aria-label': ariaLabel,
		children
	}: Props = $props();

	const base =
		'inline-flex shrink-0 cursor-pointer items-center justify-center rounded text-sm transition-colors disabled:cursor-not-allowed disabled:opacity-50';

	const sizes: Record<Size, string> = {
		default: 'w-full px-4 py-3 font-semibold tracking-wide uppercase',
		inline: 'h-auto w-auto px-2 py-1 font-normal normal-case tracking-normal',
		icon: 'h-9 w-9 min-w-9 p-0 font-semibold uppercase'
	};

	const variants: Record<Variant, string> = {
		primary: 'bg-button-primary text-white hover:brightness-110',
		secondary: 'bg-button-secondary text-white hover:brightness-110',
		ghost: 'bg-transparent text-text-muted hover:bg-bg-secondary hover:text-text',
		destructive: 'bg-error text-white hover:brightness-110'
	};

	const classes = $derived(cn(base, sizes[size], variants[variant], className));
</script>

{#if href}
	<a {href} class={classes} aria-label={ariaLabel}>
		{@render children()}
	</a>
{:else}
	<button {type} class={classes} {disabled} {onclick} aria-label={ariaLabel}>
		{@render children()}
	</button>
{/if}

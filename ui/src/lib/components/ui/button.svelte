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
		'inline-flex shrink-0 cursor-pointer items-center justify-center rounded text-sm transition-all duration-200 disabled:cursor-not-allowed disabled:opacity-50';

	const sizes: Record<Size, string> = {
		default: 'w-full px-4 py-3 font-semibold tracking-wide uppercase',
		inline: 'h-auto w-auto px-2 py-1 font-normal normal-case tracking-normal',
		icon: 'h-9 w-9 min-w-9 p-0 font-semibold uppercase'
	};

	const variantClasses = $derived.by(() => {
		switch (variant) {
			case 'primary':
				return 'bg-button-primary text-white enabled:hover:brightness-110';
			case 'secondary':
				return 'bg-button-secondary text-white enabled:hover:brightness-110';
			case 'ghost':
				return 'bg-transparent text-text-muted enabled:hover:text-text';
			case 'destructive':
				if (size === 'inline') {
					return 'bg-transparent text-error-fg enabled:hover:bg-error-bg enabled:hover:text-text';
				}

				return 'bg-error-bg text-text enabled:hover:bg-error-bg-hover';
		}
	});

	const classes = $derived(cn(base, sizes[size], variantClasses, className));
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

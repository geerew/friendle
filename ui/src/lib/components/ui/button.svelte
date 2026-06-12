<script lang="ts">
	import IconSpinner from '../icon-spinner.svelte';
	import Spinner from '../spinner.svelte';
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';

	type Variant = 'primary' | 'secondary' | 'ghost' | 'destructive';
	type Size = 'default' | 'inline' | 'icon';

	type Props = {
		href?: string;
		target?: string;
		rel?: string;
		variant?: Variant;
		size?: Size;
		class?: string;
		disabled?: boolean;
		loading?: boolean;
		type?: 'button' | 'submit' | 'reset';
		onclick?: (event: MouseEvent) => void;
		'aria-label'?: string;
		children: Snippet;
	};

	let {
		href,
		target,
		rel,
		variant = 'secondary',
		size = 'default',
		class: className = '',
		disabled = false,
		loading = false,
		type = 'button',
		onclick,
		'aria-label': ariaLabel,
		children
	}: Props = $props();

	const isDisabled = $derived(disabled || loading);
	const useIconSpinner = $derived(size === 'inline' || size === 'icon');

	const base =
		'inline-flex shrink-0 cursor-pointer items-center justify-center rounded transition-all duration-200 disabled:cursor-not-allowed disabled:opacity-50';

	const sizes: Record<Size, string> = {
		default: 'h-10 min-h-10 w-full px-4 text-sm font-semibold tracking-wide uppercase',
		inline: 'h-auto min-h-0 w-auto px-2 py-1 text-sm font-normal normal-case tracking-normal',
		icon: 'h-9 min-h-9 w-9 min-w-9 p-0 text-sm font-semibold uppercase'
	};

	const isLink = $derived(href != null && href !== '');

	const variantClasses = $derived.by(() => {
		const brightHover = isLink ? 'hover:brightness-110' : 'enabled:hover:brightness-110';

		switch (variant) {
			case 'primary':
				return cn('bg-background-primary text-white', brightHover);
			case 'secondary':
				return cn('bg-background-alt-3 text-foreground', brightHover);
			case 'ghost':
				return 'bg-transparent text-foreground-alt-2 hover:text-foreground disabled:hover:text-foreground-alt-2';
			case 'destructive':
				if (size === 'inline') {
					const inlineHover = isLink
						? 'hover:bg-background-error hover:text-foreground'
						: 'enabled:hover:bg-background-error enabled:hover:text-foreground';

					return cn('bg-transparent text-foreground-error-alt-1', inlineHover);
				}

				const solidHover = isLink
					? 'hover:bg-background-error-alt-1'
					: 'enabled:hover:bg-background-error-alt-1';

				return cn('bg-background-error text-foreground', solidHover);
		}
	});

	const spinnerClass = $derived.by(() => {
		if (variant === 'ghost' || (variant === 'destructive' && size === 'inline')) {
			return 'size-2 bg-foreground-alt-2';
		}

		return 'size-2 bg-white/70';
	});

	const classes = $derived(cn(base, sizes[size], variantClasses, className));
</script>

{#snippet buttonContents()}
	{#if loading}
		<Spinner class={spinnerClass} />
	{:else}
		{@render children()}
	{/if}
{/snippet}

{#if href}
	<a {href} {target} {rel} class={classes} aria-label={ariaLabel} {onclick}>
		{@render buttonContents()}
	</a>
{:else}
	<button
		{type}
		class={classes}
		disabled={isDisabled}
		{onclick}
		aria-label={ariaLabel}
		aria-busy={loading}
	>
		{@render buttonContents()}
	</button>
{/if}

<script lang="ts">
	import type { Snippet } from 'svelte';

	type Variant = 'primary' | 'secondary' | 'ghost' | 'danger';

	type Props = {
		href?: string;
		variant?: Variant;
		class?: string;
		disabled?: boolean;
		type?: 'button' | 'submit' | 'reset';
		onclick?: (event: MouseEvent) => void;
		children: Snippet;
	};

	let {
		href,
		variant = 'secondary',
		class: className = '',
		disabled = false,
		type = 'button',
		onclick,
		children
	}: Props = $props();

	const base =
		'inline-flex w-full items-center justify-center rounded px-4 py-3 text-sm font-semibold tracking-wide uppercase transition-colors disabled:cursor-not-allowed disabled:opacity-50';

	const variants: Record<Variant, string> = {
		primary: 'bg-button-primary text-white hover:brightness-110',
		secondary: 'bg-button-secondary text-white hover:brightness-110',
		ghost: 'bg-transparent text-text-muted hover:text-text',
		danger: 'bg-error text-white hover:brightness-110'
	};

	const classes = $derived(`${base} ${variants[variant]} ${className}`);
</script>

{#if href}
	<a {href} class={classes}>
		{@render children()}
	</a>
{:else}
	<button {type} class={classes} {disabled} {onclick}>
		{@render children()}
	</button>
{/if}

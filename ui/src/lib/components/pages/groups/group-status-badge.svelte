<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';

	type Variant = 'admin' | 'member' | 'pending' | 'rejected';

	type Props = {
		variant: Variant;
		label: string;
		icon: Snippet;
	};

	let { variant, label, icon }: Props = $props();

	const variantClass: Record<
		Variant,
		{ root: string; icon: string; divider: string; label: string }
	> = {
		admin: {
			root: 'border-background-primary-alt-2',
			icon: 'bg-background-primary-alt-2 text-foreground',
			divider: 'border-background-primary-alt-2',
			label: 'bg-background-primary-alt-1 text-foreground'
		},
		member: {
			root: 'border-foreground-alt-3',
			icon: 'bg-foreground-alt-3 text-foreground',
			divider: 'border-foreground-alt-3',
			label: 'bg-background-alt-2 text-foreground-alt-1'
		},
		pending: {
			root: 'border-foreground-alt-4',
			icon: 'bg-foreground-alt-4 text-foreground-alt-2',
			divider: 'border-foreground-alt-4',
			label: 'bg-background-alt-1 text-foreground-alt-2'
		},
		rejected: {
			root: 'border-foreground-error/30',
			icon: 'bg-background-error text-foreground-error',
			divider: 'border-foreground-error/30',
			label: 'bg-background-error-alt-1 text-foreground-error'
		}
	};

	const styles = $derived(variantClass[variant]);
</script>

<span
	class={cn(
		'inline-flex w-auto flex-row items-stretch overflow-hidden rounded-md border text-xs whitespace-nowrap',
		styles.root
	)}
>
	<span
		class={cn(
			'flex items-center self-stretch border-r px-1.5 py-0.5',
			styles.icon,
			styles.divider
		)}
	>
		{@render icon()}
	</span>
	<span class={cn('flex items-center px-2 py-0.5', styles.label)}>
		{label}
	</span>
</span>

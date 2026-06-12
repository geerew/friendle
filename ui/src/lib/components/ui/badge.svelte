<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';
	import type { HTMLAttributes } from 'svelte/elements';

	type Props = HTMLAttributes<HTMLSpanElement> & {
		class?: string;
		children?: Snippet;
	};

	let { class: className, children, ...restProps }: Props = $props();
</script>

<span
	class={cn(
		'border-foreground-alt-4 bg-background-alt-1 text-foreground-alt-2 inline-flex w-auto flex-row items-center rounded-md border px-2.5 py-0.5 text-xs whitespace-nowrap',
		className
	)}
	{...restProps}
>
	{@render children?.()}
</span>

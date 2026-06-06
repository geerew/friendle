<script lang="ts">
	import { cn } from '$lib/utils';

	type Props = {
		open?: boolean;
		class?: string;
	};

	let { open = false, class: className }: Props = $props();
</script>

<span
	class={cn(
		'relative inline-flex h-6 w-6 flex-col items-center justify-center gap-[5px]',
		className
	)}
	aria-hidden="true"
>
	<span
		class={cn(
			'h-0.5 w-full rounded-full bg-current transition-transform duration-200 ease-out',
			open && 'translate-y-[7px] rotate-45'
		)}
	></span>
	<span
		class={cn(
			'h-0.5 w-full rounded-full bg-current transition-opacity duration-200 ease-out',
			open && 'opacity-0'
		)}
	></span>
	<span
		class={cn(
			'h-0.5 w-full rounded-full bg-current transition-transform duration-200 ease-out',
			open && '-translate-y-[7px] -rotate-45'
		)}
	></span>
</span>

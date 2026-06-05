<script lang="ts">
	import { cn } from '$lib/utils';

	type Props = {
		class?: string;
		dashed?: boolean;
	};

	let { class: className = '', dashed = false }: Props = $props();
</script>

<div
	role="separator"
	class={cn(
		'shrink-0',
		dashed ? 'border-t border-dashed border-border' : 'h-px bg-border',
		className
	)}
></div>

<script lang="ts">
	import { cn } from '$lib/utils';

	type Props = {
		class?: string;
	};

	let { class: className = 'size-4' }: Props = $props();
</script>

<div
	class={cn(
		'shrink-0 animate-spin rounded-full border-2 border-current border-t-transparent',
		className
	)}
	role="status"
	aria-label="Loading"
></div>

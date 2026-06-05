<script lang="ts">
	import { cn } from '$lib/utils';

	type Props = {
		class?: string;
	};

	let { class: className = '' }: Props = $props();
</script>

<div role="separator" class={cn('h-px shrink-0 bg-border', className)}></div>

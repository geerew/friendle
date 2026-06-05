<script lang="ts">
	import { cn } from '$lib/utils';

	type Props = {
		class?: string;
	};

	let { class: className = '' }: Props = $props();

	const dotClass = $derived(
		cn('animate-bounce-dot size-[2.5em] rounded-full bg-button-primary', className)
	);
</script>

<div class="inline-flex items-end gap-[0.35em]" role="status" aria-label="Loading">
	<div class={cn(dotClass, '[animation-delay:0ms]')}></div>
	<div class={cn(dotClass, '[animation-delay:110ms]')}></div>
	<div class={cn(dotClass, '[animation-delay:220ms]')}></div>
</div>

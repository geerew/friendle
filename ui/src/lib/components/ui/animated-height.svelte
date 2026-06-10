<script lang="ts">
	import { cn } from '$lib/utils';
	import type { Snippet } from 'svelte';
	import { tick } from 'svelte';

	type Props = {
		children: Snippet;
		class?: string;
		durationMs?: number;
		refreshKey?: unknown;
	};

	let { children, class: className, durationMs = 200, refreshKey }: Props = $props();

	let innerRef = $state<HTMLDivElement | null>(null);
	let height = $state<number | undefined>(undefined);
	let transitionEnabled = $state(false);

	// syncHeight sets the outer wrapper height from the inner content
	function syncHeight(animate: boolean): void {
		const el = innerRef;
		if (!el) {
			return;
		}

		const nextHeight = el.scrollHeight;

		if (!animate) {
			transitionEnabled = false;
			height = nextHeight;
			requestAnimationFrame(() => {
				transitionEnabled = true;
			});
			return;
		}

		height = nextHeight;
	}

	// animatedHeight tracks inner content size and transitions height changes
	$effect(() => {
		const el = innerRef;
		if (!el) {
			return;
		}

		syncHeight(false);

		const observer = new ResizeObserver(() => {
			syncHeight(true);
		});
		observer.observe(el);

		return () => {
			observer.disconnect();
		};
	});

	// refreshKey re-measures after conditional content swaps
	$effect(() => {
		refreshKey;

		void tick().then(() => {
			syncHeight(true);
		});
	});
</script>

<div
	class={cn(
		'overflow-hidden',
		transitionEnabled && 'transition-[height] ease-out',
		className
	)}
	style:height={height != null ? `${height}px` : undefined}
	style:transition-duration={transitionEnabled ? `${durationMs}ms` : undefined}
>
	<div bind:this={innerRef}>
		{@render children()}
	</div>
</div>

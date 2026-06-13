<script lang="ts">
	import { cn } from '$lib/utils';

	type Props = {
		roundEndsAt: string;
		class?: string;
	};

	let { roundEndsAt, class: className }: Props = $props();

	let remainingMs = $state(0);

	$effect(() => {
		roundEndsAt;

		const tick = () => {
			remainingMs = Math.max(0, Date.parse(roundEndsAt) - Date.now());
		};

		tick();
		const interval = setInterval(tick, 1000);

		return () => {
			clearInterval(interval);
		};
	});

	const label = $derived.by(() => {
		const totalSeconds = Math.floor(remainingMs / 1000);
		const hours = Math.floor(totalSeconds / 3600);
		const minutes = Math.floor((totalSeconds % 3600) / 60);
		const seconds = totalSeconds % 60;

		return `${hours}h ${minutes}m ${seconds}s`;
	});
</script>

<p class={cn('text-background-primary text-lg tabular-nums', className)}>
	{label}
</p>

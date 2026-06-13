<script lang="ts">
	import type { SVGAttributes } from 'svelte/elements';

	let restProps: SVGAttributes<SVGElement> = $props();
</script>

<svg
	xmlns="http://www.w3.org/2000/svg"
	fill="none"
	viewBox="0 0 24 24"
	stroke="currentColor"
	{...restProps}
>
	<circle cx="12" cy="8" r="5" stroke-linecap="round" stroke-linejoin="round" />
	<path stroke-linecap="round" stroke-linejoin="round" d="M20 21a8 8 0 0 0-16 0" />
</svg>

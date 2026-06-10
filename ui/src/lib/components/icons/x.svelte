<script lang="ts">
	import type { SVGAttributes } from 'svelte/elements';

	let restProps: SVGAttributes<SVGElement> = $props();
</script>

<svg
	xmlns="http://www.w3.org/2000/svg"
	fill="none"
	viewBox="0 0 24 24"
	stroke="currentColor"
	stroke-width="2"
	stroke-linecap="round"
	stroke-linejoin="round"
	{...restProps}
>
	<path d="M18 6 6 18" />
	<path d="m6 6 12 12" />
</svg>

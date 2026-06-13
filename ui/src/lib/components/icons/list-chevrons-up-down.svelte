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
	<path d="M3 5h8" />
	<path d="M3 12h8" />
	<path d="M3 19h8" />
	<path d="m15 8 3-3 3 3" />
	<path d="m15 16 3 3 3-3" />
</svg>

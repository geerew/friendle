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
	<line x1="15" x2="15" y1="12" y2="18" stroke-linecap="round" />
	<line x1="12" x2="18" y1="15" y2="15" stroke-linecap="round" />
	<rect width="14" height="14" x="8" y="8" rx="2" ry="2" stroke-linecap="round" stroke-linejoin="round" />
	<path
		stroke-linecap="round"
		stroke-linejoin="round"
		d="M4 16c-1.1 0-2-.9-2-2V4c0-1.1.9-2 2-2h10c1.1 0 2 .9 2 2"
	/>
</svg>

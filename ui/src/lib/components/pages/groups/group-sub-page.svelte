<script lang="ts">
	import { AppShell } from '$lib/components';
	import type { BreadcrumbItem } from '$lib/components/breadcrumb.svelte';
	import type { Snippet } from 'svelte';

	type Props = {
		title: string;
		breadcrumb: BreadcrumbItem[];
		children: Snippet;
	};

	let { title, breadcrumb, children }: Props = $props();
</script>

<AppShell {breadcrumb}>
	<div class="flex flex-col gap-5">
		<h2 class="section-title">{title}</h2>

		{@render children()}
	</div>
</AppShell>

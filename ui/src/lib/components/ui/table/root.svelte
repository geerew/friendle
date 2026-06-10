<script lang="ts">
	import AppShell from '$lib/components/app-shell.svelte';
	import type { BreadcrumbItem } from '$lib/components/breadcrumb.svelte';
	import type { Snippet } from 'svelte';

	type Props = {
		title: string;
		breadcrumb?: BreadcrumbItem[];
		showTitle?: boolean;
		showMenu?: boolean;
		children: Snippet;
	};

	let { title, breadcrumb, showTitle = true, showMenu = true, children }: Props = $props();

	const breadcrumbItems = $derived(
		breadcrumb !== undefined ? breadcrumb : [{ label: title }]
	);
</script>

<AppShell breadcrumb={breadcrumbItems} {showMenu}>
	<div class="flex flex-col gap-5">
		{#if showTitle}
			<h2 class="section-title">{title}</h2>
		{/if}

		{@render children()}
	</div>
</AppShell>

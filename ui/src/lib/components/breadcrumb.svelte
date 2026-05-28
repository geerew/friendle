<script module lang="ts">
	export type BreadcrumbItem = {
		label: string;
		href?: string;
	};
</script>

<script lang="ts">
	import { HomeIcon } from '$lib/components/icons';

	type Props = {
		items: BreadcrumbItem[];
	};

	let { items }: Props = $props();
</script>

<nav aria-label="Breadcrumb" class="mb-8 flex min-w-0 items-center gap-2 text-sm">
	<a
		href="/"
		class="inline-flex shrink-0 text-text-muted transition-colors hover:text-text"
		aria-label="Home"
	>
		<HomeIcon class="h-4 w-4 stroke-2" />
	</a>

	{#each items as item, index (item.label + index)}
		<span class="shrink-0 text-text-muted" aria-hidden="true">/</span>

		{#if item.href}
			<a
				href={item.href}
				class="truncate tracking-wide text-text-muted uppercase transition-colors hover:text-text"
			>
				{item.label}
			</a>
		{:else}
			<span class="truncate font-semibold tracking-wide text-text-muted uppercase">
				{item.label}
			</span>
		{/if}
	{/each}
</nav>

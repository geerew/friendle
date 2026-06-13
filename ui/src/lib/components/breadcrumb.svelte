<script module lang="ts">
	export type BreadcrumbItem = {
		label: string;
		href?: string;
		title?: string;
	};
</script>

<script lang="ts">
	import { HomeIcon } from '$lib/components/icons';

	type Props = {
		items: BreadcrumbItem[];
	};

	let { items }: Props = $props();
</script>

<nav aria-label="Breadcrumb" class="mb-3 flex min-w-0 items-center gap-2 text-sm">
	<a
		href="/"
		class="text-foreground-alt-2 hover:text-foreground inline-flex shrink-0 transition-colors"
		aria-label="Home"
	>
		<HomeIcon class="h-4 w-4 stroke-2" />
	</a>

	{#each items as item, index (item.label + index)}
		<span class="text-foreground-alt-2 shrink-0" aria-hidden="true">/</span>

		{#if item.href}
			<a
				href={item.href}
				title={item.title}
				class="text-foreground-alt-2 hover:text-foreground truncate tracking-wide uppercase transition-colors"
			>
				{item.label}
			</a>
		{:else}
			<span
				title={item.title}
				class="text-foreground-alt-2 truncate font-semibold tracking-wide uppercase"
			>
				{item.label}
			</span>
		{/if}
	{/each}
</nav>

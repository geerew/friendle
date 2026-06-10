<script module lang="ts">
	export type BreadcrumbItem = {
		label: string;
		accentLabel?: string;
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

<nav aria-label="Breadcrumb" class="mb-8 flex min-w-0 items-center gap-2 text-sm">
	<a
		href="/"
		class="inline-flex shrink-0 text-foreground-alt-2 transition-colors hover:text-foreground"
		aria-label="Home"
	>
		<HomeIcon class="h-4 w-4 stroke-2" />
	</a>

	{#each items as item, index (item.label + index)}
		<span class="shrink-0 text-foreground-alt-2" aria-hidden="true">/</span>

		{#if item.href}
			<a
				href={item.href}
				title={item.title}
				class="flex min-w-0 max-w-full items-center gap-1.5 truncate tracking-wide text-foreground-alt-2 transition-colors hover:text-foreground"
			>
				<span class="shrink-0 uppercase">{item.label}</span>
				{#if item.accentLabel}
					<span class="text-background-primary truncate normal-case">{item.accentLabel}</span>
				{/if}
			</a>
		{:else}
			<span
				title={item.title}
				class="truncate font-semibold tracking-wide text-foreground-alt-2 uppercase"
			>
				{item.label}
			</span>
		{/if}
	{/each}
</nav>

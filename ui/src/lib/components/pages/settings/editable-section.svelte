<script lang="ts">
	import { Button } from '$lib/components/ui';
	import type { Snippet } from 'svelte';

	type EditVariant = 'ghost' | 'destructive';

	type Props = {
		title: string;
		editing: boolean;
		onEdit: () => void;
		onCancel: () => void;
		children: Snippet;
		editLabel?: string;
		editVariant?: EditVariant;
		hint?: string;
		titleDestructive?: boolean;
		cancelDisabled?: boolean;
	};

	let {
		title,
		editing,
		onEdit,
		onCancel,
		children,
		editLabel = 'Edit',
		editVariant = 'ghost',
		hint,
		titleDestructive = false,
		cancelDisabled = false
	}: Props = $props();
</script>

<section class="flex flex-col gap-3">
	<div class="flex items-center justify-between gap-3">
		<h2 class="section-title" class:section-title-destructive={titleDestructive}>{title}</h2>
		{#if editing}
			<Button
				type="button"
				variant="ghost"
				size="inline"
				disabled={cancelDisabled}
				onclick={onCancel}
			>
				Cancel
			</Button>
		{:else}
			<Button type="button" variant={editVariant} size="inline" onclick={onEdit}>
				{editLabel}
			</Button>
		{/if}
	</div>

	{#if hint}
		<p class="text-foreground-alt-2 text-sm">{hint}</p>
	{/if}

	{@render children()}
</section>

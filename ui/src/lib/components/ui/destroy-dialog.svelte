<script lang="ts">
	import { Spinner } from '$lib/components';
	import Button from './button.svelte';
	import * as Dialog from './dialog';
	import type { Snippet } from 'svelte';

	type Props = {
		open?: boolean;
		title: string;
		description?: string;
		detail?: string;
		confirmLabel?: string;
		cancelLabel?: string;
		disabled?: boolean;
		loading?: boolean;
		onConfirm: () => void | Promise<void>;
		children?: Snippet;
	};

	let {
		open = $bindable(false),
		title,
		description,
		detail,
		confirmLabel = 'Delete',
		cancelLabel = 'Cancel',
		disabled = false,
		loading = false,
		onConfirm,
		children
	}: Props = $props();

	async function handleConfirm(): Promise<void> {
		await onConfirm();
	}
</script>

<Dialog.Root bind:open>
	<Dialog.Content>
		<Dialog.Alert>
			<div class="text-text flex flex-col gap-2 text-center">
				<span class="text-lg">{title}</span>
				{#if detail}
					<span class="text-button-primary font-semibold">{detail}</span>
				{/if}
				{#if description}
					<span class="text-text-muted text-sm">{description}</span>
				{/if}
			</div>
		</Dialog.Alert>

		{#if children}
			{@render children()}
		{/if}

		<Dialog.Footer class="h-auto py-3">
			<Dialog.CloseButton class="text-xs" disabled={loading}>{cancelLabel}</Dialog.CloseButton>
			<Button
				variant="destructive"
				class="h-10 w-36 gap-2 px-0 py-0 text-xs"
				disabled={disabled || loading}
				onclick={handleConfirm}
			>
				{#if loading}
					<Spinner class="bg-text size-2" />
				{:else}
					<span class="text-xs">{confirmLabel}</span>
				{/if}
			</Button>
		</Dialog.Footer>
	</Dialog.Content>
</Dialog.Root>

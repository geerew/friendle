<script lang="ts">
	import { Dialog } from 'bits-ui';
	import type { Snippet } from 'svelte';

	type Props = Dialog.RootProps & {
		trigger?: Snippet;
		children: Snippet;
		portalProps?: Dialog.PortalProps;
	};

	let { open = $bindable(false), trigger, children, portalProps, ...restProps }: Props = $props();
</script>

<Dialog.Root bind:open {...restProps}>
	{@render trigger?.()}

	<Dialog.Portal {...portalProps}>
		<Dialog.Overlay class="fixed inset-0 z-50 bg-black/50" />

		{@render children()}
	</Dialog.Portal>
</Dialog.Root>

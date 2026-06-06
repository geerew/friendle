<script lang="ts">
	import Logo from '$lib/components/logo.svelte';

	type Props = {
		subtitle: string;
	};

	let { subtitle }: Props = $props();
</script>

<div class="flex flex-col items-center gap-2 text-center">
	<Logo />
	<p class="text-sm text-foreground-alt-2">{subtitle}</p>
</div>

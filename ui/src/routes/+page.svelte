<script lang="ts">
	import { AppShell } from '$lib/components';
	import { auth } from '$lib/auth.svelte';
</script>

<AppShell>
	<section class="flex flex-col gap-3">
		<h2 class="section-title">Welcome</h2>
		{#if auth.user}
			<p class="text-sm text-text-muted">
				Signed in as <span class="text-button-primary">{auth.user.displayName}</span>.
			</p>
		{/if}
	</section>
</AppShell>

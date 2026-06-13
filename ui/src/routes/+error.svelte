<script lang="ts">
	import { page } from '$app/state';
</script>

<div class="app-shell page-content items-center justify-center text-center">
	<h1 class="text-2xl font-bold">{page.status}</h1>
	<p class="text-foreground-alt-2">{page.error?.message ?? 'Something went wrong.'}</p>
	<a href="/" class="text-background-primary mt-4">Go home</a>
</div>

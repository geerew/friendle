<script lang="ts">
	import { page } from '$app/state';
</script>

<div class="app-shell page-content items-center justify-center text-center">
	<h1 class="text-2xl font-bold">{page.status}</h1>
	<p class="text-text-muted">{page.error?.message ?? 'Something went wrong.'}</p>
	<a href="/" class="mt-4 text-tile-correct">Go home</a>
</div>

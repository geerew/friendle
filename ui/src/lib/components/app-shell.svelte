<script lang="ts">
	import type { Snippet } from 'svelte';

	type Props = {
		title: string;
		backHref?: string;
		settingsHref?: string;
		showBack?: boolean;
		showSettings?: boolean;
		children: Snippet;
	};

	let {
		title,
		backHref = '/',
		settingsHref = '/settings/',
		showBack = true,
		showSettings = false,
		children
	}: Props = $props();
</script>

<div class="app-shell">
	<header class="sticky top-0 z-10 border-b border-border bg-bg">
		<div class="flex items-center justify-between px-3 py-3">
			<div class="w-10">
				{#if showBack}
					<a
						href={backHref}
						class="inline-flex h-10 w-10 items-center justify-center rounded text-text-muted hover:text-text"
						aria-label="Go back"
					>
						<svg viewBox="0 0 24 24" class="h-6 w-6" fill="none" stroke="currentColor" stroke-width="2">
							<path d="M15 18l-6-6 6-6" />
						</svg>
					</a>
				{/if}
			</div>

			<h1 class="truncate text-center text-base font-semibold">{title}</h1>

			<div class="w-10">
				{#if showSettings}
					<a
						href={settingsHref}
						class="inline-flex h-10 w-10 items-center justify-center rounded text-text-muted hover:text-text"
						aria-label="Settings"
					>
						<svg viewBox="0 0 24 24" class="h-5 w-5" fill="none" stroke="currentColor" stroke-width="2">
							<circle cx="12" cy="12" r="3" />
							<path
								d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"
							/>
						</svg>
					</a>
				{/if}
			</div>
		</div>
	</header>

	<main class="page-content">
		{@render children()}
	</main>
</div>

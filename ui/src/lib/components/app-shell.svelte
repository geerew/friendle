<script lang="ts">
	import { goto } from '$app/navigation';
	import { auth } from '$lib/auth.svelte';
	import { LeftChevronIcon, SettingsIcon } from '$lib/components/icons';
	import { Dropdown } from '$lib/components/ui';
	import type { Snippet } from 'svelte';

	type Props = {
		title: string;
		backHref?: string;
		showBack?: boolean;
		showMenu?: boolean;
		children: Snippet;
	};

	let {
		title,
		backHref = '/',
		showBack = true,
		showMenu = true,
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
						<LeftChevronIcon class="h-6 w-6 stroke-2" />
					</a>
				{/if}
			</div>

			<h1 class="truncate text-center text-base font-semibold">{title}</h1>

			<div class="flex w-10 justify-end">
				{#if showMenu}
					<Dropdown.Root>
						<Dropdown.Trigger aria-label="Menu">
							<SettingsIcon class="h-5 w-5 stroke-2" />
						</Dropdown.Trigger>

						<Dropdown.Content>
							<Dropdown.Item onSelect={() => goto('/settings/')}>Settings</Dropdown.Item>
							{#if auth.isAdmin}
								<Dropdown.Item onSelect={() => goto('/admin/')}>Admin</Dropdown.Item>
							{/if}
						</Dropdown.Content>
					</Dropdown.Root>
				{/if}
			</div>
		</div>
	</header>

	<main class="page-content">
		{@render children()}
	</main>
</div>

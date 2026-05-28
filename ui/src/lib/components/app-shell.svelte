<script lang="ts">
	import { goto } from '$app/navigation';
	import { logout } from '$lib/api/auth-api';
	import { auth } from '$lib/auth.svelte';
	import Breadcrumb, { type BreadcrumbItem } from '$lib/components/breadcrumb.svelte';
	import { SettingsIcon } from '$lib/components/icons';
	import Logo from '$lib/components/logo.svelte';
	import { Dropdown } from '$lib/components/ui';
	import type { Snippet } from 'svelte';

	type Props = {
		title?: string;
		breadcrumb?: BreadcrumbItem[];
		showMenu?: boolean;
		children: Snippet;
	};

	let { title, breadcrumb, showMenu = true, children }: Props = $props();

	const breadcrumbItems = $derived(breadcrumb ?? (title ? [{ label: title }] : []));

	async function handleLogout(): Promise<void> {
		await logout();
		auth.clear();
		await goto('/auth/login/');
	}
</script>

<div class="app-shell">
	<header class="sticky top-0 z-10 border-b border-border bg-bg">
		<div class="flex items-center justify-between px-3 py-3">
			<div class="w-10 shrink-0"></div>

			<div class="flex min-w-0 flex-1 justify-center px-2">
				<Logo href="/" variant="header" />
			</div>

			<div class="flex w-10 shrink-0 justify-end">
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
							<Dropdown.Separator />
							<Dropdown.Item variant="destructive" onSelect={handleLogout}>Logout</Dropdown.Item>
						</Dropdown.Content>
					</Dropdown.Root>
				{/if}
			</div>
		</div>
	</header>

	<main class="page-content">
		{#if breadcrumbItems.length > 0}
			<Breadcrumb items={breadcrumbItems} />
		{/if}

		{@render children()}
	</main>
</div>

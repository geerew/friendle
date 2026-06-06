<script lang="ts">
	import { goto } from '$app/navigation';
	import { logout } from '$lib/api/auth-api';
	import { auth } from '$lib/auth.svelte';
	import Breadcrumb, { type BreadcrumbItem } from '$lib/components/breadcrumb.svelte';
	import {
		CopyPlusIcon,
		GithubIcon,
		HomeIcon,
		LogOutIcon,
		MenuToggleIcon,
		SearchIcon,
		ShieldUserIcon,
		UserRoundIcon
	} from '$lib/components/icons';
	import Logo from '$lib/components/logo.svelte';
	import { Button, Separator } from '$lib/components/ui';
	import { cn } from '$lib/utils';
	import { Collapsible } from 'bits-ui';
	import type { Component } from 'svelte';
	import type { Snippet } from 'svelte';
	import type { SVGAttributes } from 'svelte/elements';

	type Props = {
		title?: string;
		breadcrumb?: BreadcrumbItem[];
		showMenu?: boolean;
		children: Snippet;
	};

	type NavLinkItem = {
		label: string;
		href: string;
		icon: Component<SVGAttributes<SVGElement>>;
	};

	let { title, breadcrumb, showMenu = true, children }: Props = $props();

	const breadcrumbItems = $derived(breadcrumb ?? (title ? [{ label: title }] : []));

	let navOpen = $state(false);

	const githubRepoUrl = 'https://github.com/geerew/friendle';

	const navItems = $derived.by((): NavLinkItem[] => {
		const items: NavLinkItem[] = [
			{ label: 'Home', href: '/', icon: HomeIcon },
			{ label: 'Create Group', href: '/groups/create/', icon: CopyPlusIcon },
			{ label: 'Profile', href: '/settings/', icon: UserRoundIcon }
		];

		if (auth.isAdmin) {
			items.push({ label: 'Admin', href: '/admin/', icon: ShieldUserIcon });
		}

		return items;
	});

	$effect(() => {
		if (!navOpen) {
			return;
		}

		const previousOverflow = document.body.style.overflow;
		document.body.style.overflow = 'hidden';

		return () => {
			document.body.style.overflow = previousOverflow;
		};
	});

	async function handleLogout(): Promise<void> {
		closeNav();
		await logout();
		auth.clear();
		await goto('/auth/login/');
	}

	function closeNav(): void {
		navOpen = false;
	}

	function handleNavLink(): void {
		closeNav();
	}
</script>

<div class="app-shell">
	{#if showMenu}
		<button
			type="button"
			class={cn(
				'bg-background/70 fixed inset-0 z-9 transition-opacity duration-200 ease-out',
				navOpen ? 'pointer-events-auto opacity-100' : 'pointer-events-none opacity-0'
			)}
			aria-label="Close menu"
			aria-hidden={!navOpen}
			tabindex={navOpen ? 0 : -1}
			onclick={closeNav}
		></button>
	{/if}

	<div class="bg-background sticky top-0 z-10 shrink-0">
		<Collapsible.Root bind:open={navOpen} class="border-foreground-alt-4 border-b">
			<header class="app-header">
				<div class="flex h-full items-center justify-between px-3">
					<div class="flex h-10 w-10 shrink-0 items-center justify-start">
						{#if showMenu}
							<Collapsible.Trigger
								class="text-foreground-alt-2 hover:text-foreground inline-flex h-10 w-10 cursor-pointer items-center justify-center rounded transition-colors"
								aria-label={navOpen ? 'Close menu' : 'Open menu'}
								aria-expanded={navOpen}
							>
								<MenuToggleIcon open={navOpen} />
							</Collapsible.Trigger>
						{/if}
					</div>

					<div class="flex min-w-0 flex-1 justify-center px-2">
						<Logo href="/" variant="header" />
					</div>

					<div class="flex h-10 w-10 shrink-0 items-center justify-end">
						{#if showMenu}
							<Button
								href="/groups/search/"
								variant="ghost"
								size="icon"
								aria-label="Search groups"
								class="text-foreground-alt-2 hover:text-foreground h-10 w-10"
							>
								<SearchIcon class="h-5 w-5 stroke-2" />
							</Button>
						{/if}
					</div>
				</div>
			</header>

			{#if showMenu}
				<Collapsible.Content
					class="data-[state=closed]:animate-collapsible-up data-[state=open]:animate-collapsible-down bg-background overflow-hidden"
				>
					<nav class="flex flex-col items-center py-1">
						<div class="flex w-56 max-w-full flex-col gap-2">
							{#each navItems as item, index (item.href)}
								<Button
									href={item.href}
									variant="ghost"
									size="inline"
									class="flex h-auto min-h-0 w-full items-center justify-start gap-3 rounded-none px-0 py-2 text-sm font-medium tracking-wide uppercase"
									onclick={handleNavLink}
								>
									<item.icon class="size-4.5 shrink-0 stroke-2" />
									<span>{item.label}</span>
								</Button>

								{#if index < navItems.length - 1}
									<Separator class="bg-foreground-alt-4 w-full" />
								{/if}
							{/each}

							<Separator class="bg-foreground-alt-4 w-full" />

							<div class="flex items-center justify-center gap-1 py-1">
								<Button
									href={githubRepoUrl}
									target="_blank"
									rel="noopener noreferrer"
									variant="ghost"
									size="icon"
									aria-label="GitHub repository"
									class="text-foreground-alt-2 hover:text-foreground h-10 w-10"
									onclick={closeNav}
								>
									<GithubIcon class="size-5" />
								</Button>

								<Button
									variant="ghost"
									size="icon"
									aria-label="Logout"
									class="text-foreground-alt-2 hover:text-foreground h-10 w-10"
									onclick={handleLogout}
								>
									<LogOutIcon class="size-5 stroke-2" />
								</Button>
							</div>
						</div>
					</nav>
				</Collapsible.Content>
			{/if}
		</Collapsible.Root>
	</div>

	<main class="page-content relative z-0">
		{#if breadcrumbItems.length > 0}
			<Breadcrumb items={breadcrumbItems} />
		{/if}

		{@render children()}
	</main>
</div>

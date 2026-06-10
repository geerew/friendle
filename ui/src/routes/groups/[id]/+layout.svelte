<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { getGroup } from '$lib/api/groups-api';
	import { Spinner } from '$lib/components';
	import { GROUP_PAGE_KEY, type GroupPageContext } from '$lib/context/group-page';
	import { isGroupMember } from '$lib/utils/group';
	import { setContext, type Snippet } from 'svelte';

	type Props = {
		children: Snippet;
	};

	let { children }: Props = $props();

	const groupId = $derived(page.params.id ?? '');

	const groupPage = $state<GroupPageContext>({
		group: null,
		loading: true
	});

	setContext(GROUP_PAGE_KEY, groupPage);

	let previousPathname = $state('');

	// isGroupHomePath reports whether pathname is the group detail route
	function isGroupHomePath(pathname: string, id: string): boolean {
		return pathname.replace(/\/+$/, '') === `/groups/${id}`;
	}

	$effect(() => {
		groupId;
		void loadGroupPage();
	});

	$effect(() => {
		const id = groupId;
		const pathname = page.url.pathname;
		const previous = previousPathname;
		previousPathname = pathname;

		if (
			!id ||
			!isGroupHomePath(pathname, id) ||
			groupPage.group?.id !== id ||
			previous === pathname
		) {
			return;
		}

		void loadGroupPage({ silent: true });
	});

	// loadGroupPage fetches the group and redirects non-members away from this route tree
	async function loadGroupPage(options?: { silent?: boolean }): Promise<void> {
		if (!groupId) {
			groupPage.group = null;
			groupPage.loading = false;
			await goto('/');
			return;
		}

		if (!options?.silent) {
			groupPage.loading = true;
		}

		try {
			const group = await getGroup(groupId);

			if (!isGroupMember(group)) {
				await goto('/');
				return;
			}

			groupPage.group = group;
		} catch {
			groupPage.group = null;
			await goto('/');
			return;
		} finally {
			groupPage.loading = false;
		}
	}
</script>

{#if groupPage.loading}
	<div class="flex min-h-24 items-center justify-center">
		<Spinner class="bg-foreground-alt-2 size-3" />
	</div>
{:else if groupPage.group}
	{@render children()}
{/if}

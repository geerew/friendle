<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { ApiError } from '$lib/api';
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
		loading: true,
		error: null
	});

	setContext(GROUP_PAGE_KEY, groupPage);

	$effect(() => {
		groupId;
		void loadGroupPage();
	});

	// loadGroupPage fetches the group and redirects non-members away from this route tree
	async function loadGroupPage(): Promise<void> {
		if (!groupId) {
			groupPage.group = null;
			groupPage.error = 'Group not found';
			groupPage.loading = false;
			await goto('/');
			return;
		}

		groupPage.loading = true;
		groupPage.error = null;

		try {
			const group = await getGroup(groupId);

			if (!isGroupMember(group)) {
				await goto('/');
				return;
			}

			groupPage.group = group;
		} catch (err) {
			groupPage.group = null;
			groupPage.error = err instanceof ApiError ? err.message : 'Failed to load group';
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

<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { GroupComingSoonPage } from '$lib/components/pages';
	import { GROUP_PAGE_KEY, type GroupPageContext } from '$lib/context/group-page';
	import { isGroupAdmin } from '$lib/utils/group';
	import { getContext } from 'svelte';

	const groupId = $derived(page.params.id ?? '');
	const groupPage = getContext<GroupPageContext>(GROUP_PAGE_KEY);
	const group = $derived(groupPage.group);

	$effect(() => {
		groupId;
		group;

		if (group && !isGroupAdmin(group)) {
			void goto(`/groups/${groupId}/`);
		}
	});
</script>

{#if group && isGroupAdmin(group)}
	<GroupComingSoonPage title="Rejected Requests" breadcrumbLabel="Rejected" />
{/if}

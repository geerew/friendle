<script lang="ts">
	import { GroupSubPage } from '$lib/components/pages';
	import { GROUP_PAGE_KEY, type GroupPageContext } from '$lib/context/group-page';
	import { groupChildBreadcrumb } from '$lib/utils/group';
	import { getContext } from 'svelte';

	type Props = {
		title: string;
	};

	let { title }: Props = $props();

	const groupPage = getContext<GroupPageContext>(GROUP_PAGE_KEY);
	const group = $derived(groupPage.group);
	const breadcrumb = $derived(
		group ? groupChildBreadcrumb(group.id, group.name, title) : [{ label: title }]
	);
</script>

<GroupSubPage {title} {breadcrumb}>
	<p class="text-foreground-alt-2 text-sm italic">Coming soon</p>
</GroupSubPage>

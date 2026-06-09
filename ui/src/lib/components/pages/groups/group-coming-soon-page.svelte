<script lang="ts">
	import { page } from '$app/state';
	import { AppShell } from '$lib/components';
	import { GROUP_PAGE_KEY, type GroupPageContext } from '$lib/context/group-page';
	import { getContext } from 'svelte';

	type Props = {
		title: string;
	};

	let { title }: Props = $props();

	const groupId = $derived(page.params.id ?? '');
	const groupPage = getContext<GroupPageContext>(GROUP_PAGE_KEY);
	const group = $derived(groupPage.group);
</script>

<AppShell breadcrumb={[{ label: 'Group', href: `/groups/${groupId}/` }, { label: title }]}>
	<div class="flex flex-col gap-5">
		{#if group}
			<section class="flex flex-col gap-3">
				<h2 class="section-title">Group</h2>
				<p class="text-background-primary text-2xl">{group.name}</p>
			</section>
		{/if}

		<p class="text-foreground-alt-2 text-sm italic">Coming soon</p>
	</div>
</AppShell>

<script lang="ts">
	import { GROUP_PAGE_KEY, type GroupPageContext } from '$lib/context/group-page';
	import { Separator } from '$lib/components/ui';
	import { getContext } from 'svelte';

	const groupPage = getContext<GroupPageContext>(GROUP_PAGE_KEY);
	const group = $derived(groupPage.group);
</script>

{#if group}
	<section class="flex flex-col gap-3">
		<h2 class="section-title">Group name</h2>
		<p class="text-background-primary text-2xl">{group.name}</p>
	</section>

	<Separator />
{/if}

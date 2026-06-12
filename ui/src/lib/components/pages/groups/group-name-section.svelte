<script lang="ts">
	import { GROUP_PAGE_KEY, type GroupPageContext } from '$lib/context/group-page';
	import { Button, Separator } from '$lib/components/ui';
	import { getContext } from 'svelte';

	type Props = {
		showSettings?: boolean;
	};

	let { showSettings = false }: Props = $props();

	const groupPage = getContext<GroupPageContext>(GROUP_PAGE_KEY);
	const group = $derived(groupPage.group);
</script>

{#if group}
	<section class="flex w-full items-start gap-3">
		<div class="flex min-w-0 flex-1 flex-col gap-3">
			<h2 class="section-title">Group name</h2>
			<p class="text-background-primary text-2xl">{group.name}</p>
		</div>
		{#if showSettings}
			<Button
				href="/groups/{group.id}/settings/"
				variant="secondary"
				size="inline"
				class="shrink-0"
			>
				Settings
			</Button>
		{/if}
	</section>

	<Separator />
{/if}

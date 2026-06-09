<script lang="ts">
	import { RightChevronIcon } from '$lib/components/icons';
	import { Separator } from '$lib/components/ui';
	import type { GroupModel } from '$lib/models/group-model';

	type Props = {
		groups: GroupModel[];
	};

	let { groups }: Props = $props();
</script>

<div class="flex flex-col gap-2">
	{#each groups as group, index (group.id)}
		<a
			href="/groups/{group.id}/"
			class="hover:bg-background-alt-1 flex w-full items-center gap-3 rounded-md px-2 py-3 text-left transition-all"
		>
			<div class="flex w-full items-center px-1">
				<span class="text-foreground-alt-1 min-w-0 flex-1 truncate text-base font-medium">
					{group.name}
				</span>
				<RightChevronIcon class="text-foreground-alt-2 size-5 shrink-0 stroke-2" />
			</div>
		</a>
		{#if index < groups.length - 1}
			<div class="flex w-full items-center justify-center">
				<Separator class="bg-foreground-alt-5 w-[95%]" />
			</div>
		{/if}
	{/each}
</div>

<script lang="ts">
	import { page } from '$app/state';
	import { ApiError } from '$lib/api';
	import { getGroup } from '$lib/api/groups-api';
	import { auth } from '$lib/auth.svelte';
	import { AppShell, Spinner } from '$lib/components';
	import { Separator } from '$lib/components/ui';
	import type { GroupModel } from '$lib/models/group-model';

	const groupId = $derived(page.params.id ?? '');

	let group = $state<GroupModel | null>(null);
	let loading = $state(true);
	let error = $state<string | null>(null);

	$effect(() => {
		groupId;
		void loadGroup();
	});

	async function loadGroup(): Promise<void> {
		if (!groupId) {
			group = null;
			error = 'Group not found';
			loading = false;
			return;
		}

		loading = true;
		error = null;

		try {
			group = await getGroup(groupId);
		} catch (err) {
			group = null;
			error = err instanceof ApiError ? err.message : 'Failed to load group';
		} finally {
			loading = false;
		}
	}
</script>

<AppShell breadcrumb={[{ label: 'Group' }]}>
	{#if loading}
		<div class="flex min-h-24 items-center justify-center">
			<Spinner class="bg-foreground-alt-2 size-3" />
		</div>
	{:else if error}
		<p class="text-foreground-error-alt-1 text-sm">{error}</p>
	{:else if group}
		<div class="flex flex-col gap-5">
			<section class="flex flex-col gap-3">
				<h2 class="section-title">Group name</h2>
				<p class="text-background-primary text-2xl">{group.name}</p>
			</section>

			<Separator />

			<section class="flex flex-col gap-3">
				<h2 class="section-title">
					Members ({group.memberCount})
				</h2>
			</section>
		</div>
	{/if}
</AppShell>

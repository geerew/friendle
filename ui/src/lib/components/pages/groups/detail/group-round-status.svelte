<script lang="ts">
	import { ApiError } from '$lib/api/fetch';
	import { getGroupRoundToday } from '$lib/api/groups-api';
	import type { RoundTodayModel } from '$lib/models/round_model';
	import { apiErrorMessage, withMinLoadingDelay } from '$lib/utils';
	import { toast } from 'svelte-sonner';

	type Props = {
		groupId: string;
		memberThresholdMet: boolean;
	};

	let { groupId, memberThresholdMet }: Props = $props();

	let loading = $state(true);
	let roundToday = $state<RoundTodayModel | null>(null);
	let noRoundToday = $state(false);

	$effect(() => {
		groupId;
		memberThresholdMet;
		void loadRoundToday();
	});

	async function loadRoundToday(): Promise<void> {
		if (!groupId || !memberThresholdMet) {
			roundToday = null;
			noRoundToday = false;
			loading = false;
			return;
		}

		loading = true;
		noRoundToday = false;

		try {
			roundToday = await withMinLoadingDelay(getGroupRoundToday(groupId));
		} catch (err) {
			roundToday = null;

			if (err instanceof ApiError && err.status === 404) {
				noRoundToday = true;
				return;
			}

			toast.error(apiErrorMessage(err, "Failed to load today's round"));
		} finally {
			loading = false;
		}
	}

	const statusLine = $derived.by(() => {
		if (!memberThresholdMet) {
			return 'Not enough members';
		}

		if (noRoundToday) {
			return 'No round started for today yet.';
		}

		if (!roundToday) {
			return '';
		}

		switch (roundToday.status) {
			case 'awaiting_word':
				return roundToday.isPicker
					? "You are the picker — choose a word to start today's round."
					: 'Waiting for the picker';
			case 'active':
				return roundToday.isPicker
					? 'Your word is live — waiting for others to play.'
					: "Today's round is in progress.";
			case 'completed':
				return "Today's round is complete.";
		}
	});
</script>

<section class="flex flex-col gap-3">
	<h2 class="section-title">Today&apos;s round</h2>

	{#if loading}
		<p class="text-foreground-alt-2 text-sm">Loading…</p>
	{:else}
		<p class="text-foreground-alt-1 text-base">{statusLine}</p>
	{/if}
</section>

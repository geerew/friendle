<script lang="ts">
	import { updateGroup } from '$lib/api/groups-api';
	import { EditableSection } from '$lib/components';
	import { Button, Input } from '$lib/components/ui';
	import { GROUP_PAGE_KEY, type GroupPageContext } from '$lib/context/group-page';
	import { MAX_GROUP_NAME_LENGTH } from '$lib/models/group-model';
	import { apiErrorMessage, withMinLoadingDelay } from '$lib/utils';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	type Props = {
		showSettings?: boolean;
		editable?: boolean;
	};

	let { showSettings = false, editable = false }: Props = $props();

	const groupPage = getContext<GroupPageContext>(GROUP_PAGE_KEY);
	const group = $derived(groupPage.group);

	let editingName = $state(false);
	let nameDraft = $state('');
	let savingName = $state(false);

	const canSaveName = $derived(
		nameDraft.trim() !== '' &&
			nameDraft.trim() !== group?.name &&
			!savingName
	);

	// Enter group name edit mode with the current value
	function openNameEdit(): void {
		nameDraft = group?.name ?? '';
		editingName = true;
	}

	// Close group name edit mode without saving
	function closeNameEdit(): void {
		editingName = false;
	}

	// Persist a new group name via the API
	async function saveName(): Promise<void> {
		if (!group) {
			return;
		}

		const trimmed = nameDraft.trim();
		if (!trimmed || trimmed === group.name) {
			closeNameEdit();
			return;
		}

		if ([...trimmed].length > MAX_GROUP_NAME_LENGTH) {
			toast.error(`Group name must be no more than ${MAX_GROUP_NAME_LENGTH} characters`);
			return;
		}

		savingName = true;

		try {
			const updated = await withMinLoadingDelay(updateGroup(group.id, { name: trimmed }));
			groupPage.group = updated;
			closeNameEdit();
			toast.success('Group name updated');
		} catch (err) {
			toast.error(apiErrorMessage(err, 'Failed to update group name'));
			nameDraft = group.name;
		} finally {
			savingName = false;
		}
	}
</script>

{#if group}
	{#if editable}
		<EditableSection
			title="Group name"
			editing={editingName}
			cancelDisabled={savingName}
			onEdit={openNameEdit}
			onCancel={closeNameEdit}
		>
			{#if editingName}
				<form
					class="flex flex-col gap-4"
					onsubmit={(event) => {
						event.preventDefault();
						void saveName();
					}}
				>
					<Input bind:value={nameDraft} autofocus selectOnMount required />
					<Button
						type="submit"
						variant="primary"
						class="w-auto min-w-16 self-start px-5 text-xs"
						disabled={!canSaveName}
						loading={savingName}
					>
						Save
					</Button>
				</form>
			{:else}
				<p class="text-background-primary text-2xl">{group.name}</p>
			{/if}
		</EditableSection>
	{:else}
		<section class="flex flex-col gap-3">
			<div class="flex items-center justify-between gap-3">
				<h2 class="section-title">Group name</h2>
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
			</div>
			<p class="text-background-primary text-2xl">{group.name}</p>
		</section>
	{/if}
{/if}

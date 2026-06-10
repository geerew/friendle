<script lang="ts">
	import { page } from '$app/state';
	import { listGroupMembers, updateGroupMemberRole } from '$lib/api/groups-api';
	import { EditableSection } from '$lib/components';
	import { ListChevronsDownUpIcon, ListChevronsUpDownIcon } from '$lib/components/icons';
	import { GroupNameSection } from '$lib/components/pages';
	import GroupRoleBadge from '$lib/components/pages/groups/group-role-badge.svelte';
	import { Button, AnimatedHeight, Separator, Switch, Table } from '$lib/components/ui';
	import { auth } from '$lib/auth.svelte';
	import { GROUP_PAGE_KEY, type GroupPageContext } from '$lib/context/group-page';
	import type { GroupMemberModel } from '$lib/models/group-member-model';
	import { apiErrorMessage, cn, withMinLoadingDelay } from '$lib/utils';
	import { groupChildBreadcrumb, isGroupAdmin } from '$lib/utils/group';
	import { Collapsible } from 'bits-ui';
	import { getContext } from 'svelte';
	import { toast } from 'svelte-sonner';

	const groupId = $derived(page.params.id ?? '');
	const groupPage = getContext<GroupPageContext>(GROUP_PAGE_KEY);
	const group = $derived(groupPage.group);
	const breadcrumb = $derived(
		group ? groupChildBreadcrumb(group.id, 'Members') : [{ label: 'Members' }]
	);
	const currentUserId = $derived(auth.user?.id ?? '');
	const viewerIsAdmin = $derived(group ? isGroupAdmin(group) : false);

	let members = $state<GroupMemberModel[]>([]);
	let pageNum = $state(1);
	let perPage = $state(10);
	let totalItems = $state(0);
	let loading = $state(true);
	let expandedUserId = $state<string | null>(null);
	let editingRoleUserId = $state<string | null>(null);
	let roleDraft = $state<GroupMemberModel['groupRole']>('group_user');
	let savingRole = $state(false);

	const canSaveRole = $derived.by(() => {
		if (editingRoleUserId === null || savingRole) {
			return false;
		}

		const member = members.find((item) => item.userId === editingRoleUserId);

		return member != null && member.groupRole !== roleDraft;
	});

	$effect(() => {
		groupId;
		pageNum;
		perPage;
		void loadMembers();
	});

	async function loadMembers(): Promise<void> {
		if (!groupId) {
			members = [];
			totalItems = 0;
			toast.error('Group not found');
			loading = false;
			return;
		}

		loading = true;

		try {
			const data = await withMinLoadingDelay(listGroupMembers(groupId, { page: pageNum, perPage }));
			members = data.items;
			totalItems = data.totalItems;
		} catch (err) {
			members = [];
			totalItems = 0;
			toast.error(apiErrorMessage(err, 'Failed to load members'));
		} finally {
			loading = false;
		}
	}

	function canEditMember(member: GroupMemberModel): boolean {
		return viewerIsAdmin && member.userId !== currentUserId;
	}

	function setExpanded(userId: string, open: boolean): void {
		expandedUserId = open ? userId : null;

		if (!open && editingRoleUserId === userId) {
			closeRoleEdit();
		}
	}

	function isEditingRole(userId: string): boolean {
		return editingRoleUserId === userId;
	}

	function openRoleEdit(member: GroupMemberModel): void {
		roleDraft = member.groupRole;
		editingRoleUserId = member.userId;
	}

	function closeRoleEdit(): void {
		editingRoleUserId = null;
	}

	async function saveRole(): Promise<void> {
		if (!groupId || editingRoleUserId === null) {
			return;
		}

		const member = members.find((item) => item.userId === editingRoleUserId);

		if (!member || member.groupRole === roleDraft) {
			closeRoleEdit();
			return;
		}

		savingRole = true;

		try {
			const updated = await withMinLoadingDelay(
				updateGroupMemberRole(groupId, editingRoleUserId, { groupRole: roleDraft })
			);
			members = members.map((item) => (item.userId === updated.userId ? updated : item));
			closeRoleEdit();
			toast.success('Role updated');
		} catch (err) {
			toast.error(apiErrorMessage(err, 'Failed to update role'));
		} finally {
			savingRole = false;
		}
	}

	function roleLabel(role: GroupMemberModel['groupRole']): string {
		return role === 'group_admin' ? 'Admin' : 'Member';
	}

	function setRoleDraftFromSwitch(checked: boolean): void {
		roleDraft = checked ? 'group_admin' : 'group_user';
	}
</script>

<Table.Root title="Members" {breadcrumb}>
	{#snippet header()}
		<GroupNameSection />
	{/snippet}

	<Table.PaginatedBody
		itemCount={members.length}
		{totalItems}
		bind:page={pageNum}
		bind:perPage
		{loading}
		emptyMessage="No members"
	>
		{#snippet list()}
			<Table.List>
				{#each members as member, index (member.userId)}
					{@const isExpanded = expandedUserId === member.userId}
					<Collapsible.Root
						open={isExpanded}
						onOpenChange={(open) => setExpanded(member.userId, open)}
						class="flex flex-col"
					>
						<Table.Row label={member.displayName}>
							{#snippet trailing()}
								<GroupRoleBadge role={member.groupRole} />
								{#if canEditMember(member)}
									<Collapsible.Trigger
										class={cn(
											'text-foreground-alt-2 hover:bg-background-alt-1 hover:text-foreground focus-visible:ring-background-primary ms-2 flex size-5 shrink-0 items-center justify-center rounded-md transition-colors focus-visible:ring-2 focus-visible:outline-none',
											isExpanded && 'text-foreground'
										)}
										aria-label="View {member.displayName} details"
									>
										{#if isExpanded}
											<ListChevronsDownUpIcon class="size-4 shrink-0 stroke-2" />
										{:else}
											<ListChevronsUpDownIcon class="size-4 shrink-0 stroke-2" />
										{/if}
									</Collapsible.Trigger>
								{:else if viewerIsAdmin}
									<div class="ms-2 size-5 shrink-0" aria-hidden="true"></div>
								{/if}
							{/snippet}
						</Table.Row>

						{#if canEditMember(member)}
							<Collapsible.Content
								class="data-[state=closed]:animate-collapsible-up data-[state=open]:animate-collapsible-down overflow-hidden"
							>
								<div class="px-4 pt-2 pb-3">
									<section class="bg-foreground-alt-5/15 flex flex-col gap-3 rounded-md px-4 py-3">
										<div class="flex w-fit items-stretch gap-7">
											<div class="inline-flex flex-col items-center gap-2">
												<h3 class="section-title">Times picked</h3>
												<p class="text-background-primary text-lg tabular-nums">
													{member.timesPicked}
												</p>
											</div>
											<Separator class="h-auto w-px shrink-0 self-stretch" />
											<div class="inline-flex flex-col items-center gap-2">
												<h3 class="section-title">Skips</h3>
												<p class="text-background-primary text-lg tabular-nums">
													{member.pickerSkips}
												</p>
											</div>
										</div>

										<Separator />

										<EditableSection
											title="Role"
											editing={isEditingRole(member.userId)}
											cancelDisabled={savingRole}
											onEdit={() => openRoleEdit(member)}
											onCancel={closeRoleEdit}
										>
											<AnimatedHeight refreshKey={isEditingRole(member.userId)}>
												{#if isEditingRole(member.userId)}
													<form
														class="flex flex-col gap-3"
														onsubmit={(event) => {
															event.preventDefault();
															void saveRole();
														}}
													>
														<div class="flex w-fit items-center gap-3">
															<span
																class={cn(
																	'text-sm',
																	roleDraft === 'group_user'
																		? 'text-foreground'
																		: 'text-foreground-alt-2'
																)}
															>
																Member
															</span>
															<Switch
																checked={roleDraft === 'group_admin'}
																onCheckedChange={setRoleDraftFromSwitch}
															/>
															<span
																class={cn(
																	'text-sm',
																	roleDraft === 'group_admin'
																		? 'text-foreground'
																		: 'text-foreground-alt-2'
																)}
															>
																Admin
															</span>
														</div>
														<Button
															type="submit"
															variant="primary"
															size="inline"
															class="self-start px-3 text-xs"
															disabled={!canSaveRole}
															loading={savingRole}
														>
															Save
														</Button>
													</form>
												{:else}
													<p class="text-background-primary text-lg">
														{roleLabel(member.groupRole)}
													</p>
												{/if}
											</AnimatedHeight>
										</EditableSection>
									</section>
								</div>
							</Collapsible.Content>
						{/if}
					</Collapsible.Root>

					{#if index < members.length - 1}
						<Table.Separator />
					{/if}
				{/each}
			</Table.List>
		{/snippet}
	</Table.PaginatedBody>
</Table.Root>

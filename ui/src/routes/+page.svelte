<script lang="ts">
	import { AppShell } from '$lib/components';
	import { Button, Separator } from '$lib/components/ui';
</script>

<AppShell>
	<h2 class="section-title">My Groups</h2>

	<div class="flex flex-col gap-7">
		<div class="text-text-muted flex min-h-24 items-center justify-center text-sm italic">
			No groups
		</div>

		<Separator />

		<div class="flex flex-col gap-3">
			<Button type="button" variant="secondary">Search groups</Button>
			<Button type="button" variant="primary" href="/groups/create/">Create Group</Button>
		</div>
	</div>
</AppShell>

<script lang="ts">
	import { AppShell } from '$lib/components';
	import { Button, Separator } from '$lib/components/ui';
</script>

<AppShell>
	<div class="flex flex-col gap-5 pt-4">
		<section class="flex flex-col gap-3">
			<h2 class="section-title">My Groups</h2>
			<div class="text-text-muted flex min-h-30 items-center justify-center text-sm">No groups</div>
		</section>

		<Separator />

		<div class="flex flex-col gap-3">
			<Button type="button" variant="secondary">Search groups</Button>
			<Button type="button" variant="primary">Create Group</Button>
		</div>
	</div>
</AppShell>

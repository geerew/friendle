<script lang="ts">
	import { AppShell } from '$lib/components';
	import { Button } from '$lib/components/ui';
</script>

<AppShell title="Admin">
	<div class="flex flex-col gap-3">
		<Button href="/admin/users/" variant="secondary">Users</Button>
		<Button href="/admin/groups/" variant="secondary">Groups</Button>
	</div>
</AppShell>

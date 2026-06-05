<script lang="ts">
	import { AppShell } from '$lib/components';
	import { Button } from '$lib/components/ui';
</script>

<AppShell title="Admin">
	<div class="flex flex-col gap-2">
		<Button href="/admin/users/">Users</Button>
	</div>
</AppShell>

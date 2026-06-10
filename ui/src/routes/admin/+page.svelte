<script lang="ts">
	import { Button, Table } from '$lib/components/ui';
</script>

<Table.Root title="Admin">
	<div class="flex flex-col gap-3">
		<Button href="/admin/users/" variant="secondary">Users</Button>
		<Button href="/admin/groups/" variant="secondary">Groups</Button>
	</div>
</Table.Root>

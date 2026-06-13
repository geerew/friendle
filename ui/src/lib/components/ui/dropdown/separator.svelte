<script lang="ts">
	import { cn } from '$lib/utils';
	import { DropdownMenu, type DropdownMenuSeparatorProps } from 'bits-ui';

	type Props = DropdownMenuSeparatorProps & {
		class?: string;
	};

	let { class: containerClass, ...restProps }: Props = $props();
</script>

<DropdownMenu.Separator
	class={cn('bg-foreground-alt-4 -mx-1.5 my-1 h-px', containerClass)}
	{...restProps}
/>

<script lang="ts">
	import { cn } from '$lib/utils';
	import { Switch, type WithoutChildrenOrChild } from 'bits-ui';

	type Props = WithoutChildrenOrChild<Switch.RootProps> & {
		class?: string;
	};

	let { checked = $bindable(false), class: className, ...restProps }: Props = $props();
</script>

<Switch.Root
	bind:checked
	class={cn(
		'focus-visible:ring-background-primary inline-flex h-6 w-11 shrink-0 cursor-pointer items-center rounded-full border border-foreground-alt-4 bg-background-alt-1 px-0.5 transition-colors focus-visible:ring-2 focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-50 data-[state=checked]:border-background-primary data-[state=checked]:bg-background-primary',
		className
	)}
	{...restProps}
>
	<Switch.Thumb
		class="block size-4 shrink-0 rounded-full bg-foreground shadow-sm transition-transform data-[state=checked]:translate-x-5 data-[state=unchecked]:translate-x-0"
	/>
</Switch.Root>

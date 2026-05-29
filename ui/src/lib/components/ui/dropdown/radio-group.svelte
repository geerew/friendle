<script lang="ts">
	import { DropdownMenu } from 'bits-ui';

	let {
		ref = $bindable(null),
		value = $bindable(),
		...restProps
	}: DropdownMenu.RadioGroupProps = $props();
</script>

<DropdownMenu.RadioGroup bind:ref bind:value {...restProps} />

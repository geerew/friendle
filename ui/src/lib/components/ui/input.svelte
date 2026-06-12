<script lang="ts">
	import { cn } from '$lib/utils';
	import { EyeIcon, EyeOffIcon } from '$lib/components/icons';
	import Button from '$lib/components/ui/button.svelte';
	import type { HTMLInputAttributes } from 'svelte/elements';

	type Props = HTMLInputAttributes & {
		value?: HTMLInputAttributes['value'];
		ref?: HTMLInputElement;
		class?: string;
		password?: boolean;
		selectOnMount?: boolean;
	};

	let {
		value = $bindable(),
		ref = $bindable(),
		class: className = '',
		password = false,
		type = 'text',
		autofocus = false,
		selectOnMount = false,
		...restProps
	}: Props = $props();

	let visible = $state(false);

	const inputType = $derived(password ? (visible ? 'text' : 'password') : type);

	// Focus when mounted with autofocus (deferred so it wins over the triggering button)
	$effect(() => {
		if (!autofocus || !ref) return;

		const input = ref;
		requestAnimationFrame(() => {
			input.focus();
			if (selectOnMount) {
				input.select();
			}
		});
	});

	function toggleVisibility(): void {
		visible = !visible;
	}
</script>

{#if password}
	<div class="relative w-full">
		<input
			bind:value
			bind:this={ref}
			type={inputType}
			class={cn('field pe-11', className)}
			{...restProps}
		/>

		<Button
			type="button"
			variant="ghost"
			size="icon"
			class="absolute top-1/2 right-2 h-8 w-8 min-w-8 -translate-y-1/2 normal-case"
			aria-label={visible ? 'Hide password' : 'Show password'}
			onclick={toggleVisibility}
		>
			{#if visible}
				<EyeIcon class="size-5 stroke-[1.5]" />
			{:else}
				<EyeOffIcon class="size-5 stroke-[1.5]" />
			{/if}
		</Button>
	</div>
{:else}
	<input
		bind:value
		bind:this={ref}
		type={inputType}
		class={cn('field', className)}
		{...restProps}
	/>
{/if}

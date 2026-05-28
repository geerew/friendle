<script lang="ts">
	import { cn } from '$lib/utils';
	import { EyeIcon, EyeOffIcon } from '$lib/components/icons';
	import type { HTMLInputAttributes } from 'svelte/elements';

	type Props = HTMLInputAttributes & {
		value?: HTMLInputAttributes['value'];
		ref?: HTMLInputElement;
		class?: string;
		password?: boolean;
	};

	let {
		value = $bindable(),
		ref = $bindable(),
		class: className = '',
		password = false,
		type = 'text',
		...restProps
	}: Props = $props();

	let visible = $state(false);

	const inputType = $derived(password ? (visible ? 'text' : 'password') : type);

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

		<button
			type="button"
			class="absolute top-1/2 right-2 inline-flex h-8 w-8 -translate-y-1/2 items-center justify-center rounded text-text-muted transition-colors hover:text-text"
			aria-label={visible ? 'Hide password' : 'Show password'}
			onclick={toggleVisibility}
		>
			{#if visible}
				<EyeIcon class="size-5 stroke-[1.5]" />
			{:else}
				<EyeOffIcon class="size-5 stroke-[1.5]" />
			{/if}
		</button>
	</div>
{:else}
	<input bind:value bind:this={ref} type={inputType} class={cn('field', className)} {...restProps} />
{/if}

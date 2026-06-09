<script lang="ts">
	type Props = {
		label: string;
		count: number;
		href: string;
		ariaLabel: string;
	};

	let { label, count, href, ariaLabel }: Props = $props();
</script>

<a {href} class="group flex flex-col items-center gap-3 text-center" aria-label={ariaLabel}>
	<span class="section-title">{label}</span>
	<div class="relative h-8 w-full overflow-hidden">
		<span
			class="text-background-primary absolute inset-0 flex items-center justify-center text-2xl tabular-nums transition-all duration-300 ease-out group-hover:-translate-y-full group-hover:opacity-0"
		>
			{count}
		</span>
		<span
			class="text-background-primary absolute inset-0 flex translate-y-full items-center justify-center text-sm font-semibold tracking-wide uppercase opacity-0 transition-all duration-300 ease-out group-hover:translate-y-0 group-hover:opacity-100"
		>
			View
		</span>
	</div>
</a>

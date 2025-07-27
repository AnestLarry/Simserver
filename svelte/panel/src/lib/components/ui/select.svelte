<script lang="ts">
	import { cn } from "$lib/utils";
	import type { HTMLSelectAttributes } from "svelte/elements";

	interface SelectOption {
		value: string;
		name: string;
	}

	interface SelectProps extends HTMLSelectAttributes {
		class?: string;
		items: SelectOption[];
		value?: string;
		onchange?: () => void;
	}

	let {
		class: className,
		items,
		value = $bindable(),
		onchange,
		...restProps
	}: SelectProps = $props();

	function handleChange(event: Event) {
		const target = event.target as HTMLSelectElement;
		value = target.value;
		onchange?.();
	}
</script>

<select
	class={cn(
		"flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50",
		className
	)}
	bind:value
	onchange={handleChange}
	{...restProps}
>
	{#each items as item}
		<option value={item.value}>{item.name}</option>
	{/each}
</select>
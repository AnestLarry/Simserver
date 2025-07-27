<script lang="ts">
	import { cn } from "$lib/utils";
	import type { HTMLButtonAttributes } from "svelte/elements";

	interface SwitchProps extends HTMLButtonAttributes {
		class?: string;
		checked?: boolean;
		onCheckedChange?: (checked: boolean) => void;
	}

	let {
		class: className,
		checked = $bindable(false),
		onCheckedChange,
		disabled = false,
		...restProps
	}: SwitchProps = $props();

	function toggle() {
		if (disabled) return;
		checked = !checked;
		onCheckedChange?.(checked);
	}
</script>

<button
	type="button"
	role="switch"
	aria-checked={checked}
	class={cn(
		"peer inline-flex h-6 w-11 shrink-0 cursor-pointer items-center rounded-full border-2 border-transparent transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background disabled:cursor-not-allowed disabled:opacity-50",
		checked ? "bg-primary" : "bg-input",
		className
	)}
	{disabled}
	onclick={toggle}
	{...restProps}
>
	<span
		class={cn(
			"pointer-events-none block h-5 w-5 rounded-full bg-background shadow-lg ring-0 transition-transform",
			checked ? "translate-x-5" : "translate-x-0"
		)}
	></span>
</button>
<script lang="ts">
	import { cn } from "$lib/utils";
	import type { HTMLAttributes } from "svelte/elements";

	interface ListItem {
		name: string;
		current: boolean;
		imgIndex: number;
	}

	interface ListProps extends HTMLAttributes<HTMLUListElement> {
		class?: string;
		items: ListItem[];
		active?: boolean;
		onclick?: (event: CustomEvent) => void;
	}

	let { 
		class: className, 
		items,
		active = false,
		onclick,
		children,
		...restProps 
	}: ListProps = $props();

	function handleClick(item: ListItem) {
		onclick?.(new CustomEvent('click', { detail: item }));
	}
</script>

<ul
	class={cn(
		"space-y-1",
		className
	)}
	{...restProps}
>
	{#each items as item}
		<li>
			<button
				class={cn(
					"w-full text-left px-3 py-2 rounded-md text-sm transition-colors",
					item.current 
						? "bg-primary text-primary-foreground" 
						: "hover:bg-accent hover:text-accent-foreground"
				)}
				onclick={() => handleClick(item)}
				aria-current={item.current}
			>
				{#if children}
					{@render children({ item })}
				{:else}
					{item.name}
				{/if}
			</button>
		</li>
	{/each}
</ul>
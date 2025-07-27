<script lang="ts">
	import { cn } from "$lib/utils";
	import type { HTMLAttributes } from "svelte/elements";

	interface DrawerProps extends HTMLAttributes<HTMLDivElement> {
		class?: string;
		hidden?: boolean;
		transitionType?: string;
		transitionParams?: any;
	}

	let { 
		class: className, 
		hidden = $bindable(true),
		transitionType,
		transitionParams,
		children, 
		...restProps 
	}: DrawerProps = $props();
</script>

{#if !hidden}
	<!-- Backdrop -->
	<div 
		class="fixed inset-0 bg-black/50 z-40" 
		onclick={() => hidden = true}
		role="button"
		tabindex="0"
		onkeydown={(e) => {
			if (e.key === 'Escape') hidden = true;
		}}
	></div>
	
	<!-- Drawer -->
	<div
		class={cn(
			"fixed left-0 top-0 z-50 h-screen w-80 bg-white dark:bg-gray-800 shadow-lg transform transition-transform duration-300",
			className
		)}
		{...restProps}
	>
		<div class="p-4 border-b border-gray-200 dark:border-gray-700">
			<button
				onclick={() => hidden = true}
				class="float-right text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
			>
				<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
				</svg>
			</button>
		</div>
		<div class="p-4 overflow-y-auto h-full">
			{@render children?.()}
		</div>
	</div>
{/if}
<script lang="ts">
	import { cn } from "$lib/utils";
	import type { HTMLAttributes } from "svelte/elements";

	interface CardProps extends HTMLAttributes<HTMLDivElement> {
		class?: string;
	}

	let { class: className, children, ...restProps }: CardProps = $props();
</script>

<div
	class={cn("rounded-lg border bg-card text-card-foreground shadow-sm", className)}
	{...restProps}
>
	{@render children?.()}
</div>

<script lang="ts" module>
	export { default as CardHeader } from "./card-header.svelte";
	export { default as CardTitle } from "./card-title.svelte";
	export { default as CardContent } from "./card-content.svelte";
</script>

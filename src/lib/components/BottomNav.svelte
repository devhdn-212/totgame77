<script lang="ts">
	import { cn } from "$lib/utils";
	import Receipt from "@lucide/svelte/icons/receipt";
	import BookOpen from "@lucide/svelte/icons/book-open";
	import Camera from "@lucide/svelte/icons/camera";

	let {
		active = $bindable("transaksi"),
		onTransaksiClick,
		onCameraClick,
		onBukuMimpiClick,
	}: {
		active?: string;
		onTransaksiClick?: () => void;
		onCameraClick?: () => void;
		onBukuMimpiClick?: () => void;
	} = $props();

	const items = [
		{ key: "transaksi", label: "Transaksi", href: "#/transaksi", icon: Receipt },
		{ key: "camera", label: "Camera", href: "#/camera", icon: Camera },
		{ key: "bukumimpi", label: "Buku Mimpi", href: "#/bukumimpi", icon: BookOpen },
	];
</script>

<nav
	class="fixed inset-x-4 z-50 mx-auto max-w-md rounded-3xl border border-white/40 bg-white/70 shadow-lg shadow-black/10 backdrop-blur-xl dark:border-white/10 dark:bg-white/10"
	style="bottom: calc(1rem + env(safe-area-inset-bottom));"
>
	<div class="flex">
		{#each items as item (item.key)}
			{@const Icon = item.icon}
			<a
				href={item.href}
				class={cn(
					"flex flex-1 flex-col items-center gap-1 py-3 text-xs transition-colors",
					active === item.key ? "text-primary" : "text-muted-foreground",
				)}
				onclick={(e) => {
					active = item.key;
					if (item.key === "transaksi") {
						e.preventDefault();
						onTransaksiClick?.();
					} else if (item.key === "camera") {
						e.preventDefault();
						onCameraClick?.();
					} else if (item.key === "bukumimpi") {
						e.preventDefault();
						onBukuMimpiClick?.();
					}
				}}
			>
				<Icon class="size-5" />
				{item.label}
			</a>
		{/each}
	</div>
</nav>

<script lang="ts">
	import { cn } from "$lib/utils";
	import Receipt from "@lucide/svelte/icons/receipt";
	import ShoppingBag from "@lucide/svelte/icons/shopping-bag";
	import Camera from "@lucide/svelte/icons/camera";

	let {
		active = $bindable("transaksi"),
		onTransaksiClick,
		onCameraClick,
	}: { active?: string; onTransaksiClick?: () => void; onCameraClick?: () => void } = $props();

	const items = [
		{ key: "transaksi", label: "Transaksi", href: "#/transaksi", icon: Receipt },
		{ key: "camera", label: "Camera", href: "#/camera", icon: Camera },
		{ key: "belanja", label: "Belanja", href: "#/belanja", icon: ShoppingBag },
	];
</script>

<nav
	class="bg-background fixed inset-x-0 bottom-0 z-50 border-t pb-[env(safe-area-inset-bottom)] md:hidden"
>
	<div class="mx-auto flex max-w-6xl">
		{#each items as item (item.key)}
			{@const Icon = item.icon}
			<a
				href={item.href}
				class={cn(
					"flex flex-1 flex-col items-center gap-1 py-2 text-xs",
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
					}
				}}
			>
				<Icon class="size-5" />
				{item.label}
			</a>
		{/each}
	</div>
</nav>

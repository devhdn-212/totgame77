<script lang="ts">
	import { Button } from "$lib/components/ui/button";
	import { cn } from "$lib/utils";
	import Menu from "@lucide/svelte/icons/menu";
	import X from "@lucide/svelte/icons/x";
	import ChevronDown from "@lucide/svelte/icons/chevron-down";
	import Key from "@lucide/svelte/icons/key";
	import LogOut from "@lucide/svelte/icons/log-out";

	let mobileOpen = $state(false);
	let transaksiOpen = $state(false);

	const navLinks = [
		{ label: "Rekening Saya", href: "#/bank" },
		{ label: "Transaksi Saya", href: "#/transaksi-saya" },
	];

	const transaksiLinks = [
		{ label: "Internal Transaksi - Draft", href: "#/internaltransaksi" },
	];

	function closeMobileMenu() {
		mobileOpen = false;
		transaksiOpen = false;
	}
</script>

<header class="bg-primary text-primary-foreground sticky top-0 z-50">
	<div class="mx-auto flex h-14 max-w-6xl items-center justify-between gap-2 px-4">
		<a href="#/" class="text-base font-semibold tracking-wide">CLIENT</a>

		<!-- Desktop nav -->
		<nav class="hidden items-center gap-1 md:flex">
			<div class="relative">
				<button
					class="hover:bg-primary-foreground/10 flex items-center gap-1 rounded-lg px-3 py-2 text-sm"
					onclick={() => (transaksiOpen = !transaksiOpen)}
					aria-expanded={transaksiOpen}
				>
					TRANSAKSI
					<ChevronDown class="size-4" />
				</button>
				{#if transaksiOpen}
					<div class="bg-popover text-popover-foreground absolute left-0 mt-1 min-w-56 rounded-lg border p-1 shadow-md">
						{#each transaksiLinks as link (link.href)}
							<a
								href={link.href}
								class="hover:bg-muted block rounded-md px-3 py-2 text-sm"
								onclick={closeMobileMenu}
							>
								{link.label}
							</a>
						{/each}
					</div>
				{/if}
			</div>
			{#each navLinks as link (link.href)}
				<a href={link.href} class="hover:bg-primary-foreground/10 rounded-lg px-3 py-2 text-sm">
					{link.label}
				</a>
			{/each}
		</nav>

		<div class="hidden items-center gap-2 md:flex">
			<Button variant="secondary" size="icon" aria-label="Ganti password">
				<Key />
			</Button>
			<Button variant="secondary" size="default">
				<LogOut />
				LOGOUT
			</Button>
		</div>

		<!-- Mobile toggle -->
		<button
			class="hover:bg-primary-foreground/10 rounded-lg p-2 md:hidden"
			aria-label="Toggle navigation"
			aria-expanded={mobileOpen}
			onclick={() => (mobileOpen = !mobileOpen)}
		>
			{#if mobileOpen}
				<X class="size-6" />
			{:else}
				<Menu class="size-6" />
			{/if}
		</button>
	</div>

	<!-- Mobile menu -->
	{#if mobileOpen}
		<nav class="border-primary-foreground/10 border-t px-4 pt-2 pb-4 md:hidden">
			<div class="flex flex-col">
				<button
					class="hover:bg-primary-foreground/10 flex items-center justify-between rounded-lg px-3 py-2.5 text-left text-sm"
					onclick={() => (transaksiOpen = !transaksiOpen)}
					aria-expanded={transaksiOpen}
				>
					TRANSAKSI
					<ChevronDown class={cn("size-4 transition-transform", transaksiOpen && "rotate-180")} />
				</button>
				{#if transaksiOpen}
					<div class="flex flex-col pl-3">
						{#each transaksiLinks as link (link.href)}
							<a
								href={link.href}
								class="hover:bg-primary-foreground/10 rounded-lg px-3 py-2.5 text-sm"
								onclick={closeMobileMenu}
							>
								{link.label}
							</a>
						{/each}
					</div>
				{/if}
				{#each navLinks as link (link.href)}
					<a
						href={link.href}
						class="hover:bg-primary-foreground/10 rounded-lg px-3 py-2.5 text-sm"
						onclick={closeMobileMenu}
					>
						{link.label}
					</a>
				{/each}
			</div>

			<div class="border-primary-foreground/10 mt-3 flex items-center gap-2 border-t pt-3">
				<Button variant="secondary" size="default" class="flex-1">
					<Key />
					GANTI PASSWORD
				</Button>
				<Button variant="secondary" size="default" class="flex-1">
					<LogOut />
					LOGOUT
				</Button>
			</div>
		</nav>
	{/if}
</header>

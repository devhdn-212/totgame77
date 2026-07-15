<script lang="ts">
  import Navbar from "$lib/components/Navbar.svelte";
  import BottomNav from "$lib/components/BottomNav.svelte";
  import TabFourD from "$lib/components/TabFourD.svelte";
  import TabColok from "$lib/components/TabColok.svelte";
  import Tab5050 from "$lib/components/Tab5050.svelte";
  import TabMacauKombinasi from "$lib/components/TabMacauKombinasi.svelte";
  import TabDasar from "$lib/components/TabDasar.svelte";
  import TabShio from "$lib/components/TabShio.svelte";
  import CameraOcr from "$lib/components/CameraOcr.svelte";
  import { Button } from "$lib/components/ui/button";
  import { Card, CardContent } from "$lib/components/ui/card";
  import { Tabs, TabsList, TabsTrigger, TabsContent } from "$lib/components/ui/tabs";
  import { Input } from "$lib/components/ui/input";
  import {
    AlertDialog,
    AlertDialogContent,
    AlertDialogHeader,
    AlertDialogTitle,
    AlertDialogDescription,
    AlertDialogFooter,
    AlertDialogAction,
    AlertDialogCancel,
  } from "$lib/components/ui/alert-dialog";
  import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
  } from "$lib/components/ui/dialog";
  import { formatIDR, cn } from "$lib/utils";
  import Trash2 from "@lucide/svelte/icons/trash-2";
  import ShoppingCart from "@lucide/svelte/icons/shopping-cart";
  import Search from "@lucide/svelte/icons/search";
  import X from "@lucide/svelte/icons/x";
  import Decimal from "decimal.js";
  import { scale, slide } from "svelte/transition";
  import { flip } from "svelte/animate";
  import { Tween } from "svelte/motion";
  import { cubicOut } from "svelte/easing";

  const username = "ls0999212";
  const periode = "6512";
  let credit = $state(new Decimal(1000000));
  const creditDisplay = new Tween(credit.toNumber(), { duration: 800, easing: cubicOut });

  const betTypes = ["4D/3D/2D", "Colok", "5050", "Macau/Kombinasi", "Dasar", "Shio"];

  let activeBetType = $state("4D/3D/2D");

  function handleBetTypeClick(type: string) {
    activeBetType = type;
  }

  type BetEntry = {
    id: string;
    type: string;
    number: string;
    bet: string;
    deletable?: boolean;
  };

  let bets = $state<BetEntry[]>([]);
  let listTransaksi = $state<BetEntry[]>([]);

  const MIN_BET = 500;
  let minBetAlertOpen = $state(false);

  function handleDeleteBet(id: string) {
    bets = bets.filter((entry) => entry.id !== id);
  }

  const BET_TYPE_LABELS = ["4D", "3D", "3DD", "2D", "2DD", "2DT"];

  let searchQuery = $state("");
  let activeTypeFilter = $state<string | null>(null);

  let searchFilteredBets = $derived(
    bets.filter((entry) => entry.number.replace(/\*/g, "").includes(searchQuery.trim())),
  );

  let filteredBets = $derived(
    searchFilteredBets.filter((entry) => !activeTypeFilter || entry.type === activeTypeFilter),
  );

  let betTypeCounts = $derived(
    BET_TYPE_LABELS.map((type) => ({
      type,
      count: searchFilteredBets.filter((entry) => entry.type === type).length,
    })),
  );

  function handleTypeFilterClick(type: string) {
    activeTypeFilter = type;
  }

  let totalBelanja = $derived(
    bets.reduce((sum, entry) => sum.plus(entry.bet), new Decimal(0)),
  );

  let checkoutOpen = $state(false);
  let insufficientCreditAlertOpen = $state(false);

  function handleCheckoutClick() {
    if (totalBelanja.greaterThan(credit)) {
      insufficientCreditAlertOpen = true;
      return;
    }
    checkoutOpen = true;
  }

  function handleConfirmCheckout() {
    credit = credit.minus(totalBelanja);
    creditDisplay.set(credit.toNumber());
    listTransaksi = [...bets, ...listTransaksi];
    bets = [];
    checkoutOpen = false;
  }

  function handleCancelCheckout() {
    bets = [];
    checkoutOpen = false;
  }

  let transaksiModalOpen = $state(false);
  let cameraOcrOpen = $state(false);

  let transaksiGroups = $derived(
    BET_TYPE_LABELS.map((type) => {
      const entries = listTransaksi.filter((entry) => entry.type === type);
      const sum = entries.reduce((s, entry) => s.plus(entry.bet), new Decimal(0));
      return { type, entries, sum };
    }).filter((group) => group.entries.length > 0),
  );
</script>

<main class="mx-auto max-w-6xl px-4 py-6 pb-20 md:pb-6">
  <Card>
    <CardContent class="flex items-start justify-between gap-4">
      <div>
        <p class="text-lg font-semibold">HONGKONG</p>
        <p class="text-muted-foreground text-sm">Periode {periode}</p>
      </div>
      <div class="text-right">
        <p class="text-sm font-medium">{username}</p>
        <p class="text-muted-foreground text-xs">Credit</p>
        <p class="text-lg font-semibold text-blue-600">
          {formatIDR(creditDisplay.current)}
        </p>
      </div>
    </CardContent>
  </Card>

  <div class="-mx-4 mt-4 flex scrollbar-none gap-2 overflow-x-auto px-4 pb-1">
    {#each betTypes as type (type)}
      <Button
        variant="outline"
        class={cn(
          "shrink-0 cursor-pointer rounded-full",
          activeBetType === type && "border-orange-300 bg-orange-100 text-orange-700 hover:bg-orange-200",
        )}
        onclick={() => handleBetTypeClick(type)}
      >
        {type}
      </Button>
    {/each}
  </div>

  {#if activeBetType === "4D/3D/2D"}
    <TabFourD bind:bets bind:minBetAlertOpen onCheckoutClick={handleCheckoutClick} />
  {:else if activeBetType === "Colok"}
    <TabColok />
  {:else if activeBetType === "5050"}
    <Tab5050 />
  {:else if activeBetType === "Macau/Kombinasi"}
    <TabMacauKombinasi />
  {:else if activeBetType === "Dasar"}
    <TabDasar />
  {:else if activeBetType === "Shio"}
    <TabShio />
  {/if}

  {#if bets.length > 0}
    <div
      class="bg-background sticky top-0 z-10 -mx-4 mt-4 flex items-center justify-between border-b px-4 py-3"
    >
      <div>
        <p class="text-muted-foreground text-xs">Total Belanja</p>
        <p class="text-lg font-semibold">{formatIDR(totalBelanja)}</p>
        <p class="text-muted-foreground text-xs">Total Row : {bets.length}</p>
      </div>
      <div class="flex shrink-0 items-center gap-2">
        <Button
          variant="outline"
          class="cursor-pointer text-destructive"
          onclick={() => (bets = [])}
        >
          <Trash2 />
          Clear
        </Button>
        <Button
          class="cursor-pointer bg-blue-600 text-white hover:bg-blue-700"
          onclick={handleCheckoutClick}
        >
          <ShoppingCart />
          Checkout
        </Button>
      </div>
    </div>

    <div class="relative mt-3">
      <Search class="text-muted-foreground pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2" />
      <Input
        type="text"
        inputmode="numeric"
        placeholder="Cari nomor"
        class="pl-8"
        bind:value={searchQuery}
      />
    </div>

    <div class="mt-2 -mx-4 flex scrollbar-none items-center gap-2 overflow-x-auto px-4 pb-1">
      {#each betTypeCounts as { type, count } (type)}
        <button
          type="button"
          class={cn(
            "shrink-0 cursor-pointer rounded-full border px-2 py-0.5 text-xs backdrop-blur-sm",
            activeTypeFilter === type
              ? "border-green-800 bg-green-700 text-white"
              : "border-green-200 bg-green-100/70 text-green-700",
          )}
          onclick={() => handleTypeFilterClick(type)}
        >
          {type} :
          {#key count}
            <span class="inline-block" in:scale={{ duration: 250, start: 0.4 }}>{count}</span>
          {/key}
        </button>
      {/each}
    </div>

    {#if activeTypeFilter}
      <button
        type="button"
        aria-label="Bersihkan filter"
        class="mt-2 w-fit cursor-pointer rounded-full border border-red-200 bg-red-50 p-1 text-red-600 hover:bg-red-100"
        onclick={() => (activeTypeFilter = null)}
      >
        <X class="size-3.5" />
      </button>
    {/if}

    <div class="mt-3 flex flex-col gap-2">
      {#each filteredBets as entry (entry.id)}
        <div
          class="flex items-center justify-between rounded-lg border p-3"
          out:slide={{ duration: 200 }}
          animate:flip={{ duration: 200 }}
        >
          <div>
            <p class="text-sm font-medium">{entry.type} - {entry.number.replace(/\*/g, "")}</p>
            <p class="text-xs text-blue-600">Bet : {formatIDR(entry.bet)}</p>
          </div>
          {#if entry.deletable !== false}
            <Button
              size="icon"
              variant="ghost"
              class="text-destructive cursor-pointer"
              aria-label="Hapus"
              onclick={() => handleDeleteBet(entry.id)}
            >
              <Trash2 />
            </Button>
          {/if}
        </div>
      {/each}
    </div>
  {/if}
</main>

<AlertDialog bind:open={checkoutOpen}>
  <AlertDialogContent>
    <AlertDialogHeader>
      <AlertDialogTitle>Konfirmasi Checkout</AlertDialogTitle>
      <AlertDialogDescription>
        Apakah anda yakin ingin belanja dengan total bet {formatIDR(totalBelanja)}?
      </AlertDialogDescription>
    </AlertDialogHeader>
    <AlertDialogFooter>
      <AlertDialogCancel class="cursor-pointer" onclick={handleCancelCheckout}>
        Dibatalkan
      </AlertDialogCancel>
      <Button
        variant="outline"
        class="cursor-pointer"
        onclick={() => (checkoutOpen = false)}
      >
        Lanjut Belanja
      </Button>
      <AlertDialogAction class="cursor-pointer" onclick={handleConfirmCheckout}>
        Ya
      </AlertDialogAction>
    </AlertDialogFooter>
  </AlertDialogContent>
</AlertDialog>

<AlertDialog bind:open={minBetAlertOpen}>
  <AlertDialogContent>
    <AlertDialogHeader>
      <AlertDialogTitle>Bet Tidak Valid</AlertDialogTitle>
      <AlertDialogDescription>Minimal bet {formatIDR(MIN_BET)}</AlertDialogDescription>
    </AlertDialogHeader>
    <AlertDialogFooter>
      <AlertDialogAction class="cursor-pointer" onclick={() => (minBetAlertOpen = false)}>
        OK
      </AlertDialogAction>
    </AlertDialogFooter>
  </AlertDialogContent>
</AlertDialog>

<AlertDialog bind:open={insufficientCreditAlertOpen}>
  <AlertDialogContent>
    <AlertDialogHeader>
      <AlertDialogTitle>Informasi</AlertDialogTitle>
      <AlertDialogDescription>
        Saldo tidak cukup untuk melakukan transaksi
      </AlertDialogDescription>
    </AlertDialogHeader>
    <AlertDialogFooter>
      <AlertDialogAction
        class="cursor-pointer"
        onclick={() => (insufficientCreditAlertOpen = false)}
      >
        OK
      </AlertDialogAction>
    </AlertDialogFooter>
  </AlertDialogContent>
</AlertDialog>

<Dialog bind:open={transaksiModalOpen}>
  <DialogContent>
    <DialogHeader>
      <DialogTitle>Transaksi</DialogTitle>
      <DialogDescription>Daftar transaksi yang sudah checkout.</DialogDescription>
    </DialogHeader>
    {#if transaksiGroups.length === 0}
      <p class="text-muted-foreground py-4 text-center text-sm">Belum ada transaksi.</p>
    {:else}
      <Tabs value={transaksiGroups[0].type} class="min-w-0">
        <div class="min-w-0 scrollbar-none -mx-1 overflow-x-auto px-1">
          <TabsList class="w-max">
            {#each transaksiGroups as group (group.type)}
              <TabsTrigger value={group.type} class="cursor-pointer whitespace-nowrap">
                {group.type} - {formatIDR(group.sum)}
              </TabsTrigger>
            {/each}
          </TabsList>
        </div>
        {#each transaksiGroups as group (group.type)}
          <TabsContent value={group.type}>
            <div class="flex max-h-72 flex-col gap-2 overflow-y-auto">
              {#each group.entries as entry (entry.id)}
                <div class="rounded-lg border p-3">
                  <p class="text-sm font-medium">
                    {entry.type} - {entry.number.replace(/\*/g, "")}
                  </p>
                  <p class="text-xs text-blue-600">Bet : {formatIDR(entry.bet)}</p>
                </div>
              {/each}
            </div>
          </TabsContent>
        {/each}
      </Tabs>
    {/if}
  </DialogContent>
</Dialog>

<CameraOcr bind:open={cameraOcrOpen} bind:bets bind:minBetAlertOpen />

<BottomNav
  onTransaksiClick={() => (transaksiModalOpen = true)}
  onCameraClick={() => (cameraOcrOpen = true)}
/>

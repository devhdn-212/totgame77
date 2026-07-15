<script lang="ts">
  import Navbar from "$lib/components/Navbar.svelte";
  import BottomNav from "$lib/components/BottomNav.svelte";
  import { Button } from "$lib/components/ui/button";
  import { Card, CardContent } from "$lib/components/ui/card";
  import { Tabs, TabsList, TabsTrigger, TabsContent } from "$lib/components/ui/tabs";
  import { Input } from "$lib/components/ui/input";
  import { RadioGroup, RadioGroupItem } from "$lib/components/ui/radio-group";
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
  import Plus from "@lucide/svelte/icons/plus";
  import Trash2 from "@lucide/svelte/icons/trash-2";
  import ShoppingCart from "@lucide/svelte/icons/shopping-cart";
  import Dices from "@lucide/svelte/icons/dices";
  import Search from "@lucide/svelte/icons/search";
  import X from "@lucide/svelte/icons/x";
  import { z } from "zod";
  import Decimal from "decimal.js";
  import { scale } from "svelte/transition";
  import { untrack } from "svelte";

  const username = "ls0999212";
  const periode = "6512";
  const credit = "1000000";

  const betTypes = ["4D/3D/2D", "Colok", "5050", "Macau/Kombinasi", "Dasar", "Shio"];

  const fourDSubTabs = [
    "4D/3D/2D",
    "4D/3D/2D SET",
    "BOLAK BALIK",
    "WAP",
    "QUICK 2D",
    "POLA TARUNG",
    "3D DEPAN",
    "2D DEPAN",
    "2D TENGAH",
  ];

  let showFourDTabs = $state(true);
  let activeBetType = $state("4D/3D/2D");

  function handleBetTypeClick(type: string) {
    activeBetType = type;
    if (type === "4D/3D/2D") {
      showFourDTabs = !showFourDTabs;
    }
  }

  let numberInput = $state("");
  let betInput = $state("500");
  let qtyInput = $state("1");
  let formError = $state("");

  function handleNumberInput(e: Event) {
    const target = e.currentTarget as HTMLInputElement;
    const filtered = target.value.replace(/[^0-9*]/g, "").slice(0, 4);
    numberInput = filtered;
    target.value = filtered;
    formError = "";
  }

  const NUMBER_PATTERNS = ["DDDD", "DDD", "DDD*", "DD", "DD**", "*DD*"];

  function handleGenerateNumber() {
    const pattern = NUMBER_PATTERNS[Math.floor(Math.random() * NUMBER_PATTERNS.length)];
    numberInput = pattern
      .split("")
      .map((char) => (char === "D" ? String(Math.floor(Math.random() * 10)) : char))
      .join("");
    formError = "";
  }

  function handleBetInput(e: Event) {
    const target = e.currentTarget as HTMLInputElement;
    const digitsOnly = target.value.replace(/\D/g, "");
    betInput = digitsOnly;
    target.value = digitsOnly;
    formError = "";
  }

  function handleQtyInput(e: Event) {
    const target = e.currentTarget as HTMLInputElement;
    const digitsOnly = target.value.replace(/\D/g, "").slice(0, 2);
    qtyInput = digitsOnly;
    target.value = digitsOnly;
    formError = "";
  }

  const NUMBER_PATTERN = /^(?:\d{4}|\d{3}\*|\d{3}|\d{2}\*{2}|\*\d{2}\*|\d{2})$/;

  const betEntrySchema = z.object({
    number: z.string().regex(NUMBER_PATTERN, "Format nomor tidak valid"),
    bet: z.string().regex(/^[1-9]\d*$/, "Bet wajib diisi dengan angka"),
    qty: z.string().regex(/^(?:[1-9]|[1-9][0-9])$/, "Qty harus antara 1-99"),
  });

  function classifyBetType(number: string) {
    if (/^\d{4}$/.test(number)) return "4D";
    if (/^\d{3}$/.test(number)) return "3D";
    if (/^\d{3}\*$/.test(number)) return "3DD";
    if (/^\d{2}$/.test(number)) return "2D";
    if (/^\d{2}\*{2}$/.test(number)) return "2DD";
    return "2DT";
  }

  type BetEntry = { id: string; type: string; number: string; bet: string };

  let bets = $state<BetEntry[]>([]);
  let listTransaksi = $state<BetEntry[]>([]);

  const MIN_BET = 500;
  let minBetAlertOpen = $state(false);

  function handleAddBet() {
    const result = betEntrySchema.safeParse({ number: numberInput, bet: betInput, qty: qtyInput });
    if (!result.success) {
      formError = result.error.issues[0]?.message ?? "Input tidak valid";
      return;
    }
    if (Number(result.data.bet) < MIN_BET) {
      minBetAlertOpen = true;
      return;
    }
    formError = "";
    const type = classifyBetType(result.data.number);
    const newEntries = Array.from({ length: Number(result.data.qty) }, () => ({
      id: crypto.randomUUID(),
      type,
      number: result.data.number,
      bet: result.data.bet,
    }));
    bets = [...newEntries, ...bets];
    numberInput = "";
    betInput = "500";
    qtyInput = "1";
  }

  function handleDeleteBet(id: string) {
    bets = bets.filter((entry) => entry.id !== id);
  }

  const BET_TYPE_LABELS = ["4D", "3D", "3DD", "2D", "2DD", "2DT"];

  let setNumberInput = $state("");
  let setBetInputs = $state<Record<string, string>>(
    Object.fromEntries(BET_TYPE_LABELS.map((type) => [type, ""])),
  );
  let setFormError = $state("");

  function handleSetNumberInput(e: Event) {
    const target = e.currentTarget as HTMLInputElement;
    const digitsOnly = target.value.replace(/\D/g, "").slice(0, 4);
    setNumberInput = digitsOnly;
    target.value = digitsOnly;
    setFormError = "";
  }

  function handleSetBetInput(type: string, e: Event) {
    const target = e.currentTarget as HTMLInputElement;
    const digitsOnly = target.value.replace(/\D/g, "");
    setBetInputs[type] = digitsOnly;
    target.value = digitsOnly;
    setFormError = "";
  }

  function deriveSetNumber(type: string, base: string) {
    switch (type) {
      case "4D":
        return base;
      case "3D":
        return base.slice(1);
      case "3DD":
        return base.slice(0, 3);
      case "2D":
        return base.slice(2);
      case "2DD":
        return base.slice(0, 2);
      default:
        return base.slice(1, 3);
    }
  }

  function handleAddSetBet() {
    if (!/^\d{4}$/.test(setNumberInput)) {
      setFormError = "Nomor harus 4 digit";
      return;
    }

    const filled = BET_TYPE_LABELS.filter((type) => Number(setBetInputs[type]) > 0);
    if (filled.length === 0) {
      setFormError = "Isi minimal satu bet";
      return;
    }

    for (const type of filled) {
      if (Number(setBetInputs[type]) < MIN_BET) {
        minBetAlertOpen = true;
        return;
      }
    }

    setFormError = "";
    const newEntries = filled.map((type) => ({
      id: crypto.randomUUID(),
      type,
      number: deriveSetNumber(type, setNumberInput),
      bet: setBetInputs[type],
    }));
    bets = [...newEntries, ...bets];

    setNumberInput = "";
    setBetInputs = Object.fromEntries(BET_TYPE_LABELS.map((type) => [type, ""]));
  }

  let bbNumberInput = $state("");
  let bbBetInput = $state("500");
  let bbFormError = $state("");
  let bbWarning = $state("");
  let bbKembarMode = $state("tidak-kembar");

  function handleBbNumberInput(e: Event) {
    const target = e.currentTarget as HTMLInputElement;
    const digitsOnly = target.value.replace(/\D/g, "").slice(0, 15);
    bbNumberInput = digitsOnly;
    target.value = digitsOnly;
    bbFormError = "";
    bbWarning = "";
  }

  function handleBbBetInput(e: Event) {
    const target = e.currentTarget as HTMLInputElement;
    const digitsOnly = target.value.replace(/\D/g, "");
    bbBetInput = digitsOnly;
    target.value = digitsOnly;
    bbFormError = "";
  }

  function handleGenerateBbNumber() {
    const length = 4 + Math.floor(Math.random() * 3);
    const digits = Array.from({ length: 10 }, (_, i) => String(i)).sort(() => Math.random() - 0.5);
    bbNumberInput = digits.slice(0, length).join("");
    bbFormError = "";
    bbWarning = "";
  }

  function cartesianCombos(digits: string[], width: number): string[] {
    let combos = [""];
    for (let i = 0; i < width; i++) {
      combos = combos.flatMap((prefix) => digits.map((d) => prefix + d));
    }
    return combos;
  }

  function permutationsWithoutRepeat(digits: string[], width: number): string[] {
    const n = digits.length;
    if (n < width) return [];
    const results: string[] = [];
    const used = new Array(n).fill(false);
    function permute(current: string[]) {
      if (current.length === width) {
        results.push(current.join(""));
        return;
      }
      for (let i = 0; i < n; i++) {
        if (used[i]) continue;
        used[i] = true;
        current.push(digits[i]);
        permute(current);
        current.pop();
        used[i] = false;
      }
    }
    permute([]);
    return results;
  }

  const BB_LABEL_GROUPS: Record<number, string[]> = {
    4: ["4D"],
    3: ["3D", "3DD"],
    2: ["2D", "2DD", "2DT"],
  };

  function generateBB(nomor: string, allowKembar: boolean): { pasaran: string; nomor: string }[] {
    const output: { pasaran: string; nomor: string }[] = [];
    const uniqueDigits = [...new Set(nomor.split(""))];
    for (const width of [4, 3, 2]) {
      const combos = allowKembar
        ? cartesianCombos(uniqueDigits, width)
        : permutationsWithoutRepeat(uniqueDigits, width);
      for (const label of BB_LABEL_GROUPS[width]) {
        for (const combo of combos) {
          output.push({ pasaran: label, nomor: combo });
        }
      }
    }
    return output;
  }

  function handleAddBolakBalik() {
    if (bbNumberInput.length === 0) {
      bbFormError = "Nomor wajib diisi";
      return;
    }
    if (!/^[1-9]\d*$/.test(bbBetInput)) {
      bbFormError = "Bet wajib diisi dengan angka";
      return;
    }
    if (Number(bbBetInput) < MIN_BET) {
      minBetAlertOpen = true;
      return;
    }

    const allowKembar = bbKembarMode === "kembar";
    const results = generateBB(bbNumberInput, allowKembar);
    const uniqueDigitCount = new Set(bbNumberInput.split("")).size;

    const emptyLabels = !allowKembar
      ? Object.entries(BB_LABEL_GROUPS)
          .filter(([width]) => uniqueDigitCount < Number(width))
          .flatMap(([, labels]) => labels)
      : [];
    bbWarning =
      emptyLabels.length > 0
        ? `Pasaran ${emptyLabels.join(", ")} kosong karena digit unik kurang (unik: ${uniqueDigitCount})`
        : "";

    bbFormError = "";
    const newEntries = results.map((result) => ({
      id: crypto.randomUUID(),
      type: result.pasaran,
      number: result.nomor,
      bet: bbBetInput,
    }));
    bets = [...newEntries, ...bets];

    bbNumberInput = "";
    bbBetInput = "500";
  }

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

  function handleConfirmCheckout() {
    listTransaksi = [...bets, ...listTransaksi];
    bets = [];
    checkoutOpen = false;
  }

  function handleCancelCheckout() {
    bets = [];
    checkoutOpen = false;
  }

  let activeSubTab = $state(fourDSubTabs[0]);
  let lastValidSubTab = fourDSubTabs[0];
  let bolakBalikGuardOpen = $state(false);

  $effect(() => {
    if (activeSubTab === "BOLAK BALIK" && untrack(() => bets.length > 0)) {
      activeSubTab = lastValidSubTab;
      bolakBalikGuardOpen = true;
    } else {
      lastValidSubTab = activeSubTab;
    }
  });

  function handleGuardCheckout() {
    bolakBalikGuardOpen = false;
    checkoutOpen = true;
  }

  function handleGuardCancel() {
    bolakBalikGuardOpen = false;
  }

  let transaksiModalOpen = $state(false);

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
        <p class="text-primary text-lg font-semibold">{formatIDR(credit)}</p>
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

  {#if showFourDTabs}
    <Tabs bind:value={activeSubTab} class="mt-4">
      <div class="-mx-4 scrollbar-none overflow-x-auto px-4">
        <TabsList class="w-max">
          {#each fourDSubTabs as tab (tab)}
            <TabsTrigger value={tab} class="cursor-pointer whitespace-nowrap">{tab}</TabsTrigger>
          {/each}
        </TabsList>
      </div>
      {#each fourDSubTabs as tab (tab)}
        <TabsContent value={tab} class="text-sm">
          {#if tab === "4D/3D/2D"}
            <div class="flex flex-col gap-2">
              <div class="flex items-center gap-2">
                <Input
                  type="text"
                  inputmode="text"
                  pattern="[0-9*]*"
                  maxlength={4}
                  placeholder="0000"
                  class="text-center"
                  value={numberInput}
                  oninput={handleNumberInput}
                />
                <Button
                  size="icon"
                  variant="outline"
                  class="shrink-0 cursor-pointer"
                  aria-label="Generate nomor acak"
                  onclick={handleGenerateNumber}
                >
                  <Dices />
                </Button>
                <Input
                  type="text"
                  inputmode="numeric"
                  pattern="[0-9]*"
                  maxlength={2}
                  placeholder="Qty"
                  class="w-14 shrink-0 text-center"
                  value={qtyInput}
                  oninput={handleQtyInput}
                />
                <Input
                  type="text"
                  inputmode="numeric"
                  pattern="[0-9]*"
                  placeholder="Bet"
                  class="text-center"
                  value={betInput}
                  oninput={handleBetInput}
                />
              </div>
              <Button
                class="w-full cursor-pointer bg-blue-100 text-blue-700 hover:bg-blue-200"
                onclick={handleAddBet}
              >
                <Plus />
                Tambah
              </Button>
            </div>
            {#if formError}
              <p class="text-destructive mt-1 text-xs">{formError}</p>
            {/if}
          {:else if tab === "4D/3D/2D SET"}
            <div class="flex flex-col gap-2">
              <Input
                type="text"
                inputmode="numeric"
                pattern="[0-9]*"
                maxlength={4}
                placeholder="Nomor 4 digit"
                class="text-center"
                value={setNumberInput}
                oninput={handleSetNumberInput}
              />
              <div class="grid grid-cols-3 gap-2">
                {#each BET_TYPE_LABELS as type (type)}
                  <div class="flex flex-col gap-1">
                    <span class="text-muted-foreground text-center text-xs">{type}</span>
                    <Input
                      type="text"
                      inputmode="numeric"
                      pattern="[0-9]*"
                      placeholder="Bet"
                      class="text-center"
                      value={setBetInputs[type]}
                      oninput={(e) => handleSetBetInput(type, e)}
                    />
                  </div>
                {/each}
              </div>
              <Button
                class="w-full cursor-pointer bg-blue-100 text-blue-700 hover:bg-blue-200"
                onclick={handleAddSetBet}
              >
                <Plus />
                Tambah
              </Button>
            </div>
            {#if setFormError}
              <p class="text-destructive mt-1 text-xs">{setFormError}</p>
            {/if}
          {:else if tab === "BOLAK BALIK"}
            <div class="flex flex-col gap-2">
              <div class="flex items-center gap-2">
                <Input
                  type="text"
                  inputmode="numeric"
                  pattern="[0-9]*"
                  maxlength={15}
                  placeholder="Nomor"
                  class="text-center"
                  value={bbNumberInput}
                  oninput={handleBbNumberInput}
                />
                <Button
                  size="icon"
                  variant="outline"
                  class="shrink-0 cursor-pointer"
                  aria-label="Generate nomor acak"
                  onclick={handleGenerateBbNumber}
                >
                  <Dices />
                </Button>
                <Input
                  type="text"
                  inputmode="numeric"
                  pattern="[0-9]*"
                  placeholder="Bet"
                  class="text-center"
                  value={bbBetInput}
                  oninput={handleBbBetInput}
                />
              </div>
              <RadioGroup bind:value={bbKembarMode} class="flex justify-center gap-4">
                <label class="flex items-center gap-2 text-sm">
                  <RadioGroupItem value="tidak-kembar" />
                  Tidak Kembar
                </label>
                <label class="flex items-center gap-2 text-sm">
                  <RadioGroupItem value="kembar" />
                  Kembar
                </label>
              </RadioGroup>
              <Button
                class="w-full cursor-pointer bg-blue-100 text-blue-700 hover:bg-blue-200"
                onclick={handleAddBolakBalik}
              >
                <Plus />
                Tambah
              </Button>
            </div>
            {#if bbFormError}
              <p class="text-destructive mt-1 text-xs">{bbFormError}</p>
            {/if}
            {#if bbWarning}
              <p class="mt-1 text-xs text-amber-600">{bbWarning}</p>
            {/if}
          {:else}
            <p class="text-muted-foreground">{tab}</p>
          {/if}
        </TabsContent>
      {/each}
    </Tabs>

    {#if bets.length > 0}
      <div class="mt-4 flex items-center justify-between">
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
            onclick={() => (checkoutOpen = true)}
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
          <div class="flex items-center justify-between rounded-lg border p-3">
            <div>
              <p class="text-sm font-medium">{entry.type} - {entry.number.replace(/\*/g, "")}</p>
              <p class="text-muted-foreground text-xs">Bet : {formatIDR(entry.bet)}</p>
            </div>
            <Button
              size="icon"
              variant="ghost"
              class="text-destructive cursor-pointer"
              aria-label="Hapus"
              onclick={() => handleDeleteBet(entry.id)}
            >
              <Trash2 />
            </Button>
          </div>
        {/each}
      </div>
    {/if}
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

<AlertDialog bind:open={bolakBalikGuardOpen}>
  <AlertDialogContent onCloseAutoFocus={(e: Event) => e.preventDefault()}>
    <AlertDialogHeader>
      <AlertDialogTitle>Informasi</AlertDialogTitle>
      <AlertDialogDescription>
        Harap melakukan checkout sebelum melakukan belanja bolak balik.
      </AlertDialogDescription>
    </AlertDialogHeader>
    <AlertDialogFooter>
      <AlertDialogCancel class="cursor-pointer" onclick={handleGuardCancel}>
        Batalkan
      </AlertDialogCancel>
      <AlertDialogAction class="cursor-pointer" onclick={handleGuardCheckout}>
        Checkout
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
                  <p class="text-muted-foreground text-xs">Bet : {formatIDR(entry.bet)}</p>
                </div>
              {/each}
            </div>
          </TabsContent>
        {/each}
      </Tabs>
    {/if}
  </DialogContent>
</Dialog>

<BottomNav onTransaksiClick={() => (transaksiModalOpen = true)} />

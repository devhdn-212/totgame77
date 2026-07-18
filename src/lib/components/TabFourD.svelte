<script lang="ts">
  import Form432 from "$lib/components/Form432.svelte";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Textarea } from "$lib/components/ui/textarea";
  import { RadioGroup, RadioGroupItem } from "$lib/components/ui/radio-group";
  import { Select, SelectTrigger, SelectContent, SelectItem } from "$lib/components/ui/select";
  import { Tabs, TabsList, TabsTrigger, TabsContent } from "$lib/components/ui/tabs";
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
  import Plus from "@lucide/svelte/icons/plus";
  import Dices from "@lucide/svelte/icons/dices";
  import { untrack } from "svelte";
  import { formatIDR } from "$lib/utils";
  import Decimal from "decimal.js";

  type BetEntry = {
    id: string;
    type: string;
    number: string;
    bet: string;
    deletable?: boolean;
    kombinasi?: string;
  };

  const KOMBINASI_OPTIONS = ["DISC", "FULL", "BB"];

  let {
    bets = $bindable(),
    minBetAlertOpen = $bindable(false),
    onCheckoutClick,
  }: {
    bets: BetEntry[];
    minBetAlertOpen: boolean;
    onCheckoutClick: () => void;
  } = $props();

  const MIN_BET = 500;
  const BET_TYPE_LABELS = ["4D", "3D", "3DD", "2D", "2DD", "2DT"];

  const fourDSubTabs = ["4D/3D/2D", "4D/3D/2D SET", "BOLAK BALIK", "WAP", "QUICK 2D"];

  let activeSubTab = $state(fourDSubTabs[0]);
  let lastValidSubTab = fourDSubTabs[0];
  let bolakBalikGuardOpen = $state(false);
  let leaveBolakBalikGuardOpen = $state(false);

  let totalBelanja = $derived(bets.reduce((sum, entry) => sum.plus(entry.bet), new Decimal(0)));

  $effect(() => {
    const enteringBolakBalik = activeSubTab === "BOLAK BALIK" && lastValidSubTab !== "BOLAK BALIK";
    const leavingBolakBalik = activeSubTab !== "BOLAK BALIK" && lastValidSubTab === "BOLAK BALIK";

    if (enteringBolakBalik && untrack(() => bets.length > 0)) {
      activeSubTab = lastValidSubTab;
      bolakBalikGuardOpen = true;
    } else if (leavingBolakBalik && untrack(() => bets.length > 0)) {
      activeSubTab = "BOLAK BALIK";
      leaveBolakBalikGuardOpen = true;
    } else {
      lastValidSubTab = activeSubTab;
    }
  });

  function handleGuardCheckout() {
    bolakBalikGuardOpen = false;
    onCheckoutClick();
  }

  function handleGuardCancel() {
    bolakBalikGuardOpen = false;
  }

  function handleLeaveGuardCheckout() {
    leaveBolakBalikGuardOpen = false;
    onCheckoutClick();
  }

  function handleLeaveGuardCancel() {
    bets = [];
    leaveBolakBalikGuardOpen = false;
  }

  let setNumberInput = $state("");
  let setBetInputs = $state<Record<string, string>>(
    Object.fromEntries(BET_TYPE_LABELS.map((type) => [type, ""])),
  );
  let setKombinasiInput = $state("DISC");
  let setFormError = $state("");

  function handleSetNumberInput(e: Event) {
    const target = e.currentTarget as HTMLInputElement;
    const digitsOnly = target.value.replace(/\D/g, "").slice(0, 4);
    setNumberInput = digitsOnly;
    target.value = digitsOnly;
    setFormError = "";
  }

  function handleGenerateSetNumber() {
    const digits = Array.from({ length: 10 }, (_, i) => String(i)).sort(() => Math.random() - 0.5);
    setNumberInput = digits.slice(0, 4).join("");
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
      kombinasi: setKombinasiInput,
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
  let bbKombinasiInput = $state("DISC");

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
      deletable: false,
      kombinasi: bbKombinasiInput,
    }));
    bets = [...newEntries, ...bets];

    bbNumberInput = "";
    bbBetInput = "500";
  }

  let wapInput = $state("");
  let wapKombinasiInput = $state("DISC");
  let wapFormError = $state("");

  function handleWapInput(e: Event) {
    const target = e.currentTarget as HTMLTextAreaElement;
    const filtered = target.value.replace(/[^0-9*#,]/g, "");
    wapInput = filtered;
    target.value = filtered;
    wapFormError = "";
  }

  function classifyByLength(number: string) {
    if (number.length === 4) return "4D";
    if (number.length === 3) return "3D";
    if (number.length === 2) return "2D";
    return null;
  }

  function handleAddWap() {
    const raw = wapInput.trim();
    if (raw.length === 0) {
      wapFormError = "Nomor wajib diisi";
      return;
    }

    const groups = raw.split(",").filter((g) => g.trim() !== "");
    if (groups.length === 0) {
      wapFormError = "Nomor wajib diisi";
      return;
    }

    type ParsedEntry = { type: string; number: string; bet: string };
    const parsed: ParsedEntry[] = [];

    for (const group of groups) {
      const [numbersPart, betPart] = group.split("#");

      if (!betPart || !/^[1-9]\d*$/.test(betPart)) {
        wapFormError = `Bet wajib diisi (format nomor#bet) pada "${group}"`;
        return;
      }
      const bet = betPart;
      if (Number(bet) < MIN_BET) {
        minBetAlertOpen = true;
        return;
      }

      const numbers = (numbersPart ?? "").split("*").filter((n) => n !== "");
      if (numbers.length === 0) {
        wapFormError = `Nomor tidak valid pada "${group}"`;
        return;
      }

      for (const number of numbers) {
        const type = classifyByLength(number);
        if (!type) {
          wapFormError = `Nomor "${number}" harus 2-4 digit`;
          return;
        }
        parsed.push({ type, number, bet });
      }
    }

    wapFormError = "";
    const newEntries = parsed.map((entry) => ({
      id: crypto.randomUUID(),
      type: entry.type,
      number: entry.number,
      bet: entry.bet,
      kombinasi: wapKombinasiInput,
    }));
    bets = [...newEntries, ...bets];

    wapInput = "";
  }

  const QUICK_KONDISI_OPTIONS = ["BESAR", "KECIL", "GANJIL", "GENAP"];
  const QUICK_PASARAN_OPTIONS = ["2D", "2DD", "2DT"];

  let quickKondisi = $state("BESAR");
  let quickPasaran = $state("2D");
  let quickKombinasiInput = $state("DISC");
  let quickBetInput = $state("500");
  let quickFormError = $state("");

  function handleQuickBetInput(e: Event) {
    const target = e.currentTarget as HTMLInputElement;
    const digitsOnly = target.value.replace(/\D/g, "");
    quickBetInput = digitsOnly;
    target.value = digitsOnly;
    quickFormError = "";
  }

  function handleAddQuick2D() {
    if (!/^[1-9]\d*$/.test(quickBetInput)) {
      quickFormError = "Bet wajib diisi dengan angka";
      return;
    }
    if (Number(quickBetInput) < MIN_BET) {
      minBetAlertOpen = true;
      return;
    }

    const numbers: string[] = [];
    for (let n = 0; n <= 99; n++) {
      const matches =
        (quickKondisi === "BESAR" && n > 49) ||
        (quickKondisi === "KECIL" && n < 50) ||
        (quickKondisi === "GANJIL" && n % 2 === 1) ||
        (quickKondisi === "GENAP" && n % 2 === 0);
      if (matches) numbers.push(String(n).padStart(2, "0"));
    }

    quickFormError = "";
    const newEntries = numbers.map((number) => ({
      id: crypto.randomUUID(),
      type: quickPasaran,
      number,
      bet: quickBetInput,
      kombinasi: quickKombinasiInput,
    }));
    bets = [...newEntries, ...bets];

    quickBetInput = "500";
  }
</script>

<Tabs bind:value={activeSubTab} class="mt-4">
  <div class="-mx-4 scrollbar-none overflow-x-auto px-4">
    <TabsList class="w-max">
      {#each fourDSubTabs as tab (tab)}
        <TabsTrigger
          value={tab}
          class="cursor-pointer whitespace-nowrap data-active:border-orange-300! data-active:bg-orange-100! data-active:text-orange-700!"
        >
          {tab}
        </TabsTrigger>
      {/each}
    </TabsList>
  </div>
  {#each fourDSubTabs as tab (tab)}
    <TabsContent value={tab} class="text-sm">
      {#if tab === "4D/3D/2D"}
        <Form432 bind:bets bind:minBetAlertOpen />
      {:else if tab === "4D/3D/2D SET"}
        <div class="flex flex-col gap-2">
          <div class="flex items-center gap-2">
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
            <Button
              size="icon"
              variant="outline"
              class="shrink-0 cursor-pointer"
              aria-label="Generate nomor acak"
              onclick={handleGenerateSetNumber}
            >
              <Dices />
            </Button>
          </div>
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
          <Select type="single" bind:value={setKombinasiInput}>
            <SelectTrigger class="w-full">
              {setKombinasiInput}
            </SelectTrigger>
            <SelectContent>
              {#each KOMBINASI_OPTIONS as option (option)}
                <SelectItem value={option}>{option}</SelectItem>
              {/each}
            </SelectContent>
          </Select>
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
          <Select type="single" bind:value={bbKombinasiInput}>
            <SelectTrigger class="w-full">
              {bbKombinasiInput}
            </SelectTrigger>
            <SelectContent>
              {#each KOMBINASI_OPTIONS as option (option)}
                <SelectItem value={option}>{option}</SelectItem>
              {/each}
            </SelectContent>
          </Select>
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
      {:else if tab === "WAP"}
        <div class="flex flex-col gap-2">
          <Textarea
            placeholder="contoh : 1234*234*34#1000,34*235*35#5000"
            class="min-h-24"
            value={wapInput}
            oninput={handleWapInput}
          />
          <Select type="single" bind:value={wapKombinasiInput}>
            <SelectTrigger class="w-full">
              {wapKombinasiInput}
            </SelectTrigger>
            <SelectContent>
              {#each KOMBINASI_OPTIONS as option (option)}
                <SelectItem value={option}>{option}</SelectItem>
              {/each}
            </SelectContent>
          </Select>
          <Button
            class="w-full cursor-pointer bg-blue-100 text-blue-700 hover:bg-blue-200"
            onclick={handleAddWap}
          >
            <Plus />
            Tambah
          </Button>
        </div>
        {#if wapFormError}
          <p class="text-destructive mt-1 text-xs">{wapFormError}</p>
        {/if}
      {:else if tab === "QUICK 2D"}
        <div class="flex flex-col gap-2">
          <div class="flex items-center gap-2">
            <Select type="single" bind:value={quickKondisi}>
              <SelectTrigger class="w-full">
                {quickKondisi}
              </SelectTrigger>
              <SelectContent>
                {#each QUICK_KONDISI_OPTIONS as option (option)}
                  <SelectItem value={option}>{option}</SelectItem>
                {/each}
              </SelectContent>
            </Select>
            <Select type="single" bind:value={quickPasaran}>
              <SelectTrigger class="w-full">
                {quickPasaran}
              </SelectTrigger>
              <SelectContent>
                {#each QUICK_PASARAN_OPTIONS as option (option)}
                  <SelectItem value={option}>{option}</SelectItem>
                {/each}
              </SelectContent>
            </Select>
            <Select type="single" bind:value={quickKombinasiInput}>
              <SelectTrigger class="w-full">
                {quickKombinasiInput}
              </SelectTrigger>
              <SelectContent>
                {#each KOMBINASI_OPTIONS as option (option)}
                  <SelectItem value={option}>{option}</SelectItem>
                {/each}
              </SelectContent>
            </Select>
          </div>
          <Input
            type="text"
            inputmode="numeric"
            pattern="[0-9]*"
            placeholder="Bet"
            class="text-center"
            value={quickBetInput}
            oninput={handleQuickBetInput}
          />
          <Button
            class="w-full cursor-pointer bg-blue-100 text-blue-700 hover:bg-blue-200"
            onclick={handleAddQuick2D}
          >
            <Plus />
            Tambah
          </Button>
        </div>
        {#if quickFormError}
          <p class="text-destructive mt-1 text-xs">{quickFormError}</p>
        {/if}
      {:else}
        <p class="text-muted-foreground">{tab}</p>
      {/if}
    </TabsContent>
  {/each}
</Tabs>

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

<AlertDialog bind:open={leaveBolakBalikGuardOpen}>
  <AlertDialogContent onCloseAutoFocus={(e: Event) => e.preventDefault()}>
    <AlertDialogHeader>
      <AlertDialogTitle>Harap melakukan pembelian</AlertDialogTitle>
      <AlertDialogDescription>
        Total Belanja: {formatIDR(totalBelanja)}
      </AlertDialogDescription>
    </AlertDialogHeader>
    <AlertDialogFooter>
      <AlertDialogCancel class="cursor-pointer" onclick={handleLeaveGuardCancel}>
        Dibatalkan
      </AlertDialogCancel>
      <Button
        variant="outline"
        class="cursor-pointer"
        onclick={() => (leaveBolakBalikGuardOpen = false)}
      >
        Lanjut Belanja
      </Button>
      <AlertDialogAction class="cursor-pointer" onclick={handleLeaveGuardCheckout}>
        Checkout
      </AlertDialogAction>
    </AlertDialogFooter>
  </AlertDialogContent>
</AlertDialog>

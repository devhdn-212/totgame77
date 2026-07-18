<script lang="ts">
  import { Input } from "$lib/components/ui/input";
  import { Button } from "$lib/components/ui/button";
  import { Select, SelectTrigger, SelectContent, SelectItem } from "$lib/components/ui/select";
  import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
  } from "$lib/components/ui/dialog";
  import Plus from "@lucide/svelte/icons/plus";
  import Dices from "@lucide/svelte/icons/dices";
  import Info from "@lucide/svelte/icons/info";
  import { z } from "zod";
  import { BET_TYPE_LIMITS, formatIDR } from "$lib/utils";

  type BetEntry = {
    id: string;
    type: string;
    number: string;
    bet: string;
    deletable?: boolean;
    kombinasi?: string;
  };

  let {
    bets = $bindable(),
    minBetAlertOpen = $bindable(false),
    minBetRequired = $bindable(500),
  }: {
    bets: BetEntry[];
    minBetAlertOpen: boolean;
    minBetRequired: number;
  } = $props();

  const MIN_BET = 500;

  const KOMBINASI_OPTIONS = ["DISC", "FULL", "BB"];
  const QTY_OPTIONS = Array.from({ length: 10 }, (_, i) => String(i + 1));

  let numberInput = $state("");
  let betInput = $state("500");
  let qtyInput = $state("1");
  let kombinasiInput = $state("DISC");
  let formError = $state("");
  let infoOpen = $state(false);

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

  function handleAddBet() {
    const result = betEntrySchema.safeParse({ number: numberInput, bet: betInput, qty: qtyInput });
    if (!result.success) {
      formError = result.error.issues[0]?.message ?? "Input tidak valid";
      return;
    }
    if (Number(result.data.bet) < MIN_BET) {
      minBetRequired = MIN_BET;
      minBetAlertOpen = true;
      return;
    }
    const type = classifyBetType(result.data.number);
    const maxBet = BET_TYPE_LIMITS[type]?.maxBet;
    if (maxBet !== undefined && Number(result.data.bet) > maxBet) {
      formError = `Bet melebihi maksimal ${formatIDR(maxBet)} untuk ${type}`;
      return;
    }

    formError = "";
    const newEntries = Array.from({ length: Number(result.data.qty) }, () => ({
      id: crypto.randomUUID(),
      type,
      number: result.data.number,
      bet: result.data.bet,
      kombinasi: kombinasiInput,
    }));
    bets = [...newEntries, ...bets];
    numberInput = "";
    betInput = "500";
    qtyInput = "1";
  }
</script>

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
    <Select type="single" bind:value={qtyInput}>
      <SelectTrigger class="w-16 shrink-0">
        {qtyInput}
      </SelectTrigger>
      <SelectContent>
        {#each QTY_OPTIONS as option (option)}
          <SelectItem value={option}>{option}</SelectItem>
        {/each}
      </SelectContent>
    </Select>
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
  <Select type="single" bind:value={kombinasiInput}>
    <SelectTrigger class="w-full">
      {kombinasiInput}
    </SelectTrigger>
    <SelectContent>
      {#each KOMBINASI_OPTIONS as option (option)}
        <SelectItem value={option}>{option}</SelectItem>
      {/each}
    </SelectContent>
  </Select>
  <div class="flex items-center gap-2">
    <Button
      class="flex-1 cursor-pointer bg-blue-100 text-blue-700 hover:bg-blue-200"
      onclick={handleAddBet}
    >
      <Plus />
      Tambah
    </Button>
    <Button
      size="icon"
      variant="outline"
      class="shrink-0 cursor-pointer"
      aria-label="Informasi"
      onclick={() => (infoOpen = true)}
    >
      <Info />
    </Button>
  </div>
</div>
{#if formError}
  <p class="text-destructive mt-1 text-xs">{formError}</p>
{/if}

<Dialog bind:open={infoOpen}>
  <DialogContent>
    <DialogHeader>
      <DialogTitle>Informasi</DialogTitle>
      <DialogDescription>Cara penulisan nomor</DialogDescription>
    </DialogHeader>
    <pre class="rounded-lg border bg-muted/50 p-3 font-mono text-sm leading-relaxed">4D  : 1234
3D  : 234
3DD : 234*
2D  : 34
2DD : 23**
2DT : *23*</pre>
  </DialogContent>
</Dialog>

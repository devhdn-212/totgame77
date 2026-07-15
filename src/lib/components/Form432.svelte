<script lang="ts">
  import { Input } from "$lib/components/ui/input";
  import { Button } from "$lib/components/ui/button";
  import Plus from "@lucide/svelte/icons/plus";
  import Dices from "@lucide/svelte/icons/dices";
  import { z } from "zod";

  type BetEntry = {
    id: string;
    type: string;
    number: string;
    bet: string;
    deletable?: boolean;
  };

  let {
    bets = $bindable(),
    minBetAlertOpen = $bindable(false),
  }: {
    bets: BetEntry[];
    minBetAlertOpen: boolean;
  } = $props();

  const MIN_BET = 500;

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

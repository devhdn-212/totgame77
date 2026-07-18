<script lang="ts">
  import { Input } from "$lib/components/ui/input";
  import { Button } from "$lib/components/ui/button";
  import { Select, SelectTrigger, SelectContent, SelectItem } from "$lib/components/ui/select";
  import Plus from "@lucide/svelte/icons/plus";
  import Dices from "@lucide/svelte/icons/dices";
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
    type,
    bets = $bindable(),
    minBetAlertOpen = $bindable(false),
    minBetRequired = $bindable(500),
  }: {
    type: "Bebas" | "Macau" | "Naga" | "JITU";
    bets: BetEntry[];
    minBetAlertOpen: boolean;
    minBetRequired: number;
  } = $props();

  const MIN_BET = 1000;

  const DIGIT_OPTIONS = ["0", "1", "2", "3", "4", "5", "6", "7", "8", "9"];
  const POSITION_OPTIONS = ["AS", "KOP", "KEPALA", "EKOR"];
  const DIGIT_COUNT: Record<string, number> = { Bebas: 1, Macau: 2, Naga: 3, JITU: 1 };
  const TYPE_LABEL: Record<string, string> = {
    Bebas: "COLOK_BEBAS",
    Macau: "COLOK_MACAU",
    Naga: "COLOK_NAGA",
    JITU: "COLOK_JITU",
  };

  let digitSelects = $state<string[]>(Array(DIGIT_COUNT[type]).fill("0"));
  let positionInput = $state("AS");
  let betInput = $state("1000");
  let formError = $state("");

  $effect(() => {
    // Reset digit slots whenever the active Colok sub-tab changes.
    digitSelects = Array(DIGIT_COUNT[type]).fill("0");
    formError = "";
  });

  function handleGenerateDigit(index: number) {
    digitSelects[index] = DIGIT_OPTIONS[Math.floor(Math.random() * DIGIT_OPTIONS.length)];
    formError = "";
  }

  function handleBetInput(e: Event) {
    const target = e.currentTarget as HTMLInputElement;
    const digitsOnly = target.value.replace(/\D/g, "");
    betInput = digitsOnly;
    target.value = digitsOnly;
    formError = "";
  }

  function handleAddBet() {
    if (!/^[1-9]\d*$/.test(betInput)) {
      formError = "Bet wajib diisi dengan angka";
      return;
    }
    if (Number(betInput) < MIN_BET) {
      minBetRequired = MIN_BET;
      minBetAlertOpen = true;
      return;
    }
    const maxBet = BET_TYPE_LIMITS[TYPE_LABEL[type]]?.maxBet;
    if (maxBet !== undefined && Number(betInput) > maxBet) {
      formError = `Bet melebihi maksimal ${formatIDR(maxBet)}`;
      return;
    }

    formError = "";
    const number =
      type === "JITU" ? `${digitSelects[0]}_${positionInput}` : digitSelects.join("");

    const newEntry: BetEntry = {
      id: crypto.randomUUID(),
      type: TYPE_LABEL[type],
      number,
      bet: betInput,
      kombinasi: "DISC",
    };
    bets = [newEntry, ...bets];

    digitSelects = Array(DIGIT_COUNT[type]).fill("0");
    betInput = "1000";
  }
</script>

<div class="flex flex-col gap-2">
  <div class="flex items-center gap-2">
    {#each digitSelects as digit, i (i)}
      <Select type="single" bind:value={digitSelects[i]}>
        <SelectTrigger class="w-full">
          {digit}
        </SelectTrigger>
        <SelectContent>
          {#each DIGIT_OPTIONS as option (option)}
            <SelectItem value={option}>{option}</SelectItem>
          {/each}
        </SelectContent>
      </Select>
      <Button
        size="icon"
        variant="outline"
        class="shrink-0 cursor-pointer"
        aria-label="Generate digit acak"
        onclick={() => handleGenerateDigit(i)}
      >
        <Dices />
      </Button>
    {/each}
    {#if type === "JITU"}
      <Select type="single" bind:value={positionInput}>
        <SelectTrigger class="w-full">
          {positionInput}
        </SelectTrigger>
        <SelectContent>
          {#each POSITION_OPTIONS as option (option)}
            <SelectItem value={option}>{option}</SelectItem>
          {/each}
        </SelectContent>
      </Select>
    {/if}
  </div>
  <Input
    type="text"
    inputmode="numeric"
    pattern="[0-9]*"
    placeholder="Bet"
    class="text-center"
    value={betInput}
    oninput={handleBetInput}
  />
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

<script lang="ts">
  import { Input } from "$lib/components/ui/input";
  import { Button } from "$lib/components/ui/button";
  import { Select, SelectTrigger, SelectContent, SelectItem } from "$lib/components/ui/select";
  import Plus from "@lucide/svelte/icons/plus";
  import { BET_TYPE_LIMITS, formatIDR } from "$lib/utils";

  type BetEntry = {
    id: string;
    type: string;
    number: string;
    bet: string;
    deletable?: boolean;
    kombinasi?: string;
    compactDisplay?: boolean;
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

  const MIN_BET = 1000;

  const DASAR_OPTIONS = ["GANJIL", "GENAP", "BESAR", "KECIL"];

  let dasarInput = $state(DASAR_OPTIONS[0]);
  let betInput = $state("1000");
  let formError = $state("");

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
    const maxBet = BET_TYPE_LIMITS.DASAR?.maxBet;
    if (maxBet !== undefined) {
      if (Number(betInput) > maxBet) {
        formError = `Bet melebihi maximal Bet dan maximal bet ${formatIDR(maxBet)}`;
        return;
      }
      const existingTotal = bets
        .filter((entry) => entry.type === "DASAR" && entry.number === dasarInput)
        .reduce((sum, entry) => sum + Number(entry.bet), 0);
      if (existingTotal + Number(betInput) > maxBet) {
        const remaining = Math.max(maxBet - existingTotal, 0);
        formError = `Bet melebihi limit total ${formatIDR(maxBet)} untuk nomor ini (sisa ${formatIDR(remaining)})`;
        return;
      }
    }

    formError = "";
    const newEntry: BetEntry = {
      id: crypto.randomUUID(),
      type: "DASAR",
      number: dasarInput,
      bet: betInput,
      kombinasi: "DISC",
    };
    bets = [newEntry, ...bets];

    betInput = "1000";
  }
</script>

<div class="flex flex-col gap-2">
  <div class="flex items-center gap-2">
    <Select type="single" bind:value={dasarInput}>
      <SelectTrigger class="w-full">
        {dasarInput}
      </SelectTrigger>
      <SelectContent>
        {#each DASAR_OPTIONS as option (option)}
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

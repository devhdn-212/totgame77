<script lang="ts">
  import { Input } from "$lib/components/ui/input";
  import { Button } from "$lib/components/ui/button";
  import { Select, SelectTrigger, SelectContent, SelectItem } from "$lib/components/ui/select";
  import Plus from "@lucide/svelte/icons/plus";

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

  const POSISI_OPTIONS = ["BELAKANG", "TENGAH", "DEPAN"];
  const BESAR_KECIL_OPTIONS = ["BESAR", "KECIL"];
  const GENAP_GANJIL_OPTIONS = ["GENAP", "GANJIL"];

  let posisiInput = $state(POSISI_OPTIONS[0]);
  let besarKecilInput = $state(BESAR_KECIL_OPTIONS[0]);
  let genapGanjilInput = $state(GENAP_GANJIL_OPTIONS[0]);
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

    formError = "";
    const newEntry: BetEntry = {
      id: crypto.randomUUID(),
      type: "MACAU_KOMBINASI",
      number: `${posisiInput}_${besarKecilInput}_${genapGanjilInput}`,
      bet: betInput,
      kombinasi: "DISC",
      compactDisplay: true,
    };
    bets = [newEntry, ...bets];

    betInput = "1000";
  }
</script>

<div class="flex flex-col gap-2">
  <div class="flex items-center gap-2">
    <Select type="single" bind:value={posisiInput}>
      <SelectTrigger class="w-full">
        {posisiInput}
      </SelectTrigger>
      <SelectContent>
        {#each POSISI_OPTIONS as option (option)}
          <SelectItem value={option}>{option}</SelectItem>
        {/each}
      </SelectContent>
    </Select>
    <Select type="single" bind:value={besarKecilInput}>
      <SelectTrigger class="w-full">
        {besarKecilInput}
      </SelectTrigger>
      <SelectContent>
        {#each BESAR_KECIL_OPTIONS as option (option)}
          <SelectItem value={option}>{option}</SelectItem>
        {/each}
      </SelectContent>
    </Select>
    <Select type="single" bind:value={genapGanjilInput}>
      <SelectTrigger class="w-full">
        {genapGanjilInput}
      </SelectTrigger>
      <SelectContent>
        {#each GENAP_GANJIL_OPTIONS as option (option)}
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

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
    type,
    bets = $bindable(),
    minBetAlertOpen = $bindable(false),
    minBetRequired = $bindable(500),
  }: {
    type: "Umum" | "Special" | "Kombinasi";
    bets: BetEntry[];
    minBetAlertOpen: boolean;
    minBetRequired: number;
  } = $props();

  const MIN_BET = 1000;

  const UMUM_OPTIONS = ["BESAR", "KECIL", "GENAP", "GANJIL", "TENGAH", "TEPI"];
  const SPECIAL_POSISI_OPTIONS = ["AS", "KOP", "KEPALA", "EKOR"];
  const SPECIAL_KONDISI_OPTIONS = ["BESAR", "KECIL", "GENAP", "GANJIL"];
  const KOMBINASI_POSISI_OPTIONS = ["BELAKANG", "TENGAH", "DEPAN"];
  const KOMBINASI_JENIS_OPTIONS = ["MONO", "STEREO", "KEMBANG", "KEMPIS", "KEMBAR"];

  const TYPE_LABEL: Record<string, string> = {
    Umum: "50_50_UMUM",
    Special: "50_50_SPECIAL",
    Kombinasi: "50_50_KOMBINASI",
  };

  let umumInput = $state(UMUM_OPTIONS[0]);
  let specialPosisiInput = $state(SPECIAL_POSISI_OPTIONS[0]);
  let specialKondisiInput = $state(SPECIAL_KONDISI_OPTIONS[0]);
  let kombinasiPosisiInput = $state(KOMBINASI_POSISI_OPTIONS[0]);
  let kombinasiJenisInput = $state(KOMBINASI_JENIS_OPTIONS[0]);
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
    const betType = TYPE_LABEL[type];
    const number =
      type === "Special"
        ? `${specialPosisiInput}_${specialKondisiInput}`
        : type === "Kombinasi"
          ? `${kombinasiPosisiInput}_${kombinasiJenisInput}`
          : umumInput;

    const maxBet = BET_TYPE_LIMITS[betType]?.maxBet;
    if (maxBet !== undefined) {
      if (Number(betInput) > maxBet) {
        formError = `Bet melebihi maximal Bet dan maximal bet ${formatIDR(maxBet)}`;
        return;
      }
      const existingTotal = bets
        .filter((entry) => entry.type === betType && entry.number === number)
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
      type: TYPE_LABEL[type],
      number,
      bet: betInput,
      kombinasi: "DISC",
    };
    bets = [newEntry, ...bets];

    betInput = "1000";
  }
</script>

<div class="flex flex-col gap-2">
  {#if type === "Umum"}
    <div class="flex items-center gap-2">
      <Select type="single" bind:value={umumInput}>
        <SelectTrigger class="w-full">
          {umumInput}
        </SelectTrigger>
        <SelectContent>
          {#each UMUM_OPTIONS as option (option)}
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
  {:else if type === "Special"}
    <div class="flex items-center gap-2">
      <Select type="single" bind:value={specialPosisiInput}>
        <SelectTrigger class="w-full">
          {specialPosisiInput}
        </SelectTrigger>
        <SelectContent>
          {#each SPECIAL_POSISI_OPTIONS as option (option)}
            <SelectItem value={option}>{option}</SelectItem>
          {/each}
        </SelectContent>
      </Select>
      <Select type="single" bind:value={specialKondisiInput}>
        <SelectTrigger class="w-full">
          {specialKondisiInput}
        </SelectTrigger>
        <SelectContent>
          {#each SPECIAL_KONDISI_OPTIONS as option (option)}
            <SelectItem value={option}>{option}</SelectItem>
          {/each}
        </SelectContent>
      </Select>
    </div>
  {:else}
    <div class="flex items-center gap-2">
      <Select type="single" bind:value={kombinasiPosisiInput}>
        <SelectTrigger class="w-full">
          {kombinasiPosisiInput}
        </SelectTrigger>
        <SelectContent>
          {#each KOMBINASI_POSISI_OPTIONS as option (option)}
            <SelectItem value={option}>{option}</SelectItem>
          {/each}
        </SelectContent>
      </Select>
      <Select type="single" bind:value={kombinasiJenisInput}>
        <SelectTrigger class="w-full">
          {kombinasiJenisInput}
        </SelectTrigger>
        <SelectContent>
          {#each KOMBINASI_JENIS_OPTIONS as option (option)}
            <SelectItem value={option}>{option}</SelectItem>
          {/each}
        </SelectContent>
      </Select>
    </div>
  {/if}
  {#if type !== "Umum"}
    <Input
      type="text"
      inputmode="numeric"
      pattern="[0-9]*"
      placeholder="Bet"
      class="text-center"
      value={betInput}
      oninput={handleBetInput}
    />
  {/if}
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

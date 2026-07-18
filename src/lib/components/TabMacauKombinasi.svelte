<script lang="ts">
  import Formmacaukombinasi from "$lib/components/Formmacaukombinasi.svelte";

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
</script>

<div class="mt-4">
  <Formmacaukombinasi bind:bets bind:minBetAlertOpen bind:minBetRequired />
</div>

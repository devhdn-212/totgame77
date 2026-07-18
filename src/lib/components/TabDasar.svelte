<script lang="ts">
  import Formdasar from "$lib/components/Formdasar.svelte";

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
  <Formdasar bind:bets bind:minBetAlertOpen bind:minBetRequired />
</div>

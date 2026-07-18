<script lang="ts">
  import Formshio from "$lib/components/Formshio.svelte";

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
</script>

<div class="mt-4">
  <Formshio bind:bets bind:minBetAlertOpen bind:minBetRequired />
</div>

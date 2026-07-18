<script lang="ts">
  import Form5050 from "$lib/components/Form5050.svelte";
  import { Tabs, TabsList, TabsTrigger, TabsContent } from "$lib/components/ui/tabs";

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

  const fiftyFiftySubTabs = ["Umum", "Special", "Kombinasi"] as const;
  let activeSubTab = $state<(typeof fiftyFiftySubTabs)[number]>(fiftyFiftySubTabs[0]);
</script>

<Tabs bind:value={activeSubTab} class="mt-4">
  <div class="-mx-4 scrollbar-none overflow-x-auto px-4">
    <TabsList class="w-max">
      {#each fiftyFiftySubTabs as tab (tab)}
        <TabsTrigger
          value={tab}
          class="cursor-pointer whitespace-nowrap data-active:border-orange-300! data-active:bg-orange-100! data-active:text-orange-700!"
        >
          {tab}
        </TabsTrigger>
      {/each}
    </TabsList>
  </div>
  {#each fiftyFiftySubTabs as tab (tab)}
    <TabsContent value={tab} class="text-sm">
      <Form5050 type={tab} bind:bets bind:minBetAlertOpen bind:minBetRequired />
    </TabsContent>
  {/each}
</Tabs>

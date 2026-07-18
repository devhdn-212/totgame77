<script lang="ts">
  import Formcolok from "$lib/components/Formcolok.svelte";
  import { Tabs, TabsList, TabsTrigger, TabsContent } from "$lib/components/ui/tabs";

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
    minBetRequired = $bindable(500),
  }: {
    bets: BetEntry[];
    minBetAlertOpen: boolean;
    minBetRequired: number;
  } = $props();

  const colokSubTabs = ["Bebas", "Macau", "Naga", "JITU"] as const;
  let activeSubTab = $state<(typeof colokSubTabs)[number]>(colokSubTabs[0]);
</script>

<Tabs bind:value={activeSubTab} class="mt-4">
  <div class="-mx-4 scrollbar-none overflow-x-auto px-4">
    <TabsList class="w-max">
      {#each colokSubTabs as tab (tab)}
        <TabsTrigger
          value={tab}
          class="cursor-pointer whitespace-nowrap data-active:border-orange-300! data-active:bg-orange-100! data-active:text-orange-700!"
        >
          {tab}
        </TabsTrigger>
      {/each}
    </TabsList>
  </div>
  {#each colokSubTabs as tab (tab)}
    <TabsContent value={tab} class="text-sm">
      <Formcolok type={tab} bind:bets bind:minBetAlertOpen bind:minBetRequired />
    </TabsContent>
  {/each}
</Tabs>

<script lang="ts">
  import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
  } from "$lib/components/ui/dialog";

  let {
    open = $bindable(false),
    activeBetType,
  }: {
    open: boolean;
    activeBetType: string;
  } = $props();

  type InfoTable = { columns: string[]; rows: { label: string; values: string[] }[] };

  const INFO_TABLES: Record<string, InfoTable | null> = {
    "4D/3D/2D": {
      columns: ["4D", "3D", "3DD", "2D", "2DD", "2DT"],
      rows: [
        { label: "Hadiah Diskon", values: ["3000X", "400X", "355X", "70X", "65X", "65X"] },
        { label: "Hadiah Full", values: ["10000X", "950X", "945X", "95X", "90X", "90X"] },
        { label: "Hadiah BB Match", values: ["4000X", "400X", "400X", "70X", "65X", "65X"] },
        { label: "Hadiah BB", values: ["200X", "110X", "100X", "20X", "15X", "15X"] },
        { label: "Diskon", values: ["66%", "59%", "56%", "29%", "26%", "26%"] },
      ],
    },
    Colok: null,
    "5050": null,
    "Macau/Kombinasi": null,
    Dasar: null,
    Shio: null,
  };

  let table = $derived(INFO_TABLES[activeBetType] ?? null);
</script>

<Dialog bind:open>
  <DialogContent class="max-w-lg">
    <DialogHeader>
      <DialogTitle>Informasi</DialogTitle>
      <DialogDescription>Tabel hadiah dan diskon untuk permainan {activeBetType}</DialogDescription>
    </DialogHeader>

    {#if table}
      <div class="overflow-x-auto">
        <table class="w-full min-w-max border-collapse text-xs">
          <thead>
            <tr>
              <th class="border-b p-2 text-left font-medium"></th>
              {#each table.columns as col (col)}
                <th class="border-b p-2 text-center font-medium">{col}</th>
              {/each}
            </tr>
          </thead>
          <tbody>
            {#each table.rows as row (row.label)}
              <tr>
                <td class="text-muted-foreground border-b p-2 whitespace-nowrap">{row.label}</td>
                {#each row.values as value, i (table.columns[i])}
                  <td class="border-b p-2 text-center">{value}</td>
                {/each}
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {:else}
      <p class="text-muted-foreground py-4 text-center text-sm">
        Belum ada informasi untuk permainan {activeBetType}.
      </p>
    {/if}
  </DialogContent>
</Dialog>

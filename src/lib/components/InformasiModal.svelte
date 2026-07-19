<script lang="ts">
  import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
  } from "$lib/components/ui/dialog";
  import { formatIDR } from "$lib/utils";

  let {
    open = $bindable(false),
    activeBetType,
  }: {
    open: boolean;
    activeBetType: string;
  } = $props();

  type InfoContent =
    | { kind: "table"; columns: string[]; rows: { label: string; values: string[] }[] }
    | { kind: "sections"; sections: { title: string; items: { label: string; value: string }[] }[] };

  const INFO_CONTENT: Record<string, InfoContent | null> = {
    "4D/3D/2D": {
      kind: "table",
      columns: ["4D", "3D", "3DD", "2D", "2DD", "2DT"],
      rows: [
        { label: "Hadiah Diskon", values: ["3000X", "400X", "355X", "70X", "65X", "65X"] },
        { label: "Hadiah Full", values: ["10000X", "950X", "945X", "95X", "90X", "90X"] },
        { label: "Hadiah BB Match", values: ["4000X", "400X", "400X", "70X", "65X", "65X"] },
        { label: "Hadiah BB", values: ["200X", "110X", "100X", "20X", "15X", "15X"] },
        { label: "Diskon", values: ["66%", "59%", "56%", "29%", "26%", "26%"] },
        {
          label: "Minimal Bet",
          values: [500, 500, 500, 500, 500, 500].map((v) => formatIDR(v)),
        },
        {
          label: "Maksimal Bet Disc",
          values: [50000, 200000, 200000, 1000000, 1000000, 1000000].map((v) => formatIDR(v)),
        },
        {
          label: "Maksimal Bet Full",
          values: [50000, 200000, 200000, 1000000, 1000000, 1000000].map((v) => formatIDR(v)),
        },
        {
          label: "Maksimal Bet BB",
          values: [50000, 200000, 200000, 1000000, 1000000, 1000000].map((v) => formatIDR(v)),
        },
        {
          label: "Limit Total",
          values: [50000, 200000, 200000, 1000000, 1000000, 1000000].map((v) => formatIDR(v)),
        },
        {
          label: "Limit Global",
          values: [50000, 200000, 200000, 1000000, 1000000, 1000000].map((v) => formatIDR(v)),
        },
      ],
    },
    Colok: {
      kind: "sections",
      sections: [
        {
          title: "Colok Bebas",
          items: [
            { label: "Hadiah", value: "1.53x" },
            { label: "Disc", value: "6%" },
            { label: "Minimal Bet", value: formatIDR(1000) },
            { label: "Maksimal Bet", value: formatIDR(3000000) },
            { label: "Limit Total", value: formatIDR(3000000) },
            { label: "Limit Global", value: formatIDR(3000000) },
          ],
        },
        {
          title: "Colok Macau",
          items: [
            { label: "Hadiah 2 Digit", value: "6.5x" },
            { label: "Hadiah 3 Digit", value: "11x" },
            { label: "Hadiah 4 Digit", value: "18x" },
            { label: "Disc", value: "10%" },
            { label: "Minimal Bet", value: formatIDR(1000) },
            { label: "Maksimal Bet", value: formatIDR(3000000) },
            { label: "Limit Total", value: formatIDR(3000000) },
            { label: "Limit Global", value: formatIDR(3000000) },
          ],
        },
        {
          title: "Colok Naga",
          items: [
            { label: "Hadiah 3 Digit", value: "25x" },
            { label: "Hadiah 4 Digit", value: "37x" },
            { label: "Disc", value: "10%" },
            { label: "Minimal Bet", value: formatIDR(1000) },
            { label: "Maksimal Bet", value: formatIDR(3000000) },
            { label: "Limit Total", value: formatIDR(3000000) },
            { label: "Limit Global", value: formatIDR(3000000) },
          ],
        },
        {
          title: "Colok Jitu",
          items: [
            { label: "Hadiah As", value: "8x" },
            { label: "Hadiah Kop", value: "8x" },
            { label: "Hadiah Kepala", value: "8x" },
            { label: "Hadiah Ekor", value: "8x" },
            { label: "Disc", value: "5%" },
            { label: "Minimal Bet", value: formatIDR(1000) },
            { label: "Maksimal Bet", value: formatIDR(3000000) },
            { label: "Limit Total", value: formatIDR(3000000) },
            { label: "Limit Global", value: formatIDR(3000000) },
          ],
        },
      ],
    },
    "5050": {
      kind: "sections",
      sections: [
        {
          title: "Umum",
          items: [
            { label: "Hadiah", value: "1x" },
            { label: "Disc", value: "2%" },
            { label: "Kei", value: "-1.5%" },
            { label: "Minimal Bet", value: formatIDR(1000) },
            { label: "Maksimal Bet", value: formatIDR(1000000) },
            { label: "Limit Total", value: formatIDR(1000000) },
            { label: "Limit Global", value: formatIDR(1000000) },
          ],
        },
        {
          title: "Special",
          items: [
            { label: "Hadiah", value: "1x" },
            { label: "Disc", value: "1.5%" },
            { label: "Kei", value: "-2.5%" },
            { label: "Minimal Bet", value: formatIDR(1000) },
            { label: "Maksimal Bet", value: formatIDR(1000000) },
            { label: "Limit Total", value: formatIDR(1000000) },
            { label: "Limit Global", value: formatIDR(1000000) },
          ],
        },
        {
          title: "Kombinasi",
          items: [
            { label: "Hadiah", value: "1x" },
            { label: "Disc", value: "1.5%" },
            { label: "Kei", value: "-2.5%" },
            { label: "Minimal Bet", value: formatIDR(1000) },
            { label: "Maksimal Bet", value: formatIDR(1000000) },
            { label: "Limit Total", value: formatIDR(1000000) },
            { label: "Limit Global", value: formatIDR(1000000) },
          ],
        },
      ],
    },
    "Macau/Kombinasi": {
      kind: "sections",
      sections: [
        {
          title: "Macau/Kombinasi",
          items: [
            { label: "Hadiah", value: "2.7x" },
            { label: "Disc", value: "5%" },
            { label: "Minimal Bet", value: formatIDR(1000) },
            { label: "Maksimal Bet", value: formatIDR(1000000) },
            { label: "Limit Total", value: formatIDR(1000000) },
            { label: "Limit Global", value: formatIDR(1000000) },
          ],
        },
      ],
    },
    Dasar: {
      kind: "sections",
      sections: [
        {
          title: "Dasar",
          items: [
            { label: "Hadiah", value: "1x" },
            { label: "Disc Besar", value: "0%" },
            { label: "Disc Kecil", value: "0%" },
            { label: "Disc Genap", value: "0%" },
            { label: "Disc Ganjil", value: "0%" },
            { label: "Kei Besar", value: "-25%" },
            { label: "Kei Kecil", value: "0%" },
            { label: "Kei Genap", value: "-25%" },
            { label: "Kei Ganjil", value: "0%" },
            { label: "Minimal Bet", value: formatIDR(1000) },
            { label: "Maksimal Bet", value: formatIDR(1000000) },
            { label: "Limit Total", value: formatIDR(1000000) },
            { label: "Limit Global", value: formatIDR(1000000) },
          ],
        },
      ],
    },
    Shio: {
      kind: "sections",
      sections: [
        {
          title: "Shio",
          items: [
            { label: "Hadiah", value: "9x" },
            { label: "Disc", value: "8%" },
            { label: "Minimal Bet", value: formatIDR(1000) },
            { label: "Maksimal Bet", value: formatIDR(1000000) },
            { label: "Limit Total", value: formatIDR(1000000) },
            { label: "Limit Global", value: formatIDR(1000000) },
          ],
        },
      ],
    },
  };

  const INFO_DESCRIPTION: Record<string, string> = {
    "4D/3D/2D":
      "Tabel informasi minimal bet, maximal bet, hadiah, diskon, limit total dan limit global untuk permainan 4D/3D/2D",
    Colok:
      "Tabel informasi minimal bet, maximal bet, hadiah, diskon, limit total dan limit global untuk permainan Colok",
    "5050":
      "Tabel informasi minimal bet, maximal bet, hadiah, diskon, limit total dan limit global untuk permainan 5050",
    "Macau/Kombinasi":
      "Tabel informasi minimal bet, maximal bet, hadiah, diskon, limit total dan limit global untuk permainan Macau/Kombinasi",
    Dasar:
      "Tabel informasi minimal bet, maximal bet, hadiah, diskon, limit total dan limit global untuk permainan Dasar",
    Shio: "Tabel informasi minimal bet, maximal bet, hadiah, diskon, limit total dan limit global untuk permainan Shio",
  };

  let content = $derived(INFO_CONTENT[activeBetType] ?? null);
  let description = $derived(
    INFO_DESCRIPTION[activeBetType] ?? `Tabel hadiah dan diskon untuk permainan ${activeBetType}`,
  );
</script>

<Dialog bind:open>
  <DialogContent class="max-w-lg">
    <DialogHeader>
      <DialogTitle>Informasi</DialogTitle>
      <DialogDescription>{description}</DialogDescription>
    </DialogHeader>

    {#if content?.kind === "table"}
      <div class="overflow-x-auto">
        <table class="w-full min-w-max border-collapse text-xs">
          <thead>
            <tr>
              <th class="border-b p-2 text-left font-medium"></th>
              {#each content.columns as col (col)}
                <th class="border-b p-2 text-right font-medium">{col}</th>
              {/each}
            </tr>
          </thead>
          <tbody>
            {#each content.rows as row (row.label)}
              <tr>
                <td class="text-muted-foreground border-b p-2 whitespace-nowrap">{row.label}</td>
                {#each row.values as value, i (content.columns[i])}
                  <td class="border-b p-2 text-right">{value}</td>
                {/each}
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    {:else if content?.kind === "sections"}
      <div class="flex max-h-96 flex-col gap-3 overflow-y-auto">
        {#each content.sections as section (section.title)}
          <div class="rounded-lg border p-3">
            <p class="mb-2 text-sm font-medium">{section.title}</p>
            <div class="flex flex-col gap-1">
              {#each section.items as item (item.label)}
                <div class="flex items-center justify-between text-xs">
                  <span class="text-muted-foreground">{item.label}</span>
                  <span class="font-medium">{item.value}</span>
                </div>
              {/each}
            </div>
          </div>
        {/each}
      </div>
    {:else}
      <p class="text-muted-foreground py-4 text-center text-sm">
        Belum ada informasi untuk permainan {activeBetType}.
      </p>
    {/if}
  </DialogContent>
</Dialog>

<script lang="ts">
  import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
  } from "$lib/components/ui/dialog";
  import { Button } from "$lib/components/ui/button";
  import { Input } from "$lib/components/ui/input";
  import { Select, SelectTrigger, SelectContent, SelectItem } from "$lib/components/ui/select";
  import Search from "@lucide/svelte/icons/search";
  import LoaderCircle from "@lucide/svelte/icons/loader-circle";

  let { open = $bindable(false) }: { open: boolean } = $props();

  const TIPE_OPTIONS = ["ALL", "4D", "3D", "2D"];

  type BukuMimpiRow = { id: string; type: string; name: string; nomor: string };

  let tipe = $state("ALL");
  let namaQuery = $state("");
  let results = $state<BukuMimpiRow[]>([]);
  let loading = $state(false);
  let errorMessage = $state("");

  // The same text input both filters what's already loaded (instant, as you
  // type) and becomes the "nama" sent to the server when Search is clicked.
  let filteredResults = $derived(
    results.filter((row) => {
      if (!namaQuery) return true;
      const q = namaQuery.toLowerCase();
      return row.name.toLowerCase().includes(q) || row.nomor.toLowerCase().includes(q);
    }),
  );

  async function fetchData(tipeParam: string, nama: string) {
    loading = true;
    errorMessage = "";
    try {
      const res = await fetch("/api/bukumimpi", {
        method: "POST",
        headers: { "content-type": "application/json" },
        body: JSON.stringify({ tipe: tipeParam, nama }),
      });
      const body = await res.json();

      if (!res.ok) {
        throw new Error(body?.message ?? `request failed (${res.status})`);
      }

      results = (body.record ?? []).map(
        (row: { bukumimpi_type: string; bukumimpi_name: string; bukumimpi_nomor: string }) => ({
          id: crypto.randomUUID(),
          type: row.bukumimpi_type.trim(),
          name: row.bukumimpi_name,
          nomor: row.bukumimpi_nomor,
        }),
      );
    } catch (err) {
      console.error(err);
      results = [];
      errorMessage = "Gagal memuat data buku mimpi.";
    } finally {
      loading = false;
    }
  }

  function handleSearchClick() {
    fetchData(tipe === "ALL" ? "" : tipe, namaQuery);
  }

  $effect(() => {
    if (open) {
      tipe = "ALL";
      namaQuery = "";
      fetchData("", "");
    }
  });
</script>

<Dialog bind:open>
  <DialogContent>
    <DialogHeader>
      <DialogTitle>Buku Mimpi</DialogTitle>
      <DialogDescription>Cari arti mimpi dan nomor terkait.</DialogDescription>
    </DialogHeader>

    <div class="-mx-6 max-h-[70vh] overflow-y-auto px-6">
      <div class="bg-background sticky top-0 z-10 flex flex-col gap-2 pb-3">
        <div class="flex gap-2">
          <Select type="single" bind:value={tipe}>
            <SelectTrigger class="w-28 shrink-0">
              {tipe}
            </SelectTrigger>
            <SelectContent>
              {#each TIPE_OPTIONS as option (option)}
                <SelectItem value={option}>{option}</SelectItem>
              {/each}
            </SelectContent>
          </Select>
          <Input
            type="text"
            placeholder="Cari mimpi..."
            bind:value={namaQuery}
            onkeydown={(e) => e.key === "Enter" && handleSearchClick()}
          />
        </div>
        <Button
          class="w-full cursor-pointer bg-blue-100 text-blue-700 hover:bg-blue-200"
          onclick={handleSearchClick}
        >
          <Search />
          Search
        </Button>
      </div>

      {#if loading}
        <div class="flex flex-col items-center gap-2 py-8">
          <LoaderCircle class="animate-spin" />
          <p class="text-muted-foreground text-xs">Memuat data...</p>
        </div>
      {:else if errorMessage}
        <p class="text-destructive py-4 text-center text-sm">{errorMessage}</p>
      {:else if filteredResults.length === 0}
        <p class="text-muted-foreground py-4 text-center text-sm">Tidak ada hasil.</p>
      {:else}
        <div class="flex flex-col gap-2 pt-1 pb-2">
          {#each filteredResults as row (row.id)}
            <div class="rounded-lg border p-3">
              <p class="text-sm font-medium">{row.name}</p>
              <p class="text-muted-foreground text-xs">{row.type} · {row.nomor}</p>
            </div>
          {/each}
        </div>
      {/if}
    </div>
  </DialogContent>
</Dialog>

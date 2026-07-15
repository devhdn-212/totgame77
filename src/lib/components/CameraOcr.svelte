<script lang="ts">
  import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
    DialogDescription,
  } from "$lib/components/ui/dialog";
  import { Button } from "$lib/components/ui/button";
  import { formatIDR } from "$lib/utils";
  import Camera from "@lucide/svelte/icons/camera";
  import LoaderCircle from "@lucide/svelte/icons/loader-circle";
  import Plus from "@lucide/svelte/icons/plus";
  import Trash2 from "@lucide/svelte/icons/trash-2";
  import { createWorker } from "tesseract.js";

  type BetEntry = {
    id: string;
    type: string;
    number: string;
    bet: string;
    deletable?: boolean;
  };

  type ParsedRow = { id: string; number: string; bet: string };

  let {
    open = $bindable(false),
    bets = $bindable(),
    minBetAlertOpen = $bindable(false),
  }: {
    open: boolean;
    bets: BetEntry[];
    minBetAlertOpen: boolean;
  } = $props();

  const MIN_BET = 500;

  let fileInput: HTMLInputElement | null = $state(null);
  let imageUrl = $state("");
  let status = $state<"idle" | "recognizing" | "done" | "error">("idle");
  let progress = $state(0);
  let rawText = $state("");
  let parsedRows = $state<ParsedRow[]>([]);
  let errorMessage = $state("");

  function reset() {
    imageUrl = "";
    status = "idle";
    progress = 0;
    rawText = "";
    parsedRows = [];
    errorMessage = "";
  }

  function handleTriggerCapture() {
    reset();
    fileInput?.click();
  }

  function deriveTypes(number: string): { type: string; number: string }[] {
    if (number.length >= 4) {
      const base = number.slice(-4);
      return [
        { type: "4D", number: base },
        { type: "3D", number: base.slice(1) },
        { type: "2D", number: base.slice(2) },
      ];
    }
    if (number.length === 3) {
      return [
        { type: "3D", number },
        { type: "2D", number: number.slice(1) },
      ];
    }
    return [{ type: "2D", number }];
  }

  function parseRawText(text: string): ParsedRow[] {
    const rows: ParsedRow[] = [];
    const matches = text.matchAll(/(\d{2,4})\D+(\d{2,})/g);
    for (const match of matches) {
      rows.push({ id: crypto.randomUUID(), number: match[1], bet: match[2] });
    }
    return rows;
  }

  async function handleFileSelected(e: Event) {
    const target = e.currentTarget as HTMLInputElement;
    const file = target.files?.[0];
    target.value = "";
    if (!file) return;

    reset();
    imageUrl = URL.createObjectURL(file);
    status = "recognizing";
    errorMessage = "";

    try {
      const worker = await createWorker("eng", 1, {
        logger: (m) => {
          if (m.status === "recognizing text") progress = Math.round(m.progress * 100);
        },
      });
      const { data } = await worker.recognize(file);
      await worker.terminate();

      rawText = data.text;
      parsedRows = parseRawText(data.text);
      status = "done";
      if (parsedRows.length === 0) {
        errorMessage = "Tidak ada nomor yang terbaca. Coba foto ulang dengan tulisan lebih jelas.";
      }
    } catch (err) {
      status = "error";
      errorMessage = "Gagal membaca gambar. Coba lagi.";
    }
  }

  function handleRemoveRow(id: string) {
    parsedRows = parsedRows.filter((row) => row.id !== id);
  }

  function classifyByLength(number: string) {
    if (number.length >= 4) return "4D";
    if (number.length === 3) return "3D";
    return "2D";
  }

  function handleConfirmAdd() {
    const invalidBet = parsedRows.find((row) => Number(row.bet) < MIN_BET);
    if (invalidBet) {
      minBetAlertOpen = true;
      return;
    }

    // A single detected number/bet pair auto-expands into 4D/3D/2D.
    // Multiple detected pairs are each their own distinct number, added as-is.
    const newEntries: BetEntry[] =
      parsedRows.length === 1
        ? deriveTypes(parsedRows[0].number).map((t) => ({
            id: crypto.randomUUID(),
            type: t.type,
            number: t.number,
            bet: parsedRows[0].bet,
          }))
        : parsedRows.map((row) => ({
            id: crypto.randomUUID(),
            type: classifyByLength(row.number),
            number: row.number,
            bet: row.bet,
          }));
    bets = [...newEntries, ...bets];

    open = false;
    reset();
  }
</script>

<input
  bind:this={fileInput}
  type="file"
  accept="image/*"
  capture="environment"
  class="hidden"
  onchange={handleFileSelected}
/>

<Dialog bind:open>
  <DialogContent>
    <DialogHeader>
      <DialogTitle>Scan Nomor</DialogTitle>
      <DialogDescription>
        Foto tulisan nomor & bet di kertas, lalu review sebelum ditambahkan ke list.
      </DialogDescription>
    </DialogHeader>

    {#if status === "idle"}
      <Button class="w-full cursor-pointer" onclick={handleTriggerCapture}>
        <Camera />
        Ambil Foto
      </Button>
    {:else}
      {#if imageUrl}
        <img src={imageUrl} alt="Hasil foto" class="max-h-48 w-full rounded-lg border object-contain" />
      {/if}

      {#if status === "recognizing"}
        <div class="flex flex-col items-center gap-2 py-4">
          <LoaderCircle class="animate-spin" />
          <p class="text-muted-foreground text-xs">Membaca gambar... {progress}%</p>
        </div>
      {:else if status === "error"}
        <p class="text-destructive text-xs">{errorMessage}</p>
        <Button variant="outline" class="w-full cursor-pointer" onclick={handleTriggerCapture}>
          Coba Lagi
        </Button>
      {:else if status === "done"}
        {#if errorMessage}
          <p class="text-destructive text-xs">{errorMessage}</p>
        {/if}
        {#if parsedRows.length > 0}
          <div class="flex max-h-56 flex-col gap-2 overflow-y-auto">
            {#each parsedRows as row (row.id)}
              <div class="flex items-center justify-between rounded-lg border p-2">
                <div class="text-sm">
                  <span class="font-medium">{row.number}</span>
                  <span class="text-muted-foreground"> - Bet {formatIDR(row.bet)}</span>
                </div>
                <Button
                  size="icon"
                  variant="ghost"
                  class="text-destructive cursor-pointer"
                  aria-label="Hapus"
                  onclick={() => handleRemoveRow(row.id)}
                >
                  <Trash2 />
                </Button>
              </div>
            {/each}
          </div>
        {/if}
        <div class="flex gap-2">
          <Button variant="outline" class="flex-1 cursor-pointer" onclick={handleTriggerCapture}>
            Foto Ulang
          </Button>
          <Button
            class="flex-1 cursor-pointer bg-blue-100 text-blue-700 hover:bg-blue-200"
            disabled={parsedRows.length === 0}
            onclick={handleConfirmAdd}
          >
            <Plus />
            Tambahkan
          </Button>
        </div>
      {/if}
    {/if}
  </DialogContent>
</Dialog>

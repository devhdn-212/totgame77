const { chromium } = require('playwright');
(async () => {
  const browser = await chromium.launch();
  const errors = [];
  const ctx = await browser.newContext({ viewport: { width: 390, height: 844 } });
  const page = await ctx.newPage();
  page.on('console', (msg) => { if (msg.type() === 'error') errors.push(msg.text()); });
  page.on('pageerror', (err) => errors.push('pageerror: ' + err.message));
  await page.goto('http://localhost:5176');
  await page.waitForSelector('text=HONGKONG');

  await page.locator('a:has-text("Camera")').click();
  await page.waitForTimeout(200);
  console.log('DIALOG_OPEN:', await page.locator('text=Scan Nomor').isVisible());
  console.log('AMBIL_FOTO_BTN_VISIBLE:', await page.locator('button:has-text("Ambil Foto")').isVisible());
  console.log('OLD_BUKA_KAMERA_GONE:', !(await page.locator('button:has-text("Buka Kamera")').isVisible().catch(() => false)));
  console.log('GALLERY_LINK_VISIBLE:', await page.locator('text=atau pilih foto dari galeri').isVisible());

  // Verify the camera input has the capture attribute (native camera app trigger)
  const captureAttr = await page.locator('input[type=file]').first().getAttribute('capture');
  console.log('CAMERA_INPUT_CAPTURE_ATTR:', captureAttr);

  // Simulate file selection (what happens after native camera app returns a photo)
  const fileInput = page.locator('input[type=file]').first();
  await fileInput.setInputFiles('C:/Users/PC/AppData/Local/Temp/claude/c--Users-PC-Documents-PROGRAMMING-2026-golang-project-2026-gotomodern-client-gocu-client-web2026/00f62140-a99b-4ee2-871f-358a01daeae1/scratchpad/test-ocr-image2.png');

  await page.waitForFunction(() => {
    return document.body.innerText.includes('Tambahkan') || document.body.innerText.includes('Tidak ada nomor') || document.body.innerText.includes('Gagal membaca');
  }, null, { timeout: 90000 });

  console.log('IMAGE_PREVIEW_VISIBLE:', await page.locator('img[alt="Hasil foto"]').isVisible());
  const rowCount = await page.locator('div.flex.max-h-56 > div').count();
  console.log('PARSED_ROW_COUNT:', rowCount);
  const rowTexts = await page.locator('div.flex.max-h-56 > div').allInnerTexts();
  console.log('PARSED_ROWS:', JSON.stringify(rowTexts));

  await page.screenshot({ path: 'C:/Users/PC/AppData/Local/Temp/claude/c--Users-PC-Documents-PROGRAMMING-2026-golang-project-2026-gotomodern-client-gocu-client-web2026/00f62140-a99b-4ee2-871f-358a01daeae1/scratchpad/camera-native-result.png' });

  // Confirm add works
  await page.locator('button:has-text("Tambahkan")').click();
  await page.waitForTimeout(300);
  console.log('DIALOG_CLOSED:', !(await page.locator('text=Scan Nomor').isVisible().catch(() => false)));
  await page.screenshot({ path: 'C:/Users/PC/AppData/Local/Temp/claude/c--Users-PC-Documents-PROGRAMMING-2026-golang-project-2026-gotomodern-client-gocu-client-web2026/00f62140-a99b-4ee2-871f-358a01daeae1/scratchpad/camera-native-added.png' });

  console.log('ERRORS:', JSON.stringify(errors));
  await browser.close();
})();

import path from "path";
import puppeteer from "puppeteer";

export async function openPDF(pdfPath: string) {
  const browser = await puppeteer.launch({
    args: ["--no-sandbox", "--font-render-hinting=none"],
    headless: false,
  });
  const page = await browser.newPage();
  const absolutePath = path.resolve(pdfPath);
  const fileUrl = `http://localhost:5555/resume_1YOE.pdf`;
  await page.goto(fileUrl);
  setTimeout(async () => {
    const pageCount = await page.evaluate(
      () => document.querySelectorAll(".page").length
    );
    for (let pageNumber = 1; pageNumber <= pageCount; pageNumber++) {
      await page.goto(`${fileUrl}#page=${pageNumber}`);
      await page.waitForSelector(".page");

      const screenshotPath = path.resolve(`page_${pageNumber}.png`);
      await page.screenshot({ path: screenshotPath });
      console.log(`Screenshot captured for page ${pageNumber}`);
    }
  }, 3000);
  setTimeout(async () => {
    await browser.close();
  }, 30000);
}

import { expect, test } from "@playwright/test";

test("Feature flags page renders", async ({ page }) => {
  await page.goto("/feature-flags");
  await expect(page.getByRole("heading", { name: "Feature Flags" })).toBeVisible();
});

test("WebRTC gray simulator computes percent", async ({ page }) => {
  await page.goto("/webrtc/gray-sim");
  await expect(page.getByRole("heading", { name: "WebRTC 灰度模拟器" })).toBeVisible();
  await expect(page.locator("text=/命中：/")).toBeVisible();
});

import { test, expect } from "@playwright/test";

test.beforeEach(async ({ page }) => {
	await page.goto("http://localhost:3173/");
});

test.describe("mapbox", () => {
	test("has copyright text", async ({ page }) => {
		// Expect a title "to contain" a substring.
		const text = page.getByText("© Mapbox");
		await expect(text).toBeInViewport();
	});

	test("mapbox should render", async ({ page }) => {
		const mapbox = page.locator(".mapboxgl-map");
		await mapbox.waitFor();
		await expect(mapbox).toBeInViewport();
	});

	test("has markers", async ({ page }) => {
		// Expect a title "to contain" a substring.
		const mapbox = page.locator(".mapboxgl-map");
		await mapbox.waitFor();
		await expect(mapbox).toBeInViewport();

		const markerLoc = mapbox.locator("[aria-label='Map marker']");
		await expect(markerLoc).toHaveCount(2);
	});

	test("drawer should open", async ({ page }) => {
		const icon = page.getByTestId("search_btn");
		await icon.click();
		await expect(icon).toBeVisible();

		const drawer = page.locator(".drawer");
		await expect(drawer).toBeVisible();
	});
});

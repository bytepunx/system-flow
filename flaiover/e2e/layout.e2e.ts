// S-0458: the nav rail. From Tailwind's md width the rail is docked at the left of every page,
// beside the main area; below md it is a drawer that the menu button at the start of the top bar
// opens, that choosing a page or pressing Escape closes, and that gives focus back to the button.
// At 390px and at 1440px no page scrolls sideways, the drawer open or closed. The pages are ones
// that render signed in with no project connected, so no fake flai is needed. S-0484: Charts in the
// rail expands to list the charts, on arrival at a chart or when pressed, and choosing one opens it,
// marks it current, and closes the drawer.
import { expect, test, type Page } from '@playwright/test';
import { readFileSync } from 'node:fs';
import { join } from 'node:path';

const PAGES = [
	'/',
	'/board',
	'/settings',
	'/license',
	'/docs/design/system/README.md',
	'/charts/cycle-time'
];

// what playwright.config.ts made for the run and gave the server
const dir = process.env.FLAIOVER_E2E_DIR ?? '';
const token = () => readFileSync(join(dir, 'token'), 'utf8').trim();

test.beforeEach(async ({ context, baseURL }) => {
	await context.addCookies([{ name: 'flaiover_session', value: token(), url: baseURL! }]);
});

/** How far the page reaches past the window sideways: 0 when it never scrolls sideways. */
const overflow = (page: Page) =>
	page.evaluate(() => document.documentElement.scrollWidth - window.innerWidth);

async function arrive(page: Page, path: string) {
	await page.goto(path);
	await expect(page.locator('[data-top-bar]')).toBeVisible();
}

test.describe('at 1440px', () => {
	test.use({ viewport: { width: 1440, height: 900 } });

	for (const path of PAGES)
		test(`${path}: the rail is docked at the left of the main area, and nothing scrolls sideways`, async ({
			page
		}) => {
			await arrive(page, path);
			const docked = page.locator('[data-nav-docked]');
			await expect(docked).toBeVisible();
			await expect(docked.locator('nav[aria-label="Site"][data-nav-rail]')).toBeVisible();
			await expect(page.locator('[data-nav-toggle]')).toBeHidden();

			const rail = (await docked.boundingBox())!;
			const main = (await page.locator('main').boundingBox())!;
			expect(rail.x).toBeCloseTo(0, 0);
			expect(rail.x + rail.width).toBeLessThanOrEqual(main.x);
			expect(await overflow(page)).toBeLessThanOrEqual(0);
		});

	test('arriving on a chart expands Charts in the rail and marks the chart current', async ({
		page
	}) => {
		await arrive(page, '/charts/cycle-time');
		const docked = page.locator('[data-nav-docked]');
		const charts = docked.locator('[data-rail-page=charts]');
		await expect(charts).toHaveAttribute('aria-expanded', 'true');
		const shown = docked.locator('[data-rail-chart=cycle-time]');
		await expect(shown).toHaveAttribute('aria-current', 'page');
	});

	test('pressing Charts lists the charts, and choosing one opens it and marks it current', async ({
		page
	}) => {
		await arrive(page, '/board');
		const docked = page.locator('[data-nav-docked]');
		const charts = docked.locator('[data-rail-page=charts]');
		await expect(charts).toHaveAttribute('aria-expanded', 'false');
		await expect(docked.locator('[data-rail-chart]')).toHaveCount(0);
		await charts.click();
		await expect(charts).toHaveAttribute('aria-expanded', 'true');
		await docked.locator('[data-rail-chart=throughput]').click();

		await expect(page).toHaveURL(/\/charts\/throughput(\?|$)/);
		const shown = docked.locator('[data-rail-chart=throughput]');
		await expect(shown).toHaveAttribute('aria-current', 'page');
	});
});

test.describe('at 390px', () => {
	test.use({ viewport: { width: 390, height: 844 } });

	for (const path of PAGES)
		test(`${path}: the rail is a drawer, and nothing scrolls sideways, open or closed`, async ({
			page
		}) => {
			await arrive(page, path);
			await expect(page.locator('[data-nav-docked]')).toBeHidden();
			const toggle = page.locator('[data-top-bar] > :first-child');
			await expect(toggle).toHaveAttribute('data-nav-toggle', '');
			await expect(toggle).toBeVisible();
			await expect(toggle).toHaveAttribute('aria-expanded', 'false');
			expect(await overflow(page)).toBeLessThanOrEqual(0);

			await toggle.click();
			const drawer = page.locator('[data-nav-drawer]');
			await expect(drawer).toHaveAttribute('open', '');
			await expect(drawer.locator('nav[aria-label="Site"][data-nav-rail]')).toBeVisible();
			await expect(toggle).toHaveAttribute('aria-expanded', 'true');
			expect(await overflow(page)).toBeLessThanOrEqual(0);

			await page.keyboard.press('Escape');
			await expect(drawer).not.toHaveAttribute('open');
			await expect(toggle).toBeFocused();
		});

	test('choosing a page in the drawer goes there, closes it, and gives focus back to the button', async ({
		page
	}) => {
		await arrive(page, '/');
		const toggle = page.locator('[data-nav-toggle]');
		const drawer = page.locator('[data-nav-drawer]');
		await toggle.click();
		await expect(drawer).toHaveAttribute('open', '');

		const flow = drawer.locator('[data-rail-group=flow]');
		if ((await flow.getAttribute('aria-expanded')) !== 'true') await flow.click();
		await drawer.locator('[data-rail-page=board]').click();

		await expect(page).toHaveURL(/\/board(\?|$)/);
		await expect(drawer).not.toHaveAttribute('open');
		await expect(toggle).toBeFocused();
		await expect(toggle).toHaveAttribute('aria-expanded', 'false');

		// opened again, the rail marks the page it went to
		await toggle.click();
		await expect(drawer.locator('[data-rail-page=board]')).toHaveAttribute('aria-current', 'page');
		await page.keyboard.press('Escape');
		await expect(drawer).not.toHaveAttribute('open');
		await expect(toggle).toBeFocused();
	});

	test('choosing a chart in the drawer opens it, closes the drawer, and marks it current', async ({
		page
	}) => {
		await arrive(page, '/board');
		const toggle = page.locator('[data-nav-toggle]');
		const drawer = page.locator('[data-nav-drawer]');
		await toggle.click();
		await expect(drawer).toHaveAttribute('open', '');

		const charts = drawer.locator('[data-rail-page=charts]');
		await expect(charts).toHaveAttribute('aria-expanded', 'false');
		await charts.click();
		await expect(charts).toHaveAttribute('aria-expanded', 'true');
		await expect(drawer).toHaveAttribute('open', '');
		await drawer.locator('[data-rail-chart=throughput]').click();

		await expect(page).toHaveURL(/\/charts\/throughput(\?|$)/);
		await expect(drawer).not.toHaveAttribute('open');
		await expect(toggle).toBeFocused();

		// opened again, the rail lists the charts and marks the one it went to
		await toggle.click();
		await expect(charts).toHaveAttribute('aria-expanded', 'true');
		const shown = drawer.locator('[data-rail-chart=throughput]');
		await expect(shown).toHaveAttribute('aria-current', 'page');
		await page.keyboard.press('Escape');
		await expect(drawer).not.toHaveAttribute('open');
	});
});

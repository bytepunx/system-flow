// S-0557, I-0180: the top bar's project switcher with two long-named projects, one connected and
// one not, so the "n not connected: why?" link shows too. Its select once grew to fit its longest
// option, and the page was wider than a 390px window. At 390px and at 1440px, on / and /board, the
// select shows the current project, ends inside the window, and the page never scrolls sideways.
// The projects are answered in the browser, as board-layout.e2e.ts answers the board: a fake flai
// would make a project known to the server for every test file that runs after it, and the pages
// layout.e2e.ts checks render with none.
import { expect, test, type Page } from '@playwright/test';
import { readFileSync } from 'node:fs';
import { join } from 'node:path';

const PAGES = ['/', '/board'];
const CURRENT = 'e2e-board';

// what playwright.config.ts made for the run and gave the server
const dir = process.env.FLAIOVER_E2E_DIR ?? '';
const token = () => readFileSync(join(dir, 'token'), 'utf8').trim();

/** /api/projects' answer: two projects whose labels are as long as I-0180's, or longer. */
const projects = {
	projects: [
		{
			key: CURRENT,
			name: 'E2E board of a project with a long name',
			connected: true,
			unsigned: ['flai', 'dashboard']
		},
		{
			key: 'e2e-quay',
			name: 'E2E quay of another project with a long name',
			connected: false,
			served: true,
			lastError: 'dial: connection refused'
		}
	]
};
const LABEL = 'E2E board of a project with a long name (unsigned allowed: flai and dashboard)';

/** /api/board's answer: an empty board, so /board renders its lanes. */
const board = {
	wip_limits: {},
	order: [],
	writable: true,
	columns: Object.fromEntries(
		['backlog', 'ready', 'in-progress', 'review', 'done', 'cancelled'].map((s) => [s, []])
	)
};

test.beforeEach(async ({ context, baseURL }) => {
	await context.addCookies([{ name: 'flaiover_session', value: token(), url: baseURL! }]);
	await context.route(/\/api\/projects(\?|$)/, (route) => route.fulfill({ json: projects }));
	await context.route(/\/api\/board(\?|$)/, (route) => route.fulfill({ json: board }));
});

/** How far the page reaches past the window sideways: 0 when it never scrolls sideways. */
const overflow = (page: Page) =>
	page.evaluate(() => document.documentElement.scrollWidth - window.innerWidth);

async function expectFits(page: Page, path: string, width: number) {
	await page.goto(`${path}?project=${CURRENT}`);
	await expect(page.locator('[data-top-bar]')).toBeVisible();
	const select = page.locator('[data-testid="project-switcher"] select');
	await expect(select).toBeVisible();
	await expect(select).toHaveValue(CURRENT);
	await expect(select.locator('option:checked')).toHaveText(LABEL);
	await expect(page.locator('[data-testid="project-switcher-why"]')).toHaveText(
		'1 not connected: why?'
	);
	const box = (await select.boundingBox())!;
	expect(box.x).toBeGreaterThanOrEqual(0);
	expect(box.x + box.width).toBeLessThanOrEqual(width);
	expect(await overflow(page)).toBeLessThanOrEqual(0);
}

test.describe('at 390px', () => {
	test.use({ viewport: { width: 390, height: 844 } });

	for (const path of PAGES)
		test(`${path}: two long-named projects leave the page no wider than the window`, async ({
			page
		}) => {
			await expectFits(page, path, 390);
		});
});

test.describe('at 1440px', () => {
	test.use({ viewport: { width: 1440, height: 900 } });

	for (const path of PAGES)
		test(`${path}: two long-named projects leave the page no wider than the window`, async ({
			page
		}) => {
			await expectFits(page, path, 1440);
		});
});

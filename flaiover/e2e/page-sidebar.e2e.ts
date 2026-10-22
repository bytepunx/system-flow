// S-0461: the page sidebar on the review page and the documents page, at a phone's width and a
// desktop's. Below Tailwind's md width the sidebar is a row above the page's content, its entries
// listed in it, and no page scrolls sideways; from md it is a column at the left of the content,
// and still no page scrolls sideways. A file's entry on the review page opens its panel and brings
// it into view, and a heading's entry on the documents page brings its heading to the top. The
// charts page has no sidebar since S-0484: the nav rail lists the charts, and layout.e2e.ts checks
// it. Each page's own requests are answered in the browser with enough to lay it out as a real
// project would: long paths, a diff with a long line, long headings, a wide table, and a long line
// of code. A fake flai would make a project known to the server for every test file that runs
// after it, and the pages layout.e2e.ts checks render with none.
import { expect, test, type Page, type Route } from '@playwright/test';
import { readFileSync } from 'node:fs';
import { join } from 'node:path';

// what playwright.config.ts made for the run and gave the server
const dir = process.env.FLAIOVER_E2E_DIR ?? '';
const token = () => readFileSync(join(dir, 'token'), 'utf8').trim();

const STORY = 'S-9001';
const NOW = '2026-10-21T09:00:00Z';

/** A line far wider than any phone, as a diff or a code block holds one. */
const LONG_LINE =
	'export const aConfigurationValueWithAVeryLongName = computeSomethingFrom(' +
	Array.from({ length: 24 }, (_, n) => `argumentNumber${n}`).join(', ') +
	');';

const patch = (lines: number) =>
	[
		`@@ -1,${lines} +1,${lines + 1} @@`,
		' // the start of the file',
		`-${LONG_LINE.replace('computeSomethingFrom', 'computeSomethingElseFrom')}`,
		`+${LONG_LINE}`,
		...Array.from({ length: lines }, (_, n) => ` const line${n} = ${n};`)
	].join('\n');

/** The story's changed files: a long path among them, and a patch with a long line in each. */
const FILES = [
	'flaiover/src/lib/components/PageSidebar.svelte',
	'flaiover/src/routes/review/[id]/a-deeply/nested/directory/structure/that/goes/on/AndOnWithAVeryLongComponentFileNameIndeed.svelte',
	'design/system/flaiover-dashboard.md',
	'flaiover/src/routes/charts/[kind]/+page.svelte',
	'flai/internal/server/handlers_for_the_review_page_and_its_sidebar_with_a_long_name_test.go'
];
const fileId = (path: string) => `review-file-${path.replace(/[^A-Za-z0-9_-]+/g, '-')}`;

const diff = {
	branch: `story/${STORY}`,
	base: 'main',
	commits: 4,
	additions: FILES.length * 13,
	deletions: FILES.length,
	truncated: false,
	files: FILES.map((path) => ({
		path,
		status: 'modified',
		additions: 13,
		deletions: 1,
		binary: false,
		truncated: false,
		patch: patch(12)
	}))
};

const item = {
	id: STORY,
	type: 'story',
	title:
		'The review page keeps Accept and its blockers on top and scrolls to its sections and each changed file from the sidebar',
	status: 'review',
	path: `wip/kanban/stories/${STORY}.md`,
	body: [
		'## Goal',
		'',
		'A sidebar on the pages that hold several sections.',
		'',
		'## Acceptance criteria',
		'',
		'- [x] The review page lists its sections and each changed file in a sidebar, and choosing one scrolls to it.',
		'- [x] The charts page lists its groups and kinds in the sidebar in place of the rows above the chart.',
		'- [ ] The documents page holds the document tree and the open document’s headings in the sidebar.',
		'- [ ] Below Tailwind’s `md` width no page scrolls sideways.',
		''
	].join('\n')
};

const narrative = [
	'## Current state',
	'',
	`Working in \`.flai-cache/worktrees/${STORY}\`: the sidebar is on three pages, and the end-to-end test of \`flaiover/src/lib/components/PageSidebar.svelte\` at a phone's width is next.`,
	'',
	'```sh',
	`scripts/flaiover-e2e.sh e2e/page-sidebar.e2e.ts # ${LONG_LINE}`,
	'```',
	'',
	'## Next steps',
	'',
	'1. Run the end-to-end tests.',
	'2. Move the story to review.',
	''
].join('\n');

const preview = {
	branch: `story/${STORY}`,
	blockers: [
		`the story's narrative wip/agents/${STORY}.md has an open question in its Open questions section that names flaiover/src/lib/components/PageSidebar.svelte`
	],
	plan: {
		level: 'minor',
		commits: ['feat: [S-9001] the page sidebar'],
		steps: [
			{
				component: { name: 'flaiover' },
				delivered: false,
				level: 'minor',
				from: '1.41.0',
				to: '1.42.0',
				files: FILES
			}
		]
	}
};

const thread = (on: string, anchor: { path: string; item?: string }) => ({
	id: 'TH-0901',
	title: `Should ${on}'s sidebar keep a heading's entry current while a long section is in view, or the one after it?`,
	anchor,
	status: 'open',
	participants: ['operator', `agent-${STORY}`],
	updated: NOW,
	entries: [
		{
			at: NOW,
			author: `agent-${STORY}`,
			text: `The entry stays current while its section is in view: see flaiover/src/lib/pagesidebar.svelte.ts and ${FILES[1]}.`
		}
	]
});

/** The document open on the documents page: headings of two levels, a wide table, a long line. */
const DOC = 'design/system/a-design-document-with-a-rather-long-file-name-for-the-tree.md';
const filler = (n: number) =>
	Array.from(
		{ length: n },
		(_, i) =>
			`Paragraph ${i + 1}: the sidebar lists the headings of the open document, each a link that scrolls to its heading, and the tree of documents heads it.`
	).join('\n\n');
const docBody = [
	'# The page sidebar and the sections it scrolls to on the pages that hold several',
	'',
	filler(2),
	'',
	'## Where the sidebar sits at a phone’s width and at a desktop’s width',
	'',
	filler(3),
	'',
	'### The row of entries above the content scrolls within itself',
	'',
	'| Page | Entries | Kind | What choosing one does | Where the content is | Notes on the layout at each width |',
	'|------|---------|------|------------------------|----------------------|-----------------------------------|',
	'| /review/S-9001 | Criteria, State, Threads, Changes, each changed file | scroll | scrolls to the section and opens a file’s panel | beside the sidebar from md | the Acceptance section stays above the sidebar |',
	'| /charts/cycle-time | flow, usage, planning, strategic, each chart | link | goes to the chart | beside the sidebar from md | the group is a label over its charts |',
	'',
	filler(3),
	'',
	'### A long line of code in a block of its own',
	'',
	'```ts',
	LONG_LINE,
	'```',
	'',
	filler(4),
	'',
	'## The entries of a document with a heading that is longer than the sidebar is wide',
	'',
	filler(6),
	'',
	'### The last heading, followed by enough text to scroll it to the top of the window',
	'',
	filler(30)
].join('\n');

const tree = [
	{
		name: 'design',
		path: 'design',
		kind: 'dir',
		children: [
			{
				name: 'system',
				path: 'design/system',
				kind: 'dir',
				children: [
					{
						name: 'README.md',
						path: 'design/system/README.md',
						kind: 'file',
						title: 'The design'
					},
					{
						name: DOC.split('/').pop(),
						path: DOC,
						kind: 'file',
						title: 'The page sidebar and the sections it scrolls to'
					}
				]
			},
			{
				name: 'conventions-with-a-long-directory-name-for-the-tree',
				path: 'design/conventions-with-a-long-directory-name-for-the-tree',
				kind: 'dir',
				children: [
					{
						name: 'README.md',
						path: 'design/conventions-with-a-long-directory-name-for-the-tree/README.md',
						kind: 'file'
					}
				]
			}
		]
	},
	{ name: 'wip', path: 'wip', kind: 'dir', children: [] }
];

/** Each page's own requests, answered here; the rest, the event stream among them, go on. */
async function answer(route: Route) {
	const url = new URL(route.request().url());
	const at = url.pathname;
	const json = (body: unknown, status = 200) => route.fulfill({ status, json: body });
	if (route.request().method() !== 'GET') return json({ error: 'not in this test' }, 405);
	if (at === '/api/board') return json({ writable: true, wip_limits: {}, order: [], columns: {} });
	if (at === '/api/items') return json([]);
	if (at === `/api/items/${STORY}`) return json({ item });
	if (at === `/api/items/${STORY}/diff`) return json(diff);
	if (at === `/api/items/${STORY}/acceptance`) return json(preview);
	if (at === `/api/items/${STORY}/checks`) return json({ checks_enabled: false });
	if (at === `/api/items/${STORY}/verify`) return json({ report: null });
	if (at === `/api/items/${STORY}/forge`) return json({ pull_request: null });
	if (at === '/api/components') return json({ error: 'no catalogue in this test' }, 404);
	if (at === '/api/queue') return json({ enabled: false });
	if (at === '/api/issues') return json([]);
	if (at === '/api/threads') {
		const on = url.searchParams.get('on') ?? '';
		return json([thread(on, on === STORY ? { path: item.path, item: STORY } : { path: on })]);
	}
	if (at === '/api/docs/tree') return json(tree);
	if (at === '/api/docs/file') {
		const path = url.searchParams.get('path');
		if (path === `wip/agents/${STORY}.md`)
			return json({ path, frontMatter: { story: STORY }, body: narrative });
		if (path === DOC)
			return json({
				path,
				frontMatter: {
					title: 'The page sidebar',
					topics: ['flaiover', 'layout', 'navigation'],
					links: [`flaiover/src/lib/components/PageSidebar.svelte`, FILES[1]]
				},
				body: docBody
			});
		return json({ error: 'no such document' }, 404);
	}
	return route.fallback();
}

test.beforeEach(async ({ context, baseURL }) => {
	await context.addCookies([{ name: 'flaiover_session', value: token(), url: baseURL! }]);
	await context.route(/\/api\//, answer);
});

/** How far the page reaches past the window sideways: 0 when it never scrolls sideways. */
const overflow = (page: Page) =>
	page.evaluate(() => document.documentElement.scrollWidth - window.innerWidth);

/**
 * What makes the page wider than the window, for the message of a failed overflow: from the body
 * down, the child that reaches furthest past the window's right edge, to the deepest. A box that
 * scrolls or clips sideways reaches as far as its own edge; any other as far as its content does,
 * a long word in its text among it.
 */
const overflowing = (page: Page) =>
	page.evaluate(() => {
		// a box of display: contents has no edge of its own, and reaches as far as its children
		const reach = (el: Element): number => {
			const style = getComputedStyle(el);
			if (style.display === 'contents') return Math.max(0, ...[...el.children].map(reach));
			const r = el.getBoundingClientRect();
			return style.overflowX === 'visible' ? Math.max(r.right, r.left + el.scrollWidth) : r.right;
		};
		const name = (el: Element) =>
			[
				el.tagName.toLowerCase(),
				el.id && `#${el.id}`,
				...[...el.attributes].filter((a) => a.name.startsWith('data-')).map((a) => `[${a.name}]`)
			].join('');
		const trail: string[] = [];
		let at: Element | undefined = document.body;
		while (at) {
			let furthest: Element | undefined;
			for (const child of at.children)
				if (reach(child) > window.innerWidth + 0.5 && (!furthest || reach(child) > reach(furthest)))
					furthest = child;
			if (furthest) trail.push(`${name(furthest)} to ${Math.round(reach(furthest))}px`);
			at = furthest;
		}
		return `past ${window.innerWidth}px: ${trail.slice(-4).join(' > ')}`;
	});

/** The page does not scroll sideways; when it does, the message names what reaches past. */
async function expectNoOverflow(page: Page) {
	expect(await overflow(page), await overflowing(page)).toBeLessThanOrEqual(0);
}

const nav = (page: Page) => page.locator('[data-sidebar-nav]');
const content = (page: Page) => page.locator('[data-sidebar-content]');
const entry = (page: Page, id: string) => page.locator(`[data-sidebar-entry="${id}"]`);

/** The review page, with its diff drawn and the acceptance worked out. */
async function arriveReview(page: Page) {
	await page.goto(`/review/${STORY}`);
	await expect(page.locator('[data-testid="diff-toggle"]')).toHaveCount(FILES.length);
	await expect(page.locator('#review-release')).toContainText('Tag a minor release');
	await expect(page.locator('#review-state')).toContainText('Current state');
}

/** The documents page, the tree and the open document's headings drawn. */
async function arriveDocs(page: Page) {
	await page.goto(`/docs/${DOC}`);
	await expect(page.locator('[data-doc-tree] a[title]').first()).toBeVisible();
	await expect(page.locator('article.prose h3').last()).toBeVisible();
	await expect(entry(page, 'On this page')).toBeVisible();
}

/** The sidebar is a row above the content: its bottom at or above the content's top. */
async function expectAbove(page: Page) {
	const n = (await nav(page).boundingBox())!;
	const c = (await content(page).boundingBox())!;
	expect(n.y + n.height, 'the sidebar above the content').toBeLessThanOrEqual(c.y + 0.5);
}

/** The sidebar is a column left of the content: its right edge at or left of the content's left. */
async function expectBeside(page: Page) {
	const n = (await nav(page).boundingBox())!;
	const c = (await content(page).boundingBox())!;
	expect(n.x + n.width, 'the sidebar left of the content').toBeLessThanOrEqual(c.x + 0.5);
	expect(n.y, 'the sidebar beside the content, not above it').toBeLessThan(c.y + c.height);
}

/** Every entry named is listed in the sidebar. */
async function expectListed(page: Page, ids: string[]) {
	for (const id of ids) {
		await expect(nav(page).locator(`[data-sidebar-entry="${id}"]`), id).toHaveCount(1);
		await expect(nav(page).locator(`[data-sidebar-entry="${id}"]`), id).toBeVisible();
	}
}

const REVIEW_ENTRIES = [
	'review-criteria',
	'review-state',
	'review-threads',
	'review-changes',
	...FILES.map(fileId),
	'review-release',
	'review-checks'
];
const DOC_HEADINGS = (page: Page) =>
	page
		.locator('article.prose')
		.locator('h1, h2, h3')
		.evaluateAll((hs) => hs.map((h) => h.id).filter(Boolean));

test.describe('at 375px', () => {
	test.use({ viewport: { width: 375, height: 800 } });

	test('the review page lists its sections and files above the content, and does not scroll sideways', async ({
		page
	}) => {
		await arriveReview(page);
		await expectListed(page, REVIEW_ENTRIES);
		await expectAbove(page);
		await expectNoOverflow(page);

		// a file's panel, with its long line, open
		await entry(page, fileId(FILES[1])).click();
		const row = page.locator(`[id="${fileId(FILES[1])}"]`);
		await expect(row.locator('[data-testid="diff-panel"]')).toBeVisible();
		await expectNoOverflow(page);
	});

	test('the documents page lists the tree and the headings above the document, and does not scroll sideways', async ({
		page
	}) => {
		await arriveDocs(page);
		await expect(nav(page).locator('[data-sidebar-header] [data-doc-tree]')).toBeVisible();
		await expect(
			nav(page).locator(`[data-sidebar-header] a[title="${DOC}"]`),
			'the open document in the tree'
		).toBeVisible();
		const ids = await DOC_HEADINGS(page);
		expect(ids.length).toBeGreaterThanOrEqual(6);
		await expectListed(page, ['On this page', ...ids]);
		await expectAbove(page);
		await expectNoOverflow(page);
	});
});

test.describe('at 1280px', () => {
	test.use({ viewport: { width: 1280, height: 800 } });

	test('the review page has its sidebar left of the sections, and a file entry opens its panel in view', async ({
		page
	}) => {
		await arriveReview(page);
		await expectListed(page, REVIEW_ENTRIES);
		await expectBeside(page);
		await expectNoOverflow(page);

		const last = FILES[FILES.length - 1];
		const row = page.locator(`[id="${fileId(last)}"]`);
		await expect(row.locator('[data-testid="diff-panel"]')).toHaveCount(0);
		await entry(page, fileId(last)).click();
		await expect(row.locator('[data-testid="diff-panel"]')).toBeVisible();
		await expect(row).toBeInViewport();
		await expect(page).toHaveURL(new RegExp(`#${fileId(last)}$`));
		await expectNoOverflow(page);
	});

	test('the documents page has the tree and headings left of the document, and a heading entry scrolls to it', async ({
		page
	}) => {
		await arriveDocs(page);
		const ids = await DOC_HEADINGS(page);
		await expectListed(page, ['On this page', ...ids]);
		await expectBeside(page);
		await expectNoOverflow(page);

		const target = ids[ids.length - 1];
		await entry(page, target).click();
		const heading = page.locator(`article.prose [id="${target}"]`);
		await expect
			.poll(async () => (await heading.boundingBox())!.y, { message: 'the heading at the top' })
			.toBeLessThan(80);
		expect((await heading.boundingBox())!.y).toBeGreaterThanOrEqual(-1);
		await expect(entry(page, target)).toHaveAttribute('aria-current', 'true');
		await expectNoOverflow(page);
	});
});

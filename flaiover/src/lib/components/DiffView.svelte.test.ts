import { afterEach, describe, expect, it } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import DiffView from './DiffView.svelte';

const file = (over: Record<string, unknown>) => ({
	path: 'docs/a.md',
	status: 'modified',
	additions: 0,
	deletions: 0,
	binary: false,
	truncated: false,
	patch: '',
	...over
});
const DIFF = {
	branch: 'story/S-0164',
	base: 'abc1234',
	commits: 2,
	additions: 4,
	deletions: 3,
	truncated: false,
	files: [
		file({
			additions: 4,
			deletions: 3,
			patch:
				'@@ -1,5 +1,6 @@\n one\n-two\n-\tthree\n+2\n+3\n+3.5\n four\n-five\n+5\n\\ No newline at end of file\n'
		}),
		file({ path: 'static/logo.png', status: 'added', binary: true }),
		file({
			path: 'cli/main.go',
			status: 'modified',
			additions: 1,
			patch: '@@ -1 +1,2 @@\n a\n+b\n'
		})
	]
};

const toggles = () => [
	...document.querySelectorAll<HTMLButtonElement>('[data-testid="diff-toggle"]')
];
const panels = () => [...document.querySelectorAll<HTMLElement>('[data-testid="diff-panel"]')];
const click = (b: HTMLElement) => {
	b.click();
	flushSync();
};
/** Each run of a panel as its kind and its lines, a line as its margin and its text. */
const runs = (panel: HTMLElement) =>
	[...panel.querySelectorAll<HTMLElement>('[data-run]')].map((run) => ({
		kind: run.dataset.run,
		classes: run.className,
		lines: [...run.querySelectorAll<HTMLElement>('[data-line]')].map((l) => ({
			sign: l.querySelector('[data-sign]')!.textContent,
			text: l.lastElementChild!.textContent,
			classes: l.className
		}))
	}));

describe('DiffView (S-0164)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.body.innerHTML = '';
	});
	const show = (diff: unknown = DIFF) => {
		c = mount(DiffView, { target: document.body, props: { diff: diff as never } });
		flushSync();
	};

	it('has a toggle control on every summary line that reveals and hides its panel', () => {
		show();
		expect(toggles()).toHaveLength(3);
		expect(panels()).toHaveLength(0);
		const [first, , third] = toggles();
		expect(first.textContent).toContain('docs/a.md');
		expect(first.textContent).toContain('▸');
		expect(first.getAttribute('aria-expanded')).toBe('false');

		click(first);
		expect(first.getAttribute('aria-expanded')).toBe('true');
		expect(first.textContent).toContain('▾');
		expect(panels()).toHaveLength(1);
		// the line names the panel it opens
		expect(panels()[0].id).toBe(first.getAttribute('aria-controls'));
		expect(panels()[0].textContent).toContain('three');

		// each file has a panel of its own
		click(third);
		expect(panels()).toHaveLength(2);
		expect(third.getAttribute('aria-controls')).not.toBe(first.getAttribute('aria-controls'));
		expect(panels()[1].id).toBe(third.getAttribute('aria-controls'));

		click(first);
		expect(first.getAttribute('aria-expanded')).toBe('false');
		expect(first.textContent).toContain('▸');
		expect(panels()).toHaveLength(1);
		expect(panels()[0].textContent).not.toContain('three');
	});

	it('puts - in the margin of every removed line and + in that of every added line', () => {
		show();
		click(toggles()[0]);
		const lines = runs(panels()[0]).flatMap((r) => r.lines.map((l) => ({ kind: r.kind, ...l })));
		const of = (kind: string) => lines.filter((l) => l.kind === kind);
		expect(of('del').map((l) => [l.sign, l.text])).toEqual([
			['-', 'two'],
			['-', '\tthree'],
			['-', 'five']
		]);
		expect(of('add').map((l) => [l.sign, l.text])).toEqual([
			['+', '2'],
			['+', '3'],
			['+', '3.5'],
			['+', '5']
		]);
		// nothing else has a sign, and a hunk header and a note keep their text
		expect(lines.filter((l) => l.kind !== 'add' && l.kind !== 'del').map((l) => l.sign)).toEqual([
			'',
			'',
			'',
			''
		]);
		expect(of('hunk')[0].text).toBe('@@ -1,5 +1,6 @@');
		expect(of('note')[0].text).toBe('\\ No newline at end of file');
	});

	it('gives a line its dim background and consecutive lines one outline', () => {
		show();
		click(toggles()[0]);
		const got = runs(panels()[0]);
		expect(got.map((r) => [r.kind, r.lines.length])).toEqual([
			['hunk', 1],
			['context', 1],
			['del', 2],
			['add', 3],
			['context', 1],
			['del', 1],
			['add', 1],
			['note', 1]
		]);
		for (const run of got) {
			const outlined = /\bborder-(good|danger)\b/.exec(run.classes)?.[1];
			expect(outlined).toBe({ add: 'good', del: 'danger' }[run.kind as string]);
			for (const l of run.lines) {
				const behind = /\bbg-(good|danger)-soft\b/.exec(l.classes)?.[1];
				expect(behind).toBe(outlined);
			}
		}
	});

	it('says why a file has no lines to show', () => {
		show({
			...DIFF,
			truncated: true,
			files: [
				DIFF.files[1],
				file({ path: 'big.txt', additions: 9000, truncated: true }),
				file({ path: 'mode.sh' }),
				file({ path: 'long.txt', additions: 1, truncated: true, patch: '@@ -0,0 +1 @@\n+a\n' })
			]
		});
		toggles().forEach(click);
		expect(panels().map((p) => p.textContent?.replace(/\s+/g, ' ').trim())).toEqual([
			'Binary file, no hunks.',
			'Left out for size.',
			'No textual changes.',
			'@@ -0,0 +1 @@+a Cut for size.'
		]);
	});
});

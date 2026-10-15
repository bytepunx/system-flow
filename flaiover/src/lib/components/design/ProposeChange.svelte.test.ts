// S-0384: the Design page's request form makes a story under the epic chosen, with topic design,
// the operator's words as its goal, fixed criteria, and the token file and the stylesheet in its
// touches, and moves it to ready; a refusal of either step is shown in flai's words.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import ProposeChange from './ProposeChange.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));
vi.mock('$app/paths', () => ({
	resolve: (route: string, p?: { id: string }) => route.replace('[id]', p?.id ?? '')
}));

const settle = async () => {
	for (let i = 0; i < 6; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const ok = (body: unknown) => ({ ok: true, statusText: 'OK', json: async () => body });
const refused = (error: string, findings?: unknown[]) => ({
	ok: false,
	statusText: 'Unprocessable',
	json: async () => ({ error, findings })
});

const EPICS = [
	{ id: 'E-0020', title: 'Design system' },
	{ id: 'E-0006', title: 'Workbench' }
];
const TOKENS = 'design/system/tokens.json';
const STYLESHEET = 'flaiover/src/theme.css';

// The calls the form made, and what create and move answer.
function backend(
	create: (body: Record<string, unknown>) => unknown = () => ok({ item: { id: 'S-0400' } }),
	move: (id: string, to: string) => unknown = (id, to) => ok({ id, status: to })
) {
	const calls: { url: string; body: Record<string, unknown> }[] = [];
	api.mockImplementation(async (url: string, init?: { method?: string; body?: string }) => {
		const body = JSON.parse(init?.body ?? '{}') as Record<string, unknown>;
		calls.push({ url, body });
		const moving = url.match(/^\/api\/items\/([^/]+)\/move$/);
		if (moving) return move(moving[1], body.to as string);
		if (url === '/api/items' && init?.method === 'POST') return create(body);
		throw new Error('unexpected ' + url);
	});
	return calls;
}

const q = <T extends HTMLElement>(id: string) =>
	document.querySelector<T>(`[data-testid="${id}"]`)!;
const type = (el: HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement, value: string) => {
	el.value = value;
	el.dispatchEvent(
		new Event(el instanceof HTMLSelectElement ? 'change' : 'input', { bubbles: true })
	);
	flushSync();
};
const button = () => document.querySelector<HTMLButtonElement>('form button')!;
const words = (e: Element) => e.textContent!.replace(/\s+/g, ' ').trim();

describe('ProposeChange (S-0384)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		document.body.innerHTML = '';
	});

	function render(props: Record<string, unknown> = {}) {
		const oncreated = vi.fn();
		c = mount(ProposeChange, {
			target: document.body,
			props: { epics: EPICS, tokens: TOKENS, oncreated, ...props }
		});
		flushSync();
		return oncreated;
	}
	function fill(epic = 'E-0020') {
		type(q<HTMLSelectElement>('propose-epic'), epic);
		type(q<HTMLInputElement>('propose-title'), 'Warmer primary colour');
		type(q<HTMLTextAreaElement>('propose-request'), '  Make the primary colour warmer.\n');
	}
	async function submit() {
		q<HTMLFormElement>('propose-change').dispatchEvent(
			new Event('submit', { bubbles: true, cancelable: true })
		);
		await settle();
	}

	it('offers the epics given and is disabled until an epic, a title, and the request are filled', () => {
		backend();
		render();
		const options = [...q<HTMLSelectElement>('propose-epic').options].map((o) => o.value);
		expect(options).toEqual(['', 'E-0020', 'E-0006']);
		expect(button().disabled).toBe(true);
		type(q<HTMLSelectElement>('propose-epic'), 'E-0020');
		expect(button().disabled).toBe(true);
		type(q<HTMLInputElement>('propose-title'), 'Warmer primary colour');
		expect(button().disabled).toBe(true);
		type(q<HTMLTextAreaElement>('propose-request'), '   ');
		expect(button().disabled).toBe(true);
		type(q<HTMLTextAreaElement>('propose-request'), 'Make it warmer.');
		expect(button().disabled).toBe(false);
		type(q<HTMLInputElement>('propose-title'), '');
		expect(button().disabled).toBe(true);
	});

	it('chooses a lone epic for the operator', () => {
		backend();
		render({ epics: [EPICS[1]] });
		expect(q<HTMLSelectElement>('propose-epic').value).toBe('E-0006');
	});

	it('makes a story with topic design, the goal in their words, the fixed criteria, and the token file and stylesheet in its touches, then moves it to ready', async () => {
		const calls = backend();
		const oncreated = render({ stylesheet: STYLESHEET });
		expect(words(q('propose-about'))).toContain(`changes ${TOKENS} and ${STYLESHEET}, moved`);
		fill('E-0006');
		type(q<HTMLSelectElement>('propose-nature'), 'improvement');
		await submit();
		expect(calls.map((c) => c.url)).toEqual(['/api/items', '/api/items/S-0400/move']);
		expect(calls[0].body).toEqual({
			type: 'story',
			title: 'Warmer primary colour',
			nature: 'improvement',
			parent: 'E-0006',
			topics: ['design'],
			touches: [TOKENS, STYLESHEET],
			body:
				'## Goal\n\nMake the primary colour warmer.\n\n' +
				'## Acceptance criteria\n\n' +
				'- [ ] The tokens are changed as the goal asks.\n' +
				'- [ ] The generated stylesheet is in step with the tokens: `flai design css --write` is run, and `flai check` reports no `design.stale`.\n' +
				"- [ ] The UI's tests pass.\n\n" +
				'## Tasks\n\n## Notes\n'
		});
		expect(calls[1].body).toEqual({ to: 'ready' });
		expect(oncreated).toHaveBeenCalledWith('S-0400');
	});

	it('touches only the token file when no stylesheet is given, and starts as a feature', async () => {
		const calls = backend();
		render();
		expect(words(q('propose-about'))).toContain(`changes ${TOKENS}, moved`);
		fill();
		await submit();
		expect(calls[0].body.touches).toEqual([TOKENS]);
		expect(calls[0].body.nature).toBe('feature');
	});

	it('shows the new story with a link to it and clears the form', async () => {
		backend();
		render();
		fill();
		await submit();
		const created = q('propose-created');
		expect(words(created)).toBe('S-0400 is ready for a story’s agent to pull.');
		expect(created.querySelector('a')!.getAttribute('href')).toBe('/items/S-0400');
		expect(q<HTMLInputElement>('propose-title').value).toBe('');
		expect(q<HTMLTextAreaElement>('propose-request').value).toBe('');
		expect(q<HTMLSelectElement>('propose-epic').value).toBe('E-0020');
		expect(document.querySelector('[data-testid="propose-error"]')).toBeNull();
		expect(button().disabled).toBe(true);
	});

	it("shows flai's refusal to create, with its findings, keeps the words, and moves nothing", async () => {
		const calls = backend(() =>
			refused('flai check has 1 finding(s) with this story; nothing was created', [
				{ level: 'error', rule: 'item.parent', path: 'wip/x.md', line: 3, message: 'no such epic' }
			])
		);
		const oncreated = render();
		fill();
		await submit();
		const error = words(q('propose-error'));
		expect(error).toContain(
			'The story was not created: flai check has 1 finding(s) with this story; nothing was created'
		);
		expect(error).toContain('wip/x.md:3: error: item.parent: no such epic');
		expect(calls.map((c) => c.url)).toEqual(['/api/items']);
		expect(q<HTMLInputElement>('propose-title').value).toBe('Warmer primary colour');
		expect(q<HTMLTextAreaElement>('propose-request').value).toBe(
			'  Make the primary colour warmer.\n'
		);
		expect(document.querySelector('[data-testid="propose-created"]')).toBeNull();
		expect(oncreated).not.toHaveBeenCalled();
		expect(button().disabled).toBe(false);
	});

	it("says a story whose move was refused is in backlog, names it, and shows flai's error", async () => {
		backend(undefined, () => refused('the in-progress limit is reached'));
		const oncreated = render();
		fill();
		await submit();
		expect(words(q('propose-error'))).toBe(
			'S-0400 was created and is in backlog, but moving it to ready was refused: the in-progress limit is reached'
		);
		const created = q('propose-created');
		expect(words(created)).toBe('S-0400 is in backlog.');
		expect(created.querySelector('a')!.getAttribute('href')).toBe('/items/S-0400');
		// It exists, so its words are cleared and it is not made again.
		expect(q<HTMLInputElement>('propose-title').value).toBe('');
		expect(button().disabled).toBe(true);
		expect(oncreated).toHaveBeenCalledWith('S-0400');
	});

	it('shows a request that did not reach the dashboard', async () => {
		api.mockRejectedValue(new Error('Failed to fetch'));
		render();
		fill();
		await submit();
		expect(words(q('propose-error'))).toBe('Failed to fetch');
	});
});

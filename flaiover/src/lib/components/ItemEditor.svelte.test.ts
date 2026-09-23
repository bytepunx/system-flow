import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import ItemEditor from './ItemEditor.svelte';

const api = vi.fn();
vi.mock('$lib/api', () => ({ api: (...args: unknown[]) => api(...args) }));

const settle = async () => {
	for (let i = 0; i < 6; i++) await new Promise((r) => setTimeout(r, 0));
	flushSync();
};
const view = {
	id: 'S-0007',
	type: 'story',
	status: 'ready',
	title: 'A plain story',
	nature: 'feature',
	tags: ['cli'],
	touches: ['flai/cmd'],
	parent: 'E-0001',
	body: '## Goal\nx\n\n## Acceptance criteria\n- [ ] it works\n',
	path: 'wip/kanban/stories/S-0007-a-plain-story.md',
	hash: 'a'.repeat(64),
	editable: true,
	natures: ['feature', 'improvement', 'remediation', 'research', 'experiment'],
	parents: [
		{ id: 'E-0001', title: 'First' },
		{ id: 'E-0002', title: 'Second' }
	]
};
const answer = (status: number, body: unknown) => ({
	ok: status < 400,
	status,
	statusText: String(status),
	json: async () => body
});
const type = (selector: string, value: string) => {
	const el = document.querySelector<HTMLInputElement | HTMLTextAreaElement>(selector)!;
	el.value = value;
	el.dispatchEvent(new Event('input', { bubbles: true }));
	flushSync();
};
const submit = async () => {
	document
		.querySelector('form')!
		.dispatchEvent(new Event('submit', { bubbles: true, cancelable: true }));
	await settle();
};
const sent = () => {
	const put = api.mock.calls.filter((c) => c[1]?.method === 'PUT').at(-1)!;
	return JSON.parse(put[1].body as string);
};

describe('ItemEditor (S-0085)', () => {
	let c: ReturnType<typeof mount> | undefined;
	const saved = vi.fn();
	const mountIt = async () => {
		c = mount(ItemEditor, {
			target: document.body,
			props: { id: 'S-0007', oncancel: () => {}, onsaved: saved }
		});
		await settle();
	};
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		api.mockReset();
		saved.mockReset();
		document.body.innerHTML = '';
	});

	it('shows the fields, says what stays flai’s, and sends only what changed, with the hash', async () => {
		api.mockImplementation((_url: string, init?: { method?: string }) =>
			Promise.resolve(
				init?.method === 'PUT' ? answer(200, { changed: ['title', 'tags'] }) : answer(200, view)
			)
		);
		await mountIt();
		expect(api).toHaveBeenCalledWith('/api/items/S-0007/edit');
		const text = document.body.textContent!.replace(/\s+/g, ' ');
		expect(text).toContain('S-0007');
		expect(text).toContain('the status changes by moving the card');
		const save = document.querySelector<HTMLButtonElement>('[data-testid="edit-save"]')!;
		expect(save.disabled).toBe(true); // nothing changed yet
		type('[data-testid="edit-title"]', 'A better name');
		type('input[placeholder="comma separated"]', 'cli, dashboard');
		expect(save.disabled).toBe(false);
		await submit();
		expect(sent()).toEqual({
			title: 'A better name',
			tags: ['cli', 'dashboard'],
			hash: view.hash
		});
		expect(saved).toHaveBeenCalledWith(['title', 'tags']);
	});

	it('keeps the text and lists the findings when the check refuses', async () => {
		api.mockImplementation((_url: string, init?: { method?: string }) =>
			Promise.resolve(
				init?.method === 'PUT'
					? answer(422, {
							error: 'flai check has 1 finding(s) with this change; nothing was changed',
							findings: [
								{
									path: view.path,
									line: 1,
									level: 'error',
									rule: 'story.criteria',
									message: 'a ready story needs acceptance criteria'
								}
							]
						})
					: answer(200, view)
			)
		);
		await mountIt();
		type('[data-testid="edit-body"]', '## Goal\nNo criteria.\n');
		await submit();
		const err = document.querySelector('[data-testid="edit-error"]')!.textContent!;
		expect(err).toContain('story.criteria');
		expect(err).toContain('your text is still here');
		expect(
			document.querySelector<HTMLTextAreaElement>('[data-testid="edit-body"]')!.value
		).toContain('No criteria.');
		expect(saved).not.toHaveBeenCalled();
	});

	it('shows a conflict and lets the designer load the current version or save over it', async () => {
		let puts = 0;
		api.mockImplementation((_url: string, init?: { method?: string }) => {
			if (init?.method !== 'PUT') return Promise.resolve(answer(200, view));
			puts++;
			return Promise.resolve(
				puts === 1
					? answer(409, { error: 'conflict', hash: 'b'.repeat(64), current: '...' })
					: answer(200, { changed: ['title'] })
			);
		});
		await mountIt();
		type('[data-testid="edit-title"]', 'Mine');
		await submit();
		const box = document.querySelector('[data-testid="edit-conflict"]')!;
		expect(box.textContent).toContain('changed after you opened it');
		expect(saved).not.toHaveBeenCalled();
		[...box.querySelectorAll('button')].find((b) => b.textContent!.includes('over it'))!.click();
		await settle();
		expect(sent()).toEqual({ title: 'Mine', hash: 'b'.repeat(64) });
		expect(saved).toHaveBeenCalledWith(['title']);
	});

	it('says why an item cannot be edited', async () => {
		api.mockResolvedValue(
			answer(200, {
				...view,
				editable: false,
				reason: 'S-0007 is done; a closed item is not edited'
			})
		);
		await mountIt();
		expect(document.body.textContent).toContain('a closed item is not edited');
		expect(document.querySelector('form')).toBeNull();
	});

	// S-0103: the story's own agent is shown, and a change replaces it; emptied, it is removed
	it("edits the story's agent and sends it whole, or null when emptied", async () => {
		const withAgent = {
			...view,
			agent: { harness: 'claude-code', model: 'claude-opus-5-5', config: { effort: 'high' } },
			default_agent: { harness: 'claude-code', model: 'claude-opus-5-5' }
		};
		api.mockImplementation((_url: string, init?: { method?: string }) =>
			Promise.resolve(
				init?.method === 'PUT' ? answer(200, { changed: ['agent'] }) : answer(200, withAgent)
			)
		);
		await mountIt();
		const model = document.querySelector<HTMLInputElement>('[data-testid="agent-model"]')!;
		expect(model.value).toBe('claude-opus-5-5');
		expect(document.querySelector<HTMLTextAreaElement>('[data-testid="agent-config"]')!.value).toBe(
			'effort=high'
		);
		type('[data-testid="agent-model"]', 'claude-sonnet-5');
		await submit();
		expect(sent()).toEqual({
			agent: { harness: 'claude-code', model: 'claude-sonnet-5', config: { effort: 'high' } },
			hash: view.hash
		});
		type('[data-testid="agent-harness"]', '');
		type('[data-testid="agent-model"]', '');
		type('[data-testid="agent-config"]', '');
		await submit();
		expect(sent()).toEqual({ agent: null, hash: view.hash });
	});
});

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
		expect(err).toContain(
			`${view.path}:1: error: story.criteria: a ready story needs acceptance criteria`
		);
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

	// S-0130: a story's after: is shown and sent as a list; emptied, it is sent empty and removed
	it('edits the stories a story waits for', async () => {
		api.mockImplementation((_url: string, init?: { method?: string }) =>
			Promise.resolve(
				init?.method === 'PUT'
					? answer(200, { changed: ['after'] })
					: answer(200, { ...view, after: ['S-0005'] })
			)
		);
		await mountIt();
		const field = document.querySelector<HTMLInputElement>('[data-testid="edit-after"]')!;
		expect(field.value).toBe('S-0005');
		expect(document.body.textContent).toContain('held, until each of these is done');
		type('[data-testid="edit-after"]', 'S-0005, S-0006');
		await submit();
		expect(sent()).toEqual({ after: ['S-0005', 'S-0006'], hash: view.hash });
		type('[data-testid="edit-after"]', '');
		await submit();
		expect(sent()).toEqual({ after: [], hash: view.hash });
	});

	it('has no after field for an epic, and sends none for a story whose flai does not know it', async () => {
		api.mockImplementation((_url: string, init?: { method?: string }) =>
			Promise.resolve(
				init?.method === 'PUT' ? answer(200, { changed: ['title'] }) : answer(200, view)
			)
		);
		await mountIt();
		type('[data-testid="edit-title"]', 'Renamed');
		await submit();
		expect(sent()).toEqual({ title: 'Renamed', hash: view.hash });
		unmount(c!);
		api.mockResolvedValue(answer(200, { ...view, type: 'epic', parent: undefined }));
		await mountIt();
		expect(document.querySelector('[data-testid="edit-after"]')).toBeNull();
	});

	// S-0135: a story's or epic's topics are shown next to its tags and sent as a list; emptied,
	// they are sent empty and removed
	it('edits the topics of a story and of an epic', async () => {
		api.mockImplementation((_url: string, init?: { method?: string }) =>
			Promise.resolve(
				init?.method === 'PUT'
					? answer(200, { changed: ['topics'] })
					: answer(200, { ...view, topics: ['logging'] })
			)
		);
		await mountIt();
		const field = document.querySelector<HTMLInputElement>('[data-testid="edit-topics"]')!;
		expect(field.value).toBe('logging');
		expect(document.body.textContent).toContain('beyond the components its tags and touches reach');
		type('[data-testid="edit-topics"]', 'logging, release');
		await submit();
		expect(sent()).toEqual({ topics: ['logging', 'release'], hash: view.hash });
		type('[data-testid="edit-topics"]', '');
		await submit();
		expect(sent()).toEqual({ topics: [], hash: view.hash });
		unmount(c!);
		api.mockReset();
		api.mockImplementation((_url: string, init?: { method?: string }) =>
			Promise.resolve(
				init?.method === 'PUT'
					? answer(200, { changed: ['topics'] })
					: answer(200, { ...view, type: 'epic', parent: undefined })
			)
		);
		await mountIt();
		type('[data-testid="edit-topics"]', 'release');
		await submit();
		expect(sent()).toEqual({ topics: ['release'], hash: view.hash });
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

	// S-0204: an epic's or story's cost of delay inputs in a closed panel; a save sends the changed
	// inputs, an emptied one as "", and never the planner's value
	describe('the cost of delay panel', () => {
		const planned = {
			...view,
			cost_of_delay: {
				inputs: { revenue_per_week: 1200, penalty_per_week: 300, time_lost_per_cycle: '6h' },
				value: 900,
				by: 'alex',
				at: '2026-10-01T09:00:00Z'
			},
			currency: 'EUR'
		};
		const panel = () => document.querySelector<HTMLDetailsElement>('[data-testid="cod-panel"]');

		it('is closed, says the inputs, and sends only those changed', async () => {
			api.mockImplementation((_url: string, init?: { method?: string }) =>
				Promise.resolve(
					init?.method === 'PUT'
						? answer(200, { changed: ['cost_of_delay'] })
						: answer(200, planned)
				)
			);
			await mountIt();
			expect(panel()!.open).toBe(false);
			expect(document.querySelector('[data-testid="cod-summary"]')!.textContent).toBe(
				'revenue 1,200 EUR/week · penalty 300 EUR/week · 6h lost per cycle'
			);
			expect(document.querySelector<HTMLInputElement>('[data-testid="cod-revenue"]')!.value).toBe(
				'1200'
			);
			const save = document.querySelector<HTMLButtonElement>('[data-testid="edit-save"]')!;
			expect(save.disabled).toBe(true);
			type('[data-testid="cod-revenue"]', '1500');
			expect(save.disabled).toBe(false);
			await submit();
			expect(sent()).toEqual({ cost_of_delay: { revenue_per_week: '1500' }, hash: view.hash });
		});

		it('clears every input as "" and leaves the value', async () => {
			api.mockImplementation((_url: string, init?: { method?: string }) =>
				Promise.resolve(
					init?.method === 'PUT'
						? answer(200, { changed: ['cost_of_delay'] })
						: answer(200, planned)
				)
			);
			await mountIt();
			type('[data-testid="cod-revenue"]', '');
			type('[data-testid="cod-penalty"]', '');
			type('[data-testid="cod-time-lost"]', '');
			await submit();
			expect(sent()).toEqual({
				cost_of_delay: { revenue_per_week: '', penalty_per_week: '', time_lost_per_cycle: '' },
				hash: view.hash
			});
		});

		it('is on an epic in USD when the view names no currency, and not on a task', async () => {
			api.mockImplementation((_url: string, init?: { method?: string }) =>
				Promise.resolve(
					init?.method === 'PUT'
						? answer(200, { changed: ['cost_of_delay'] })
						: answer(200, { ...view, type: 'epic', parent: undefined })
				)
			);
			await mountIt();
			expect(panel()!.textContent).toContain('revenue per week, USD');
			expect(document.querySelector('[data-testid="cod-summary"]')).toBeNull();
			type('[data-testid="cod-time-lost"]', '1h30m');
			await submit();
			expect(sent()).toEqual({
				cost_of_delay: { time_lost_per_cycle: '1h30m' },
				hash: view.hash
			});
			unmount(c!);
			api.mockReset();
			api.mockResolvedValue(answer(200, { ...view, id: 'T-0001', type: 'task' }));
			await mountIt();
			expect(document.querySelector('[data-testid="item-editor"]')).not.toBeNull();
			expect(panel()).toBeNull();
		});
	});

	// S-0201: a draft story says it cannot go to ready and is finalized from the form
	describe('a draft story', () => {
		const draft = { ...view, status: 'backlog', draft: true };
		const finalized = { ...view, status: 'backlog', draft: false, hash: 'c'.repeat(64) };
		const box = () => document.querySelector('[data-testid="editor-draft"]');
		const finalizeIt = async () => {
			document.querySelector<HTMLButtonElement>('[data-testid="editor-finalize"]')!.click();
			await settle();
		};

		it('is warned about with a Finalize button, and no other item is', async () => {
			api.mockResolvedValue(answer(200, draft));
			await mountIt();
			expect(box()!.textContent).toContain('cannot be moved to ready until you finalize it');
			expect(document.querySelector('[data-testid="editor-finalize"]')).not.toBeNull();
			unmount(c!);
			api.mockResolvedValue(answer(200, { ...view, draft: false }));
			await mountIt();
			expect(box()).toBeNull();
			unmount(c!);
			api.mockResolvedValue(answer(200, view));
			await mountIt();
			expect(box()).toBeNull();
			unmount(c!);
			api.mockResolvedValue(answer(200, { ...draft, type: 'epic', parent: undefined }));
			await mountIt();
			expect(box()).toBeNull();
			expect(document.querySelector('[data-testid="editor-finalize"]')).toBeNull();
		});

		it('is finalized in place, keeping what is typed, and a save carries the new hash', async () => {
			let shown = draft;
			api.mockImplementation((_url: string, init?: { method?: string }) => {
				if (init?.method === 'POST') {
					shown = finalized;
					return Promise.resolve(answer(200, { id: 'S-0007' }));
				}
				if (init?.method === 'PUT') return Promise.resolve(answer(200, { changed: ['title'] }));
				return Promise.resolve(answer(200, shown));
			});
			await mountIt();
			type('[data-testid="edit-title"]', 'Typed before finalizing');
			await finalizeIt();
			expect(api).toHaveBeenCalledWith(
				'/api/items/S-0007/finalize',
				expect.objectContaining({ method: 'POST' })
			);
			expect(box()).toBeNull();
			expect(document.querySelector('[data-testid="edit-error"]')).toBeNull();
			expect(document.querySelector<HTMLInputElement>('[data-testid="edit-title"]')!.value).toBe(
				'Typed before finalizing'
			);
			expect(saved).not.toHaveBeenCalled();
			await submit();
			expect(sent()).toEqual({ title: 'Typed before finalizing', hash: finalized.hash });
		});

		it('keeps the old hash when its words changed meanwhile, so the save reports the conflict', async () => {
			let shown = draft;
			api.mockImplementation((_url: string, init?: { method?: string }) => {
				if (init?.method === 'POST') {
					shown = { ...finalized, title: 'Renamed by an agent' };
					return Promise.resolve(answer(200, { id: 'S-0007' }));
				}
				if (init?.method === 'PUT') return Promise.resolve(answer(200, { changed: ['body'] }));
				return Promise.resolve(answer(200, shown));
			});
			await mountIt();
			await finalizeIt();
			expect(box()).toBeNull();
			type('[data-testid="edit-body"]', '## Goal\ny\n');
			await submit();
			expect(sent().hash).toBe(view.hash);
		});

		it('shows a refusal as an error and stays a draft', async () => {
			api.mockImplementation((_url: string, init?: { method?: string }) =>
				Promise.resolve(
					init?.method === 'POST'
						? answer(400, { error: 'S-0007 is not a draft' })
						: answer(200, draft)
				)
			);
			await mountIt();
			await finalizeIt();
			expect(document.querySelector('[data-testid="edit-error"]')!.textContent).toContain(
				'S-0007 is not a draft'
			);
			expect(box()).not.toBeNull();
		});
	});
});

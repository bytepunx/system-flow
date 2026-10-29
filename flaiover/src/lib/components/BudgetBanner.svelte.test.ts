// S-0534: while flai serve holds every agent start on a provider out of budget, every page names
// the provider, the limit, and until when, one line each, and says when a probe is running.
import { afterEach, describe, expect, it } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import BudgetBanner from './BudgetBanner.svelte';
import { hostFlai } from '$lib/hostflai.svelte';
import { localTime } from '$lib/localtime';

const CONNECTED = { configured: true, connected: true, flai: '1.60.0' };
const CLAUDE = {
	provider: 'claude-code',
	limit: 'usage_limit_reached',
	status: 429,
	since: '2026-10-28T05:30:00Z',
	until: '2026-10-28T09:00:00Z'
};
const CODEX = {
	provider: 'codex',
	since: '2026-10-28T07:55:00Z',
	until: '2026-10-28T08:00:00Z',
	probing: true
};

describe('BudgetBanner', () => {
	afterEach(() => {
		hostFlai.status = null;
		document.body.innerHTML = '';
	});
	const banner = () => document.querySelector('[data-budget-banner]');
	const line = (provider: string) =>
		document
			.querySelector(`[data-budget-provider="${provider}"]`)
			?.textContent?.replace(/\s+/g, ' ')
			.trim();

	it('says nothing until the status is known, nor while no provider is held', () => {
		const c = mount(BudgetBanner, { target: document.body });
		expect(banner()).toBeNull();
		hostFlai.status = CONNECTED;
		flushSync();
		expect(banner()).toBeNull();
		hostFlai.status = { ...CONNECTED, budget: [] };
		flushSync();
		expect(banner()).toBeNull();
		unmount(c);
	});

	it('names each provider held, its limit, and until when in the viewer’s zone', () => {
		const c = mount(BudgetBanner, { target: document.body });
		hostFlai.status = { ...CONNECTED, budget: [CLAUDE] };
		flushSync();
		expect(banner()?.getAttribute('role')).toBe('status');
		expect(line('claude-code')).toBe(
			`claude-code is out of budget (usage_limit_reached) until ${localTime(CLAUDE.until)}; flai serve holds every agent start on it until then.`
		);
		expect(banner()?.textContent).not.toContain('probe');
		unmount(c);
	});

	it('says a probe is running, and leaves out a limit the provider did not name', () => {
		const c = mount(BudgetBanner, { target: document.body });
		hostFlai.status = { ...CONNECTED, budget: [CLAUDE, CODEX] };
		flushSync();
		expect(document.querySelectorAll('[data-budget-provider]')).toHaveLength(2);
		expect(line('codex')).toBe(
			`codex is out of budget until ${localTime(CODEX.until)}; flai serve holds every agent start on it until then. A probe is running to learn whether the limit has reset.`
		);
		expect(line('claude-code')).not.toContain('probe');
		unmount(c);
	});

	it('goes once the hold clears', () => {
		const c = mount(BudgetBanner, { target: document.body });
		hostFlai.status = { ...CONNECTED, budget: [CLAUDE] };
		flushSync();
		expect(banner()).not.toBeNull();
		hostFlai.status = CONNECTED;
		flushSync();
		expect(banner()).toBeNull();
		unmount(c);
	});
});

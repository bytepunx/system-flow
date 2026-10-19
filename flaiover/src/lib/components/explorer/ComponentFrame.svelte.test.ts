// S-0449: the frame takes only the answers about the component it last selected, so a late
// rendered or error about the pick before it changes nothing it shows; it is as tall as the
// harness's size message says, at the width chosen, and its wrapper clips nothing.
import { afterEach, describe, expect, it, vi } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import ComponentFrame from './ComponentFrame.svelte';

const HARNESS = 'http://127.0.0.1:5199/__flai/components/';
const DEV = 'http://127.0.0.1:5199';
const BUTTON = 'src/lib/components/Button.svelte';
const SWATCH = 'src/lib/components/design/Swatch.svelte';

type Fake = { postMessage: ReturnType<typeof vi.fn> };

describe('ComponentFrame (S-0449)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.body.innerHTML = '';
	});

	const q = (testid: string) => document.querySelector<HTMLElement>(`[data-testid="${testid}"]`);
	const frame = () => q('component-frame') as HTMLIFrameElement;
	const box = () => q('component-frame-box')!;

	/** Frames BUTTON, with props the test changes as the page would, and a fake frame window. */
	function show(width: number | null = null) {
		const props = $state({ server: HARNESS, component: BUTTON, theme: 'light' as const, width });
		c = mount(ComponentFrame, { target: document.body, props });
		flushSync();
		const w: Fake = { postMessage: vi.fn() };
		Object.defineProperty(frame(), 'contentWindow', { value: w, configurable: true });
		return { props, w };
	}
	/** The harness in the frame answers, from the dev server's origin. */
	const harness = (from: Fake, m: Record<string, unknown>) => {
		const e = new MessageEvent('message', {
			data: { source: 'flai-components', version: 1, ...m },
			origin: DEV
		});
		Object.defineProperty(e, 'source', { value: from });
		window.dispatchEvent(e);
		flushSync();
	};
	const lastPosted = (w: Fake) => w.postMessage.mock.lastCall![0] as Record<string, unknown>;
	/** Frames BUTTON, then picks SWATCH once the harness is ready. */
	function pickedAgain() {
		const { props, w } = show();
		harness(w, { type: 'ready' });
		props.component = SWATCH;
		flushSync();
		expect(lastPosted(w)).toMatchObject({ type: 'select', component: SWATCH });
		return w;
	}

	it('drops an error about the component picked before', () => {
		const w = pickedAgain();
		harness(w, { type: 'error', component: BUTTON, message: 'label is not defined' });
		expect(q('component-error')).toBeNull();
		expect(box().hidden).toBe(false);
		harness(w, { type: 'error', component: SWATCH, message: 'count is not a number' });
		expect(q('component-error')!.textContent).toContain(`${SWATCH} did not render:`);
		expect(q('component-error')!.textContent).toContain('count is not a number');
		expect(box().hidden).toBe(true);
	});

	it("keeps the current component's error through a rendered about the one before", () => {
		const w = pickedAgain();
		harness(w, { type: 'error', component: SWATCH, message: 'count is not a number' });
		harness(w, { type: 'rendered', component: BUTTON });
		expect(q('component-error')!.textContent).toContain('count is not a number');
		expect(box().hidden).toBe(true);
		harness(w, { type: 'rendered', component: SWATCH });
		expect(q('component-error')).toBeNull();
		expect(box().hidden).toBe(false);
	});

	it('is as tall as the size message says, at the width chosen, and clips none of it', () => {
		const { w } = show(480);
		expect(frame().style.height).toBe('160px');
		harness(w, { type: 'size', width: 1200, height: 517.3 });
		expect(frame().style.height).toBe('518px');
		// the frame is 480 pixels inside the box's 1-pixel border, however wide the content is
		expect(box().style.width).toBe('482px');
		expect(box().className).not.toMatch(/overflow-(hidden|clip)/);
	});
});

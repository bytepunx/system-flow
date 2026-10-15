import { afterEach, describe, expect, it } from 'vitest';
import { flushSync, mount, unmount } from 'svelte';
import PreviewFrame from './PreviewFrame.svelte';
import type { PreviewToken } from '$lib/design/preview';

const colour = (name: string, value: string): PreviewToken => ({
	path: `color.${name}`,
	type: 'color',
	value: { colorSpace: 'srgb', hex: value },
	css: value,
	variables: [{ name: `--color-${name}`, value }]
});

const LIGHT = [colour('background', '#ffffff'), colour('primary', '#2563eb')];
const DARK = [colour('background', '#0f172a'), colour('primary', '#60a5fa')];

describe('PreviewFrame (S-0382, ADR-0149)', () => {
	let c: ReturnType<typeof mount> | undefined;
	afterEach(() => {
		if (c) unmount(c);
		c = undefined;
		document.body.innerHTML = '';
		document.head.innerHTML = '';
		document.documentElement.removeAttribute('data-theme');
		document.documentElement.removeAttribute('style');
	});

	function render(props: { light: PreviewToken[]; dark: PreviewToken[] }) {
		c = mount(PreviewFrame, { target: document.body, props });
		flushSync();
		const frame = document.querySelector<HTMLIFrameElement>('[data-testid="design-preview"]')!;
		return { frame, doc: frame.contentDocument! };
	}

	it('mounts the gallery in its blank frame, once per theme, each panel carrying its theme', () => {
		const { frame, doc } = render({ light: LIGHT, dark: DARK });
		expect(frame.getAttribute('src')).toBeNull();
		expect(frame.getAttribute('title')).toBe('Design system preview');
		// Standards mode, so that a table's text inherits its panel's colour.
		expect(doc.compatMode).toBe('CSS1Compat');
		expect(doc.documentElement.lang).toBe('en');
		const panels = [...doc.querySelectorAll<HTMLElement>('[data-ds-theme]')];
		expect(panels.map((p) => p.dataset.dsTheme)).toEqual(['light', 'dark']);
		for (const p of panels) {
			const gallery = p.querySelector('[data-gallery]')!;
			expect(gallery).not.toBeNull();
			for (const name of ['Primary', 'Secondary', 'Ghost', 'Danger', 'Disabled'])
				expect(
					[...gallery.querySelectorAll('button')].some((b) => b.textContent?.trim() === name)
				).toBe(true);
			expect(gallery.querySelector('form')).not.toBeNull();
			expect(gallery.querySelector('table caption')).not.toBeNull();
			expect(gallery.querySelectorAll('[role="tab"]').length).toBe(3);
		}
		expect(panels[0].querySelector('[data-gallery] p')!.textContent).toBe('Light');
		expect(panels[1].querySelector('[data-gallery] p')!.textContent).toBe('Dark');
		// Two galleries in one document: every label still names its own input.
		const ids = [...doc.querySelectorAll('[id]')].map((e) => e.id);
		expect(new Set(ids).size).toBe(ids.length);
		for (const label of doc.querySelectorAll('label[for]'))
			expect(doc.getElementById(label.getAttribute('for')!)).not.toBeNull();
		// Nothing of the gallery is mounted in the dashboard's own document.
		expect(document.querySelector('[data-gallery]')).toBeNull();
	});

	it("writes each theme's tokens as --ds-* values into the frame's one token style", () => {
		const { doc } = render({ light: LIGHT, dark: DARK });
		const tokens = doc.querySelector('style[data-ds-tokens]')!.textContent!;
		expect(tokens).toContain(
			"[data-ds-theme='light'] {\n\tcolor-scheme: light;\n\t--ds-color-background: #ffffff;\n\t--ds-color-primary: #2563eb;\n}"
		);
		expect(tokens).toContain(
			"[data-ds-theme='dark'] {\n\tcolor-scheme: dark;\n\t--ds-color-background: #0f172a;\n\t--ds-color-primary: #60a5fa;\n}"
		);
		expect(doc.querySelector('style[data-ds-gallery]')).not.toBeNull();
	});

	it('follows a change of tokens at once, rewriting that style and keeping the mounted galleries', () => {
		const props = $state({ light: LIGHT, dark: DARK });
		const { doc } = render(props);
		const gallery = doc.querySelector('[data-ds-theme="light"] [data-gallery]');
		const style = doc.querySelector('style[data-ds-tokens]');
		props.light = [colour('background', '#ffffff'), colour('primary', '#ff0000')];
		flushSync();
		expect(doc.querySelector('style[data-ds-tokens]')).toBe(style);
		expect(style!.textContent).toContain('--ds-color-primary: #ff0000;');
		expect(style!.textContent).not.toContain('#2563eb');
		expect(style!.textContent).toContain('--ds-color-primary: #60a5fa;');
		expect(doc.querySelector('[data-ds-theme="light"] [data-gallery]')).toBe(gallery);
	});

	it('handles events inside the frame: a sample tab is selected when clicked', () => {
		const { doc } = render({ light: LIGHT, dark: DARK });
		const panel = doc.querySelector('[data-ds-theme="dark"]')!;
		const settings = [...panel.querySelectorAll<HTMLButtonElement>('[role="tab"]')].find(
			(b) => b.textContent === 'Settings'
		)!;
		expect(settings.getAttribute('aria-selected')).toBe('false');
		settings.click();
		flushSync();
		expect(settings.getAttribute('aria-selected')).toBe('true');
		expect(panel.querySelector('[data-tab-panel]')!.textContent).toBe('Settings');
		// The light panel's tabs are its own.
		const light = doc.querySelector('[data-ds-theme="light"] [data-tab-panel]')!;
		expect(light.textContent).toBe('Overview');
		// A sample form's submit does not navigate the frame away.
		const form = panel.querySelector('form')!;
		const submit = new (doc.defaultView as Window & typeof globalThis).SubmitEvent('submit', {
			cancelable: true
		});
		form.dispatchEvent(submit);
		expect(submit.defaultPrevented).toBe(true);
	});

	it("keeps the dashboard's theme out of the frame", () => {
		document.documentElement.setAttribute('data-theme', 'dark');
		document.documentElement.style.setProperty('--t-ground', '#191615');
		const sheet = document.createElement('style');
		sheet.textContent = ":root { --t-primary: #054a91; } [data-theme='dark'] * { color: red; }";
		document.head.append(sheet);
		const { doc } = render({ light: LIGHT, dark: DARK });
		expect(doc.documentElement.hasAttribute('data-theme')).toBe(false);
		expect(doc.querySelector('[data-theme]')).toBeNull();
		const styles = [...doc.querySelectorAll('style, link')];
		expect(styles.map((s) => Object.keys((s as HTMLElement).dataset))).toEqual([
			['dsGallery'],
			['dsTokens']
		]);
		for (const s of styles) expect(s.textContent).not.toMatch(/--t-|data-theme/);
		const computed = doc.defaultView!.getComputedStyle(doc.documentElement);
		expect(computed.getPropertyValue('--t-ground')).toBe('');
		expect(computed.getPropertyValue('--t-primary')).toBe('');
		// And the preview's values stay out of the dashboard's document.
		expect(document.querySelector('style[data-ds-tokens], style[data-ds-gallery]')).toBeNull();
	});

	it('leaves out, and says so, a value that cannot go into a stylesheet', () => {
		const bad: PreviewToken = {
			...colour('brand', '#000000'),
			variables: [{ name: '--color-brand', value: 'url(https://example.com/x)' }]
		};
		const { doc } = render({ light: [...LIGHT, bad], dark: DARK });
		expect(doc.querySelector('style[data-ds-tokens]')!.textContent).not.toContain('url(');
		const said = document.querySelector('[data-testid="design-preview-refused"]')!;
		expect(said.textContent).toContain('--color-brand');
		expect(said.textContent).toContain('(light): its value contains url(');
	});

	it('unmounts the galleries when it is destroyed', () => {
		const { doc } = render({ light: LIGHT, dark: DARK });
		const panels = [...doc.querySelectorAll('[data-ds-theme]')];
		unmount(c!);
		c = undefined;
		for (const p of panels) expect(p.querySelector('[data-gallery]')).toBeNull();
	});
});

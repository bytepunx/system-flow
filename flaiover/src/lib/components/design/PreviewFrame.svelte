<!--
	The design preview (S-0382, ADR-0149 § 1): a blank frame, same-origin with the page, that holds
	the gallery's own stylesheet, one <style> of the project's tokens as --ds-* values, and the
	gallery mounted twice, light and dark side by side. A new document inherits nothing, so the
	dashboard's theme (its data-theme, layout.css, and --t-* properties) never reaches the gallery.
	A change of tokens rewrites that one <style>: no remount and no compile.
-->
<script lang="ts">
	import { mount, onMount, unmount } from 'svelte';
	import galleryCss from '$lib/design/gallery.css?inline';
	import {
		PREVIEW_THEMES,
		THEME_ATTRIBUTE,
		previewStyle,
		type PreviewToken
	} from '$lib/design/preview';
	import Gallery from './Gallery.svelte';

	let {
		light,
		dark,
		title = 'Design system preview'
	}: {
		/** The light theme's tokens, as design.get answers them. */
		light: readonly PreviewToken[];
		/** The dark theme's tokens, as design.get answers them. */
		dark: readonly PreviewToken[];
		/** The frame's accessible name. */
		title?: string;
	} = $props();

	const CAPTIONS = { light: 'Light', dark: 'Dark' } as const;

	const style = $derived(previewStyle({ light, dark }));
	let frame: HTMLIFrameElement;
	let tokenStyle = $state<HTMLStyleElement>();
	let height = $state(0);

	onMount(() => {
		const doc = frame.contentDocument;
		if (!doc)
			throw new Error('The design preview could not reach its frame: it is not same-origin.');
		// A blank frame's document is in quirks mode, where a table's text does not inherit its
		// colour; write it again with a doctype so the gallery renders in standards mode.
		doc.open();
		doc.write('<!doctype html><html lang="en"><head></head><body></body></html>');
		doc.close();
		const base = doc.createElement('style');
		base.dataset.dsGallery = '';
		base.textContent = galleryCss;
		const tokens = doc.createElement('style');
		tokens.dataset.dsTokens = '';
		doc.head.append(base, tokens);
		tokenStyle = tokens;

		const panels = doc.createElement('div');
		panels.className = 'ds-panels';
		doc.body.append(panels);
		const galleries = PREVIEW_THEMES.map((theme) => {
			const panel = doc.createElement('div');
			panel.setAttribute(THEME_ATTRIBUTE, theme);
			panels.append(panel);
			return mount(Gallery, { target: panel, props: { caption: CAPTIONS[theme] } });
		});

		// Size the frame to the panels, and its own border, with the frame's own ResizeObserver,
		// which sees its document.
		const fit = () =>
			(height =
				Math.ceil(panels.getBoundingClientRect().height) + frame.offsetHeight - frame.clientHeight);
		fit();
		const Observer = (doc.defaultView as (Window & typeof globalThis) | null)?.ResizeObserver;
		const observer = Observer ? new Observer(fit) : undefined;
		observer?.observe(panels);

		return () => {
			observer?.disconnect();
			for (const g of galleries) unmount(g);
		};
	});

	$effect(() => {
		if (tokenStyle) tokenStyle.textContent = style.css;
	});
</script>

<iframe
	bind:this={frame}
	{title}
	class="block w-full rounded-md border border-line"
	style:height={height ? `${height}px` : undefined}
	data-testid="design-preview"
></iframe>
{#if style.refused.length}
	<div class="mt-2 text-sm text-danger" role="status" data-testid="design-preview-refused">
		<p>The preview left out values it cannot write into a stylesheet:</p>
		<ul class="list-disc pl-5">
			{#each style.refused as r (`${r.theme} ${r.name}`)}
				<li><code>{r.name}</code> ({r.theme}): {r.reason}</li>
			{/each}
		</ul>
	</div>
{/if}

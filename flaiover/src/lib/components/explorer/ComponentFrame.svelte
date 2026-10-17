<!--
	One theme's frame of the component harness (S-0439, ADR-0163): the harness page from a
	sub-project's dev server, told the dashboard's origin, sent the theme and the component picked
	by message, and as tall as the harness says its content is. The frame loads once per dev server;
	a later pick or theme is a message, so neither reloads it. A harness error is shown in the
	frame's place, with the frame kept loaded but hidden, so that the next hot update's render
	brings it back.
-->
<script lang="ts">
	import { untrack } from 'svelte';
	import {
		connect,
		harnessUrl,
		type Client,
		type RenderError,
		type Theme
	} from '$lib/explorer/protocol';

	let {
		server,
		component,
		theme,
		width = null
	}: {
		/** The harness page, as components.status answers it. */
		server: string;
		/** The component's path in its sub-project. */
		component: string;
		/** The theme the harness document is given. */
		theme: Theme;
		/** The frame's width in pixels; null for the width it is laid out in. */
		width?: number | null;
	} = $props();

	const CAPTIONS = { light: 'Light', dark: 'Dark' } as const;
	// The frame's height until the harness says how tall its content is.
	const UNSIZED = 160;

	let iframe: HTMLIFrameElement | undefined = $state();
	let client = $state.raw<Client | null>(null);
	let height = $state(0);
	let error = $state.raw<RenderError | null>(null);

	// The component picked when the frame loads renders first; later picks are select messages.
	const first = () => untrack(() => component);
	const src = $derived(harnessUrl(server, location.origin, first()));

	$effect(() => {
		const c = connect(() => iframe?.contentWindow, server, {
			rendered: () => (error = null),
			error: (m) => (error = m),
			size: (m) => (height = Math.ceil(m.height))
		});
		client = c;
		return () => {
			c.close();
			client = null;
		};
	});

	$effect(() => {
		client?.theme(theme);
	});

	$effect(() => {
		if (!client) return;
		// an error is the component's it was raised for: the next one starts clear
		error = null;
		client.select(component);
	});

	// file:line:column, as far as the harness knows them
	const place = $derived(
		[error?.file, error?.line, error?.column].filter((v) => v !== undefined).join(':')
	);
	const caption = $derived(width ? `${CAPTIONS[theme]} · ${width}px` : CAPTIONS[theme]);
</script>

<figure class="min-w-0" data-testid="component-frame-{theme}">
	<figcaption class="mb-1 text-xs text-muted">{caption}</figcaption>
	{#if error}
		<div
			class="rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
			style:width={width ? `${width}px` : undefined}
			role="alert"
			data-testid="component-error"
		>
			<p class="font-medium">{component} did not render:</p>
			<pre class="mt-1 font-mono text-xs whitespace-pre-wrap">{error.message}</pre>
			{#if place}
				<p class="mt-1 font-mono text-xs" data-testid="component-error-place">at {place}</p>
			{/if}
		</div>
	{/if}
	<!-- the border is the wrapper's, so that the frame's height is the harness content's own -->
	<div
		class="overflow-hidden rounded-md border border-line"
		style:width={width ? `${width + 2}px` : undefined}
		hidden={!!error}
		data-testid="component-frame-box"
	>
		<iframe
			bind:this={iframe}
			{src}
			title="{component} ({theme})"
			class="block w-full border-0"
			style:height="{height || UNSIZED}px"
			data-testid="component-frame"
		></iframe>
	</div>
</figure>

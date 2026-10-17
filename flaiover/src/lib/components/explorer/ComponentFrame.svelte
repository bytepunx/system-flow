<!--
	One theme's frame of the component harness (S-0439, ADR-0163): the harness page from a
	sub-project's dev server, told the dashboard's origin, sent the theme and the component picked
	by message, and as tall as the harness says its content is. The frame loads once per dev server;
	a later pick or theme is a message, so neither reloads it. A pick is sent with its props, and a
	change of props alone as a props message, which the harness gives the mounted component without
	remounting it (S-0440); a function prop's calls come back as events. A harness error, or a
	required prop without a value, is shown in the frame's place, with the frame kept loaded but
	hidden, so that the next hot update's render, or the value, brings it back. Tried token values
	are sent as they change, and an empty tokens message once they are gone (S-0441); the client
	sends the last of each message again on every ready, so a reloaded harness keeps them.
-->
<script lang="ts">
	import { untrack } from 'svelte';
	import {
		connect,
		harnessUrl,
		type Called,
		type Client,
		type Props,
		type RenderError,
		type Theme,
		type Tokens
	} from '$lib/explorer/protocol';

	let {
		server,
		component,
		theme,
		width = null,
		props = {},
		functions = [],
		snippets = [],
		missing = [],
		tokens = null,
		onevent
	}: {
		/** The harness page, as components.status answers it. */
		server: string;
		/** The component's path in its sub-project. */
		component: string;
		/** The theme the harness document is given. */
		theme: Theme;
		/** The frame's width in pixels; null for the width it is laid out in. */
		width?: number | null;
		/** Every prop the component is given; one left out takes its default. */
		props?: Props;
		/** The function props the harness gives a logger. */
		functions?: string[];
		/** The props whose text the harness renders as a snippet. */
		snippets?: string[];
		/** The required props without a value, which keep the component from being shown. */
		missing?: string[];
		/** The tried token values laid over the stylesheet's, light and dark; null for none. */
		tokens?: { light: Tokens; dark: Tokens } | null;
		/** Called with each call of a logged function prop of the component framed. */
		onevent?: (m: Called) => void;
	} = $props();

	const CAPTIONS = { light: 'Light', dark: 'Dark' } as const;
	// The frame's height until the harness says how tall its content is.
	const UNSIZED = 160;

	let iframe: HTMLIFrameElement | undefined = $state();
	let client = $state.raw<Client | null>(null);
	let height = $state(0);
	let error = $state.raw<RenderError | null>(null);
	// The component the client last selected, and the props it was last sent, as JSON: a catalogue
	// read again gives equal props anew, which are not sent again.
	let selected: string | null = null;
	let sent = '';
	// The tokens last sent, as JSON; '' when none are laid over the stylesheet.
	let tinted = '';

	// The component picked when the frame loads renders first; later picks are select messages.
	const first = () => untrack(() => component);
	const src = $derived(harnessUrl(server, location.origin, first()));

	$effect(() => {
		const c = connect(() => iframe?.contentWindow, server, {
			rendered: () => (error = null),
			error: (m) => (error = m),
			size: (m) => (height = Math.ceil(m.height)),
			// a call answered after another pick is the component's before it
			event: (m) => {
				if (m.component === component) onevent?.(m);
			}
		});
		selected = null;
		tinted = '';
		client = c;
		return () => {
			c.close();
			client = null;
		};
	});

	$effect(() => {
		client?.theme(theme);
	});

	// A new component is a select, carrying its props; a change of props alone is a props message.
	$effect(() => {
		if (!client) return;
		const given = JSON.stringify([props, functions, snippets]);
		if (component !== selected) {
			selected = component;
			sent = given;
			// an error is the component's it was raised for: the next one starts clear
			error = null;
			client.select(component, props, functions, snippets);
		} else if (given !== sent) {
			sent = given;
			client.props(props, functions, snippets);
		}
	});

	// Tried tokens as they change; once they are gone, an empty message removes them.
	$effect(() => {
		if (!client) return;
		const given = tokens ? JSON.stringify(tokens) : '';
		if (given === tinted) return;
		tinted = given;
		if (tokens) client.tokens(tokens.light, tokens.dark);
		else client.tokens({}, {});
	});

	// file:line:column, as far as the harness knows them
	const place = $derived(
		[error?.file, error?.line, error?.column].filter((v) => v !== undefined).join(':')
	);
	const caption = $derived(width ? `${CAPTIONS[theme]} · ${width}px` : CAPTIONS[theme]);
</script>

<figure class="min-w-0" data-testid="component-frame-{theme}">
	<figcaption class="mb-1 text-xs text-muted">{caption}</figcaption>
	{#if missing.length}
		<div
			class="rounded border border-danger bg-danger-soft p-3 text-sm text-danger"
			style:width={width ? `${width}px` : undefined}
			role="alert"
			data-testid="component-missing"
		>
			<p class="font-medium">{component} needs a value for {missing.join(', ')}.</p>
			<p class="mt-1 text-xs">The frame shows it once each has one in the props panel.</p>
		</div>
	{:else if error}
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
		hidden={!!error || missing.length > 0}
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

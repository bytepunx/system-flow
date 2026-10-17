// The page's side of the component harness's protocol (S-0439, ADR-0163): what the Components
// page posts to the harness it frames from a sub-project's dev server, and what it takes back.
// Every message is {source: 'flai-components', version: 1, type, ...}; the page posts only to the
// dev server's origin and takes answers only from that origin and its own frame.

export const SOURCE = 'flai-components';
export const VERSION = 1;

export type Theme = 'light' | 'dark';
/** CSS custom properties by name, with or without their leading `--`. */
export type Tokens = Record<string, string | number>;
export type Props = Record<string, unknown>;

/** What the page posts to the harness. */
export type PageMessage =
	| { type: 'select'; component: string; props?: Props; functions?: string[] }
	| { type: 'props'; props: Props; functions?: string[] }
	| { type: 'theme'; theme: Theme }
	| { type: 'tokens'; light: Tokens; dark: Tokens };

export type Ready = { type: 'ready'; component?: string };
export type Rendered = { type: 'rendered'; component: string };
export type RenderError = {
	type: 'error';
	component: string;
	message: string;
	file?: string;
	line?: number;
	column?: number;
};
export type Called = { type: 'event'; component: string; name: string; args: unknown[] };
export type Size = { type: 'size'; width: number; height: number };
/** What the harness posts to the page. */
export type HarnessMessage = Ready | Rendered | RenderError | Called | Size;

/** What the page does with each of the harness's answers. */
export type Handlers = {
	ready?: (m: Ready) => void;
	rendered?: (m: Rendered) => void;
	error?: (m: RenderError) => void;
	event?: (m: Called) => void;
	size?: (m: Size) => void;
};

/** The frame's window, as far as the client posts to it. */
export type FrameWindow = { postMessage(message: unknown, targetOrigin: string): void };
/** The page's window, as far as the client listens on it. */
export type PageWindow = {
	addEventListener(type: 'message', listener: (e: MessageEvent) => void): void;
	removeEventListener(type: 'message', listener: (e: MessageEvent) => void): void;
};

/** A connection to one harness frame. */
export type Client = {
	/** Mounts the component at path, with props and a logger for each name in functions. */
	select(component: string, props?: Props, functions?: string[]): void;
	/** Gives the mounted component props without remounting it. */
	props(props: Props, functions?: string[]): void;
	/** Sets the harness document's theme. */
	theme(theme: Theme): void;
	/** Lays CSS custom properties over the stylesheet's; empty objects clear them. */
	tokens(light: Tokens, dark: Tokens): void;
	/** Stops listening and sending. */
	close(): void;
};

/** The harness frame's src: the harness page at url, as components.status answers it, told the
 * page's origin and, when given, the component to render first. */
export function harnessUrl(url: string, pageOrigin: string, component?: string): string {
	const u = new URL(url);
	u.searchParams.set('origin', new URL(pageOrigin).origin);
	if (component) u.searchParams.set('c', component);
	return u.href;
}

const isString = (v: unknown): v is string => typeof v === 'string';
const isNumber = (v: unknown): v is number => typeof v === 'number';
const optional = <T>(v: unknown, is: (v: unknown) => v is T) => v === undefined || is(v);

// The fields each answer must carry, by its type.
const shapes: Record<HarnessMessage['type'], (m: Record<string, unknown>) => boolean> = {
	ready: (m) => optional(m.component, isString),
	rendered: (m) => isString(m.component),
	error: (m) =>
		isString(m.component) &&
		isString(m.message) &&
		optional(m.file, isString) &&
		optional(m.line, isNumber) &&
		optional(m.column, isNumber),
	event: (m) => isString(m.component) && isString(m.name) && Array.isArray(m.args),
	size: (m) => isNumber(m.width) && isNumber(m.height)
};

/** data as a harness message, or null when it is not one of this protocol's. */
export function harnessMessage(data: unknown): HarnessMessage | null {
	if (typeof data !== 'object' || data === null) return null;
	const m = data as Record<string, unknown>;
	if (m.source !== SOURCE || m.version !== VERSION) return null;
	if (typeof m.type !== 'string' || !Object.hasOwn(shapes, m.type)) return null;
	return shapes[m.type as HarnessMessage['type']](m) ? (m as HarnessMessage) : null;
}

// Theme and tokens go before the component mounts, and its props after it.
const order: PageMessage['type'][] = ['theme', 'tokens', 'select', 'props'];

/** Talks to the harness in the frame whose window frame answers, served from server (the dev
 * server's URL or origin), answering handlers; page is the window it listens on. */
export function connect(
	frame: () => FrameWindow | null | undefined,
	server: string,
	handlers: Handlers,
	page: PageWindow = window
): Client {
	const origin = new URL(server).origin;
	// The last of each kind, sent on ready and again on each later one, as the harness reloads.
	const last = new Map<PageMessage['type'], PageMessage>();
	let ready = false;
	let closed = false;

	function post(m: PageMessage) {
		frame()?.postMessage({ source: SOURCE, version: VERSION, ...m }, origin);
	}

	function send(m: PageMessage) {
		if (closed) return;
		// A select carries its component's props, so props given for the one before it are dropped.
		if (m.type === 'select') last.delete('props');
		last.set(m.type, m);
		if (ready) post(m);
	}

	function receive(e: MessageEvent) {
		if (e.origin !== origin) return;
		const w = frame();
		if (!w || e.source !== w) return;
		const m = harnessMessage(e.data);
		if (!m) return;
		switch (m.type) {
			case 'ready':
				ready = true;
				for (const type of order) {
					const queued = last.get(type);
					if (queued) post(queued);
				}
				handlers.ready?.(m);
				break;
			case 'rendered':
				handlers.rendered?.(m);
				break;
			case 'error':
				handlers.error?.(m);
				break;
			case 'event':
				handlers.event?.(m);
				break;
			case 'size':
				handlers.size?.(m);
				break;
		}
	}

	page.addEventListener('message', receive);
	return {
		select: (component, props, functions) => send({ type: 'select', component, props, functions }),
		props: (props, functions) => send({ type: 'props', props, functions }),
		theme: (theme) => send({ type: 'theme', theme }),
		tokens: (light, dark) => send({ type: 'tokens', light, dark }),
		close() {
			closed = true;
			last.clear();
			page.removeEventListener('message', receive);
		}
	};
}

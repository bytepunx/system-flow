// Client-side markdown rendering for the documentation explorer. Runs in
// the browser from raw markdown served by /api/docs/file. Mermaid and shiki
// are loaded lazily by enhance() after the HTML is in the DOM, and each
// repository markdown file the HTML links or names gets an edit control (S-0378).
import MarkdownIt from 'markdown-it';
import anchor from 'markdown-it-anchor';
import taskLists from 'markdown-it-task-lists';
import type { Attachment } from 'svelte/attachments';
import { editLink, fileRef, type EditLink, type LinkOptions } from './filelinks';

/** The id a heading gets, markdown-it-anchor's own rule: what a link to a document's heading ends with. */
export const slug = (heading: string): string =>
	encodeURIComponent(heading.trim().toLowerCase().replace(/\s+/g, '-'));

const md = new MarkdownIt({ html: true, linkify: true, typographer: false });
md.use(anchor, {
	slugify: slug,
	permalink: anchor.permalink.headerLink({ safariReaderFix: true })
});
md.use(taskLists, { enabled: false, label: true });

// Fenced mermaid blocks become <pre class="mermaid"> for the client renderer;
// other fences keep their language on the <code> element for shiki.
const defaultFence = md.renderer.rules.fence!;
md.renderer.rules.fence = (tokens, idx, options, env, self) => {
	const token = tokens[idx];
	const lang = token.info.trim().split(/\s+/)[0];
	if (lang === 'mermaid') {
		return `<pre class="mermaid">${md.utils.escapeHtml(token.content)}</pre>\n`;
	}
	return defaultFence(tokens, idx, options, env, self);
};

/** Resolve a relative link against the directory of the current document. */
export function resolveRelative(docPath: string, href: string): string {
	const base = docPath.split('/').slice(0, -1);
	const parts = href.split('/');
	for (const p of parts) {
		if (p === '..') base.pop();
		else if (p !== '.' && p !== '') base.push(p);
	}
	return base.join('/');
}

/**
 * Rewrite links: relative .md links open in the explorer, other relative
 * links stay repository-relative under /docs, absolute URLs are untouched.
 */
export function rewriteLink(docPath: string, href: string): string {
	if (/^[a-z]+:/i.test(href) || href.startsWith('#') || href.startsWith('/')) return href;
	const [target, hash] = href.split('#');
	const resolved = resolveRelative(docPath, target);
	const withHash = hash ? `#${hash}` : '';
	return `/docs/${resolved}${withHash}`;
}

/** Render markdown to HTML with links rewritten for docPath. */
export function render(markdown: string, docPath: string): string {
	const env = { docPath };
	const defaultLink =
		md.renderer.rules.link_open ??
		((tokens, idx, options, _env, self) => self.renderToken(tokens, idx, options));
	md.renderer.rules.link_open = (tokens, idx, options, e, self) => {
		const href = tokens[idx].attrGet('href');
		if (href)
			tokens[idx].attrSet('href', rewriteLink((e as { docPath: string }).docPath, String(href)));
		return defaultLink(tokens, idx, options, e, self);
	};
	return md.render(markdown, env);
}

/** The repository path of a rendered link to a markdown file, which render() pointed at the explorer. */
function linkedPath(a: HTMLAnchorElement): string | null {
	const m = /^\/docs\/([^#?]+\.md)(?:#.*)?$/.exec(a.getAttribute('href') ?? '');
	return m ? decodeURIComponent(m[1]) : null;
}

// A pencil, drawn rather than written, so that a control adds nothing to the text around it.
const PENCIL =
	'<svg viewBox="0 0 16 16" width="12" height="12" fill="none" stroke="currentColor" stroke-width="1.5" aria-hidden="true"><path d="M11 2.5 13.5 5 5 13.5H2.5V11z"/></svg>';

/** Put a control that opens the editor right after an element, once. */
function placeControl(el: Element, edit: EditLink): void {
	const next = el.nextSibling;
	if (next?.nodeType === 1 && (next as Element).hasAttribute('data-edit')) return;
	const button = el.ownerDocument.createElement('button');
	button.type = 'button';
	button.className = 'ml-1 rounded px-1 text-xs text-muted hover:text-ink';
	button.dataset.edit = edit.path;
	if (edit.line) button.dataset.editLine = String(edit.line);
	button.setAttribute('aria-label', edit.label);
	button.title = edit.label;
	button.innerHTML = PENCIL;
	button.addEventListener('click', edit.open);
	el.after(button);
}

/**
 * Give each link to a repository markdown file, and each inline code span that names one with an
 * optional `:line`, an edit control beside it that opens the editor at that line; following the
 * link still shows the rendered document. Files the editor cannot edit get none.
 */
export function addEditControls(root: HTMLElement, opts: LinkOptions = {}): void {
	for (const a of root.querySelectorAll<HTMLAnchorElement>('a[href]')) {
		const path = linkedPath(a);
		const edit = path && editLink({ path }, opts);
		if (edit) placeControl(a, edit);
	}
	for (const code of root.querySelectorAll<HTMLElement>('code')) {
		if (code.closest('pre, a')) continue;
		const ref = fileRef(code.textContent ?? '', opts.layout);
		const edit = ref && editLink(ref, opts);
		if (edit) placeControl(code, edit);
	}
}

/** An attachment that renders markdown into its element, with edit controls, for markdown shown without enhance(). */
export function markdownInto(
	markdown: string,
	docPath: string,
	opts: LinkOptions = {}
): Attachment<HTMLElement> {
	return (node) => {
		node.innerHTML = render(markdown, docPath);
		addEditControls(node, opts);
	};
}

/**
 * After the HTML is in the DOM: add edit controls, unless `edit` is false, then render mermaid
 * diagrams and highlight code.
 */
export async function enhance(
	root: HTMLElement,
	dark: boolean,
	edit: LinkOptions | false = {}
): Promise<void> {
	if (edit) addEditControls(root, edit);
	const diagrams = root.querySelectorAll<HTMLElement>('pre.mermaid');
	if (diagrams.length) {
		const mermaid = (await import('mermaid')).default;
		mermaid.initialize({
			startOnLoad: false,
			theme: dark ? 'dark' : 'default',
			securityLevel: 'strict'
		});
		await mermaid.run({ nodes: diagrams });
	}
	const blocks = [...root.querySelectorAll<HTMLElement>('pre > code[class*="language-"]')];
	if (blocks.length) {
		const { codeToHtml } = await import('shiki');
		for (const code of blocks) {
			const lang = code.className.replace(/.*language-(\S+).*/, '$1');
			try {
				const html = await codeToHtml(code.textContent ?? '', {
					lang,
					theme: dark ? 'github-dark' : 'github-light'
				});
				code.parentElement!.outerHTML = html;
			} catch {
				// unknown language: leave the plain block
			}
		}
	}
}

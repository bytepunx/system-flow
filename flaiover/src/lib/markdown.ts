// Client-side markdown rendering for the documentation explorer. Runs in
// the browser from raw markdown served by /api/docs/file. Mermaid and shiki
// are loaded lazily by enhance() after the HTML is in the DOM.
import MarkdownIt from 'markdown-it';
import anchor from 'markdown-it-anchor';
import taskLists from 'markdown-it-task-lists';

const md = new MarkdownIt({ html: true, linkify: true, typographer: false });
md.use(anchor, { permalink: anchor.permalink.headerLink({ safariReaderFix: true }) });
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

/** After the HTML is in the DOM: render mermaid diagrams and highlight code. */
export async function enhance(root: HTMLElement, dark: boolean): Promise<void> {
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

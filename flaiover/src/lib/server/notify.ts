// Webhook for new inbox entries (S-0042): when system-flow.yaml sets
// dashboard.notify_url, POST each entry that appears after the server started.
// One attempt, a short timeout, a warning on failure. The token and file
// contents are never sent: an entry is its key, kind, title, link, and time.
import type { Repo } from './repo';
import { inbox, type InboxEntry } from './inbox';
import { log } from './log';

export type NotifyBody = {
	project: string;
	entry: Pick<InboxEntry, 'key' | 'kind' | 'title' | 'href' | 'at'>;
};

export type Notifier = {
	stop: () => void;
	/** resolves when a pending check has finished; for tests */ idle: () => Promise<void>;
};

const TIMEOUT_MS = 5000;
const DEBOUNCE_MS = 500;

/** The webhook URL from the manifest, or null when unset or not http(s). */
export async function notifyUrl(repo: Repo): Promise<string | null> {
	const m = (await repo.manifest()) as { dashboard?: { notify_url?: unknown } };
	const url = m.dashboard?.notify_url;
	if (typeof url !== 'string' || !url.trim()) return null;
	try {
		const u = new URL(url.trim());
		return u.protocol === 'http:' || u.protocol === 'https:' ? u.toString() : null;
	} catch {
		return null;
	}
}

async function post(url: string, body: NotifyBody): Promise<void> {
	const controller = new AbortController();
	const timer = setTimeout(() => controller.abort(), TIMEOUT_MS);
	try {
		const r = await fetch(url, {
			method: 'POST',
			headers: { 'content-type': 'application/json', 'user-agent': 'flaiover' },
			body: JSON.stringify(body),
			signal: controller.signal
		});
		if (!r.ok) throw new Error(`webhook answered ${r.status}`);
	} finally {
		clearTimeout(timer);
	}
}

/**
 * Watch the repository and POST new inbox entries. Returns null, and does
 * nothing, when the manifest sets no notify_url. Entries present when the
 * notifier starts are not news.
 */
export async function startNotifier(
	repo: Repo,
	debounceMs = DEBOUNCE_MS
): Promise<Notifier | null> {
	const url = await notifyUrl(repo);
	if (!url) return null;
	const project = String(((await repo.manifest()) as { name?: unknown }).name ?? '');
	let known = new Set((await inbox(repo)).entries.map((e) => e.key));
	let timer: ReturnType<typeof setTimeout> | null = null;
	let running: Promise<void> = Promise.resolve();

	const check = async () => {
		try {
			const entries = (await inbox(repo)).entries;
			const fresh = entries.filter((e) => !known.has(e.key));
			known = new Set(entries.map((e) => e.key));
			for (const e of fresh) {
				try {
					await post(url, {
						project,
						entry: { key: e.key, kind: e.kind, title: e.title, href: e.href, at: e.at }
					});
				} catch (err) {
					// the URL may carry a secret in its path or query: log the host only
					log().warn(
						{ component: 'notify', host: new URL(url).host, key: e.key, err: String(err) },
						'inbox webhook failed'
					);
				}
			}
		} catch (err) {
			log().warn({ component: 'notify', err: String(err) }, 'inbox webhook check failed');
		}
	};
	const onChange = () => {
		if (timer) clearTimeout(timer);
		timer = setTimeout(() => {
			timer = null;
			running = running.then(check);
		}, debounceMs);
	};
	repo.on('change', onChange);
	log().info({ component: 'notify', host: new URL(url).host }, 'inbox webhook enabled');
	return {
		stop: () => {
			repo.off('change', onChange);
			if (timer) clearTimeout(timer);
		},
		idle: async () => {
			while (timer) await new Promise((r) => setTimeout(r, 10));
			await running;
		}
	};
}

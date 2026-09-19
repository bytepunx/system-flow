// Who is working on what, derived from the narratives (S-0042). Nothing is
// written for presence: an agent that stops logging simply ages out.
import { readdir, readFile } from 'node:fs/promises';
import { join } from 'node:path';
import { splitFrontMatter, type Item, type Repo } from './repo';

export type StreamActivity = {
	stream: string;
	title: string;
	agent: string;
	session: string;
	updated: string;
	age_seconds: number;
	status: string; // the story's state; "unknown" when the story is not found
	blocked: boolean; // the story or a task of it has an open blocked interval
	task?: { id: string; title: string }; // the task in progress
	last_log?: { at: string; text: string };
	path: string;
};

const LOG_HEADING = /^### (\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z)/;

/** The last entry of a narrative's `## Log`: its timestamp and the lines under it. */
export function lastLogEntry(body: string): { at: string; text: string } | undefined {
	const lines = body.split('\n');
	const start = lines.findIndex((l) => l.trim() === '## Log');
	let last: { at: string; text: string[] } | undefined;
	for (const line of lines.slice(start + 1)) {
		if (/^##\s/.test(line)) break;
		const m = LOG_HEADING.exec(line);
		if (m) last = { at: m[1], text: [] };
		else last?.text.push(line);
	}
	return last ? { at: last.at, text: last.text.join('\n').trim() } : undefined;
}

const stamp = (v: unknown): string =>
	v instanceof Date ? v.toISOString().replace(/\.\d{3}Z$/, 'Z') : String(v ?? '');

export async function activity(
	repo: Repo,
	now = new Date()
): Promise<{ streams: StreamActivity[] }> {
	const layout = await repo.layout();
	const dir = `${layout.wip}/agents`;
	const abs = repo.resolveInside(dir);
	const names = (await readdir(abs).catch(() => [] as string[])).filter(
		(n) => n.endsWith('.md') && n !== 'index.md' && n !== 'README.md'
	);
	const items = (await repo.items()).filter((it) => !it.archived);
	const byId = new Map<string, Item>(items.map((it) => [it.id, it]));
	const open = (it: Item | undefined) => (it?.blocked ?? []).some((b) => !b.until);
	const streams: StreamActivity[] = [];
	for (const name of names) {
		const { frontMatter, body } = splitFrontMatter(await readFile(join(abs, name), 'utf8'));
		if (!frontMatter) continue;
		const fm = frontMatter as Record<string, unknown>;
		const stream = String(fm.stream ?? name.replace(/\.md$/, ''));
		const story = byId.get(stream);
		const tasks = items.filter((it) => it.type === 'task' && it.parent === stream);
		const current = tasks.find((t) => t.status === 'in-progress');
		const updated = stamp(fm.updated);
		streams.push({
			stream,
			title: String(fm.title ?? story?.title ?? ''),
			agent: String(fm.agent ?? ''),
			session: String(fm.session ?? ''),
			updated,
			age_seconds: Math.max(0, Math.round((now.getTime() - Date.parse(updated)) / 1000)) || 0,
			status: story?.status ?? 'unknown',
			blocked: open(story) || tasks.some(open),
			task: current ? { id: current.id, title: current.title } : undefined,
			last_log: lastLogEntry(body),
			path: `${dir}/${name}`
		});
	}
	streams.sort((a, b) => b.updated.localeCompare(a.updated) || a.stream.localeCompare(b.stream));
	return { streams };
}

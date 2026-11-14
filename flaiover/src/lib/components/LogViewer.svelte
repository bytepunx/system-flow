<script lang="ts">
	// One log on the host by its path (S-0785): the path, then the log's contents in a LogFeed that
	// reads them through GET /api/log, flai's log.read, from the start, then from where the last read
	// stopped about every second, so a growing log is followed live. Another path starts the log
	// over. A log that is not there says so, and shows once it is written; a path flai does not read,
	// or a host flai too old to read one, shows flai's words.
	import LogFeed, { type LogPart } from '$lib/components/LogFeed.svelte';
	import { api } from '$lib/api';
	import type { LogRead } from '$lib/logs';

	let {
		path,
		open = $bindable(true)
	}: {
		/** The log's path as flai names it: absolute, ~/-prefixed, or from the project's root. */
		path: string;
		/** Whether the log's contents are open; the reader opens and closes them. */
		open?: boolean;
	} = $props();

	/** Reads the log from an offset; throws flai's refusal word for word. */
	async function read(offset: number): Promise<LogPart> {
		const r = await api(`/api/log?path=${encodeURIComponent(path)}&offset=${offset}`);
		const body = await r.json().catch(() => ({}));
		if (!r.ok)
			throw new Error(body.error || r.statusText || `the log could not be read (${r.status})`);
		return body as LogRead;
	}
</script>

<section class="text-sm" aria-label="Log" data-testid="log-viewer" data-path={path}>
	<p>
		<span class="text-muted">Log</span>
		<code class="font-mono break-all" data-testid="log-viewer-path">{path}</code>
	</p>
	<div class="mt-2">
		<LogFeed step={path} label="Contents" follow {read} bind:open />
	</div>
</section>

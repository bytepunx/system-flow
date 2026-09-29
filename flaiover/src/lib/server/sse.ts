// Server-sent events over the Repo watcher. The stream tears itself down
// on client cancel or request abort, whichever comes first, and never
// touches the controller after it is closed (a closed controller throws,
// and an unhandled throw here would take the whole server down).
import type { Repo } from './repo';

export function sseStream(
	r: Repo,
	signal: AbortSignal,
	pingMs = 25000
): ReadableStream<Uint8Array> {
	const enc = new TextEncoder();
	let closed = false;
	let cleanup = () => {};
	return new ReadableStream<Uint8Array>({
		start(controller) {
			const safe = (chunk: string) => {
				if (closed) return;
				try {
					controller.enqueue(enc.encode(chunk));
				} catch {
					finish();
				}
			};
			const send = (path: string, kind?: string) =>
				safe(`event: change\ndata: ${JSON.stringify({ path, kind })}\n\n`);
			const agent = (story: string) => safe(`event: agent\ndata: ${JSON.stringify({ story })}\n\n`);
			const ping = setInterval(() => safe(': ping\n\n'), pingMs);
			const finish = () => {
				if (closed) return;
				closed = true;
				clearInterval(ping);
				r.off('change', send);
				r.off('agent', agent);
				signal.removeEventListener('abort', finish);
				try {
					controller.close();
				} catch {
					// already closed by cancel
				}
			};
			cleanup = finish;
			r.on('change', send);
			r.on('agent', agent);
			signal.addEventListener('abort', finish);
			safe('event: ready\ndata: {}\n\n');
		},
		cancel() {
			cleanup();
		}
	});
}

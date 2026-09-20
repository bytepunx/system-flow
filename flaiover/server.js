// The server entry of the image: adapter-node's handler on a plain HTTP
// server, plus the one thing SvelteKit cannot do itself, the WebSocket
// upgrade for /agent that flai on the host dials (ADR-0029). The hub lives
// in the application (src/lib/server/agent.ts) and is found through a global
// it sets when the server initialises.
import http from 'node:http';
import { noteScheme, schemeHeader } from './scheme.js';

// adapter-node reads this when its handler is loaded, so it is set first and
// the handler is imported after it. Without it every request is taken for
// https (S-0083). ORIGIN, when an operator sets it, still overrides both.
process.env.PROTOCOL_HEADER = schemeHeader();
const { handler } = await import('./build/handler.js');

const host = process.env.HOST ?? '0.0.0.0';
const port = Number(process.env.PORT ?? 3000);
const shutdownMs = Number(process.env.SHUTDOWN_TIMEOUT ?? 5) * 1000;

const server = http.createServer((req, res) => {
	noteScheme(req);
	handler(req, res, () => {
		res.statusCode = 404;
		res.end('Not found');
	});
});

server.on('upgrade', (req, socket, head) => {
	const upgrade = globalThis.__flaioverAgentUpgrade;
	if (upgrade) upgrade(req, socket, head);
	else socket.destroy();
});

server.listen(port, host, () => {
	console.log(
		JSON.stringify({
			level: 'INFO',
			ts: new Date().toISOString(),
			service: 'flaiover',
			component: 'server',
			msg: `listening on http://${host}:${port}`
		})
	);
});

// Stop listening, give requests a moment, then close what is left: an open
// event stream must not keep a server that no longer listens alive (I-0025).
let stopping = false;
function shutdown() {
	if (stopping) return;
	stopping = true;
	server.close(() => process.exit(0));
	server.closeIdleConnections();
	setTimeout(() => {
		server.closeAllConnections();
		process.exit(0);
	}, shutdownMs).unref();
}
process.on('SIGTERM', shutdown);
process.on('SIGINT', shutdown);

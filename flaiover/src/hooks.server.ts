// Server hooks: one span, one log line, and the golden-signal metrics per
// request (design/conventions/logging.md and telemetry.md). Errors are
// logged once, here, with the trace id.
import type { Handle, HandleServerError, ServerInit } from '@sveltejs/kit';
import { SpanStatusCode } from '@opentelemetry/api';
import { log } from '$lib/server/log';
import {
	inFlight,
	requestDuration,
	requestsTotal,
	routeLabel,
	statusClass,
	version,
	commit
} from '$lib/server/metrics';
import { startTracing, traceId, withRequestSpan } from '$lib/server/otel';
import { authenticate, decide, initAuth } from '$lib/server/auth';
import { repo } from '$lib/server/repo';
import { startNotifier } from '$lib/server/notify';
import { projectIdentity, setIdentityHeaders } from '$lib/server/project';
import { agent, exposeAgentUpgrade } from '$lib/server/agent';
import { redirect, json } from '@sveltejs/kit';

export const init: ServerInit = async () => {
	const auth = initAuth();
	const tracing = await startTracing();
	// The channel flai on the host dials (ADR-0029): the server entry hands /agent upgrades to the hub.
	exposeAgentUpgrade();
	const channel = agent().status().configured;
	// flai on the host says which files changed; the caches, /api/events, and the notifier listen (S-0073).
	const r = repo();
	await r.watch();
	// Webhook for new inbox entries, only when system-flow.yaml asks for it (S-0042). The manifest is
	// asked of flai, which may not have connected yet: try now, and again whenever one connects.
	let notifying = false;
	const startNotifying = async () => {
		if (notifying) return;
		try {
			notifying = (await startNotifier(r)) !== null;
		} catch (err) {
			log().debug({ component: 'notify', err: String(err) }, 'inbox webhook not started yet');
		}
	};
	agent().on('connected', () => void startNotifying());
	void startNotifying();
	log().info(
		{
			component: 'server',
			version,
			commit,
			project_dir: process.env.PROJECT_DIR ?? process.cwd(),
			tracing,
			auth,
			channel
		},
		'server started'
	);
};

const QUIET = new Set(['/_health', '/_ready', '/metrics']);

export const handle: Handle = async ({ event, resolve }) => {
	const method = event.request.method;
	const path = event.url.pathname;
	const route = routeLabel(event.route.id, path);
	const startedAt = process.hrtime.bigint();
	inFlight.inc();
	try {
		return await withRequestSpan(method, path, event.request.headers, async (span) => {
			const tid = traceId(span);
			const requestId = tid ? undefined : crypto.randomUUID();
			event.locals.traceId = tid ?? requestId;
			let status = 500;
			try {
				const auth = authenticate(event.request.headers);
				event.locals.auth = auth;
				const decision = decide(path, method, auth, event.request.headers);
				if (decision.kind === 'unauthorized') {
					const response = json({ error: 'unauthorized' }, { status: 401 });
					status = 401;
					return response;
				}
				if (decision.kind === 'forbidden') {
					const response = json({ error: 'forbidden' }, { status: 403 });
					status = 403;
					return response;
				}
				if (decision.kind === 'login') {
					status = 303;
					redirect(303, `/login?next=${encodeURIComponent(decision.next)}`);
				}
				const response = await resolve(event);
				status = response.status;
				// Every API and MCP answer names its project (ADR-0024); refusals above do not,
				// so an unauthenticated caller learns nothing about what is served here.
				if (path.startsWith('/api/') || path === '/mcp')
					return setIdentityHeaders(response, await projectIdentity(repo()));
				return response;
			} finally {
				const seconds = Number(process.hrtime.bigint() - startedAt) / 1e9;
				span.setAttributes({ 'http.route': route, 'http.response.status_code': status });
				if (status >= 500) span.setStatus({ code: SpanStatusCode.ERROR });
				requestsTotal.labels(method, route, statusClass(status)).inc();
				requestDuration.labels(method, route).observe(seconds);
				if (!QUIET.has(path) && !path.startsWith('/_app/')) {
					log()[status >= 500 ? 'error' : 'info'](
						{
							component: 'http',
							trace_id: tid,
							request_id: requestId,
							method,
							route,
							path,
							status,
							duration_ms: Math.round(seconds * 1000)
						},
						'request handled'
					);
				}
			}
		});
	} finally {
		inFlight.dec();
	}
};

export const handleError: HandleServerError = ({ error, event, status, message }) => {
	log().error(
		{
			component: 'http',
			trace_id: event.locals.traceId,
			route: event.route.id,
			path: event.url.pathname,
			status,
			err: error instanceof Error ? (error.stack ?? error.message) : String(error)
		},
		'request failed'
	);
	return { message };
};

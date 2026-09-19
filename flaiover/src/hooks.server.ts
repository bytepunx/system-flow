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
import { redirect, json } from '@sveltejs/kit';

export const init: ServerInit = async () => {
	const auth = initAuth();
	const tracing = await startTracing();
	// Webhook for new inbox entries, only when system-flow.yaml asks for it (S-0042).
	// A project that cannot be read yet must not stop the server from starting.
	try {
		const r = repo();
		if (await startNotifier(r)) await r.watch();
	} catch (err) {
		log().warn({ component: 'notify', err: String(err) }, 'inbox webhook not started');
	}
	log().info(
		{
			component: 'server',
			version,
			commit,
			project_dir: process.env.PROJECT_DIR ?? process.cwd(),
			tracing,
			auth
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

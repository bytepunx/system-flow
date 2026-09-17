// OpenTelemetry traces: exported over OTLP/HTTP only when
// OTEL_EXPORTER_OTLP_ENDPOINT is set; otherwise the API is a no-op. One
// span per request is created in hooks.server.ts with W3C context
// propagated from the incoming headers, so the trace_id in logs joins the
// spans of the caller.
import { context, propagation, trace, type Span } from '@opentelemetry/api';

let started = false;

/** Start the SDK once per process when an endpoint is configured. */
export async function startTracing(): Promise<boolean> {
	if (started) return true;
	const endpoint =
		process.env.OTEL_EXPORTER_OTLP_ENDPOINT ?? process.env.OTEL_EXPORTER_OTLP_TRACES_ENDPOINT;
	if (!endpoint) return false;
	const [
		{ NodeSDK },
		{ OTLPTraceExporter },
		{ resourceFromAttributes },
		{ ATTR_SERVICE_NAME, ATTR_SERVICE_VERSION }
	] = await Promise.all([
		import('@opentelemetry/sdk-node'),
		import('@opentelemetry/exporter-trace-otlp-http'),
		import('@opentelemetry/resources'),
		import('@opentelemetry/semantic-conventions')
	]);
	const { version } = await import('./metrics');
	const sdk = new NodeSDK({
		resource: resourceFromAttributes({
			[ATTR_SERVICE_NAME]: process.env.OTEL_SERVICE_NAME ?? 'flaiover',
			[ATTR_SERVICE_VERSION]: version
		}),
		traceExporter: new OTLPTraceExporter()
	});
	sdk.start();
	started = true;
	const stop = () => void sdk.shutdown().catch(() => {});
	process.once('SIGTERM', stop);
	process.once('SIGINT', stop);
	return true;
}

export const tracer = trace.getTracer('flaiover');

/** Run fn inside a server span whose parent comes from the request headers. */
export async function withRequestSpan<T>(
	method: string,
	path: string,
	headers: Headers,
	fn: (span: Span) => Promise<T>
): Promise<T> {
	const carrier: Record<string, string> = {};
	headers.forEach((v, k) => (carrier[k.toLowerCase()] = v));
	const parent = propagation.extract(context.active(), carrier);
	return tracer.startActiveSpan(
		`${method} ${path}`,
		{ kind: 1 /* SERVER */, attributes: { 'http.request.method': method, 'url.path': path } },
		parent,
		async (span) => {
			try {
				return await fn(span);
			} finally {
				span.end();
			}
		}
	);
}

/** The active trace id, for the request log line. */
export function traceId(span: Span): string | undefined {
	const id = span.spanContext().traceId;
	return id && /[^0]/.test(id) ? id : undefined;
}

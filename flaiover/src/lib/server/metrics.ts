// The four golden signals plus build_info, per design/conventions/telemetry.md,
// exposed at /metrics in Prometheus format. Names are flaiover_ prefixed,
// snake_case, with unit suffixes; labels are a declared low-cardinality set.
import { Counter, Gauge, Histogram, Registry, collectDefaultMetrics } from '@prometheus-io/client';

export const registry = new Registry();

export const requestsTotal = new Counter({
	name: 'flaiover_http_requests_total',
	help: 'HTTP requests handled, by method, route, and status class',
	labelNames: ['method', 'route', 'status'] as const,
	registers: [registry]
});

export const requestDuration = new Histogram({
	name: 'flaiover_http_request_duration_seconds',
	help: 'HTTP request latency in seconds, by method and route',
	labelNames: ['method', 'route'] as const,
	buckets: [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5],
	registers: [registry]
});

export const inFlight = new Gauge({
	name: 'flaiover_http_requests_in_flight',
	help: 'HTTP requests currently being handled',
	registers: [registry]
});

export const buildInfo = new Gauge({
	name: 'flaiover_build_info',
	help: 'Build metadata; always 1',
	labelNames: ['version', 'commit'] as const,
	registers: [registry]
});

// Version: the release tag baked into the image (FLAIOVER_VERSION), else the
// package version Vite compiles in, so build_info is right in both cases.
declare const __FLAIOVER_VERSION__: string | undefined;
export const version: string =
	process.env.FLAIOVER_VERSION ??
	(typeof __FLAIOVER_VERSION__ === 'string' ? __FLAIOVER_VERSION__ : '0.0.0');
export const commit = process.env.FLAIOVER_COMMIT ?? 'unknown';
buildInfo.labels(version, commit).set(1);

collectDefaultMetrics({ register: registry, prefix: 'flaiover_' });

/** Route label: the matched SvelteKit route id, so IDs and paths never become labels. */
export function routeLabel(routeId: string | null, path: string): string {
	if (routeId) return routeId;
	return path.startsWith('/_app/') ? '/_app/*' : 'unmatched';
}

/** Status class label: 2xx, 3xx, 4xx, 5xx. */
export function statusClass(status: number): string {
	return `${Math.floor(status / 100)}xx`;
}

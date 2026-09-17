import { describe, expect, it } from 'vitest';
import {
	registry,
	requestsTotal,
	requestDuration,
	inFlight,
	routeLabel,
	statusClass,
	version
} from './metrics';

describe('metrics', () => {
	it('exposes the golden signals and build_info with the flaiover prefix', async () => {
		requestsTotal.labels('GET', '/api/items', '2xx').inc();
		requestDuration.labels('GET', '/api/items').observe(0.02);
		inFlight.inc();
		inFlight.dec();
		const text = await registry.metrics();
		for (const name of [
			'flaiover_http_requests_total',
			'flaiover_http_request_duration_seconds_bucket',
			'flaiover_http_requests_in_flight',
			'flaiover_build_info',
			'flaiover_process_cpu_seconds_total'
		]) {
			expect(text, name).toContain(name);
		}
		expect(text).toContain(`flaiover_build_info{version="${version}"`);
		expect(text).toContain(
			'flaiover_http_requests_total{method="GET",route="/api/items",status="2xx"} 1'
		);
		expect(registry.contentType).toContain('text/plain');
	});
	it('labels routes by route id and statuses by class', () => {
		expect(routeLabel('/api/items/[id]', '/api/items/S-001')).toBe('/api/items/[id]');
		expect(routeLabel(null, '/_app/immutable/x.js')).toBe('/_app/*');
		expect(routeLabel(null, '/nope')).toBe('unmatched');
		expect(statusClass(204)).toBe('2xx');
		expect(statusClass(503)).toBe('5xx');
	});
});

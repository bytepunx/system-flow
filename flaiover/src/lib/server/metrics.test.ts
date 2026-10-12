import { describe, expect, it } from 'vitest';
import {
	registry,
	requestsTotal,
	requestDuration,
	inFlight,
	routeLabel,
	statusClass,
	version,
	commit,
	recordBuildInfo
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
	it('labels build_info signed by whether the release stamp verified, one series at a time (S-0235)', async () => {
		const series = (signed: string) =>
			`flaiover_build_info{version="${version}",commit="${commit}",signed="${signed}"} 1`;
		expect(await registry.metrics()).toContain(series('false'));
		recordBuildInfo(true);
		let text = await registry.metrics();
		expect(text).toContain(series('true'));
		expect(text).not.toContain(series('false'));
		recordBuildInfo(false);
		text = await registry.metrics();
		expect(text).toContain(series('false'));
		expect(text).not.toContain(series('true'));
	});
	it('labels routes by route id and statuses by class', () => {
		expect(routeLabel('/api/items/[id]', '/api/items/S-001')).toBe('/api/items/[id]');
		expect(routeLabel(null, '/_app/immutable/x.js')).toBe('/_app/*');
		expect(routeLabel(null, '/nope')).toBe('unmatched');
		expect(statusClass(204)).toBe('2xx');
		expect(statusClass(503)).toBe('5xx');
	});
});

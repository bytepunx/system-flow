import { describe, expect, it } from 'vitest';
import { context, trace } from '@opentelemetry/api';
import {
	BasicTracerProvider,
	InMemorySpanExporter,
	SimpleSpanProcessor
} from '@opentelemetry/sdk-trace-base';
import { W3CTraceContextPropagator } from '@opentelemetry/core';
import { propagation } from '@opentelemetry/api';
import { AsyncLocalStorageContextManager } from '@opentelemetry/context-async-hooks';
import { startTracing, traceId, withRequestSpan } from './otel';

describe('tracing', () => {
	it('is a no-op without an endpoint', async () => {
		delete process.env.OTEL_EXPORTER_OTLP_ENDPOINT;
		expect(await startTracing()).toBe(false);
		await withRequestSpan('GET', '/x', new Headers(), async (span) => {
			expect(traceId(span)).toBeUndefined();
		});
	});
	it('creates a server span with the incoming W3C parent and exposes the trace id', async () => {
		const exporter = new InMemorySpanExporter();
		const provider = new BasicTracerProvider({
			spanProcessors: [new SimpleSpanProcessor(exporter)]
		});
		trace.setGlobalTracerProvider(provider);
		propagation.setGlobalPropagator(new W3CTraceContextPropagator());
		context.setGlobalContextManager(new AsyncLocalStorageContextManager().enable());
		const headers = new Headers({
			traceparent: '00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01'
		});
		let tid: string | undefined;
		await withRequestSpan('GET', '/api/items', headers, async (span) => {
			tid = traceId(span);
			span.setAttributes({ 'http.route': '/api/items', 'http.response.status_code': 200 });
		});
		expect(tid).toBe('0af7651916cd43dd8448eb211c80319c');
		const spans = exporter.getFinishedSpans();
		expect(spans).toHaveLength(1);
		expect(spans[0].name).toBe('GET /api/items');
		expect(spans[0].parentSpanContext?.spanId).toBe('b7ad6b7169203331');
		expect(spans[0].attributes['http.route']).toBe('/api/items');
		expect(spans[0].kind).toBe(1);
		trace.disable();
	});
});

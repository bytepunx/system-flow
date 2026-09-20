import { afterEach, describe, expect, it } from 'vitest';
// the server entry's helper is plain JavaScript beside server.js, outside src
import { noteScheme, schemeHeader } from '../../../scheme.js';

type Req = {
	headers: Record<string, string | string[] | undefined>;
	socket?: { encrypted?: boolean };
};
const req = (headers: Req['headers'] = {}, encrypted = false): Req => ({
	headers,
	socket: { encrypted }
});

// S-0083: adapter-node takes every request for https unless a header says otherwise, and a cookie
// marked Secure over plain HTTP is kept by a browser on localhost only.
describe('the scheme a request came over', () => {
	afterEach(() => {
		delete process.env.PROTOCOL_HEADER;
	});

	it('is what the connection is when no proxy said', () => {
		const plain = req();
		expect(noteScheme(plain)).toBe('http');
		expect(plain.headers['x-forwarded-proto']).toBe('http');
		expect(noteScheme(req({}, true))).toBe('https');
	});

	it('is what a TLS-terminating proxy said, the first when it names a chain', () => {
		expect(noteScheme(req({ 'x-forwarded-proto': 'https' }))).toBe('https');
		expect(noteScheme(req({ 'x-forwarded-proto': 'HTTPS, http' }))).toBe('https');
		expect(noteScheme(req({ 'x-forwarded-proto': ['http', 'https'] }))).toBe('http');
	});

	it('replaces anything else, which adapter-node would refuse with an error per request', () => {
		for (const odd of ['javascript:', 'https://evil.example', 'ws', '', ' ']) {
			const r = req({ 'x-forwarded-proto': odd });
			expect(noteScheme(r)).toBe('http');
			expect(r.headers['x-forwarded-proto']).toBe('http');
		}
	});

	it('follows PROTOCOL_HEADER when an operator names another header', () => {
		process.env.PROTOCOL_HEADER = 'X-Scheme';
		expect(schemeHeader()).toBe('x-scheme');
		const r = req({ 'x-scheme': 'https', 'x-forwarded-proto': 'http' });
		expect(noteScheme(r)).toBe('https');
	});
});

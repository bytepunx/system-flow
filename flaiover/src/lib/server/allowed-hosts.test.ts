// ADR-0193: the hosts a browser's Origin may name beside the request's own, from FLAIOVER_ALLOWED_HOSTS.
import { describe, expect, it } from 'vitest';
import { allowedHosts, allowsOrigin, parseAllowedHosts } from './allowed-hosts';

describe('the allowed hosts (ADR-0193)', () => {
	it('reads a comma-separated list, trimmed and lower-cased, the empty entries dropped', () => {
		expect(parseAllowedHosts('flow.example.com,Board.Example.NET:8443')).toEqual([
			'flow.example.com',
			'board.example.net:8443'
		]);
		expect(parseAllowedHosts('  flow.example.com , ,[::1]:8443,')).toEqual([
			'flow.example.com',
			'[::1]:8443'
		]);
		expect(parseAllowedHosts('')).toEqual([]);
		expect(parseAllowedHosts(' , ')).toEqual([]);
		expect(parseAllowedHosts(undefined)).toEqual([]);
	});

	it('takes the list from FLAIOVER_ALLOWED_HOSTS, none when it is unset', () => {
		expect(allowedHosts({ FLAIOVER_ALLOWED_HOSTS: 'flow.example.com' })).toEqual([
			'flow.example.com'
		]);
		expect(allowedHosts({})).toEqual([]);
	});

	it('allows an origin that names a listed host, whatever its case', () => {
		const hosts = parseAllowedHosts('Flow.Example.com');
		expect(allowsOrigin(new URL('https://flow.example.com'), hosts)).toBe(true);
		expect(allowsOrigin(new URL('https://FLOW.example.COM'), hosts)).toBe(true);
		expect(allowsOrigin(new URL('http://flow.example.com'), hosts)).toBe(true);
		expect(allowsOrigin(new URL('https://elsewhere.example'), hosts)).toBe(false);
		expect(allowsOrigin(new URL('https://flow.example.com.elsewhere.example'), hosts)).toBe(false);
		expect(allowsOrigin(new URL('https://board.flow.example.com'), hosts)).toBe(false);
	});

	it('matches the port: an entry with none means the default port of the origin', () => {
		const hosts = parseAllowedHosts(
			'flow.example.com,board.example.net:8443,[::1]:9000,edge.example:443'
		);
		expect(allowsOrigin(new URL('https://flow.example.com:8443'), hosts)).toBe(false);
		expect(allowsOrigin(new URL('https://flow.example.com:443'), hosts)).toBe(true);
		expect(allowsOrigin(new URL('https://board.example.net:8443'), hosts)).toBe(true);
		expect(allowsOrigin(new URL('https://board.example.net'), hosts)).toBe(false);
		expect(allowsOrigin(new URL('http://[::1]:9000'), hosts)).toBe(true);
		expect(allowsOrigin(new URL('http://[::1]'), hosts)).toBe(false);
		expect(allowsOrigin(new URL('https://edge.example'), hosts)).toBe(true);
		expect(allowsOrigin(new URL('http://edge.example'), hosts)).toBe(false);
	});

	it('allows nothing when no host is listed', () => {
		expect(allowsOrigin(new URL('https://flow.example.com'), [])).toBe(false);
	});
});

import { afterEach, describe, expect, it } from 'vitest';
import { mkdtempSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join } from 'node:path';
import {
	authenticate,
	clearCookie,
	decide,
	initAuth,
	sessionCookie,
	setToken,
	tokenMatches
} from './auth';

const H = (h: Record<string, string> = {}) => new Headers(h);

describe('auth', () => {
	afterEach(() => setToken(null));

	it('loads the token from the file and refuses to start without one', () => {
		const dir = mkdtempSync(join(tmpdir(), 'flaiover-auth-'));
		const file = join(dir, 'dashboard.token');
		writeFileSync(file, 'abc123\n');
		expect(initAuth({ FLAIOVER_TOKEN_FILE: file })).toBe('token');
		expect(tokenMatches('abc123')).toBe(true);
		expect(() => initAuth({})).toThrow(/FLAIOVER_TOKEN_FILE/);
		expect(() => initAuth({ FLAIOVER_AUTH: 'off', NODE_ENV: 'production' })).toThrow();
		expect(initAuth({ FLAIOVER_AUTH: 'off', NODE_ENV: 'development' })).toBe('off');
		expect(authenticate(H())).toBe('off');
	});

	it('accepts bearer and cookie, rejects missing, wrong, and rotated tokens', () => {
		setToken('secret-one');
		expect(authenticate(H({ authorization: 'Bearer secret-one' }))).toBe('bearer');
		expect(authenticate(H({ authorization: 'bearer secret-one' }))).toBe('bearer');
		expect(authenticate(H({ cookie: 'x=1; flaiover_session=secret-one' }))).toBe('cookie');
		expect(authenticate(H())).toBeNull();
		expect(authenticate(H({ authorization: 'Bearer secret-two' }))).toBeNull();
		expect(authenticate(H({ authorization: 'Basic secret-one' }))).toBeNull();
		expect(authenticate(H({ cookie: 'flaiover_session=secret-on' }))).toBeNull();
		setToken('secret-rotated');
		expect(authenticate(H({ cookie: 'flaiover_session=secret-one' }))).toBeNull();
		expect(authenticate(H({ authorization: 'Bearer secret-rotated' }))).toBe('bearer');
	});

	it('decides per path, method, and credential', () => {
		expect(decide('/_health', 'GET', null, H())).toEqual({ kind: 'allow' });
		expect(decide('/_ready', 'GET', null, H())).toEqual({ kind: 'allow' });
		expect(decide('/login', 'GET', null, H())).toEqual({ kind: 'allow' });
		expect(decide('/api/login', 'POST', null, H())).toEqual({ kind: 'allow' });
		expect(decide('/_app/immutable/x.js', 'GET', null, H())).toEqual({ kind: 'allow' });
		expect(decide('/metrics', 'GET', null, H(), {})).toEqual({ kind: 'unauthorized' });
		expect(decide('/metrics', 'GET', null, H(), { FLAIOVER_METRICS_PUBLIC: 'true' })).toEqual({
			kind: 'allow'
		});
		expect(decide('/api/items', 'GET', null, H())).toEqual({ kind: 'unauthorized' });
		expect(decide('/board', 'GET', null, H())).toEqual({ kind: 'login', next: '/board' });
		expect(decide('/api/items', 'GET', 'bearer', H())).toEqual({ kind: 'allow' });
		expect(decide('/api/items/S-1/move', 'POST', 'bearer', H())).toEqual({ kind: 'allow' });
		expect(decide('/api/items/S-1/move', 'POST', 'cookie', H())).toEqual({ kind: 'forbidden' });
		expect(
			decide('/api/items/S-1/move', 'POST', 'cookie', H({ 'x-requested-with': 'flaiover' }))
		).toEqual({ kind: 'allow' });
		expect(decide('/api/items', 'GET', 'off', H())).toEqual({ kind: 'allow' });
	});

	it('builds cookies', () => {
		expect(sessionCookie('a b', false)).toBe(
			'flaiover_session=a%20b; Path=/; HttpOnly; SameSite=Lax; Max-Age=2592000'
		);
		expect(sessionCookie('t', true)).toContain('; Secure');
		expect(clearCookie()).toContain('Max-Age=0');
	});
});

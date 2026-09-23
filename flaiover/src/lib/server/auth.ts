// Dashboard authentication (ADR-0018): one per-project token, read once from
// FLAIOVER_TOKEN_FILE. Agents send it as a bearer header; browsers hold it in
// an HttpOnly cookie set by /api/login from the login page's URL fragment.
// The token is never logged, never echoed, and compared in constant time.
import { createHash, timingSafeEqual } from 'node:crypto';
import { readFileSync } from 'node:fs';
import { log } from './log';

export const COOKIE = 'flaiover_session';
export const CSRF_HEADER = 'x-requested-with';
export const CSRF_VALUE = 'flaiover';
const COOKIE_MAX_AGE = 30 * 24 * 3600;

export type AuthKind = 'bearer' | 'cookie' | 'off';

let token: string | null = null;
let disabled = false;

/** Load the token from FLAIOVER_TOKEN_FILE. Without a file the server refuses to start unless FLAIOVER_AUTH=off outside production. */
export function initAuth(env: NodeJS.ProcessEnv = process.env): AuthKind | 'token' {
	const file = env.FLAIOVER_TOKEN_FILE;
	if (file) {
		const value = readFileSync(file, 'utf8').trim();
		if (!value) throw new Error(`${file} is empty; run flai dashboard token`);
		token = value;
		disabled = false;
		return 'token';
	}
	if (env.FLAIOVER_AUTH === 'off' && env.NODE_ENV !== 'production') {
		disabled = true;
		token = null;
		log().warn(
			{ component: 'auth' },
			'authentication is off (FLAIOVER_AUTH=off); development only'
		);
		return 'off';
	}
	throw new Error(
		'FLAIOVER_TOKEN_FILE is not set; flai dashboard mounts the token, or set FLAIOVER_AUTH=off for development'
	);
}

/** For tests: install a token directly. */
export function setToken(value: string | null, off = false): void {
	token = value;
	disabled = off;
}

export function authDisabled(): boolean {
	return disabled;
}

/** Constant-time comparison over digests, so length never leaks. */
export function tokenMatches(candidate: string | null | undefined): boolean {
	if (!token || !candidate) return false;
	const a = createHash('sha256').update(candidate).digest();
	const b = createHash('sha256').update(token).digest();
	return timingSafeEqual(a, b);
}

export function bearerOf(headers: Headers): string | null {
	const h = headers.get('authorization');
	if (!h) return null;
	const [scheme, value] = h.split(' ', 2);
	return scheme?.toLowerCase() === 'bearer' && value ? value.trim() : null;
}

export function cookieOf(headers: Headers, name = COOKIE): string | null {
	const raw = headers.get('cookie');
	if (!raw) return null;
	for (const part of raw.split(';')) {
		const [k, ...v] = part.trim().split('=');
		if (k === name) return decodeURIComponent(v.join('='));
	}
	return null;
}

/** Who the request is: bearer, cookie, off (auth disabled), or null when unauthenticated. */
export function authenticate(headers: Headers): AuthKind | null {
	if (disabled) return 'off';
	if (tokenMatches(bearerOf(headers))) return 'bearer';
	if (tokenMatches(cookieOf(headers))) return 'cookie';
	return null;
}

export type Decision =
	| { kind: 'allow' }
	| { kind: 'unauthorized' }
	| { kind: 'login'; next: string }
	| { kind: 'forbidden' };

/**
 * Policy for one request. Health and readiness are open; /login and the app
 * bundle are open so the login page can render; /metrics is open only when
 * FLAIOVER_METRICS_PUBLIC=true; everything else needs a token. Cookie-backed
 * writes must also carry the CSRF header.
 */
export function decide(
	path: string,
	method: string,
	auth: AuthKind | null,
	headers: Headers,
	env: NodeJS.ProcessEnv = process.env
): Decision {
	if (
		path === '/_health' ||
		path === '/_ready' ||
		path === '/login' ||
		path === '/api/login' ||
		path.startsWith('/_app/') ||
		path === '/robots.txt'
	) {
		return { kind: 'allow' };
	}
	if (path === '/metrics' && env.FLAIOVER_METRICS_PUBLIC === 'true') return { kind: 'allow' };
	// MCP moved to flai on the host (ADR-0030). What is left here only says so, to anyone.
	if (path === '/mcp') return { kind: 'allow' };
	if (!auth) {
		if (method === 'GET' && !path.startsWith('/api/') && path !== '/metrics') {
			return { kind: 'login', next: path };
		}
		return { kind: 'unauthorized' };
	}
	if (
		auth === 'cookie' &&
		method !== 'GET' &&
		method !== 'HEAD' &&
		headers.get(CSRF_HEADER) !== CSRF_VALUE
	) {
		return { kind: 'forbidden' };
	}
	return { kind: 'allow' };
}

/** Set-Cookie value for a session; Secure whenever the page is served over HTTPS. */
export function sessionCookie(value: string, secure: boolean): string {
	const parts = [
		`${COOKIE}=${encodeURIComponent(value)}`,
		'Path=/',
		'HttpOnly',
		'SameSite=Lax',
		`Max-Age=${COOKIE_MAX_AGE}`
	];
	if (secure) parts.push('Secure');
	return parts.join('; ');
}

export function clearCookie(): string {
	return `${COOKIE}=; Path=/; HttpOnly; SameSite=Lax; Max-Age=0`;
}

export function isSecure(url: URL, headers: Headers): boolean {
	return url.protocol === 'https:' || headers.get('x-forwarded-proto') === 'https';
}

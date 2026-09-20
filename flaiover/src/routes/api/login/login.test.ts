import { afterEach, describe, expect, it } from 'vitest';
import { setToken } from '$lib/server/auth';
import { POST } from './+server';

type Handler = (event: never) => Promise<Response>;
const login = (url: string, headers: Record<string, string> = {}, token = 'the-token') =>
	(POST as Handler)({
		url: new URL(url),
		request: new Request(url, {
			method: 'POST',
			headers: { 'content-type': 'application/json', ...headers },
			body: JSON.stringify({ token })
		})
	} as never);

// S-0083: a cookie marked Secure over plain HTTP is kept by a browser on http://localhost and
// dropped everywhere else, so a login by any other address came straight back to the login page.
describe('the session cookie follows the scheme the page is served over', () => {
	afterEach(() => setToken(null));

	it('is not Secure over plain HTTP, by any address', async () => {
		setToken('the-token');
		for (const url of ['http://192.168.0.162:4242/api/login', 'http://localhost:4242/api/login']) {
			const cookie = (await login(url)).headers.get('set-cookie')!;
			expect(cookie).toContain('flaiover_session=the-token');
			expect(cookie).toContain('HttpOnly');
			expect(cookie).toContain('SameSite=Lax');
			expect(cookie).not.toContain('Secure');
		}
	});

	it('is Secure over HTTPS, and behind a proxy that terminates TLS and says so', async () => {
		setToken('the-token');
		expect((await login('https://dash.example/api/login')).headers.get('set-cookie')).toContain(
			'; Secure'
		);
		const forwarded = await login('http://127.0.0.1:4242/api/login', {
			'x-forwarded-proto': 'https'
		});
		expect(forwarded.headers.get('set-cookie')).toContain('; Secure');
	});

	it('sets nothing for a wrong token', async () => {
		setToken('the-token');
		const r = await login('http://192.168.0.162:4242/api/login', {}, 'nope');
		expect(r.status).toBe(401);
		expect(r.headers.get('set-cookie')).toBeNull();
	});
});

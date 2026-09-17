import { afterEach, describe, expect, it } from 'vitest';
import { Writable } from 'node:stream';
import { createLogger, log, resetLog } from './log';

describe('log', () => {
	afterEach(() => {
		resetLog();
		delete process.env.LOG_LEVEL;
	});
	it('emits JSON events with ts, level, component, and msg', () => {
		process.env.LOG_LEVEL = 'debug';
		const lines: string[] = [];
		const dest = new Writable({
			write(chunk, _enc, cb) {
				lines.push(String(chunk));
				cb();
			}
		});
		const l = createLogger(dest);
		expect(l.level).toBe('debug');
		l.child({ component: 'test' }).info({ route: '/x' }, 'request handled');
		const ev = JSON.parse(lines.join('').trim().split('\n').at(-1)!);
		expect(ev).toMatchObject({
			level: 'INFO',
			service: 'flaiover',
			component: 'test',
			msg: 'request handled',
			route: '/x'
		});
		expect((lines.join('').match(/"component"/g) ?? []).length).toBe(1);
		expect(ev.ts).toMatch(/^\d{4}-\d{2}-\d{2}T/);
	});
	it('falls back to info for an unknown level', () => {
		process.env.LOG_LEVEL = 'loud';
		expect(log().level).toBe('info');
	});
});

// Structured logging per design/conventions/logging.md: JSON by default,
// key-value text when LOG_FORMAT=text (local development), one static
// message per event with data in fields, levels from LOG_LEVEL.
import pino, { type DestinationStream, type Logger } from 'pino';

const LEVELS = new Set(['fatal', 'error', 'warn', 'info', 'debug', 'trace']);

/** Build a logger; tests pass a destination stream to capture output. */
export function createLogger(dest?: DestinationStream): Logger {
	const level = (process.env.LOG_LEVEL ?? 'info').toLowerCase();
	const text = (process.env.LOG_FORMAT ?? '').toLowerCase() === 'text';
	const base = {
		level: LEVELS.has(level) ? level : 'info',
		base: { service: 'flaiover' },
		timestamp: () => `,"ts":"${new Date().toISOString()}"`,
		messageKey: 'msg',
		formatters: { level: (label: string) => ({ level: label.toUpperCase() }) }
	};
	if (dest) return pino(base, dest);
	if (text && process.env.NODE_ENV !== 'production') {
		return pino({
			...base,
			transport: {
				target: 'pino-pretty',
				options: { colorize: true, translateTime: false, messageKey: 'msg' }
			}
		});
	}
	return pino(base);
}

let shared: Logger | null = null;
/** The process logger. Child loggers add `component` per module. */
export function log(): Logger {
	if (!shared) shared = createLogger();
	return shared;
}

/** For tests: rebuild with the current environment. */
export function resetLog() {
	shared = null;
}

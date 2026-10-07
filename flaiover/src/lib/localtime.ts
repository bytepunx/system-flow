/**
 * Times as the dashboard shows them: flai records every time in UTC (`YYYY-MM-DDTHH:MM:SSZ`),
 * and each is shown in the browser's time zone, or in `zone` when one is given.
 */

const formats = new Map<string, Intl.DateTimeFormat>();

function format(zone: string | undefined): Intl.DateTimeFormat {
	const key = zone ?? '';
	let f = formats.get(key);
	if (!f) {
		f = new Intl.DateTimeFormat('en-US', {
			timeZone: zone,
			year: 'numeric',
			month: '2-digit',
			day: '2-digit',
			hour: '2-digit',
			minute: '2-digit',
			hourCycle: 'h23',
			timeZoneName: 'short'
		});
		formats.set(key, f);
	}
	return f;
}

function parts(at: string, zone?: string): Record<string, string> | undefined {
	if (typeof at !== 'string' || !/^\d{4}-\d\d-\d\dT\d\d:\d\d/.test(at)) return undefined;
	const ms = Date.parse(at);
	if (Number.isNaN(ms)) return undefined;
	const out: Record<string, string> = {};
	for (const p of format(zone).formatToParts(ms)) out[p.type] = p.value;
	return out;
}

/** A recorded time to the minute in the local zone with its short name, as 2026-10-07 15:35 EDT; what it was given when that does not parse. */
export function localTime(at: string, zone?: string): string {
	const p = parts(at, zone);
	if (!p) return at;
	return `${p.year}-${p.month}-${p.day} ${p.hour}:${p.minute} ${p.timeZoneName}`;
}

/** A text flai wrote, such as a refusal, with each time it names as flai records one shown by localTime. */
export function localTimes(text: string, zone?: string): string {
	return text.replace(/\d{4}-\d\d-\d\dT\d\d:\d\d:\d\dZ/g, (at) => localTime(at, zone));
}

/** A recorded time's calendar date in the local zone, YYYY-MM-DD; what it was given when that does not parse. */
export function localDate(at: string, zone?: string): string {
	const p = parts(at, zone);
	if (!p) return at;
	return `${p.year}-${p.month}-${p.day}`;
}

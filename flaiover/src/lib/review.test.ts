import { describe, expect, it } from 'vitest';
import { criteriaOf, patchLines, readNdjson, sectionOf } from './review';

const STORY = `# S-0041 Review

## Goal
Review in one place.

## Acceptance criteria
- [x] A review page per story
- [ ] Accept runs as the designer
* [X] Failures are shown verbatim

## Tasks
- T-0185 Living design
- [ ] not a criterion, it is under Tasks

## Notes
n
`;

describe('review helpers', () => {
	it('reads the acceptance criteria and their ticks, and nothing from other sections', () => {
		expect(criteriaOf(STORY)).toEqual([
			{ checked: true, text: 'A review page per story' },
			{ checked: false, text: 'Accept runs as the designer' },
			{ checked: true, text: 'Failures are shown verbatim' }
		]);
		expect(criteriaOf('# no criteria\n')).toEqual([]);
	});
	it('returns a section up to the next heading', () => {
		expect(sectionOf(STORY, 'Goal')).toBe('Review in one place.');
		expect(sectionOf(STORY, 'goal')).toBe('Review in one place.');
		expect(sectionOf(STORY, 'Notes')).toBe('n');
		expect(sectionOf(STORY, 'Missing')).toBe('');
	});
	it('classifies patch lines', () => {
		const lines = patchLines('@@ -1,2 +1,2 @@\n one\n-two\n+2\n\\ No newline at end of file');
		expect(lines.map((l) => l.kind)).toEqual(['hunk', 'context', 'del', 'add', 'note']);
		expect(patchLines('')).toEqual([]);
	});
	it('delivers NDJSON lines as they arrive, across chunk boundaries', async () => {
		const enc = new TextEncoder();
		const chunks = [
			'{"event":"progress","st',
			'ep":"merged"}\n{"event":"pro',
			'gress","step":"done"}\nnot json\n{"event":"done"}'
		];
		const body = new ReadableStream<Uint8Array>({
			start(c) {
				for (const ch of chunks) c.enqueue(enc.encode(ch));
				c.close();
			}
		});
		const got: Record<string, unknown>[] = [];
		await readNdjson(new Response(body), (l) => got.push(l));
		expect(got).toEqual([
			{ event: 'progress', step: 'merged' },
			{ event: 'progress', step: 'done' },
			{ event: 'done' }
		]);
	});
});

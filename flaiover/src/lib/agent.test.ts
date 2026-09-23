import { describe, expect, it } from 'vitest';
import { agentFrom, agentLine, configText, parseConfig, sameAgent } from './agent';

describe('agent helpers (S-0103)', () => {
	it('reads and writes the config as key=value lines', () => {
		expect(parseConfig('effort=high\n\n max_turns = 50 \n')).toEqual({
			config: { effort: 'high', max_turns: '50' }
		});
		expect(parseConfig('no equals sign')).toEqual({ error: '"no equals sign" is not key=value' });
		expect(configText({ max_turns: '50', effort: 'high' })).toBe('effort=high\nmax_turns=50');
	});
	it('gives an agent only for what is set, and compares by what it says', () => {
		expect(agentFrom(' ', '', {})).toBeUndefined();
		expect(agentFrom('claude-code', '', { effort: 'high' })).toEqual({
			harness: 'claude-code',
			config: { effort: 'high' }
		});
		expect(sameAgent({ model: 'm', config: {} }, { model: 'm' })).toBe(true);
		expect(sameAgent({ model: 'm' }, { model: 'n' })).toBe(false);
		expect(
			agentLine({ harness: 'claude-code', model: 'claude-opus-5-5', config: { effort: 'high' } })
		).toBe('claude-code, claude-opus-5-5, effort=high');
		expect(agentLine(undefined)).toBe('none');
	});
});

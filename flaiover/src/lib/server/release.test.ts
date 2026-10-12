// S-0232: the dashboard carries the release key's public half, the same PEM as flai's.
import { describe, expect, it } from 'vitest';
import { createHash, createPublicKey } from 'node:crypto';
import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';
import { RELEASE_PUBLIC_KEY } from './release';

describe('the release public key (S-0232)', () => {
	it('parses as an ECDSA P-256 public key with the fingerprint the runbook records', () => {
		const key = createPublicKey(RELEASE_PUBLIC_KEY);
		expect(key.type).toBe('public');
		expect(key.asymmetricKeyType).toBe('ec');
		expect(key.asymmetricKeyDetails?.namedCurve).toBe('prime256v1');
		const der = key.export({ type: 'spki', format: 'der' });
		expect(createHash('sha256').update(der).digest('hex')).toBe(
			'ee6df72313c9dbac1cf2270ed7f01cc523786d37bfce2dec12210215d6c5da40'
		);
	});

	it('is the same PEM flai embeds, read from a checkout', () => {
		const flai = readFileSync(resolve('..', 'flai', 'internal', 'buildinfo', 'cosign.pub'), 'utf8');
		expect(RELEASE_PUBLIC_KEY).toBe(flai);
	});
});

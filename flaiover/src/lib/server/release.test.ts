// S-0232: the dashboard carries the release key's public half, the same PEM as flai's.
import { describe, expect, it } from 'vitest';
import { createHash, createPublicKey, generateKeyPairSync, sign } from 'node:crypto';
import { mkdtempSync, readFileSync, writeFileSync } from 'node:fs';
import { tmpdir } from 'node:os';
import { join, resolve } from 'node:path';
import {
	RELEASE_PUBLIC_KEY,
	RELEASE_STAMP_FILE,
	checkReleaseStamp,
	isDevelopmentVersion,
	readStampFile,
	releaseStampPath,
	stampWarrantsWarning,
	verifyReleaseStamp,
	type ReleaseStamp,
	type StampBuild
} from './release';

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

// S-0235: the image keeps a stamp, `flaiover <version> <commit>` and its base64 ASN.1 signature, in a
// file the server reads and verifies once at start with the key it carries.
function keyPair() {
	const { publicKey, privateKey } = generateKeyPairSync('ec', { namedCurve: 'P-256' });
	return { pem: publicKey.export({ type: 'spki', format: 'pem' }).toString(), privateKey };
}
const ours = keyPair();
const theirs = keyPair();
const build: StampBuild = { version: '1.2.3', commit: 'abc1234', development: false };

function stampText(statement: string, key = ours.privateKey): string {
	return `${statement}\n${sign('sha256', Buffer.from(statement), key).toString('base64')}\n`;
}
function verifyText(text: string, keys = [ours.pem]): ReleaseStamp {
	return verifyReleaseStamp({ kind: 'read', path: '/stamp', text }, build, keys);
}

describe('the release stamp (S-0235)', () => {
	it('is signed when the release key signed this component, version, and commit', () => {
		const text = stampText('flaiover 1.2.3 abc1234');
		const stamp = verifyText(text);
		expect(stamp).toEqual({
			state: 'signed',
			reason: null,
			error: null,
			path: '/stamp',
			statement: 'flaiover 1.2.3 abc1234',
			signature: text.split('\n')[1],
			development: false
		});
		expect(verifyText(text.trimEnd()).state).toBe('signed');
		expect(verifyText(text, [theirs.pem, ours.pem]).state).toBe('signed');
	});

	it('is invalid when another key signed it', () => {
		const stamp = verifyText(stampText('flaiover 1.2.3 abc1234', theirs.privateKey));
		expect(stamp.state).toBe('invalid');
		expect(stamp.reason).toBe('bad-signature');
		expect(stamp.statement).toBe('flaiover 1.2.3 abc1234');
	});

	it('is invalid when the statement was changed after signing', () => {
		const signature = stampText('flaiover 1.2.3 abc1234').split('\n')[1];
		expect(verifyText(`flaiover 1.2.4 abc1234\n${signature}\n`).reason).toBe('bad-signature');
	});

	it('names the component, version, and commit the image reports', () => {
		expect(verifyText(stampText('flai 1.2.3 abc1234')).reason).toBe('wrong-component');
		expect(verifyText(stampText('flaiover 1.2.4 abc1234')).reason).toBe('version-mismatch');
		expect(verifyText(stampText('flaiover 1.2.3 def5678')).reason).toBe('commit-mismatch');
	});

	it('is malformed unless it is two lines, three fields, and a base64 signature', () => {
		const signature = stampText('flaiover 1.2.3 abc1234').split('\n')[1];
		for (const text of [
			'',
			'flaiover 1.2.3 abc1234\n',
			`flaiover 1.2.3 abc1234\n${signature}\nmore\n`,
			'flaiover 1.2.3 abc1234\nnot base64!\n',
			stampText('flaiover 1.2.3'),
			stampText('flaiover  1.2.3 abc1234')
		]) {
			expect(verifyText(text).reason, JSON.stringify(text)).toBe('malformed');
		}
	});

	it('is missing when there is no file, and invalid when the file cannot be read', () => {
		const dir = mkdtempSync(join(tmpdir(), 'flaiover-stamp-'));
		expect(readStampFile(join(dir, 'none'))).toEqual({ kind: 'absent', path: join(dir, 'none') });
		const missing = verifyReleaseStamp(readStampFile(join(dir, 'none')), build, [ours.pem]);
		expect(missing).toMatchObject({ state: 'missing', reason: null, statement: null });
		const unreadable = verifyReleaseStamp(readStampFile(dir), build, [ours.pem]);
		expect(unreadable).toMatchObject({ state: 'invalid', reason: 'unreadable' });
		expect(unreadable.error).toContain('EISDIR');
	});

	it('reads the file FLAIOVER_RELEASE_STAMP names, else /app/release-stamp', () => {
		expect(RELEASE_STAMP_FILE).toBe('/app/release-stamp');
		expect(releaseStampPath({})).toBe('/app/release-stamp');
		const dir = mkdtempSync(join(tmpdir(), 'flaiover-stamp-'));
		const path = join(dir, 'release-stamp');
		writeFileSync(path, stampText('flaiover 1.2.3 abc1234'));
		const env = { FLAIOVER_RELEASE_STAMP: path, FLAIOVER_VERSION: '1.2.3' };
		expect(releaseStampPath(env)).toBe(path);
		const stamp = checkReleaseStamp({ version: '1.2.3', commit: 'abc1234' }, env, [ours.pem]);
		expect(stamp).toMatchObject({ state: 'signed', path, development: false });
		const unset = { FLAIOVER_RELEASE_STAMP: join(dir, 'none') };
		const dev = checkReleaseStamp({ version: '1.2.3', commit: 'abc1234' }, unset, [ours.pem]);
		expect(dev).toMatchObject({ state: 'missing', development: true });
	});

	it('calls a build development when its version is unset, 0.0.0, or 0.0.0-<commit>', () => {
		expect(isDevelopmentVersion(undefined)).toBe(true);
		expect(isDevelopmentVersion('')).toBe(true);
		expect(isDevelopmentVersion('0.0.0')).toBe(true);
		expect(isDevelopmentVersion('0.0.0-abc1234')).toBe(true);
		expect(isDevelopmentVersion('0.42.4')).toBe(false);
		expect(isDevelopmentVersion('0.0.1')).toBe(false);
	});

	it('warns for a versioned build without a valid stamp, never for a development or signed one', () => {
		const stamp = (state: ReleaseStamp['state'], development: boolean): ReleaseStamp => ({
			state,
			development,
			reason: state === 'invalid' ? 'bad-signature' : null,
			error: null,
			path: '/stamp',
			statement: null,
			signature: null
		});
		expect(stampWarrantsWarning(stamp('missing', false))).toBe(true);
		expect(stampWarrantsWarning(stamp('invalid', false))).toBe(true);
		expect(stampWarrantsWarning(stamp('missing', true))).toBe(false);
		expect(stampWarrantsWarning(stamp('invalid', true))).toBe(false);
		expect(stampWarrantsWarning(stamp('signed', false))).toBe(false);
	});

	it('verifies with the embedded release key by default, which no test key satisfies', () => {
		const text = stampText('flaiover 1.2.3 abc1234');
		const stamp = verifyReleaseStamp({ kind: 'read', path: '/stamp', text }, build);
		expect(stamp.reason).toBe('bad-signature');
	});
});

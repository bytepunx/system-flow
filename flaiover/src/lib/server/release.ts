import { verify } from 'node:crypto';
import { readFileSync } from 'node:fs';

// The release key's public half (ADR-0070, S-0232): the ECDSA P-256 key whose private half signs
// flai's checksums.txt and the dashboard's digest list in CI. flai carries the same PEM in
// flai/internal/buildinfo/cosign.pub, and the test keeps the two the same; the release-key runbook in
// docs/operators records its fingerprint.
export const RELEASE_PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEXGMlBGpRPBSMMOcQrk81afzRbx/M
izI2Hxp0X5OBfTx1U0q1kqKv2oMeL9d0FRsHo7kNz/pNi/gEs82LtNPbJA==
-----END PUBLIC KEY-----
`;

// The release stamp (S-0235, release-signing.md § A release stamp in hello): CI signs the line
// `flaiover <version> <commit>` with the release key, and the image keeps it as a file of two
// lines, the statement and the base64 of its ASN.1 DER signature, as cosign sign-blob writes it.

/** Where the image keeps its release stamp unless FLAIOVER_RELEASE_STAMP names another file. */
export const RELEASE_STAMP_FILE = '/app/release-stamp';

/** The component name a dashboard's stamp must carry. */
export const STAMP_COMPONENT = 'flaiover';

export type StampState = 'signed' | 'missing' | 'invalid';

export type StampReason =
	| 'unreadable'
	| 'malformed'
	| 'bad-signature'
	| 'wrong-component'
	| 'version-mismatch'
	| 'commit-mismatch';

/** What the dashboard knows of its own release stamp, checked once at start. */
export type ReleaseStamp = {
	state: StampState;
	/** Why the stamp is invalid; null unless state is invalid. */
	reason: StampReason | null;
	/** The read error, for an unreadable file; null otherwise. */
	error: string | null;
	path: string;
	statement: string | null;
	signature: string | null;
	/** A development build: one whose FLAIOVER_VERSION is unset, 0.0.0, or 0.0.0-<anything>. */
	development: boolean;
};

/** What reading the stamp file gave: its text, nothing there, or an error. */
export type StampFile =
	| { kind: 'read'; path: string; text: string }
	| { kind: 'absent'; path: string }
	| { kind: 'unreadable'; path: string; error: string };

/** The build the stamp must name, and whether it is a development build. */
export type StampBuild = { version: string; commit: string; development: boolean };

/** The stamp file's path: FLAIOVER_RELEASE_STAMP, else the image's default. */
export function releaseStampPath(env: Record<string, string | undefined> = process.env): string {
	return env.FLAIOVER_RELEASE_STAMP || RELEASE_STAMP_FILE;
}

/** Whether a FLAIOVER_VERSION names a development build rather than a release. */
export function isDevelopmentVersion(version: string | undefined): boolean {
	return !version || version === '0.0.0' || version.startsWith('0.0.0-');
}

/** Read the stamp file; a file that is not there is absent, any other failure unreadable. */
export function readStampFile(path: string): StampFile {
	try {
		return { kind: 'read', path, text: readFileSync(path, 'utf8') };
	} catch (err) {
		if ((err as NodeJS.ErrnoException).code === 'ENOENT') return { kind: 'absent', path };
		return { kind: 'unreadable', path, error: String(err) };
	}
}

const BASE64 = /^[A-Za-z0-9+/]+={0,2}$/;

/** Verify a stamp file's content against the build and the release keys; no I/O. */
export function verifyReleaseStamp(
	file: StampFile,
	build: StampBuild,
	keys: readonly string[] = [RELEASE_PUBLIC_KEY]
): ReleaseStamp {
	const base = {
		path: file.path,
		development: build.development,
		reason: null,
		error: null,
		statement: null,
		signature: null
	};
	if (file.kind === 'absent') return { ...base, state: 'missing' };
	if (file.kind === 'unreadable')
		return { ...base, state: 'invalid', reason: 'unreadable', error: file.error };
	const lines = file.text.replace(/\r?\n$/, '').split(/\r?\n/);
	if (lines.length !== 2 || !lines[0] || !BASE64.test(lines[1]))
		return { ...base, state: 'invalid', reason: 'malformed' };
	const [statement, signature] = lines;
	const invalid = (reason: StampReason): ReleaseStamp => ({
		...base,
		state: 'invalid',
		reason,
		statement,
		signature
	});
	const fields = statement.split(' ');
	if (fields.length !== 3) return invalid('malformed');
	const signed = keys.some((key) =>
		verify('sha256', Buffer.from(statement), key, Buffer.from(signature, 'base64'))
	);
	if (!signed) return invalid('bad-signature');
	const [component, version, commit] = fields;
	if (component !== STAMP_COMPONENT) return invalid('wrong-component');
	if (version !== build.version) return invalid('version-mismatch');
	if (commit !== build.commit) return invalid('commit-mismatch');
	return { ...base, state: 'signed', statement, signature };
}

/** Read and verify this dashboard's own stamp, for the version and commit it reports. */
export function checkReleaseStamp(
	build: { version: string; commit: string },
	env: Record<string, string | undefined> = process.env,
	keys: readonly string[] = [RELEASE_PUBLIC_KEY]
): ReleaseStamp {
	const development = isDevelopmentVersion(env.FLAIOVER_VERSION);
	return verifyReleaseStamp(readStampFile(releaseStampPath(env)), { ...build, development }, keys);
}

/** Whether the server warns at start: a build with a version and no valid stamp. */
export function stampWarrantsWarning(stamp: ReleaseStamp): boolean {
	return !stamp.development && stamp.state !== 'signed';
}

// The release key's public half (ADR-0070, S-0232): the ECDSA P-256 key whose private half signs
// flai's checksums.txt and the dashboard's digest list in CI. flai carries the same PEM in
// flai/internal/buildinfo/cosign.pub, and the test keeps the two the same; the release-key runbook in
// docs/operators records its fingerprint.
export const RELEASE_PUBLIC_KEY = `-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEXGMlBGpRPBSMMOcQrk81afzRbx/M
izI2Hxp0X5OBfTx1U0q1kqKv2oMeL9d0FRsHo7kNz/pNi/gEs82LtNPbJA==
-----END PUBLIC KEY-----
`;

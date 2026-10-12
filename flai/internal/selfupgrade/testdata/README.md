# selfupgrade test keys

Two ECDSA P-256 key pairs that the package's tests sign `checksums.txt` with, as cosign does for a release. They are test fixtures only: no release is signed with them, and no flai trusts them outside a test.

| Files | Stands for |
|-------|------------|
| `trusted.key`, `trusted.pub` | the release key the binary trusts, which the tests pass in `Options.Keys` |
| `unknown.key`, `unknown.pub` | a release key the binary does not know, such as one that replaced it |

The private keys are SEC1 PEM (`EC PRIVATE KEY`) and the public keys PKIX PEM (`PUBLIC KEY`), as `openssl ecparam -name prime256v1 -genkey -noout` and `openssl ec -pubout` write them. The real release key's public half is `flai/internal/buildinfo/cosign.pub`.

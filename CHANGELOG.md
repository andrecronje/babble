# Change Log

## UNRELEASED

SECURITY:

* keygen: Use plain-text representation of private key, rather that X509/PEM.

FEATURES:

* crypto: Use secp256k1 curve which is also used by Ethereum and Bitcoin.

IMPROVEMENTS:

BUG FIXES:

## v0.4.1 (January 28, 2019)

IMPROVEMENTS:

* Bug fixes and refactoring pulled from develop branch

## v0.4.0 (October 14, 2018)

SECURITY:

* keygen: write keys to files instead of tty.

FEATURES:

* proxy: Introduced in-memory proxy.
* cmd: Enable reading config from file (babble.toml, .json, or .yaml)

IMPROVEMENTS:

* node: major refactoring of configuration and initialization of Babble node.
* node: Node ID is calculated from public key rather than from sorting the
peers.json file.

## v0.3.0 (September 4, 2018)

FEATURES:

* hashgraph: Replaced Leemon Baird's original "Fair" ordering method with
Lamport timestamps.
* hashgraph: Introduced the concept of Frames and Roots to enable initializing a
hashgraph from a "non-zero" state.
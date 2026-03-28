# Security Policy

## Supported Versions

Only the latest release of `threatctl` receives security updates — older
releases are not backported.

## Reporting a Vulnerability

Please do not open a public GitHub issue for security vulnerabilities.
To report a vulnerability, use one of the following methods:

- **GitHub Private Security Advisory** — open a private advisory in this
  repository (preferred).
- **Email** — contact the maintainers via the address listed in the
  repository's GitHub profile.

### What to include in your report

- A clear description of the vulnerability and the potential impact.
- Steps to reproduce or a proof-of-concept.
- The version(s) affected.
- Any suggested mitigation or fix (optional).

### What to expect

| Timeline        | Action |
| --------------- | ------ |
| Within 48 hours | Acknowledgement of the report. |
| Within 7 days   | Initial assessment and severity classification. |
| Within 30 days  | Fix released or a mitigation plan communicated. |

If the vulnerability is accepted, a CVE will be requested and a patched
release will be published along with a public disclosure in the GitHub
Security Advisories tab.

If the vulnerability is declined, you will receive an explanation of why
it was not considered a security issue.

## Scope

The following are considered in scope:

- Code within this repository (`cmd/`, `internal/`, `pkg/`).
- Official release artifacts produced by this project.

Dependencies (Go modules, base images) that are vulnerable should be
reported to their respective upstream projects. We keep dependencies up to
date via Dependabot.

## Out of scope

- Issues in third-party services or cloud providers.
- Vulnerabilities in dependencies where an upstream fix is required (we
  will coordinate with maintainers when possible).

## Security policy contact

Create a private security advisory on this GitHub repository or contact
the maintainers by email.

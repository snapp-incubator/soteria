# Change log

# Unreleased - (yyyy-mm-dd)

---

## July 17, 2021

### Added

- Add superuser token support

### Changed

- Refactor internal packages

### Fixed

## May 31, 2021

### Added

- Add chat topic in acl logics

### Changed

### Fixed

## May 29, 2021

### Added

- Add location sharing topics acl logic.
- Ignore tokens which contain `is_superuser` flag in their claims.
- Tracing with jeager. jeager configuration is added

| **Variable Name**              | **Type** | **Description**                               |
| ------------------------------ | -------- | --------------------------------------------- |
| `SOTERIA_TRACER_SAMPLER_TYPE`  | string   | client sampler type e.g. const, probabilistic |
| `SOTERIA_TRACER_SAMPLER_PARAM` | float    | client sampler paramer e.g. 1, 0              |

For more on sampler please refer to [this](https://www.jaegertracing.io/docs/1.22/sampling/).

### Changed

### Fixed

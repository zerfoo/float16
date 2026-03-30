# Changelog

## [1.0.0](https://github.com/zerfoo/float16/compare/v0.2.1...v1.0.0) (2026-03-30)


### Features

* **bfloat16:** add ArithmeticMode support and NaN propagation ([5ddc6a1](https://github.com/zerfoo/float16/commit/5ddc6a1df2b42ec7814adfe5b08c20156bbff8e9))
* **bfloat16:** add error handling infrastructure ([d5fa65f](https://github.com/zerfoo/float16/commit/d5fa65f5164eda4e1dbc541edff71de50d9d83eb))
* **bfloat16:** add Phase 4 math functions ([eb88d1e](https://github.com/zerfoo/float16/commit/eb88d1e3be1d7f84cf1d5cbd314ffba9eed57298))
* **bfloat16:** add Phase 5 parse and format functions ([a89cb77](https://github.com/zerfoo/float16/commit/a89cb77f01280d7ec4634e33f70e1d07cbc7cc62))


### Miscellaneous Chores

* release 1.0.0 ([376291d](https://github.com/zerfoo/float16/commit/376291df54e1ccbb17e045658d972c8f082ad2b3))

## [0.2.1](https://github.com/zerfoo/float16/compare/v0.2.0...v0.2.1) (2026-03-13)


### Bug Fixes

* **ci:** checkout tag ref in goreleaser to match tag commit ([4fa4e8b](https://github.com/zerfoo/float16/commit/4fa4e8b4c7503620d0302ab6fa31a2e43ae41b57))
* **ci:** skip release-please GitHub release, let GoReleaser handle it ([62cb9b0](https://github.com/zerfoo/float16/commit/62cb9b05631d3c2db3fec8c3cd952213f3f8cc32))

## [0.2.0](https://github.com/zerfoo/float16/compare/v0.1.0...v0.2.0) (2026-03-13)


### Features

* Add comprehensive tests to increase coverage ([faaf17e](https://github.com/zerfoo/float16/commit/faaf17e83dbdcad5b509e01992fb32beadd88f41))
* **bfloat16:** Implement core infrastructure for BFloat16 (Phase 1 complete) ([f007fc6](https://github.com/zerfoo/float16/commit/f007fc6d8b95ac7ac8d18d411a222affcebd84f2))


### Bug Fixes

* Correct ToFloat16 subnormal conversion ([33ab34d](https://github.com/zerfoo/float16/commit/33ab34d5adf46bff0ec904317f9be2913a8a4bdc))
* Correct ToFloat64 conversion ([38c084e](https://github.com/zerfoo/float16/commit/38c084e57f29792d2aba2a3f3f37566812262a7b))

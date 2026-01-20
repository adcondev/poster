# Changelog

All notable changes to Poster will be documented in this file.

## [4.4.0](https://github.com/adcondev/poster/compare/v4.3.0...v4.4.0) (2026-01-20)


### 🤖 Continuous Integration

* **gh-actions:** enhance Github Actions workflows ([#91](https://github.com/adcondev/poster/issues/91)) ([480b2ea](https://github.com/adcondev/poster/commit/480b2eaa8b7daadfb734cb7a10a29564ebb66396)), closes [#92](https://github.com/adcondev/poster/issues/92)


### ♻️ Code Refactoring

* **linters:** improve golangci-lint configuration for better analysis ([ee2d80b](https://github.com/adcondev/poster/commit/ee2d80b553763810fb95f6c08b7e0e150bb4994c))
* **poster:** optimize command buffer allocation and improve readability ([1a6e148](https://github.com/adcondev/poster/commit/1a6e14832e42bbe88eaff01ea92d4cf3656ff6bd))


### 🐛 Bug Fixes

* **github:** update release workflow and improve changelog management ([096db5d](https://github.com/adcondev/poster/commit/096db5de66d71db6610ab8cab0c9163a05efcf3c))
* **go.mod:** remove version suffix from module path ([1f85a12](https://github.com/adcondev/poster/commit/1f85a12b578c1189ec0d1797a63a7b363c4d9604))
* **go.mod:** update module path to include version v4 ([e113ed3](https://github.com/adcondev/poster/commit/e113ed3abd0fadd9c8b6af486962b0b0efe075f4))
* **npm:** restore package version from null to 4.3.0 ([79508d3](https://github.com/adcondev/poster/commit/79508d3fc3c51c09f8e0767ff64ff48d2d6c10c1))
* **npm:** restore package version from null to 4.3.0 ([9517910](https://github.com/adcondev/poster/commit/95179103055e1b64caa3b210b289729e87663195))
* **npm:** update package name and dependencies for poster library ([#94](https://github.com/adcondev/poster/issues/94)) ([ac2ce69](https://github.com/adcondev/poster/commit/ac2ce69b3ff279e73ee962673afa752876e9d80f)), closes [#92](https://github.com/adcondev/poster/issues/92)
* **Taskfile:** update golangci-lint command to use absolute path ([9492552](https://github.com/adcondev/poster/commit/9492552e545340ef09aa6c3529b44643d394c688))


### ✨ Features

* **connection:** enhance printer listing functionality for Windows ([1794548](https://github.com/adcondev/poster/commit/1794548ad27c2bb2671e752aee24545eb2503667))
* **connection:** enhance printer listing functionality for Windows ([#100](https://github.com/adcondev/poster/issues/100)) ([48118b8](https://github.com/adcondev/poster/commit/48118b8c117374e68e41518ded6c6e09ee837604)), closes [/#diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34R22-R26](https://github.com/adcondev/poster/issues/diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34R22-R26) [/#diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34R76](https://github.com/adcondev/poster/issues/diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34R76) [/#diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34L100-R150](https://github.com/adcondev/poster/issues/diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34L100-R150) [/#diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34L141-R165](https://github.com/adcondev/poster/issues/diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34L141-R165) [/#diff-7a68c862ed13ecb99f59c4f61a92bbbc265afe66afa76b17bb9739a8cce7cab1L2](https://github.com/adcondev/poster/issues/diff-7a68c862ed13ecb99f59c4f61a92bbbc265afe66afa76b17bb9739a8cce7cab1L2) [/#diff-7a68c862ed13ecb99f59c4f61a92bbbc265afe66afa76b17bb9739a8cce7cab1L18-R45](https://github.com/adcondev/poster/issues/diff-7a68c862ed13ecb99f59c4f61a92bbbc265afe66afa76b17bb9739a8cce7cab1L18-R45) [/#diff-02effedb378056307ac3c7278d22cf5d4e84596b606179bd8e550ab1e95cb4deL21-R38](https://github.com/adcondev/poster/issues/diff-02effedb378056307ac3c7278d22cf5d4e84596b606179bd8e550ab1e95cb4deL21-R38) [/#diff-02effedb378056307ac3c7278d22cf5d4e84596b606179bd8e550ab1e95cb4deL54-R63](https://github.com/adcondev/poster/issues/diff-02effedb378056307ac3c7278d22cf5d4e84596b606179bd8e550ab1e95cb4deL54-R63) [/#diff-02effedb378056307ac3c7278d22cf5d4e84596b606179bd8e550ab1e95cb4deL73-R103](https://github.com/adcondev/poster/issues/diff-02effedb378056307ac3c7278d22cf5d4e84596b606179bd8e550ab1e95cb4deL73-R103) [/#diff-02effedb378056307ac3c7278d22cf5d4e84596b606179bd8e550ab1e95cb4deR148-R153](https://github.com/adcondev/poster/issues/diff-02effedb378056307ac3c7278d22cf5d4e84596b606179bd8e550ab1e95cb4deR148-R153) [/#diff-02effedb378056307ac3c7278d22cf5d4e84596b606179bd8e550ab1e95cb4deR164-R228](https://github.com/adcondev/poster/issues/diff-02effedb378056307ac3c7278d22cf5d4e84596b606179bd8e550ab1e95cb4deR164-R228) [/#diff-85b77b574161b5fe9ba5d18209e374200eae93e42abf3286dd6845915d624faaR1-R6](https://github.com/adcondev/poster/issues/diff-85b77b574161b5fe9ba5d18209e374200eae93e42abf3286dd6845915d624faaR1-R6) [/#diff-85b77b574161b5fe9ba5d18209e374200eae93e42abf3286dd6845915d624faaL15-R66](https://github.com/adcondev/poster/issues/diff-85b77b574161b5fe9ba5d18209e374200eae93e42abf3286dd6845915d624faaL15-R66)
* **connection:** implement printer enumeration and details for Windows ([e7a9735](https://github.com/adcondev/poster/commit/e7a9735d1627f462023a20a3ab849de11fd1ce60))
* **schema:** add PulseCommand and BeepCommand to document schema ([5edba5e](https://github.com/adcondev/poster/commit/5edba5e4a9d9c0422cf3f959ee987ee1f1a468fc))

## [4.3.0](https://github.com/adcondev/poster/compare/v4.2.0...v4.3.0) (2025-12-19)

### ✨ Features

* **tables:** implement automatic column width reduction for table
  overflow ([#90](https://github.com/adcondev/poster/issues/90)) ([1e70d00](https://github.com/adcondev/poster/commit/1e70d009a9fe228760c06467d25f1b9d601c4ec1))

## [4.2.0](https://github.com/adcondev/poster/compare/v4.1.0...v4.2.0) (2025-12-17)


### 📝 Documentation

* **poster:** Rename project to Poster and expand `LEARNING.md` with detailed technical architecture, new features like
  the visual emulator, and enhanced ESC/POS command
  support. ([19157f2](https://github.com/adcondev/poster/commit/19157f2de9f6073f280fca75e62f47764fbb2b4c))


### ⚡ Performance

* **emulator:** optimize ToImage method for faster bitmap
  rendering ([473ba79](https://github.com/adcondev/poster/commit/473ba799864c8c0d0b6df156b624415e050c3b6c))


### ✅ Tests

* **emulator:** add tests for AutoAdjustCursorOnScale functionality and image printing
  methods ([3fd0f80](https://github.com/adcondev/poster/commit/3fd0f80542814a3e2958df74b5dc7f1d1f8eebf1))


### ✨ Features

* **config:** add AutoAdjustCursorOnScale option and update default DPI
  settings ([0a6a3c7](https://github.com/adcondev/poster/commit/0a6a3c711353668fce69615292946301b4430f44))
* **fonts:** add caching for scaled font faces and clear cache
  functionality ([9adb600](https://github.com/adcondev/poster/commit/9adb6000dcec23ef34fb7adc6c4a81c78fc62b75))
* **graphics:** add image rendering capabilities to the
  emulator ([092b722](https://github.com/adcondev/poster/commit/092b722b3daa5e8748fcf4af685917766e971991))
* **graphics:** add image rendering capabilities to the
  emulator ([#89](https://github.com/adcondev/poster/issues/89)) ([6a0c896](https://github.com/adcondev/poster/commit/6a0c8960d27fdb14569e1574751927b4e8233a40)),
  closes [/#diff-d33d3a77b7fb25f5fa2eb04de161c22a84bc2b3cd8e7ebc4519093f1ba9e077dR1-R294](https://github.com/adcondev///issues/diff-d33d3a77b7fb25f5fa2eb04de161c22a84bc2b3cd8e7ebc4519093f1ba9e077dR1-R294) [/#diff-0cfaec7ad648f145a6bca7504542f19c11ee0ad28fd340d2534f6b51c8f8132cR32](https://github.com/adcondev///issues/diff-0cfaec7ad648f145a6bca7504542f19c11ee0ad28fd340d2534f6b51c8f8132cR32) [/#diff-0cfaec7ad648f145a6bca7504542f19c11ee0ad28fd340d2534f6b51c8f8132cR86](https://github.com/adcondev///issues/diff-0cfaec7ad648f145a6bca7504542f19c11ee0ad28fd340d2534f6b51c8f8132cR86) [/#diff-0cfaec7ad648f145a6bca7504542f19c11ee0ad28fd340d2534f6b51c8f8132cR111](https://github.com/adcondev///issues/diff-0cfaec7ad648f145a6bca7504542f19c11ee0ad28fd340d2534f6b51c8f8132cR111) [/#diff-0cfaec7ad648f145a6bca7504542f19c11ee0ad28fd340d2534f6b51c8f8132cR281-R314](https://github.com/adcondev///issues/diff-0cfaec7ad648f145a6bca7504542f19c11ee0ad28fd340d2534f6b51c8f8132cR281-R314) [/#diff-980eee076fe58ca002f2e7d1970113ad6ca0825c2fb3f29ddbe6bf2d87f02b10R1-R168](https://github.com/adcondev///issues/diff-980eee076fe58ca002f2e7d1970113ad6ca0825c2fb3f29ddbe6bf2d87f02b10R1-R168) [/#diff-980eee076fe58ca002f2e7d1970113ad6ca0825c2fb3f29ddbe6bf2d87f02b10R1-R168](https://github.com/adcondev///issues/diff-980eee076fe58ca002f2e7d1970113ad6ca0825c2fb3f29ddbe6bf2d87f02b10R1-R168) [/#diff-897f07c4bba8e42c7d53c14563f93c9eedd993d316d3173cd33e79df8596c25cL124-R136](https://github.com/adcondev///issues/diff-897f07c4bba8e42c7d53c14563f93c9eedd993d316d3173cd33e79df8596c25cL124-R136) [/#diff-687d92981816c4e009e0729b938d4f7b81aab7da50c3758e4c21c76d6a194cd0R67-R87](https://github.com/adcondev///issues/diff-687d92981816c4e009e0729b938d4f7b81aab7da50c3758e4c21c76d6a194cd0R67-R87) [/#diff-0cfaec7ad648f145a6bca7504542f19c11ee0ad28fd340d2534f6b51c8f8132cR12](https://github.com/adcondev///issues/diff-0cfaec7ad648f145a6bca7504542f19c11ee0ad28fd340d2534f6b51c8f8132cR12) [/#diff-897f07c4bba8e42c7d53c14563f93c9eedd993d316d3173cd33e79df8596c25cL6-L11](https://github.com/adcondev///issues/diff-897f07c4bba8e42c7d53c14563f93c9eedd993d316d3173cd33e79df8596c25cL6-L11)

## [4.1.0](https://github.com/adcondev/poster/compare/v4.0.0...v4.1.0) (2025-12-12)


### 📝 Documentation

* **github:** optimize
  images ([#87](https://github.com/adcondev/poster/issues/87)) ([6ed6dba](https://github.com/adcondev/poster/commit/6ed6dba5a87cdf7ee97a3d0f72116c83e3b289b3))


### ♻️ Code Refactoring

* **executor:** implement text styling and formatting
  functions ([bf334e5](https://github.com/adcondev/poster/commit/bf334e599a0d771364d43977af22907ebf26cfef))


### 🐛 Bug Fixes

* **tests:** fix formatting issues in test error
  messages ([0dadd1c](https://github.com/adcondev/poster/commit/0dadd1c5e99af92cb5cf07bbb7da5840a9f83914))


### ✅ Tests

* **builder:** add unit tests for barcode, image, QR, and raw command
  builders ([883c4c2](https://github.com/adcondev/poster/commit/883c4c20a4bfefa516704dbdfaa213fcaf244d5b))
* **builder:** enhance barcode and image builder tests with default checks and alignment
  validation ([0ac2537](https://github.com/adcondev/poster/commit/0ac25377dc71c0e3401bb1ab89df2fd54e3ac52c))
* **executor:** add barcode command structure and handler
  tests ([cadf5f7](https://github.com/adcondev/poster/commit/cadf5f7542f85338a448e665116f166ae1d7c1b2))


### ✨ Features

* **builder:** add pulse and beep command structures with default
  values ([411c497](https://github.com/adcondev/poster/commit/411c497a215d306674e996814175476450e1cb4c))
* **executor:** add pulse and beep command handlers with default
  values ([53ed6af](https://github.com/adcondev/poster/commit/53ed6af2e6e487b6934c1810c2cfa244398efd16))
* **executor:** add pulse and beep command handlers with default
  values ([#88](https://github.com/adcondev/poster/issues/88)) ([8fb04a9](https://github.com/adcondev/poster/commit/8fb04a93a03e3519d61781fbf9b34d1b5de3d77f)),
  closes [/#diff-70cac1a3fa258e56c51e09b23d54349833efdac2f84859f74226e09850e9c8e4L4-R31](https://github.com/adcondev///issues/diff-70cac1a3fa258e56c51e09b23d54349833efdac2f84859f74226e09850e9c8e4L4-R31) [/#diff-a160129cc08fbc35984cb3a8276694b1635d2940f5fac1c0f024b13311edcff7L117-L134](https://github.com/adcondev///issues/diff-a160129cc08fbc35984cb3a8276694b1635d2940f5fac1c0f024b13311edcff7L117-L134) [/#diff-36043890c52c8201a8bc84238c219be45ce07bf172d89f693c7c54ffe70d046eR384-R400](https://github.com/adcondev///issues/diff-36043890c52c8201a8bc84238c219be45ce07bf172d89f693c7c54ffe70d046eR384-R400) [/#diff-c1473fbe2c123dff34107f23da1a509dd9c8195e70fad5adeb41225f254677e5L112-R112](https://github.com/adcondev///issues/diff-c1473fbe2c123dff34107f23da1a509dd9c8195e70fad5adeb41225f254677e5L112-R112) [/#diff-c1473fbe2c123dff34107f23da1a509dd9c8195e70fad5adeb41225f254677e5R138-R156](https://github.com/adcondev///issues/diff-c1473fbe2c123dff34107f23da1a509dd9c8195e70fad5adeb41225f254677e5R138-R156) [/#diff-1158ffa9f9cfd564d24cbdfaf0dbdac24e67d3e25d89d655b20ae0deb3d00a70R1-R170](https://github.com/adcondev///issues/diff-1158ffa9f9cfd564d24cbdfaf0dbdac24e67d3e25d89d655b20ae0deb3d00a70R1-R170) [/#diff-699810127ea3e65eaaac1072424265bc55eb04b1e6e1d16c1fb9ab12f98a3167R1-R110](https://github.com/adcondev///issues/diff-699810127ea3e65eaaac1072424265bc55eb04b1e6e1d16c1fb9ab12f98a3167R1-R110) [/#diff-0b2da2a163d66decd44372f48ba3d56d6be062ac1bf43301cace620e3e4e24d9R1-R102](https://github.com/adcondev///issues/diff-0b2da2a163d66decd44372f48ba3d56d6be062ac1bf43301cace620e3e4e24d9R1-R102) [/#diff-4215faad80ec87078126ed36169d3965bc97d7729f862dcde53f0664757eab02R1-R109](https://github.com/adcondev///issues/diff-4215faad80ec87078126ed36169d3965bc97d7729f862dcde53f0664757eab02R1-R109) [/#diff-ef6fbbead24bab5e7947b465fab0421f836ffd735c1aefea7c0041afdc861f90R1-R98](https://github.com/adcondev///issues/diff-ef6fbbead24bab5e7947b465fab0421f836ffd735c1aefea7c0041afdc861f90R1-R98) [/#diff-c1765478c80d6bd2fc7368a564489ecaf1dd9776c6743a9b5cac12f7702615afL1](https://github.com/adcondev///issues/diff-c1765478c80d6bd2fc7368a564489ecaf1dd9776c6743a9b5cac12f7702615afL1) [/#diff-5a7ba45dd092108eaaf75887b78b7c56c1722fe0dabbab870280da23008306f7R50-R51](https://github.com/adcondev///issues/diff-5a7ba45dd092108eaaf75887b78b7c56c1722fe0dabbab870280da23008306f7R50-R51)
* **printer:** add profile access to PrinterActions interface and
  MockPrinter ([6195394](https://github.com/adcondev/poster/commit/61953946467dc3d237aadc2576e9d64b8845e319))
* **service:** add PrinterActions interface and MockPrinter
  implementation ([64f2165](https://github.com/adcondev/poster/commit/64f21653de943eed7cf88b65e983c540e3bde448))

## [4.0.0](https://github.com/adcondev/poster/compare/v3.6.1...v4.0.0) (2025-12-10)


### ⚠ BREAKING CHANGES

* **poster:** add ESC/POS emulator functionality with receipt generation (#86)
* **poster:** update package references from poster to poster

### 🐛 Bug Fixes

* update
  pkg/document/executor/table_handler.go ([ef8f9d5](https://github.com/adcondev/poster/commit/ef8f9d57f6c7531f6ec35868430c77cf0f3ccc3f))


### 📝 Documentation

* **poster:** update README and .gitignore for poster image
  inclusion ([ca7432f](https://github.com/adcondev/poster/commit/ca7432ffcf7b866f7b3784f31c163834a78c8f66))


### ♻️ Code Refactoring

* **constants:** introduce constants for paper dimensions and rendering
  parameters ([e4096b9](https://github.com/adcondev/poster/commit/e4096b984d7a5c49b0d5172128909ea257047ea2))
* **tables:** implement default alignment and configuration options for table
  formatting ([3c35b3c](https://github.com/adcondev/poster/commit/3c35b3c9bd1ec501bbf5366c073fb0632d01a0ac))


### ✨ Features

* **emulator:** add ESC/POS emulator functionality with receipt
  generation ([d08ccc5](https://github.com/adcondev/poster/commit/d08ccc51062ff5f690a45e04d5dc11c26cfbc348))
* **poster:** add ESC/POS emulator functionality with receipt
  generation ([#86](https://github.com/adcondev/poster/issues/86)) ([03dcf8b](https://github.com/adcondev/poster/commit/03dcf8b8653c89bb43b6b78159d9b6d17b211336)),
  closes [/#diff-33ef32bf6c23acb95f5902d7097b7a1d5128ca061167ec0716715b0b9eeaa5f6L1-R10](https://github.com/adcondev///issues/diff-33ef32bf6c23acb95f5902d7097b7a1d5128ca061167ec0716715b0b9eeaa5f6L1-R10) [/#diff-02effedb378056307ac3c7278d22cf5d4e84596b606179bd8e550ab1e95cb4deL12-R20](https://github.com/adcondev///issues/diff-02effedb378056307ac3c7278d22cf5d4e84596b606179bd8e550ab1e95cb4deL12-R20) [/#diff-d4ee3a8f9dcd4bf2bad9d891f11121d0f2f8d49c0019afc23d1c0da3682acad3L10-R11](https://github.com/adcondev///issues/diff-d4ee3a8f9dcd4bf2bad9d891f11121d0f2f8d49c0019afc23d1c0da3682acad3L10-R11) [/#diff-85b77b574161b5fe9ba5d18209e374200eae93e42abf3286dd6845915d624faaL9-R11](https://github.com/adcondev///issues/diff-85b77b574161b5fe9ba5d18209e374200eae93e42abf3286dd6845915d624faaL9-R11) [/#diff-f07392da22522af73cbf199f87fb184e3e7f097515b4569ab59462af74815658L12-R15](https://github.com/adcondev///issues/diff-f07392da22522af73cbf199f87fb184e3e7f097515b4569ab59462af74815658L12-R15) [/#diff-809c1fbe1dfa048ffe8a5d9b35ba3337d6574f4ac2f879e3c611400519cd1203L9-R16](https://github.com/adcondev///issues/diff-809c1fbe1dfa048ffe8a5d9b35ba3337d6574f4ac2f879e3c611400519cd1203L9-R16) [/#diff-023610994a3e13e717bb2f341581b98599c4ab15975a574ea0356afbdd94326fL7-R7](https://github.com/adcondev///issues/diff-023610994a3e13e717bb2f341581b98599c4ab15975a574ea0356afbdd94326fL7-R7) [/#diff-55fbc650eb407e0f997756721dd12c5acdc64a69c6bc405267900d74c7bb7c59L1-R12](https://github.com/adcondev///issues/diff-55fbc650eb407e0f997756721dd12c5acdc64a69c6bc405267900d74c7bb7c59L1-R12) [/#diff-fecc96d8fa561aa1e3e1af1ae980ac8f6d4b35bb5a2253eb4b59a4f79d5613edL1-R12](https://github.com/adcondev///issues/diff-fecc96d8fa561aa1e3e1af1ae980ac8f6d4b35bb5a2253eb4b59a4f79d5613edL1-R12) [/#diff-264ff1948e3bb2221f8eba6d0d2e85731b4f1c27211b8955c141e67a9f0d27caL10-R16](https://github.com/adcondev///issues/diff-264ff1948e3bb2221f8eba6d0d2e85731b4f1c27211b8955c141e67a9f0d27caL10-R16) [/#diff-af506b9fe4fcc35d7d61e3b6aba087ea5f4187c8a4a7025e40d5248dd0302abbL8-R13](https://github.com/adcondev///issues/diff-af506b9fe4fcc35d7d61e3b6aba087ea5f4187c8a4a7025e40d5248dd0302abbL8-R13) [/#diff-02effedb378056307ac3c7278d22cf5d4e84596b606179bd8e550ab1e95cb4deL1-R2](https://github.com/adcondev///issues/diff-02effedb378056307ac3c7278d22cf5d4e84596b606179bd8e550ab1e95cb4deL1-R2) [/#diff-02effedb378056307ac3c7278d22cf5d4e84596b606179bd8e550ab1e95cb4deL150-R150](https://github.com/adcondev///issues/diff-02effedb378056307ac3c7278d22cf5d4e84596b606179bd8e550ab1e95cb4deL150-R150) [/#diff-55fbc650eb407e0f997756721dd12c5acdc64a69c6bc405267900d74c7bb7c59L1-R12](https://github.com/adcondev///issues/diff-55fbc650eb407e0f997756721dd12c5acdc64a69c6bc405267900d74c7bb7c59L1-R12) [/#diff-fecc96d8fa561aa1e3e1af1ae980ac8f6d4b35bb5a2253eb4b59a4f79d5613edL1-R12](https://github.com/adcondev///issues/diff-fecc96d8fa561aa1e3e1af1ae980ac8f6d4b35bb5a2253eb4b59a4f79d5613edL1-R12) [/#diff-264ff1948e3bb2221f8eba6d0d2e85731b4f1c27211b8955c141e67a9f0d27caL1-R1](https://github.com/adcondev///issues/diff-264ff1948e3bb2221f8eba6d0d2e85731b4f1c27211b8955c141e67a9f0d27caL1-R1)
* **poster:** update package references from poster to
  poster ([03271f3](https://github.com/adcondev/poster/commit/03271f311d3838ba3a2ed58c35d09a751c0d4960))

### [3.6.1](https://github.com/adcondev/poster/compare/v3.6.0...v3.6.1) (2025-12-08)

### ♻️ Code Refactoring

* **tables:** implement default alignment and configuration options for table
  formatting ([#83](https://github.com/adcondev/poster/issues/83)) ([949fa9f](https://github.com/adcondev/poster/commit/949fa9f5c78312d1ac765c19c0ad798d90dd601b))

### 📝 Documentation

* **poster:** update
  readme ([#84](https://github.com/adcondev/poster/issues/84)) ([84af0d0](https://github.com/adcondev/poster/commit/84af0d0d505f29f22d17a43e3f6c593f1192faaa))

### 🐛 Bug Fixes

* **labels:** add JavaScript related label to auto-merge
  workflow ([7b10cf2](https://github.com/adcondev/poster/commit/7b10cf2e0e437926be4681b8a2d9253e13cbf5f4))

### 📦 Dependencies

* **npm:** bump the minor-and-patch group with 2
  updates ([62b290c](https://github.com/adcondev/poster/commit/62b290c336043f244e89ecb4814545897094c212))
* **npm:** bump the minor-and-patch group with 2
  updates ([#85](https://github.com/adcondev/poster/issues/85)) ([f863545](https://github.com/adcondev/poster/commit/f8635450b5bef14df4fba398e4e9960a58e679d1)),
  closes [conventional-changelog/commitlint#4542](https://github.com/conventional-changelog/commitlint/issues/4542) [conventional-changelog/commitlint#4548](https://github.com/conventional-changelog/commitlint/issues/4548) [conventional-changelog/commitlint#4540](https://github.com/conventional-changelog/commitlint/issues/4540) [conventional-changelog/commitlint#4559](https://github.com/conventional-changelog/commitlint/issues/4559) [conventional-changelog/commitlint#4548](https://github.com/conventional-changelog/commitlint/issues/4548) [conventional-changelog/commitlint#4540](https://github.com/conventional-changelog/commitlint/issues/4540) [conventional-changelog/commitlint#4559](https://github.com/conventional-changelog/commitlint/issues/4559) [conventional-changelog/commitlint#4542](https://github.com/conventional-changelog/commitlint/issues/4542) [conventional-changelog/commitlint#4548](https://github.com/conventional-changelog/commitlint/issues/4548) [conventional-changelog/commitlint#4540](https://github.com/conventional-changelog/commitlint/issues/4540) [conventional-changelog/commitlint#4559](https://github.com/conventional-changelog/commitlint/issues/4559) [conventional-changelog/commitlint#4548](https://github.com/conventional-changelog/commitlint/issues/4548) [conventional-changelog/commitlint#4540](https://github.com/conventional-changelog/commitlint/issues/4540) [conventional-changelog/commitlint#4559](https://github.com/conventional-changelog/commitlint/issues/4559) [conventional-changelog/commitlint#4553](https://github.com/conventional-changelog/commitlint/issues/4553) [conventional-changelog/commitlint#4551](https://github.com/conventional-changelog/commitlint/issues/4551) [conventional-changelog/commitlint#4553](https://github.com/conventional-changelog/commitlint/issues/4553) [conventional-changelog/commitlint#4551](https://github.com/conventional-changelog/commitlint/issues/4551)

## [3.6.0](https://github.com/adcondev/poster/compare/v3.5.0...v3.6.0) (2025-12-04)


### ✅ Tests

* **tables:** update method names to follow Go naming
  conventions ([97518ca](https://github.com/adcondev/poster/commit/97518cac6b807e1c8fd239ebc122f9866bf7d888))


### 🐛 Bug Fixes

* **poster:** correct alignment and default values in barcode and image
  builders ([95a7857](https://github.com/adcondev/poster/commit/95a7857a7704516b2ad228c997f7a43d427c4423))


### ✨ Features

* **barcode:** enhance barcode configuration and mapping
  functions ([db5e03b](https://github.com/adcondev/poster/commit/db5e03bdc49249e392cbabee96bb0ba0b806e174))
* **builder:** use centralized default version in document
  creation ([0beb10d](https://github.com/adcondev/poster/commit/0beb10d5279a6845f2894842b22b5faa71403457))
* **constants:** expand default constants for various
  configurations ([4d56b0e](https://github.com/adcondev/poster/commit/4d56b0e3288af00f90c229ff3c8b3e63fdc10290))
* **constants:** introduce centralized constants for alignment and default
  values ([fb97aed](https://github.com/adcondev/poster/commit/fb97aedbf8842bf91e78b2f0bf9aefc5f170e08e))
* **constants:** introduce centralized constants for alignment and default
  values ([#82](https://github.com/adcondev/poster/issues/82)) ([8255c69](https://github.com/adcondev/poster/commit/8255c69b9caca46837015b4df283ea1ce672f840))
* **executor:** add new executor package for handling JSON print
  documents ([e480f32](https://github.com/adcondev/poster/commit/e480f32ad65ff93e205bbaa2495feb05f410ed57))
* **executor:** implement handler registry for command
  management ([9ccb2be](https://github.com/adcondev/poster/commit/9ccb2be414cbe650937ba8ebfdf978b779ceb91e))
* **poster:** update default constants for image and QR code
  handling ([d0d638f](https://github.com/adcondev/poster/commit/d0d638f1839d486955a3c07b664d992b7517c9d2))
* **tables:** introduce TableBuilder for fluent API to build table
  commands ([5e54f02](https://github.com/adcondev/poster/commit/5e54f028fd83747cd016c8d6a14c8069dfcc0703))

## [3.5.0](https://github.com/adcondev/poster/compare/v3.4.0...v3.5.0) (2025-12-02)


### ✨ Features

* **executor:** add new executor package for handling JSON print
  documents ([#81](https://github.com/adcondev/poster/issues/81)) ([8183c79](https://github.com/adcondev/poster/commit/8183c797252db05a327b4d4c285b099c06171125))

## [3.4.0](https://github.com/adcondev/poster/compare/v3.3.3...v3.4.0) (2025-11-27)


### ✨ Features

* **document:** add support for raw command handling and
  documentation ([#80](https://github.com/adcondev/poster/issues/80)) ([dca16cf](https://github.com/adcondev/poster/commit/dca16cf535e099de3c7f922d969187bb35ca2d03))

### [3.3.3](https://github.com/adcondev/poster/compare/v3.3.2...v3.3.3) (2025-11-25)


### 🐛 Bug Fixes

* **security:** fix workflow does not contain permissions
  alert ([#79](https://github.com/adcondev/poster/issues/79)) ([ad945d3](https://github.com/adcondev/poster/commit/ad945d37cd3f764c4b9877d2727c3a38bde3b31a))

### [3.3.2](https://github.com/adcondev/poster/compare/v3.3.1...v3.3.2) (2025-11-25)


### 📦 Dependencies

* **gh-actions:** bump actions/checkout from 5 to
  6 ([#77](https://github.com/adcondev/poster/issues/77)) ([3a232ce](https://github.com/adcondev/poster/commit/3a232ceb69c1cb132f8381b750b3d66260020fc4))

### [3.3.1](https://github.com/adcondev/poster/compare/v3.3.0...v3.3.1) (2025-11-25)


### ✅ Tests

* **poster:** add unit tests for profile, connection and graphics
  packages ([#74](https://github.com/adcondev/poster/issues/74)) ([d56a1e5](https://github.com/adcondev/poster/commit/d56a1e5cf5d4deb33e1c413084d704a1e5c586d9))


### 📦 Dependencies

* **gomod:** bump github.com/stretchr/testify from 1.7.0 to
  1.11.1 ([#76](https://github.com/adcondev/poster/issues/76)) ([0f40661](https://github.com/adcondev/poster/commit/0f406613e8d716346b01b8116a3e10f72f6a765f))
* **gomod:** bump gopkg.in/yaml.v3 from 3.0.0 to
  3.0.1 ([#75](https://github.com/adcondev/poster/issues/75)) ([254c4cf](https://github.com/adcondev/poster/commit/254c4cf515b4c6f2ed22455a1225b09b42a07d86))

## [3.3.0](https://github.com/adcondev/poster/compare/v3.2.0...v3.3.0) (2025-11-25)


### ♻️ Code Refactoring

* **document:** update ticket JSON structure and improve text command
  handling ([5217e05](https://github.com/adcondev/poster/commit/5217e053db72deec249dea2c60f42922a4f90ef4))


### 📝 Documentation

* **document:** add JSON schema and documentation for POS printer document
  format ([ab93895](https://github.com/adcondev/poster/commit/ab938954a29bfdf5c596a8792bc014807c350bfb))


### 🐛 Bug Fixes

* **document:** fix panic by enhancing profile application and command
  handling ([2ce903c](https://github.com/adcondev/poster/commit/2ce903c985210c27bde12224603115bf636f122a))


### ✨ Features

* **barcode:** add BarcodeCommand structure for barcode
  generation ([6762a4c](https://github.com/adcondev/poster/commit/6762a4c897eb6abe8a49e56dac7bef3956442449))
* **barcode:** implement barcode command handling and
  configuration ([d76eef1](https://github.com/adcondev/poster/commit/d76eef1aaa6a6f8ef3c858f6a663ed305dff9fca))
* **barcode:** implement barcode command handling and
  configuration ([#73](https://github.com/adcondev/poster/issues/73)) ([2408c41](https://github.com/adcondev/poster/commit/2408c41e691c72673a76cab4534d4b1ebd60a936))
* **document:** add new JSON document examples for receipt and table
  commands ([75e5881](https://github.com/adcondev/poster/commit/75e58814079099e71d8e98d73392aa0f4847dc76))
* **document:** add QR and table command
  handling ([edac830](https://github.com/adcondev/poster/commit/edac830a1472f8a1cbad4d74b98105cef04f83cb))
* **document:** enhance text command structure with label
  support ([e389512](https://github.com/adcondev/poster/commit/e389512b505e7d624fc8b52c4b7c14df4921769c))
* **document:** introduce new text command structure with optional label and style
  support ([b129df6](https://github.com/adcondev/poster/commit/b129df6db305dfd450fcee9c33d54f8c64650308))

## [3.2.0](https://github.com/adcondev/poster/compare/v3.1.0...v3.2.0) (2025-11-19)


### 🐛 Bug Fixes

* **graphics:** improve logo handling in QR code
  generation ([68ae2c5](https://github.com/adcondev/poster/commit/68ae2c56adb96ae9b821152659163fb572886038))


### 📦 Dependencies

* **tables:** update golang.org/x/image and golang.org/x/text
  dependencies ([3969cd0](https://github.com/adcondev/poster/commit/3969cd02c4f664f9084a0f1c37c1c43705187e45))


### ✨ Features

* **ci:** add tables scope to commit message
  guidelines ([4cc7a16](https://github.com/adcondev/poster/commit/4cc7a1696fac15f15dd243cf917bb249468cbd65))
* **document:** add QR code and table command support in document
  builder ([8172b4d](https://github.com/adcondev/poster/commit/8172b4d45ca5bd071c2f9270bbb12380b6fe611d))
* **qrcode:** refactor image loading and command handling for ESC/POS
  printer ([fafcf7a](https://github.com/adcondev/poster/commit/fafcf7ac877003b31e4218d5ee5174516ae402ad))
* **tables:** add table generation and rendering for ESC/POS
  printers ([9089536](https://github.com/adcondev/poster/commit/90895364ba1e7317e55507915a2b53e0adc00db4))
* **tables:** add table generation and rendering for ESC/POS
  printers ([#72](https://github.com/adcondev/poster/issues/72)) ([261086e](https://github.com/adcondev/poster/commit/261086eff59db1e0e28ae85f32f6994b85d2b48e))

## [3.1.0](https://github.com/adcondev/poster/compare/v3.0.7...v3.1.0) (2025-11-18)


### ✨ Features

* **qrcode:** enhance QR code processing and image
  generation ([#71](https://github.com/adcondev/poster/issues/71)) ([c38f115](https://github.com/adcondev/poster/commit/c38f115377dbf861333f3ba51d83469e647a0b66))

### [3.0.7](https://github.com/adcondev/poster/compare/v3.0.6...v3.0.7) (2025-11-18)


### 📦 Dependencies

* **gomod:** bump the golang-x group with 2
  updates ([#70](https://github.com/adcondev/poster/issues/70)) ([7499533](https://github.com/adcondev/poster/commit/749953339a82e57e2d6e96e3150e64de887191ab))

### [3.0.6](https://github.com/adcondev/poster/compare/v3.0.5...v3.0.6) (2025-11-18)


### 📝 Documentation

* **github:** update DevOps and CI/CD
  documentation ([#68](https://github.com/adcondev/poster/issues/68)) ([afc4480](https://github.com/adcondev/poster/commit/afc4480267d74265ee6aa6f3b38135374acac0fc))


### 📦 Dependencies

* **npm:** bump js-yaml from 4.1.0 to
  4.1.1 ([#69](https://github.com/adcondev/poster/issues/69)) ([dc4b4a8](https://github.com/adcondev/poster/commit/dc4b4a8f4ee9fe6955167cf13382b870bef32727))

### [3.0.5](https://github.com/adcondev/poster/compare/v3.0.4...v3.0.5) (2025-11-10)


### 📦 Dependencies

* **gh-actions:** bump golangci/golangci-lint-action from 8 to
  9 ([#66](https://github.com/adcondev/poster/issues/66)) ([b964ed1](https://github.com/adcondev/poster/commit/b964ed1682e9251b1e8d41067ff511f043b2fedc))

### [3.0.4](https://github.com/adcondev/poster/compare/v3.0.3...v3.0.4) (2025-11-10)


### 📦 Dependencies

* **gh-actions:** bump github/codeql-action from 3 to
  4 ([#65](https://github.com/adcondev/poster/issues/65)) ([9b437df](https://github.com/adcondev/poster/commit/9b437df564430b352d4c064ae42f23e48b990431))

### [3.0.3](https://github.com/adcondev/poster/compare/v3.0.2...v3.0.3) (2025-11-10)


### 🤖 Continuous Integration

* **github:** enhance workflows with caching and test
  coverage ([7eae932](https://github.com/adcondev/poster/commit/7eae932e64ba76f1e9f0686e15e380e469fd8e58))
* **github:** enhance workflows with caching and test
  coverage ([#67](https://github.com/adcondev/poster/issues/67)) ([506a1ba](https://github.com/adcondev/poster/commit/506a1ba9ff4fbddb4a7f90dcdfe21ddde8216aff)),
  closes [/#diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fL78-R79](https://github.com/adcondev///issues/diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fL78-R79) [/#diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fL95-L121](https://github.com/adcondev///issues/diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fL95-L121) [/#diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fR116-R120](https://github.com/adcondev///issues/diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fR116-R120) [/#diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fL148-R147](https://github.com/adcondev///issues/diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fL148-R147) [/#diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fR156](https://github.com/adcondev///issues/diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fR156) [/#diff-8906835152921ef903f55779586f3a092362f65fca98df94d71801cf974ec95bR35](https://github.com/adcondev///issues/diff-8906835152921ef903f55779586f3a092362f65fca98df94d71801cf974ec95bR35) [/#diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34L1-L9](https://github.com/adcondev///issues/diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34L1-L9) [/#diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34R19-R26](https://github.com/adcondev///issues/diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34R19-R26) [/#diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34R42](https://github.com/adcondev///issues/diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34R42) [/#diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34L49-R68](https://github.com/adcondev///issues/diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34L49-R68) [/#diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34L97-R162](https://github.com/adcondev///issues/diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34L97-R162) [/#diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34R171-R184](https://github.com/adcondev///issues/diff-87db21a973eed4fef5f32b267aa60fcee5cbdf03c67fafdc2a9b553bb0b15f34R171-R184) [/#diff-cd2d359855d0301ce190f1ec3b4c572ea690c83747f6df61c9340720e3d2425eL6-R6](https://github.com/adcondev///issues/diff-cd2d359855d0301ce190f1ec3b4c572ea690c83747f6df61c9340720e3d2425eL6-R6) [/#diff-cd2d359855d0301ce190f1ec3b4c572ea690c83747f6df61c9340720e3d2425eL94-R94](https://github.com/adcondev///issues/diff-cd2d359855d0301ce190f1ec3b4c572ea690c83747f6df61c9340720e3d2425eL94-R94) [/#diff-327c9e81d7b5e65336eb4b7f9a2b51fb692ef83d205ff131895e1f30727260b9L7-R7](https://github.com/adcondev///issues/diff-327c9e81d7b5e65336eb4b7f9a2b51fb692ef83d205ff131895e1f30727260b9L7-R7) [/#diff-327c9e81d7b5e65336eb4b7f9a2b51fb692ef83d205ff131895e1f30727260b9L16-R46](https://github.com/adcondev///issues/diff-327c9e81d7b5e65336eb4b7f9a2b51fb692ef83d205ff131895e1f30727260b9L16-R46) [/#diff-327c9e81d7b5e65336eb4b7f9a2b51fb692ef83d205ff131895e1f30727260b9L68-R84](https://github.com/adcondev///issues/diff-327c9e81d7b5e65336eb4b7f9a2b51fb692ef83d205ff131895e1f30727260b9L68-R84) [/#diff-327c9e81d7b5e65336eb4b7f9a2b51fb692ef83d205ff131895e1f30727260b9L119-R126](https://github.com/adcondev///issues/diff-327c9e81d7b5e65336eb4b7f9a2b51fb692ef83d205ff131895e1f30727260b9L119-R126) [/#diff-327c9e81d7b5e65336eb4b7f9a2b51fb692ef83d205ff131895e1f30727260b9L148-R229](https://github.com/adcondev///issues/diff-327c9e81d7b5e65336eb4b7f9a2b51fb692ef83d205ff131895e1f30727260b9L148-R229) [/#diff-327c9e81d7b5e65336eb4b7f9a2b51fb692ef83d205ff131895e1f30727260b9L258-R262](https://github.com/adcondev///issues/diff-327c9e81d7b5e65336eb4b7f9a2b51fb692ef83d205ff131895e1f30727260b9L258-R262)


### 📦 Dependencies

* **deps:** bump actions/checkout from 4 to
  5 ([adadd53](https://github.com/adcondev/poster/commit/adadd53dbfb4f427aac6a08ca347894d9efca3cb))
* **gh-actions:** bump actions/checkout from 4 to
  5 ([#64](https://github.com/adcondev/poster/issues/64)) ([eb69625](https://github.com/adcondev/poster/commit/eb6962548c7a6a0f5d8d2011e9f79e280bf19965))

### [3.0.2](https://github.com/adcondev/poster/compare/v3.0.1...v3.0.2) (2025-11-10)


### 📦 Dependencies

* **npm:** bump @commitlint/cli from 19.8.1 to
  20.1.0 ([#63](https://github.com/adcondev/poster/issues/63)) ([1d19593](https://github.com/adcondev/poster/commit/1d1959395768b428b075502cd21c2168e21918b5))

### [3.0.1](https://github.com/adcondev/poster/compare/v3.0.0...v3.0.1) (2025-11-10)


### 🤖 Continuous Integration

* **dependabot:** enhance auto-merge workflow with timeout and success
  criteria ([5f74df6](https://github.com/adcondev/poster/commit/5f74df60c908f8361b777fe9b4337200538a3052))
* **github:** add new scopes for gomod, npm, and
  gh-actions ([39746f4](https://github.com/adcondev/poster/commit/39746f4dc5f920e163582d627415ad1fdcbf3285))


### 📦 Dependencies

* **npm:** bump @commitlint/config-conventional from 19.8.1 to
  20.0.0 ([#62](https://github.com/adcondev/poster/issues/62)) ([277a12c](https://github.com/adcondev/poster/commit/277a12ca66553a20cbd94a61edca87cba494f1f9))

## [3.0.0](https://github.com/adcondev/poster/compare/v2.3.0...v3.0.0) (2025-11-07)


### ⚠ BREAKING CHANGES

* **commands:** rename controllers to commands (#61)
* **escpos:** rename controllers to commands

### 📝 Documentation

* add LEARNING.md and update
  README.md ([21bc0f5](https://github.com/adcondev/poster/commit/21bc0f525faee77b6a32df265fd6451622fec28c))


### 🐛 Bug Fixes

* update LEARNING.md ([99a7a85](https://github.com/adcondev/poster/commit/99a7a858cce47007610f21c23448e355ced51204))


### ♻️ Code Refactoring

* **base64:** rename example file for
  clarity ([8d6533d](https://github.com/adcondev/poster/commit/8d6533db418a2a58ede14fc245a3e285a3fb0d16))
* **commands:** improve command documentation and
  comments ([764a58f](https://github.com/adcondev/poster/commit/764a58feafe8f4b68a87740fe77e2c608c5b4363))


### 🤖 Continuous Integration

* **github:** update workflow to dynamically find example
  directories ([cc643b8](https://github.com/adcondev/poster/commit/cc643b88009eb569ed2e745d9c384827daf866d2))


### ✨ Features

* **commands:** rename controllers to
  commands ([#61](https://github.com/adcondev/poster/issues/61)) ([fa786f3](https://github.com/adcondev/poster/commit/fa786f3e5e445935ae2a193ffe15ac74f24c6272)),
  closes [/#diff-8f12d5f5467d5a9e7d05741bfaa3151879b70dc9b12ed2da94e21fcbbd80a986L14-L16](https://github.com/adcondev///issues/diff-8f12d5f5467d5a9e7d05741bfaa3151879b70dc9b12ed2da94e21fcbbd80a986L14-L16) [/#diff-8f12d5f5467d5a9e7d05741bfaa3151879b70dc9b12ed2da94e21fcbbd80a986L27-R69](https://github.com/adcondev///issues/diff-8f12d5f5467d5a9e7d05741bfaa3151879b70dc9b12ed2da94e21fcbbd80a986L27-R69) [/#diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fL48-L50](https://github.com/adcondev///issues/diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fL48-L50) [/#diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fL61-R70](https://github.com/adcondev///issues/diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fL61-R70) [/#diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fR103-R118](https://github.com/adcondev///issues/diff-b803fcb7f17ed9235f1e5cb1fcd2f5d3b2838429d4368ae4c57ce4436577f03fR103-R118) [/#diff-89708b194914fd416481a81522862b10b3cd5f123be3eb550e6a1de67a01765aL6-R6](https://github.com/adcondev///issues/diff-89708b194914fd416481a81522862b10b3cd5f123be3eb550e6a1de67a01765aL6-R6) [/#diff-89708b194914fd416481a81522862b10b3cd5f123be3eb550e6a1de67a01765aL240-R242](https://github.com/adcondev///issues/diff-89708b194914fd416481a81522862b10b3cd5f123be3eb550e6a1de67a01765aL240-R242) [/#diff-89708b194914fd416481a81522862b10b3cd5f123be3eb550e6a1de67a01765aL269-R269](https://github.com/adcondev///issues/diff-89708b194914fd416481a81522862b10b3cd5f123be3eb550e6a1de67a01765aL269-R269) [/#diff-5c25f228ee796b995a56de464ea307b5cc1eeaaffb990f9329fe91a063fc3429L6-R6](https://github.com/adcondev///issues/diff-5c25f228ee796b995a56de464ea307b5cc1eeaaffb990f9329fe91a063fc3429L6-R6) [/#diff-cd2d359855d0301ce190f1ec3b4c572ea690c83747f6df61c9340720e3d2425eL26-R67](https://github.com/adcondev///issues/diff-cd2d359855d0301ce190f1ec3b4c572ea690c83747f6df61c9340720e3d2425eL26-R67) [/#diff-6179837f7df53a6f05c522b6b7bb566d484d5465d9894fb04910dd08bb40dcc9R8-R62](https://github.com/adcondev///issues/diff-6179837f7df53a6f05c522b6b7bb566d484d5465d9894fb04910dd08bb40dcc9R8-R62) [/#diff-6179837f7df53a6f05c522b6b7bb566d484d5465d9894fb04910dd08bb40dcc9R71-R83](https://github.com/adcondev///issues/diff-6179837f7df53a6f05c522b6b7bb566d484d5465d9894fb04910dd08bb40dcc9R71-R83) [/#diff-7a68c862ed13ecb99f59c4f61a92bbbc265afe66afa76b17bb9739a8cce7cab1L2-R2](https://github.com/adcondev///issues/diff-7a68c862ed13ecb99f59c4f61a92bbbc265afe66afa76b17bb9739a8cce7cab1L2-R2) [/#diff-7a68c862ed13ecb99f59c4f61a92bbbc265afe66afa76b17bb9739a8cce7cab1L12-R27](https://github.com/adcondev///issues/diff-7a68c862ed13ecb99f59c4f61a92bbbc265afe66afa76b17bb9739a8cce7cab1L12-R27)
* **escpos:** rename controllers to
  commands ([1aa49fc](https://github.com/adcondev/poster/commit/1aa49fc332c9a369d524d8f15e63dd9cfa3f1ab1))

## [2.3.0](https://github.com/adcondev/poster/compare/v2.2.1...v2.3.0) (2025-11-07)


### ✨ Features

* **escpos:** add document building and execution
  features ([#58](https://github.com/adcondev/poster/issues/58)) ([3e89d73](https://github.com/adcondev/poster/commit/3e89d734d6f86221088cb7ce9e0c64f2cf842864)),
  closes [#59](https://github.com/adcondev/poster/issues/59)

### [2.2.1](https://github.com/adcondev/poster/compare/v2.2.0...v2.2.1) (2025-11-06)


### 🐛 Bug Fixes

* **profile:** improve logging for unsupported code
  tables ([#57](https://github.com/adcondev/poster/issues/57)) ([b886820](https://github.com/adcondev/poster/commit/b886820fa250fa03a20189c54ed1da80fefb674d))

## [2.2.0](https://github.com/adcondev/poster/compare/v2.1.0...v2.2.0) (2025-11-06)


### 🐛 Bug Fixes

* **profile:** enhance encoding support and error
  handling ([b28539b](https://github.com/adcondev/poster/commit/b28539b1a312af0bdaf8431bbe8f8c985aaea6e4))


### ✨ Features

* **character:** add encoding support for character
  tables ([b6a2a6f](https://github.com/adcondev/poster/commit/b6a2a6f98c248898ff09090eb683cba73556b40d))
* **escpos:** add graphics base64 image printing and autoencoding for printer code
  tables ([#56](https://github.com/adcondev/poster/issues/56)) ([d636895](https://github.com/adcondev/poster/commit/d63689555d53a42710014cf0db064108192a4e13)),
  closes [/#diff-562156e83a675c98d4982e84e8336971ca88db93fecfcceb69cd0fc1ca6fca18L28-R45](https://github.com/adcondev///issues/diff-562156e83a675c98d4982e84e8336971ca88db93fecfcceb69cd0fc1ca6fca18L28-R45) [/#diff-562156e83a675c98d4982e84e8336971ca88db93fecfcceb69cd0fc1ca6fca18L64-R56](https://github.com/adcondev///issues/diff-562156e83a675c98d4982e84e8336971ca88db93fecfcceb69cd0fc1ca6fca18L64-R56) [/#diff-562156e83a675c98d4982e84e8336971ca88db93fecfcceb69cd0fc1ca6fca18L84-R76](https://github.com/adcondev///issues/diff-562156e83a675c98d4982e84e8336971ca88db93fecfcceb69cd0fc1ca6fca18L84-R76) [/#diff-562156e83a675c98d4982e84e8336971ca88db93fecfcceb69cd0fc1ca6fca18L116-L125](https://github.com/adcondev///issues/diff-562156e83a675c98d4982e84e8336971ca88db93fecfcceb69cd0fc1ca6fca18L116-L125) [/#diff-a1c3ae8fd8f99a288574b99ee6b6d1b0fcfc7b6291679c91a9b2cd2b26631f74L47-R50](https://github.com/adcondev///issues/diff-a1c3ae8fd8f99a288574b99ee6b6d1b0fcfc7b6291679c91a9b2cd2b26631f74L47-R50) [/#diff-a1c3ae8fd8f99a288574b99ee6b6d1b0fcfc7b6291679c91a9b2cd2b26631f74L71-R74](https://github.com/adcondev///issues/diff-a1c3ae8fd8f99a288574b99ee6b6d1b0fcfc7b6291679c91a9b2cd2b26631f74L71-R74) [/#diff-a1c3ae8fd8f99a288574b99ee6b6d1b0fcfc7b6291679c91a9b2cd2b26631f74L80-R87](https://github.com/adcondev///issues/diff-a1c3ae8fd8f99a288574b99ee6b6d1b0fcfc7b6291679c91a9b2cd2b26631f74L80-R87) [/#diff-a1c3ae8fd8f99a288574b99ee6b6d1b0fcfc7b6291679c91a9b2cd2b26631f74R204-R213](https://github.com/adcondev///issues/diff-a1c3ae8fd8f99a288574b99ee6b6d1b0fcfc7b6291679c91a9b2cd2b26631f74R204-R213) [/#diff-af506b9fe4fcc35d7d61e3b6aba087ea5f4187c8a4a7025e40d5248dd0302abbL128-R128](https://github.com/adcondev///issues/diff-af506b9fe4fcc35d7d61e3b6aba087ea5f4187c8a4a7025e40d5248dd0302abbL128-R128) [/#diff-af506b9fe4fcc35d7d61e3b6aba087ea5f4187c8a4a7025e40d5248dd0302abbR138-R151](https://github.com/adcondev///issues/diff-af506b9fe4fcc35d7d61e3b6aba087ea5f4187c8a4a7025e40d5248dd0302abbR138-R151)
* **graphics:** add base64 image loading functionality and
  example ([aa5b68a](https://github.com/adcondev/poster/commit/aa5b68aaf4475a9f8d859fc85466106dea5c03a7))
* **graphics:** update version to 2.1.0 and
  changelog ([32b7757](https://github.com/adcondev/poster/commit/32b7757fb53bbda55a0d87f0c358442d84a5ddb7))

## [2.1.0](https://github.com/adcondev/poster/compare/v2.0.1...v2.1.0) (2025-11-05)

### ✨ Features

* **graphics:** add base64 image loading functionality and
  example ([#55](https://github.com/adcondev/poster/issues/55)) ([8abc948](https://github.com/adcondev/poster/commit/8abc948da3d68f90792f796228ec8b861be42019))

### [2.0.1](https://github.com/adcondev/poster/compare/v2.0.0...v2.0.1) (2025-11-05)


### 📝 Documentation

* **readme:** update LEARNING.md and
  README.md ([#53](https://github.com/adcondev/poster/issues/53)) ([54d39a4](https://github.com/adcondev/poster/commit/54d39a4ae9136ad3ff548d43a5285e3bb221b6a7))


### 🤖 Continuous Integration

* **escpos:** update CI configuration and add pre-commit
  hooks ([#54](https://github.com/adcondev/poster/issues/54)) ([6ee16d6](https://github.com/adcondev/poster/commit/6ee16d6be8e410d287b2c232b3146a2a461060cb))


### 📦 Dependencies

* **ci:** bump actions/setup-node from 5 to
  6 ([#49](https://github.com/adcondev/poster/issues/49)) ([925e1be](https://github.com/adcondev/poster/commit/925e1bef25466a41aaec64a4899fa5f5aa12a19d))

## [2.0.0](https://github.com/AdConDev/pos-daemon/compare/v1.8.0...v2.0.0) (2025-11-04)


### ⚠ BREAKING CHANGES

* **escpos:** refactor Protocol to Commands and enhance functionality

### ✨ Features

* **arq:** introduce modular escpos architecture and printer profiles ([86d877b](https://github.com/AdConDev/pos-daemon/commit/86d877b64813613cbd437f24f9e1c8dd2c30371f))
* **escpos:** add QR code capability to ESC/POS protocol ([2d1f7f3](https://github.com/AdConDev/pos-daemon/commit/2d1f7f3a691bd7b915c96fc3b6b860eddb38a9cd))
* **escpos:** refactor Protocol to Commands and enhance functionality ([eec9582](https://github.com/AdConDev/pos-daemon/commit/eec9582bc030b068941a20d8daabb0668ca49b95))
* **graphics:** add advanced image processing engine ([0148e8c](https://github.com/AdConDev/pos-daemon/commit/0148e8c5a8e04987d384e3b5beaef4c136b1eeb9))
* **printposition:** introduce composer and refactor print position ([9b4aeb0](https://github.com/AdConDev/pos-daemon/commit/9b4aeb0104c09c2910f603570fa1005df5f9ce78))
* **qrcode:** implement QR Code generation commands ([57571fa](https://github.com/AdConDev/pos-daemon/commit/57571fa2c4bbe4a2135d200743841b4a373226de))


### 📦 Dependencies

* remove unused QR code dependency ([f479762](https://github.com/AdConDev/pos-daemon/commit/f4797628613309fb42f644ae289a1638f905790d))

## [1.8.0](https://github.com/AdConDev/pos-daemon/compare/v1.7.0...v1.8.0) (2025-10-28)


### ✨ Features

* **taskfiles:** add Task files to automate and improve some
  processes ([#48](https://github.com/adcondev/poster/issues/48)) ([54de86c](https://github.com/AdConDev/pos-daemon/commit/54de86c6bde378206dca8c8d2f86627e45d3eff5))

## [1.7.0](https://github.com/AdConDev/pos-daemon/compare/v1.6.0...v1.7.0) (2025-10-23)


### ✨ Features

* **bitimage:** implement ESC/POS commands for bit
  images ([#45](https://github.com/adcondev/poster/issues/45)) ([1b98b1f](https://github.com/AdConDev/pos-daemon/commit/1b98b1f9587ee1121961676f516ae946db0fc4e0))

## [1.6.0](https://github.com/AdConDev/pos-daemon/compare/v1.5.0...v1.6.0) (2025-10-22)


### ✨ Features

* **mechanismcontrol:** add commands for printer
  control ([#41](https://github.com/adcondev/poster/issues/41)) ([53ef71c](https://github.com/AdConDev/pos-daemon/commit/53ef71c7f59ee935a7763a3d508ccb2ac95dca3c))

## [1.5.0](https://github.com/AdConDev/pos-daemon/compare/v1.4.0...v1.5.0) (2025-10-12)


### ✨ Features

* **escpos:** implemented SetTextSize(widthMultiplier, heightMultiplier int) method to adjust text
  dimensions. ([#36](https://github.com/adcondev/poster/issues/36)) ([29c6e4e](https://github.com/AdConDev/pos-daemon/commit/29c6e4ece6ffa2767f8a19d8c51a327667ca5743))

## [1.4.0](https://github.com/AdConDev/pos-daemon/compare/v1.3.1...v1.4.0) (2025-10-07)


### 🤖 Continuous Integration

* bump actions/stale from 9 to
  10 ([#27](https://github.com/adcondev/poster/issues/27)) ([ee65c43](https://github.com/AdConDev/pos-daemon/commit/ee65c43cf84a99616a0ac5eb1934ec3483107ca3))
* bump lewagon/wait-on-check-action from 1.4.0 to
  1.4.1 ([#30](https://github.com/adcondev/poster/issues/30)) ([45bf780](https://github.com/AdConDev/pos-daemon/commit/45bf780aebbe70d48188b8f324aab391fcbf81f1))


### ✅ Tests

* **barcode:** add byte slice validation helpers ([2830158](https://github.com/AdConDev/pos-daemon/commit/283015888a1c1c2812e5b93cd370c2e41dced96b))


### ✨ Features

* **barcode:** add barcode commands and integration tests for barcode commands ([464a741](https://github.com/AdConDev/pos-daemon/commit/464a741c730267e5c6de1256b03ffc2cd8da2d0c))


### 🐛 Bug Fixes

* **barcode:** update utils/test/validation_helpers.go ([07fa11a](https://github.com/AdConDev/pos-daemon/commit/07fa11a2de8992c8bd929db1d9b21a18148b751d))

### [1.3.1](https://github.com/AdConDev/pos-daemon/compare/v1.3.0...v1.3.1) (2025-09-15)


### 🤖 Continuous Integration

* bump actions/checkout from 4 to 5 ([1e9499f](https://github.com/AdConDev/pos-daemon/commit/1e9499f554f7e4ff92a4b530e59607b54efb79bd))
* bump actions/github-script from 7 to 8 ([71bad7d](https://github.com/AdConDev/pos-daemon/commit/71bad7dd4ce67078fff9c54b17a178db0f71a023))
* bump actions/labeler from 5 to 6 ([e965a54](https://github.com/AdConDev/pos-daemon/commit/e965a5492638a9a97927630e38b83675e6916a18))
* bump actions/setup-go from 5 to 6 ([e62e0e0](https://github.com/AdConDev/pos-daemon/commit/e62e0e00eba21e2af109e7016bed19e2776d998a))
* bump actions/setup-node from 4 to 5 ([07802e7](https://github.com/AdConDev/pos-daemon/commit/07802e78c30ba8e3f66006d9fb2328ce76e2c771))
* bump lewagon/wait-on-check-action from 1.3.1 to 1.4.0 ([b39aad5](https://github.com/AdConDev/pos-daemon/commit/b39aad5ad15bb25245b36d7e55f993581d8644f0))


### 📦 Dependencies

* **deps:** bump golang.org/x/image from 0.30.0 to 0.31.0 in the golang-x
  group ([#22](https://github.com/adcondev/poster/issues/22)) ([3c7131d](https://github.com/AdConDev/pos-daemon/commit/3c7131d6b5f338f248aaa9b19b1b6559f8698ddb))

## [1.3.0](https://github.com/AdConDev/pos-daemon/compare/v1.2.1...v1.3.0) (2025-09-15)


### ✨ Features

* **udchars:** add user-defined character test example ([f1e53cc](https://github.com/AdConDev/pos-daemon/commit/f1e53cc2799228917bfb476a96443794ddb78f81))

### [1.2.1](https://github.com/AdConDev/pos-daemon/compare/v1.2.0...v1.2.1) (2025-09-15)

### 📦 Dependencies

* **deps:** bump golang.org/x/text in the golang-x
  group ([cc79970](https://github.com/AdConDev/pos-daemon/commit/cc79970113218c838146bc75c0bac88c8a624c05))

### 🤖 Continuous Integration

* **dependabot:** enhance auto-merge workflow and add PR status
  dashboard ([1add7a1](https://github.com/AdConDev/pos-daemon/commit/1add7a13707c7835ad0b0ba5616daee9003d527a))

## [1.2.0](https://github.com/AdConDev/pos-daemon/compare/v1.1.0...v1.2.0) (2025-09-08)


### ✨ Features

* **printposition:** add print position management commands ([bc25641](https://github.com/AdConDev/pos-daemon/commit/bc256411ee830abdfd4757097942acb5c68dcabb))


### ✅ Tests

* **printposition:** update tests for print position functionality ([ab7705f](https://github.com/AdConDev/pos-daemon/commit/ab7705f279c4af5845b0145ac81463baa460fe5a))
* **print:** update error messages and assertions in tests ([3932bc0](https://github.com/AdConDev/pos-daemon/commit/3932bc0b14cc08b0b99fc3a2bb52052bed0b0b4a))
* **print:** update tests for print command integration ([74d8e58](https://github.com/AdConDev/pos-daemon/commit/74d8e5835311ce7d3b61843f2f5bc6df365249ce))
* **test:** add assertion helpers for testing utilities ([e382f4f](https://github.com/AdConDev/pos-daemon/commit/e382f4f36446f4d602de02499c14857fbc68f3e7))

## [1.1.0](https://github.com/AdConDev/pos-daemon/compare/v1.0.1...v1.1.0) (2025-09-02)


### ✨ Features

* **character:** add character handling
  capabilities ([#12](https://github.com/adcondev/poster/issues/12)) ([c62dd8e](https://github.com/AdConDev/pos-daemon/commit/c62dd8eb651ee69cf9c7c92cadb2b3676bc2a344))

### [1.0.1](https://github.com/AdConDev/pos-daemon/compare/v1.0.0...v1.0.1) (2025-08-26)


### 📦 Dependencies

* **go.mod:** add missing golang.org/x/text dependency ([83d8598](https://github.com/AdConDev/pos-daemon/commit/83d859877d9b6a46d7ae9c6f65862ff6d7d09d9e))
* **go.mod:** fix go.mod
  file ([9d5b966](https://github.com/AdConDev/pos-daemon/commit/9d5b966795494d94b3fdd651fcbd03379de9da9e)),
  closes [#7](https://github.com/adcondev/poster/issues/7)

## [1.0.0](https://github.com/AdConDev/pos-daemon/compare/v0.2.0...v1.0.0) (2025-08-26)


### ⚠ BREAKING CHANGES

* **protocols:** add new architecture for command chaining and 2-layered commands.
* **escpos:** The protocol interface has been modified
to support multiple protocols and may require updates to
existing implementations.

Signed-off-by: Adrián Constante <ad_con.reload@proton.me>

### 🐛 Bug Fixes

* **errors:** standardize error variable names ([b865304](https://github.com/AdConDev/pos-daemon/commit/b865304a04e7079ca09a08bbafbb4fb00528995e))
* **escpos:** improve comments and code clarity ([bbf654c](https://github.com/AdConDev/pos-daemon/commit/bbf654c2e3d705af8e7825f1836bd39fd51c5673))


### ✅ Tests

* **escpos:** add tests for dependency injection functionality ([3d7c680](https://github.com/AdConDev/pos-daemon/commit/3d7c680391865a44ea79235cc70edd331b17ea72))
* **escpos:** add tests for line spacing functionality ([6bbfb2a](https://github.com/AdConDev/pos-daemon/commit/6bbfb2a9d8c7ace7b691e005d6a62b606ba5e9c0))
* **escpos:** add unit tests for command functionalities ([d759a33](https://github.com/AdConDev/pos-daemon/commit/d759a33be3c6cffd3e74e9b8601ec0ea86da894e))


### ✨ Features

* **barcode:** update barcode handling functions ([85558e7](https://github.com/AdConDev/pos-daemon/commit/85558e790d037f74c84e06a0aa8aa1ca0d213c30))
* **escpos:** add line spacing capabilities and refactor commands ([b0a0f84](https://github.com/AdConDev/pos-daemon/commit/b0a0f84e499d90ac6769ebc7491916553120202a))
* **escpos:** enhance printer command structure and add comments ([dd1c733](https://github.com/AdConDev/pos-daemon/commit/dd1c7333bf0e9b8dc58c7cc9136108d031ed0b58))
* **escpos:** refactor printer protocol handling ([9ce0903](https://github.com/AdConDev/pos-daemon/commit/9ce09039be23b004c1e282e8d09efd522c6d1129))
* **escpos:** refactor protocol structure and update imports ([f1840f8](https://github.com/AdConDev/pos-daemon/commit/f1840f87ef9b3cedb1f519184b53cb84bcc1dd30))
* **printer:** enhance printer configuration structure ([53c4c9c](https://github.com/AdConDev/pos-daemon/commit/53c4c9ccc93b1e33c6c5e2e27a8626af95a156bd))
* **protocols:** add new architecture for command chaining and 2-layered
  commands. ([599214f](https://github.com/AdConDev/pos-daemon/commit/599214f87e55896323056e47aa919776b2513d36)),
  closes [#4](https://github.com/adcondev/poster/issues/4)
* **protocol:** update import paths for escpos types ([599aca6](https://github.com/AdConDev/pos-daemon/commit/599aca6982e17ce3a83902b7ddf449a3c34b1d18))

## 0.2.0 (2025-08-15)

### 🤖 Continuous Integration

* bump amannn/action-semantic-pull-request from 5 to
  6 ([648be79](https://github.com/AdConDev/pos-daemon/commit/648be7999f29327db7bee9bbad30874ae27cbc64))
* bump codecov/codecov-action from 4 to
  5 ([3ce0298](https://github.com/AdConDev/pos-daemon/commit/3ce0298273748a58a796e0c90382bb9e3bc585e5))

### ✨ Features

* **escpos:** add initial implementation for ESC/POS
  commands ([f9772b4](https://github.com/AdConDev/pos-daemon/commit/f9772b47c1e4e2f8cd11910817250ef45ac472ca))
* **github:** add initial github workflows and
  files ([812b851](https://github.com/AdConDev/pos-daemon/commit/812b8513d31c12bb2eb240eb551d68bf9708c8e6))

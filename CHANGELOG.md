# Changelog

We follow the CalVer (https://calver.org/) versioning scheme: YY.MINOR.MICRO.

24.0.0 (2024-05-20)
===================

- Glide was replaced with go's built-in dependency manager
- Docker image was upgraded to use go 1.22
- Dockerfile was restructured using multi stage build, which reduced image size from 400+mb to 20mb
- Manifest v2 schema 2 was used to be compatible with newer docker versions
- Update readme

19.0.1 (2019-08-20)
===================

- Update readme

19.0.0 (2019-08-19)
===================

- Update fakeCAS for OSF token-scope relationship change
- Add OAuth revoke endpoint
- Fix OAuth profile endpoint
- Rewrite readme


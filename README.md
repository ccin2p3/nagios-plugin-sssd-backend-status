# Nagios plugin for checking SSSd backend status

## Description

This Nagios plugin aims to check `sssd` backend status by using the `sssctl` tool to query the `sssd-ifp` API.

Any `Offline` backend will be reported as a `CRITICAL` error.

## Usage

```
# /usr/lib64/nagios/plugins/check_sssd_backend_status -h
Checks that sssd backends are Online. This tool rely on sssctl.

Usage:
  check_sssd_backend_status [flags]

Flags:
      --debug            enable debug
      --domains string   domains to check (comma separated). Defaults to check all domains
  -h, --help             help for check_sssd_backend_status
```

### Check all domains

```
# /usr/lib64/nagios/plugins/check_sssd_backend_status
sssd domains are online
```

### Check specific domains

Domains to be checked must be _comma separated_:

```
# /usr/lib64/nagios/plugins/check_sssd_backend_status --domains domain1,domain2
sssd domains are online
```

## Build

```
$ CGO_ENABLED=0 GOFLAGS='-mod=vendor' go build -a -ldflags "-extldflags '-static'" -o check_sssd_backend_status .
```

## Configuration

### Requirements

This probe requires access to the `sssd-ifp` API.

As stated in the `sssd-ifp(5)` man page, by default only `root` user has access to this API.

If you want to execute this probe has another user, you **MUST** add your user _UID_ to to `allowed_uids` parameter. See `sssd-ifp(5)` for more informations.

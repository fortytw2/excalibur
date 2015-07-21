Excalibur
------

Gunning to be your next `init`. Combines a reliable, straightforward process
manager with a pluggable logger and a zombie reaper.

Slice through your daemons - all logging is seamlessly handled through stdout/err

Status
-----
Process handling is mostly finished except for tests.
The built binary simple runs `redis-server` for 10 seconds to test functionality

Todo
------
- Config file parsing + dependency tree generation
- Parallel launching of programs from the dependency tree
- file based logger

Goals
-----
- Control via `mangos`/`nanomsg` + creation of a `caliburctl` binary
- Pluggable logging backends (not just flatfiles in `/log/redis/YYYY-MM-DD`)
and symlinks `/log/redis/today`/`/log/redis/yesterday`
- Socket activation? (needs to be researched if this even has real benefits)
- Container control on top of `runc`/`rkt`/`docker`?

LICENSE
------
MIT, see LICENSE

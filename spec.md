# Shutter, an EBS snapshot tool.
Overview
---

Use Cases
---
- Ari, a systems administrator, wants to take a backup of the data drives for
a PostgreSQL server before performing a delicate migration that might fail. Using
shutter, he quickly takes a consistent snapshot of the drives, and performs the
migration. It does indeed fail, so he easily rolls back by re-attaching the old
versions of the drives from snapshots.

- To protect against data loss, Walter has setup a crontab on HBase instances
which takes a consistent snapshot of the data storage drives, so that if an
instance fails, a replacement can easily be brought back online. Stale data is a
concern with this approach, but he is willing to live with this trade-off.

*(As you can see, most of these are database-related - Shutter is primarily
used and designed for data-storage instances where the idea of a snapshot of
existing data is actually important.)*

Taking Snapshots
---
This is the most common use for Shutter, as most snapshots are fire-and-forget.
This mode is supported through the `shutter snap` command, and its associated
command-line options.

Taking snapshots is done in one of two ways:

- Running on an EC2 instance, `shutter snap` gathers all EBS drives attached to
the server and makes a snapshot of these all.

- `shutter snap` can be given a list of EBS volume IDs, each of which will then
be backed up. This mode can be done outside of an EC2 instance, but `shutter`
cannot make a necessarily consistent snapshot as it cannot freeze the
filesystem while it’s being backed up.

A basic diagram of this looks like so:


Restoring from Snapshots
---

Harakiri Mode
---

Acknowledgements
---

References
---

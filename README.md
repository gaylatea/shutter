# Shutter
## Self-contained EBS consistent snapshot program.

This little program is designed to be a standalone EBS snapshot program
that doesn't have crazy dependencies like ec2-consistent-snapshot or
ebs-consistent-snapshot do.

You can direct it to snapshot all EBS drives attached to the instance,
or a subset thereof, and give descriptions to them for easy
administration.

    shutter \
      -region us-west-1 \
      -partition /var/lib/pgdata \
      -volume vol-1a1a1a1a \
      -volume vol-2b2b2b2b \
      -volume vol-3c3c3c3c \
      -description "PostgreSQL RAID from $HOSTNAME"

Tested and built on darwin/amd64 (OS X Mavericks) and linux/amd64 (Amazon Linux).

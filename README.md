# Shutter
## Self-contained EBS consistent snapshot program.

This little program is designed to be a standalone EBS snapshot program
that doesn't have crazy dependencies like ec2-consistent-snapshot or
ebs-consistent-snapshot do.

You can direct it to snapshot all EBS drives attached to the instance,
or a subset thereof, and give descriptions to them for easy
administration.

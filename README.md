# Shutter
**A dead-simple EBS snapshot program with no external dependencies.**

**Developed @ [AdRoll](https://github.com/SemanticSugar) in SF.**

## Abstract
This little program is designed to be a standalone EBS snapshot program
that doesn't have crazy dependencies like
[ec2-consistent-snapshot](https://github.com/alestic/ec2-consistent-snapshot) or
[ebs-consistent-snapshot](https://github.com/Jd007/ebs-consistent-snapshot) do.

You can direct it to snapshot all EBS drives attached to the instance or a
subset thereof, and give descriptions to them for easy administration.

## Usage

    shutter \
      -region us-west-1 \
      -partition /var/lib/pgdata \
      -volume vol-1a1a1a1a \
      -volume vol-2b2b2b2b \
      -volume vol-3c3c3c3c \
      -description "PostgreSQL RAID from $HOSTNAME"

When running on an EC2 instance, helper methods are provided to automatically
grab the region and volume information from the AWS APIs. You do not need to
specify these parameters in this case:

    shutter \
      -partition /var/lib/pgdata \
      -description "RAID Snapshot."

*Regarding the -description[] parameter: You can specify either one description
that will be used for all the snapshots, or you must specify an equal number of
descriptions to volumes. This gets a little tricky if you use the EC2 helpers,
so it's recommended to just give a single description most of the time.*

*For instance:*

    shutter \
      -partiton /var/lib/pgdata \
      -volume vol-1a1a1a1a \
      -volume vol-2b2b2b2b \
      -description "dbserver backup"

*Will produce two snapshots with the same description, whereas:*

    shutter \
      -partition /var/lib/pgdata \
      -volume vol-1a1a1a1a \
      -volume vol-2b2b2b2b \
      -description "dbserver /dev/sda" \
      -description "dbserver /dev/sdb"

*Will produce two snapshots, vol-1a1a1a1a's with the /dev/sda description and
vol-2b2b2b2b's with the /dev/sdb description. This is handy for if you have
drives that should be re-assembled onto particular block devices.*

```shutter``` is tested and built on darwin/amd64 (OS X Mavericks) and
linux/amd64 (Amazon Linux), and the default build toolchain is configured for
these targets only.

## Come Hack with Me
Dependencies are pulled down via [gom](https://github.com/mattn/gom), the
awesome Golang package manager. You'll need that in order to hack on this.

Assuming that you have Golang already installed, you should just need:

    make setup

And all dependencies should be pulled down automatically for you to work. After that, ```make``` should work just fine.
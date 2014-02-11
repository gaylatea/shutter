// **shutter** builds EBS snapshots.
// It was written because I didn't want the extreme hassle of
// installing Perl/Python dependencies for ec2-consistent-snapshot or
// its descendents.
//
// Support for freezing and restoring XFS partitions is included.
package main

import (
	"flag"
	"fmt"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/ec2"
	"log"
	"os/exec"
	"sync"
	"time"
)

// Custom type for multiple command line-specified strings.
type stringslice []string

func (v *stringslice) String() string {
	return fmt.Sprintf("%#v", v)
}

func (v *stringslice) Set(value string) error {
	*v = append(*v, value)
	return nil
}

// Freeze the requested filesystems.
func freeze(mountpoints []string) {
	var w sync.WaitGroup
	w.Add(len(mountpoints))

	for _, mp := range mountpoints {
		go func(partition string) {
			log.Printf("freezing %s for the snapshots.", partition)
			cmd := exec.Command("xfs_freeze", "-f", partition)
			cmd.Run()

			output, _ := cmd.Output()
			log.Printf("%s debug: %s", partition, string(output))
			w.Done()
		}(mp)
	}

	w.Wait()
}

// Unfreeze the requested filesystems.
func unfreeze(mountpoints []string) {
	var w sync.WaitGroup
	w.Add(len(mountpoints))

	for _, mp := range mountpoints {
		go func(partition string) {
			log.Printf("unfreezing %s now that snapshots are complete.", partition)
			cmd := exec.Command("xfs_freeze", "-u", partition)
			cmd.Run()

			output, _ := cmd.Output()
			log.Printf("%s debug: %s", partition, string(output))
			w.Done()
		}(mp)
	}

	w.Wait()
}

func main() {
	ts := time.Now().UTC().Unix()
	log.Printf("shutter: startup at %d", ts)

	// **-region:** only required if not running on an EC2 instance.
	var region = flag.String("region", "", "EC2 region to look for EBS volumes in.")

	// **-volume[]:** Can be specified multiple times.
	// Defines the volumes that we are backing up.
	var volumes stringslice
	flag.Var(&volumes, "volume", "List of volumes to snapshot.")

	// **-description[]:** Can be specified multiple times.
	// Defines the description of the matching volume.
	//
	// If only one description is given it is applied to all snapshots.
	// You must either define only one, or match all volumes.
	// All other conditions are errors.
	var descriptions stringslice
	flag.Var(&descriptions, "description", "List of descriptions to give snapshots.")

	// **-partition[]:** XFS partitions that should be frozen and
	// unfrozen as a part of this snapshot.
	var partitions stringslice
	flag.Var(&partitions, "partition", "XFS partition to freeze/un-freeze.")

	flag.Parse()

	// Do some sanity-checking on the arguments we're given.
	// **TODO(silversupreme):** Add in support for figuring out what
	// region to use from the EC2 instance metadata.
	awsRegion, present := aws.Regions[*region]
	if !present {
		log.Fatalf("Given region %s not a supported AWS region!", *region)
	}

	log.Printf("Backing up volumes in %s region.", *region)

	// Check that the number of descriptions and volumes match up.
	numDescriptions := len(descriptions)
	numVolumes := len(volumes)

	// **TODO(silversupreme):** Make this instead use all volumes
	// attached to the currently-running instance, if we can.
	if numVolumes < 1 {
		log.Fatalf("No volumes specified!")
	}

	if numDescriptions != 1 {
		if numDescriptions != numVolumes {
			log.Fatalf("Mis-matched %d descriptions to %d volumes!", numDescriptions, numVolumes)
		}
	}

	// Gather up all of the volumes and their descriptions.
	volumesToSnapshot := map[string]string{}
	for key, value := range volumes {
		var thisDescription string

		if numDescriptions == 1 {
			thisDescription = descriptions[0]
		} else {
			thisDescription = descriptions[key]
		}

		volumesToSnapshot[value] = thisDescription
	}

	// Connect to EC2 and make sure everything is alright there.
	credentials, err := aws.GetAuth("", "")
	if err != nil {
		log.Fatalf("%s", err)
	}

	ec2_conn := ec2.New(credentials, awsRegion)

	// Freeze the requested filesystem.
	freeze(partitions)
	defer unfreeze(partitions)

	// Create snapshots of the requested volumes, in parallel.
	// Skip over any that are not actually present in EC2.
	var w sync.WaitGroup
	w.Add(numVolumes)

	for volumeId, desc := range volumesToSnapshot {
		go func(volumeId, desc string, ec2_conn *ec2.EC2) {
			// Create the snapshot if we can.
			log.Printf("Creating snapshot for %s", volumeId)

			resp, err := ec2_conn.CreateSnapshot(volumeId, desc)
			if err != nil {
				log.Printf("Could not create snapshot for %s: %s", volumeId, err)
				w.Done()
				return
			}

			log.Printf("Created snapshot %s for %s", resp.Snapshot.Id, volumeId)
			w.Done()
		}(volumeId, desc, ec2_conn)
	}

	w.Wait()

	log.Printf("shutter: done!")
}

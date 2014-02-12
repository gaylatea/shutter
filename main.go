// **shutter** builds EBS snapshots.
// It was written because I didn't want the extreme hassle of
// installing Perl/Python dependencies for ec2-consistent-snapshot or
// its descendents.
//
// Support for freezing and restoring XFS partitions is included.
package main

import (
	"./util"
	"flag"
	"fmt"
	"github.com/mitchellh/goamz/aws"
	"github.com/mitchellh/goamz/ec2"
	"os"
	"os/exec"
	"sync"
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
func freeze(mountpoints []string, logger util.Logger) {
	var w sync.WaitGroup
	w.Add(len(mountpoints))

	for _, mp := range mountpoints {
		go func(partition string) {
			logger.Debugf("freezing %s for the snapshots.", partition)
			cmd := exec.Command("xfs_freeze", "-f", partition)
			cmd.Run()

			output, _ := cmd.Output()
			logger.Debugf("%s debug: %s", partition, string(output))
			w.Done()
		}(mp)
	}

	w.Wait()
}

// Unfreeze the requested filesystems.
func unfreeze(mountpoints []string, logger util.Logger) {
	var w sync.WaitGroup
	w.Add(len(mountpoints))

	for _, mp := range mountpoints {
		go func(partition string) {
			logger.Debugf("unfreezing %s now that snapshots are complete.", partition)
			cmd := exec.Command("xfs_freeze", "-u", partition)
			cmd.Run()

			output, _ := cmd.Output()
			logger.Debugf("%s debug: %s", partition, string(output))
			w.Done()
		}(mp)
	}

	w.Wait()
}

func main() {
	logger, _ := util.NewColourizedOutputLogger(os.Stdout)

	// **-region:** only required if not running on an EC2 instance.
	var region string
	flag.StringVar(&region, "region", "", "EC2 region to look for EBS volumes in.")

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

	// Add in a debug flag to show detailed logging output for diagnosis.
	var debugMode bool
	flag.BoolVar(&debugMode, "debug", false, "Enable debug output in logs.")

	flag.Parse()

	// Enable debug mode if the user has requested it.
	logger.SetDebugOutput(debugMode)

	// Do some sanity-checking on the arguments we're given.
	// **TODO(silversupreme):** Add in support for figuring out what
	// region to use from the EC2 instance metadata.
	if region == "" {
		logger.Error("You did not specify a region!")
		return
	}

	awsRegion, present := aws.Regions[region]
	if !present {
		logger.Errorf("Given region %s not a supported AWS region!", region)
		return
	}

	logger.Debugf("Backing up volumes in %s region.", region)

	// Check that the number of descriptions and volumes match up.
	numDescriptions := len(descriptions)
	numVolumes := len(volumes)

	// **TODO(silversupreme):** Make this instead use all volumes
	// attached to the currently-running instance, if we can.
	if numVolumes < 1 {
		logger.Error("No volumes specified!")
		return
	}

	if numDescriptions != 1 {
		if numDescriptions != numVolumes {
			logger.Errorf("Mis-matched %d descriptions to %d volumes!", numDescriptions, numVolumes)
			return
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
		logger.Errorf("Auth failure: %s", err.Error())
		return
	}

	ec2_conn := ec2.New(credentials, awsRegion)

	// Freeze the requested filesystem.
	freeze(partitions, logger)
	defer unfreeze(partitions, logger)

	// Create snapshots of the requested volumes, in parallel.
	// Skip over any that are not actually present in EC2.
	var w sync.WaitGroup
	w.Add(numVolumes)

	for volumeId, desc := range volumesToSnapshot {
		go func(volumeId, desc string) {
			// Create the snapshot if we can.
			logger.Debugf("Creating snapshot for %s", volumeId)

			resp, err := ec2_conn.CreateSnapshot(volumeId, desc)
			if err != nil {
				logger.Errorf("Could not create snapshot for %s: %s", volumeId, err)
				w.Done()
				return
			}

			logger.Successf("Created snapshot %s for %s", resp.Snapshot.Id, volumeId)
			w.Done()
		}(volumeId, desc)
	}

	w.Wait()

	logger.Infof("Done!")
}

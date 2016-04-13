# BOSH Delta

Utility for comparing BOSH releases so that you as a BOSH jockey know what has
changed between releases.

This is _very_ much alpha software. Use at your own risk. Not responsibile for
kicked puppies or unicorn farts.

## Introduction

BOSH Delta will eventually support many types of artifacts found in a BOSH
release, however for now we only support comparing job properties. Specifically
it will tell you what has been removed or added for all jobs in the release.
This allow you as the operator to update your BOSH deployment manifest in a
sane and efficient way without uneeded worry. 

## Usage

1. Out of band, download release 1
2. Out of band, download release 2
3. Run `bosh-delta release1.tgz release2.tgz`

## Initial Design

1. Given the location of two releases which are tar files, extract each
to a temporary directory.
2. Read the release.MF file
3. Read each job out of the release.MF
4. For each job extract the job's tar file to a temp directory.
5. Read the job manifest
6. Read out each property in the job manifest
7. Look for new/removed properties and add to the delta result
8. Dump the delta results to stdout

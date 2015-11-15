# Compare drive or directory snapshots

Useful for figuring out where free space went - what is filling the drive (for example `/var/folders`).

### Create snapshots:

    ./snapshot
    ./tmsnapshot time_machine_drive

or create snapshots at two different points in time -
`snapshot.txt` and `snapshot_prev.txt` with just `./snapshot` (after first run use `cp snapshot.txt snapshot_prev.txt`)

    go run compare.go

or

    go run compare.go snapshot_prev.txt snapshot_current.txt

Colored diff is shown in terminal and also saved to `diff.txt`.

To sumarize some particular changeset:

    go run sum filter

For example observe the sum of all changes in png files between last two snapshots:

    go run sum png

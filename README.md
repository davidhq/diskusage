# Compare drive or directory snapshots

Useful for figuring out where free space went - what is filling the drive (for example `/var/folders`).

### Create snapshots:

    ./snapshot
    ./tmsnapshot time_machine_drive

or create snapshots at two different points in time -
`snapshot.txt` and `snapshot_prev.txt` with just `./snapshot` (after first run use `cp snapshot.txt snapshot_prev.txt`)

    go run compare.go

Colored diff is shown in terminal and also saved to `diff.txt`.


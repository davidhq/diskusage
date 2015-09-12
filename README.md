# Compare drive or directory snapshots

Useful for figuring out where free space went - what is filling the drive (for example `/var/folders`).

### Setup

    cp exclude.txt.sample exclude.txt

and update the content.

### Create snapshots:

    ./snapshot
    ./tmsnapshot time_machine_drive

or create snapshots at two different points in time - `snapshot.txt` and `snapshot_prev.txt` with `./snapshot`

    go run compare.go

Colored diff is shown in terminal and also saved to `diff.txt`.


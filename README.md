## Restic sync tool

Usage:

```
# Scan and show directories
resticsync <scan_directory> 

# Backup all directories
resticsync <repository_path> <scan_directory>
```

Application will scan `<scan_directory>` for directories containg `.resticsync` file, and backup it to restic repository specified by `<repository_path>`

when `.resticsync` is empty will generated and saved inplace.

Example `.resticsync` file

```
Id = 'BpLnfgDsc8WD2F8qNfHK'
Name = 'Restic sync'
Category = 'projects'
```

All snapshots craeted using this tool will be tagged with Id, category and name tags (slugified)

### Build

```
git clone https://github.com/orneo1212/restic-sync.git
cd restic-sync/
go build .
```

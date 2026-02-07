## Slicer box instructions

Linux sandbox for AI agents.

## Pre-reqs

Must point slicer to the box service and token file:

```
export SLICER_URL=https://box.slicervm.com
export SLICER_TOKEN_FILE=~/.slicer/gh-access-token
```

Then list VMs to find the box (1 per user):

```
slicer vm list
```

## Copying files into the box

* Always give full path when specifying the VM itself
* Always create sub-directories first before copying files into them

```
slicer vm exec --uid 1000 VM_NAME -- mkdir -p /home/ubuntu/.local/share/amp
slicer vm exec --uid 1000 VM_NAME -- mkdir -p /home/ubuntu/.amp
slicer vm cp --uid 1000 ~/.local/share/amp/secrets.json VM_NAME:/home/ubuntu/.local/share/amp/secrets.json
slicer vm cp --uid 1000 ~/.config/amp/settings.json VM_NAME:/home/ubuntu/.amp/settings.json
```

You can also copy in `--mode=tar` for a whole directory.

## Exec

Example:

```
slicer vm exec --uid 1000 VM_NAME -- /home/ubuntu/.arkade/bin/amp usage
```

### Interactive terminal

```
slicer vm exec --uid 1000 VM_NAME
```

### SSH/SCP (break glass only)

Typically no reason for this, but:

```
# Ensure SLICER_URL and SLICER_TOKEN_FILE are set
slicer vm forward VM_NAME \
    -L 127.0.0.1:2222:127.0.0.1:22
Copy a file from the VM using SCP:

scp -P 2222 ubuntu@127.0.0.1:/etc/os-release ./os-release
Run a command over SSH:

ssh -p 2222 ubuntu@127.0.0.1 uptime
```

## Reset the box

Only do this in case of a major issue or when having to test system-level changes.

# Ensure SLICER_URL and SLICER_TOKEN_FILE are set
slicer vm delete VM_NAME
Then launch a new VM:

slicer vm launch

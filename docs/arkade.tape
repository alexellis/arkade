Output arkade.gif

Set FontSize 20
Set Shell zsh
Set Width 1200
Set Height 600
Set TypingSpeed 100ms
Set WindowBar Colorful


# Downloading a CLI Binary
Type@50ms "# Download a CLI Binary with arkade get" Sleep 2s
Hide
Ctrl+C
Show

Type "arkade get jq" Enter Sleep 8s
Type "jq --version" Enter Sleep 2s
Hide
Type "clear" Enter
Show

# Specifying the version of the CLI
Type@50ms "# Download a specific version of a CLI" Sleep 2s
Hide
Ctrl+C
Show

Type "arkade get kind@v0.19.0" Enter Sleep 8s
Type "kind version" Enter Sleep 2s
Hide
Type "clear" Enter
Show

# Downloading the multiple CLIs
Type@50ms "# Download multiple CLIs" Sleep 2s
Hide
Ctrl+C
Show

Type "arkade get k3sup@0.12.12 viddy" Enter Sleep 15s
Hide
Type "clear" Enter
Show

# Downloading and Installing system binaries
Type@50ms "# Download and Install a system binary" Sleep 2s
Hide
Ctrl+C
Show

# Hack to run command with sudo on VHS tape
#
# Running a command with sudo on a VHS tape can be tricky. To work around this,
# you can edit the /etc/sudoers file to not require passwords for the user
# temporarily or create a hidden command that runs a sudo command and enters the
# password before the actual command.
Type "sudo arkade system install firecracker" Enter Sleep 8s
Type "firecracker --version" Enter Sleep 2s
Hide
Type "clear" Enter
Show

# Closing message
Type "Try arkade out now to see how easy it is to install your favourite tools." Enter Sleep 2s

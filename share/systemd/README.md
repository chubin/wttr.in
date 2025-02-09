To add **wttr.in** to systemd as a service, do the following steps.


1. **Create a systemd service file**: You’ll need to create a service file in `~/.config/systemd/user/` (for user-level) or `/etc/systemd/system/` (for system-wide) directory. Let’s create it for a user.

Create the directory if it doesn’t exist:

```bash
mkdir -p ~/.config/systemd/user/
```

Then, create the service file called `myscript.service`:

```bash
cp share/systemd/wttrin.service ~/.config/systemd/user/wttrin.service
```

2. **Reload the systemd daemon**: This will ensure systemd recognizes the new service.

```bash
systemctl --user daemon-reload
```

4. **Enable and start your service**:

To start the service immediately, run:

```bash
systemctl --user start wttrin.service
```

5. **Check the status**: To verify that your service is running correctly, you can check its status:

```bash
systemctl --user status wttrin.service
```

6. **Start service automatically**: This will ensure the service is running after reboot.

To enable it to start automatically at boot, run:

```bash
systemctl --user enable wttrin.service
```

Enable user services even if the user is not logged in (specify the user name instead of `USER`):

```bash
sudo loginctl enable-linger USER
```

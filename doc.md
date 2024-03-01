```shell
sudo systemctl stop systemd-resolved
sudo systemctl disable systemd-resolved

sudo systemctl start resolvconf.service


sudo ln -sf /run/systemd/resolve/resolv.conf /etc/resolv.conf
sudo nano /etc/systemd/resolved.conf
cat /etc/resolv.conf
```
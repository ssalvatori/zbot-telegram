# Create systemd service

Copy zbot service definition.
```
sudo cp zbot.service /etc/systemd/system/zbot.service
sudo cp zbot.default /etc/default/zbot
sudo chmod 664 /etc/systemd/system/zbot.service
```

Enable and start zbot service

```
sudo systemctl daemon-reload
sudo systemctl enable zbot
sudo systemctl start zbot
sudo systemctl status zbot
```

Create symlink

````
sudo ln -s /path/to//zbot-telegram-go-linux-amd64 /usr/local/bin/zbot
```
### get nvm

```
curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.35.3/install.sh | bash
. ~/.bashrc
nvm install 8
nvm use 8
```

### get golang

```
wget https://dl/google.com/go/go1.10.3.linux-amd64.tar.gz
# sudo tar xvaf go1.10.3.linux-amd64.tar.gz -C /usr/local/
```

### get pool src

```
git clone https://github.com/aquachain/open-aquachain-pool.git
cd open-aquachain-pool
```

### build backend

This compiles to `./build/bin/open-aquachain-pool`, i like to rename it shorter after it builds.

```
make
```

### build frontend

```
cd www

npm install -g bower ember-cli

npm install

bower install

./build.sh
```

### package dist folder

```
tar czpf dist.tgz dist
```

### send to your server (replace 127.0.0.1)

```
scp dist.tgz root@127.0.0.1:~/
scp ../build/bin/open-aquachain-pool root@127.0.0.1:~/
```

### edit config files

yes, make two config files. one for 'unlocking and payouts' and one for 'everything else'

### install redis-server

sudo apt-get update && sudo apt-get install -y redis-server

### iptables (firewall)

see attached files for iptables script

### caddy server (reverse proxy)

Download here: https://caddyserver.com

see attached files for CaddyFile

you must change the first line where it says a domain name.

the domain name must resolve, you will get automatic SSL certificate.

### Putting it all together

(as root)

```
adduser --system caddy
adduser --system pool
adduser --system aquachain
```

```
cp caddy /usr/local/bin/caddy
cp aquachain /usr/local/bin/aquachain
cp open-aquachain-pool /usr/local/bin/aquapool
cp pool.json /home/pool/pool.json
cp payouts.json /home/pool/payouts.json

```

#### crontabs for reboots

```
crontab -u caddy -e
```

Put this line: `@reboot /usr/local/bin/caddy -conf /home/caddy/CaddyFile`


```
crontab -u aquachain -e
```

Put this line: `@reboot tmux new-session -d -s tmux "$HOME/start.bash"`


```
crontab -u pool -e
```

Put this line: `@reboot tmux new-session -d -s tmux "$HOME/start.bash"`


#### start.bash

aquachain: `/usr/local/bin/aquachain -rpc`

pool:      `/usr/local/bin/aquapool /home/pool/pool.json 2>>/home/pool/miner.log 1>>/home/pool/miner.log`


### reboot

cross fingers

if it works, add your pool api to https://github.com/aquachain/aquachain.github.io/blob/master/pools.json 








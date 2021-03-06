#!/bin/sh
# Tmux script that sets up a mining harness. 
set -e
SESSION="harness"
NODES_ROOT=~/harness
RPC_USER="user"
RPC_PASS="pass"
WALLET_SEED="b280922d2cffda44648346412c5ec97f429938105003730414f10b01e1402eac"
WALLET_PASS=123
BACKUP_PASS=b@ckUp
NETWORK="simnet"
# NETWORK="testnet"
# NETWORK="mainnet"
SOLO_POOL=0
MAX_GEN_TIME=20
PAYMENT_METHOD="pplns"
LAST_N_PERIOD=300 # PPLNS range, 5 minutes.

if [ "${NETWORK}" = "simnet" ]; then
 POOL_MINING_ADDR="SspUvSyDGSzvPz2NfdZ5LW15uq6rmuGZyhL"
 PFEE_ADDR="SsVPfV8yoMu7AvF5fGjxTGmQ57pGkaY6n8z"
 CLIENT_ONE_ADDR="SsZckVrqHRBtvhJA5UqLZ3MDXpZHi5mK6uU"
 CLIENT_TWO_ADDR="Ssn23a3rJaCUxjqXiVSNwU6FxV45sLkiFpz"
fi

if [ "${NETWORK}" = "mainnet" ]; then
POOL_MINING_ADDR="B1wXSYY3k8oiUc5KqF9HLncxQhB1VtwAkBP"
PFEE_ADDR="B1wXSYY3k8oiUc5KqF9HLncxQhB1VtwAkBP"
CLIENT_ADDR="B1wXSYY3k8oiUc5KqF9HLncxQhB1VtwAkBP"
fi

if [ "${NETWORK}" = "testnet" ]; then
POOL_MINING_ADDR="TsWP4MhHZn77F6tQprHhFoJHMWifaDdh2Mc"
PFEE_ADDR="TsWP4MhHZn77F6tQprHhFoJHMWifaDdh2Mc"
CLIENT_ADDR="TsWP4MhHZn77F6tQprHhFoJHMWifaDdh2Mc"
fi


if [ -d "${NODES_ROOT}" ]; then
  rm -R "${NODES_ROOT}"
fi

BITUMD_RPC_CERT="${HOME}/.bitumd/rpc.cert"
BITUMD_RPC_KEY="${HOME}/.bitumd/rpc.key"
BITUMD_CONF="${HOME}/.bitumd/bitumd.conf"
WALLET_RPC_CERT="${HOME}/.bitumwallet/rpc.cert"
WALLET_RPC_KEY="${HOME}/.bitumwallet/rpc.key"
GUI_DIR="${NODES_ROOT}/gui"

echo "Writing node config files"
mkdir -p "${NODES_ROOT}/master"
mkdir -p "${NODES_ROOT}/wallet"
mkdir -p "${NODES_ROOT}/pool"
mkdir -p "${NODES_ROOT}/gui"
mkdir -p "${NODES_ROOT}/c1"
mkdir -p "${NODES_ROOT}/c2"

cp -r gui/public ${GUI_DIR}/public
cp -r gui/templates ${GUI_DIR}/templates

cat > "${NODES_ROOT}/c1/client.conf" <<EOF
debuglevel=trace
activenet=${NETWORK}
user=m1
address=${CLIENT_ONE_ADDR}
pool=127.0.0.1:5550
EOF

cat > "${NODES_ROOT}/c2/client.conf" <<EOF
debuglevel=trace
activenet=${NETWORK}
user=m2
address=${CLIENT_TWO_ADDR}
pool=127.0.0.1:5550
EOF

if [ "${NETWORK}" = "mainnet" ]; then
cat > "${NODES_ROOT}/bitumctl.conf" <<EOF
rpcuser=${RPC_USER}
rpcpass=${RPC_PASS}
rpccert=${BITUMD_RPC_CERT}
rpcserver=127.0.0.1:9209
EOF

cat > "${NODES_ROOT}/pool/pool.conf" <<EOF
rpcuser=${RPC_USER}
rpcpass=${RPC_PASS}
bitumdrpchost=127.0.0.1:9209
bitumdrpccert=${BITUMD_RPC_CERT}
walletgrpchost=127.0.0.1:9211
walletrpccert=${WALLET_RPC_CERT}
debuglevel=trace
maxgentime=${MAX_GEN_TIME}
solopool=${SOLO_POOL}
activenet=${NETWORK}
walletpass=${WALLET_PASS}
poolfeeaddrs=${PFEE_ADDR}
paymentmethod=${PAYMENT_METHOD}
lastnperiod=${LAST_N_PERIOD}
backuppass=${BACKUP_PASS}
guidir=${GUI_DIR}
EOF
fi

if [ "${NETWORK}" = "simnet" ]; then
cat > "${NODES_ROOT}/bitumctl.conf" <<EOF
rpcuser=${RPC_USER}
rpcpass=${RPC_PASS}
rpccert=${BITUMD_RPC_CERT}
rpcserver=127.0.0.1:19556
EOF

cat > "${NODES_ROOT}/pool/pool.conf" <<EOF
rpcuser=${RPC_USER}
rpcpass=${RPC_PASS}
bitumdrpchost=127.0.0.1:19556
bitumdrpccert=${BITUMD_RPC_CERT}
walletgrpchost=127.0.0.1:19558
walletrpccert=${WALLET_RPC_CERT}
debuglevel=trace
maxgentime=${MAX_GEN_TIME}
solopool=${SOLO_POOL}
activenet=${NETWORK}
walletpass=${WALLET_PASS}
poolfeeaddrs=${PFEE_ADDR}
paymentmethod=${PAYMENT_METHOD}
lastnperiod=${LAST_N_PERIOD}
backuppass=${BACKUP_PASS}
guidir=${GUI_DIR}
EOF

cat > "${NODES_ROOT}/bitumwctl.conf" <<EOF
rpcuser=${RPC_USER}
rpcpass=${RPC_PASS}
rpccert=${WALLET_RPC_CERT}
rpcserver=127.0.0.1:19557
EOF

cat > "${NODES_ROOT}/wallet/wallet.conf" <<EOF
username=${RPC_USER}
password=${RPC_PASS}
rpccert=${WALLET_RPC_CERT}
rpckey=${WALLET_RPC_KEY}
logdir=./log
appdata=./data
simnet=1
pass=${WALLET_PASS}
debuglevel=debug
EOF
fi

if [ "${NETWORK}" = "testnet" ]; then
cat > "${NODES_ROOT}/bitumctl.conf" <<EOF
rpcuser=${RPC_USER}
rpcpass=${RPC_PASS}
rpccert=${BITUMD_RPC_CERT}
rpcserver=127.0.0.1:19209
EOF

cat > "${NODES_ROOT}/pool/pool.conf" <<EOF
rpcuser=${RPC_USER}
rpcpass=${RPC_PASS}
bitumdrpchost=127.0.0.1:19209
bitumdrpccert=${BITUMD_RPC_CERT}
walletgrpchost=127.0.0.1:19211
walletrpccert=${WALLET_RPC_CERT}
debuglevel=trace
maxgentime=${MAX_GEN_TIME}
solopool=${SOLO_POOL}
activenet=${NETWORK}
walletpass=${WALLET_PASS}
poolfeeaddrs=${PFEE_ADDR}
paymentmethod=${PAYMENT_METHOD}
lastnperiod=${LAST_N_PERIOD}
backuppass=${BACKUP_PASS}
guidir=${GUI_DIR}
EOF
fi

cd ${NODES_ROOT} && tmux new-session -d -s $SESSION

################################################################################
# Setup the master node.
################################################################################
cat > "${NODES_ROOT}/master/ctl" <<EOF
#!/bin/zsh
bitumctl -C ../bitumctl.conf \$*
EOF
chmod +x "${NODES_ROOT}/master/ctl"

tmux rename-window -t $SESSION:0 'master'
tmux send-keys "cd ${NODES_ROOT}/master" C-m

echo "Starting ${NETWORK} bitumd"
if [ "${NETWORK}" = "mainnet" ]; then
tmux send-keys "bitumd -C ${BITUMD_CONF} \
--rpccert=${BITUMD_RPC_CERT} --rpckey=${BITUMD_RPC_KEY} \
--rpcuser=${RPC_USER} --rpcpass=${RPC_PASS} \
--miningaddr=${POOL_MINING_ADDR}" C-m
fi

if [ "${NETWORK}" = "testnet" ]; then
tmux send-keys "bitumd -C ${BITUMD_CONF} \
--rpccert=${BITUMD_RPC_CERT} --rpckey=${BITUMD_RPC_KEY} \
--rpcuser=${RPC_USER} --rpcpass=${RPC_PASS} \
--miningaddr=${POOL_MINING_ADDR} \
--testnet" C-m
fi

if [ "${NETWORK}" = "simnet" ]; then
tmux send-keys "bitumd -C ${BITUMD_CONF} \
--rpccert=${BITUMD_RPC_CERT} --rpckey=${BITUMD_RPC_KEY} \
--rpcuser=${RPC_USER} --rpcpass=${RPC_PASS} \
--miningaddr=${POOL_MINING_ADDR} \
--simnet" C-m
fi

################################################################################
# Setup the master node's bitumctl (mctl).
################################################################################
cat > "${NODES_ROOT}/master/mine" <<EOF
#!/bin/sh
  NUM=1
  case \$1 in
      ''|*[!0-9]*)  ;;
      *) NUM=\$1 ;;
  esac
  for i in \$(seq \$NUM) ; do
    bitumctl -C ../bitumctl.conf  generate 1
    sleep 0.5
  done
EOF
chmod +x "${NODES_ROOT}/master/mine"

tmux new-window -t $SESSION:1 -n 'mctl'
tmux send-keys "cd ${NODES_ROOT}/master" C-m

if [ "${NETWORK}" = "simnet" ]; then
echo "Waiting for some ${NETWORK} blocks to be mined"
sleep 2
# mine some blocks to start the chain if the active network is simnet.
tmux send-keys "./mine 2" C-m
fi

if [ "${NETWORK}" = "simnet" ]; then
################################################################################
# Setup the pool wallet.
################################################################################
cat > "${NODES_ROOT}/wallet/ctl" <<EOF
#!/bin/sh
bitumctl -C ../bitumwctl.conf --wallet \$*
EOF
chmod +x "${NODES_ROOT}/wallet/ctl"

tmux new-window -t $SESSION:2 -n 'wallet'
tmux send-keys "cd ${NODES_ROOT}/wallet" C-m
tmux send-keys "bitumwallet -C wallet.conf --create" C-m
echo "Creating ${NETWORK} wallet"
sleep 2
tmux send-keys "${WALLET_PASS}" C-m "${WALLET_PASS}" C-m "n" C-m "y" C-m
echo "Setting wallet passphrase"
sleep 1
tmux send-keys "${WALLET_SEED}" C-m C-m
tmux send-keys "bitumwallet -C wallet.conf" C-m

################################################################################
# Setup the pool wallet's bitumctl (wctl).
################################################################################
echo "Waiting for nodes to be synced"
sleep 6
# The consensus daemon must be synced for account generation to 
# work as expected.
tmux new-window -t $SESSION:3 -n 'wctl'
tmux send-keys "cd ${NODES_ROOT}/wallet" C-m
tmux send-keys "./ctl createnewaccount pfee" C-m
tmux send-keys "./ctl getnewaddress pfee" C-m
tmux send-keys "./ctl createnewaccount c1" C-m
tmux send-keys "./ctl getnewaddress c1" C-m
tmux send-keys "./ctl createnewaccount c2" C-m
tmux send-keys "./ctl getnewaddress c2" C-m
tmux send-keys "./ctl getnewaddress default" C-m
tmux send-keys "./ctl getbalance"
fi

################################################################################
# Setup bitumpool.
################################################################################
echo "Starting bitumpool"
sleep 4
tmux new-window -t $SESSION:4 -n 'pool'
tmux send-keys "cd ${NODES_ROOT}/pool" C-m
tmux send-keys "bitumpool --configfile pool.conf --homedir=${NODES_ROOT}/pool" C-m

if [ "${NETWORK}" = "simnet" ]; then
################################################################################
# Setup first mining client. 
################################################################################
echo "Starting mining client 1"
sleep 1
tmux new-window -t $SESSION:5 -n 'c1'
tmux send-keys "cd ${NODES_ROOT}/c1" C-m
tmux send-keys "miner --configfile=client.conf --homedir=${NODES_ROOT}/c1" C-m

################################################################################
# Setup another mining client. 
################################################################################
echo "Starting mining client 2"
sleep 1
tmux new-window -t $SESSION:6 -n 'c2'
tmux send-keys "cd ${NODES_ROOT}/c2" C-m
tmux send-keys "miner --configfile=client.conf --homedir=${NODES_ROOT}/c2" C-m
fi

tmux attach-session -t $SESSION
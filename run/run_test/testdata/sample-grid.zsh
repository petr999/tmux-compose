#!/usr/bin/zsh

cd /path/to/dumbclicker

tmux new -s dumbclicker-compose '
  printf '\''\033]2;%s\033\\'\'' '\''dumbclicker-compose'\''
  docker-compose up
  echo "
Commands for your consideration:
  docker-compose up
  docker-compose up -d"
  echo "Welcome to shell."
  zsh -l
' \; split-window '
  printf '\''\033]2;%s\033\\'\'' '\''dumbclicker_nginx_1'\''
  echo "Commands for your consideration:
  docker-compose up '\''nginx'\''
  docker-compose up -d '\''nginx'\''
  docker attach '\''dumbclicker_nginx_1'\''
  docker exec -it '\''dumbclicker_nginx_1'\''"
  echo "Hit Ctrl-C to break out to shell."
  echo "Attaching to container '\''dumbclicker_nginx_1'\''..."
  PID=0
  try_next=1
  trap '\''
    kill -INT $PID
    try_next=""
  '\'' SIGINT
  while [ '\''x1'\'' == "x${try_next}" ]; do
    /usr/bin/zsh -lc '\''
      docker attach dumbclicker_nginx_1
      sleep 1
    '\'' &
    PID=$!
    wait $PID
  done
  trap - SIGINT
  echo "
Commands for your consideration:
  docker-compose up '\''nginx'\''
  docker-compose up -d '\''nginx'\''
  docker attach '\''dumbclicker_nginx_1'\''
  docker exec -it '\''dumbclicker_nginx_1'\''"
  echo "Welcome to shell."
  /usr/bin/zsh -l
' \; split-window '
  printf '\''\033]2;%s\033\\'\'' '\''dumbclicker_h2o_1'\''
  echo "Commands for your consideration:
  docker-compose up '\''h2o'\''
  docker-compose up -d '\''h2o'\''
  docker attach '\''dumbclicker_h2o_1'\''
  docker exec -it '\''dumbclicker_h2o_1'\''"
  echo "Hit Ctrl-C to break out to shell."
  echo "Attaching to container '\''dumbclicker_h2o_1'\''..."
  PID=0
  try_next=1
  trap '\''
    kill -INT $PID
    try_next=""
  '\'' SIGINT
  while [ '\''x1'\'' == "x${try_next}" ]; do
    /usr/bin/zsh -lc '\''
      docker attach dumbclicker_h2o_1
      sleep 1
    '\'' &
    PID=$!
    wait $PID
  done
  trap - SIGINT
  echo "
Commands for your consideration:
  docker-compose up '\''h2o'\''
  docker-compose up -d '\''h2o'\''
  docker attach '\''dumbclicker_h2o_1'\''
  docker exec -it '\''dumbclicker_h2o_1'\''"
  echo "Welcome to shell."
  /usr/bin/zsh -l
' \; split-window '
  printf '\''\033]2;%s\033\\'\'' '\''dumbclicker_dumbclicker_1'\''
  echo "Commands for your consideration:
  docker-compose up '\''dumbclicker'\''
  docker-compose up -d '\''dumbclicker'\''
  docker attach '\''dumbclicker_dumbclicker_1'\''
  docker exec -it '\''dumbclicker_dumbclicker_1'\''"
  echo "Hit Ctrl-C to break out to shell."
  echo "Attaching to container '\''dumbclicker_dumbclicker_1'\''..."
  PID=0
  try_next=1
  trap '\''
    kill -INT $PID
    try_next=""
  '\'' SIGINT
  while [ '\''x1'\'' == "x${try_next}" ]; do
    /usr/bin/zsh -lc '\''
      docker attach dumbclicker_dumbclicker_1
      sleep 1
    '\'' &
    PID=$!
    wait $PID
  done
  trap - SIGINT
  echo "
Commands for your consideration:
  docker-compose up '\''dumbclicker'\''
  docker-compose up -d '\''dumbclicker'\''
  docker attach '\''dumbclicker_dumbclicker_1'\''
  docker exec -it '\''dumbclicker_dumbclicker_1'\''"
  echo "Welcome to shell."
  /usr/bin/zsh -l
' \; set -g pane-border-status bottom \; set-option status off \; select-layout tiled
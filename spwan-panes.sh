set -ex
addresses="127.0.0.1:3000,127.0.0.1:3001,127.0.0.1:3002"

tmux split-window -h -c "#{pane_current_path}" "tigerbeetle start --addresses=$addresses 0_0.tigerbeetle" \; \
  split-window -v -c "#{pane_current_path}" "tigerbeetle start --addresses=$addresses 0_1.tigerbeetle" \; \
  split-window -h -c "#{pane_current_path}" "tigerbeetle start --addresses=$addresses 0_2.tigerbeetle" \; \
  select-pane -t 1

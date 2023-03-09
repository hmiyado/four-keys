#!/usr/bin/env sh
make
chmod u+x four-keys
./four-keys timeSeries --since 2022-10-01 --interval month \
    | jq -r ".items[] | [.time, .deploymentFrequency, .leadTimeForChanges, .timeToRestore, .changeFailureRate] | @tsv" \
    > ./four-keys.tsv
gnuplot ./scripts/graph/draw_four_keys_graph.plt
mv *.jpg ./scripts/graph/
rm ./four-keys.tsv
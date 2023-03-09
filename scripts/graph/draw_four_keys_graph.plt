set xdata time
set timefmt "%Y-%m-%dT%H:%M:%SZ"
set xtics 60*60*24*30
set xtics format "%Y-%m-%d"
set mxtics 2
set term jpeg

set ytics 1
set output "deployment_frequency.jpg"
plot "four-keys.tsv" using 1:2 with lines title "deploymentFrequency" 

unset ytics
set yrange [0:*]
set output "lead_time_for_changes.jpg"
plot "four-keys.tsv" using 1:3 with lines title "leadTimeForChanges" 

set output "time_to_restore.jpg"
plot "four-keys.tsv" using 1:4 with lines title "timeToRestore"

set yrange [0:1]
set output "change_failure_rate.jpg"
plot "four-keys.tsv" using 1:5 with lines title "changeFailureRate"

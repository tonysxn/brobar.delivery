bash -c '
while true; do
  clear
  echo -e "\033[1;36mContainers status:\033[0m"
  docker ps -a --format "table {{.Names}}\t{{.Status}}\t{{.Image}}" | \
    awk '"'"'
      NR==1 {print "\033[1;33m" $0 "\033[0m"; next}
      /Up/ {print "\033[1;32m" $0 "\033[0m"; next}
      /Exited/ {print "\033[1;31m" $0 "\033[0m"; next}
      {print $0}
    '"'"'
  sleep 1
done
'

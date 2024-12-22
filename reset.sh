set -ex
rm *.tigerbeetle || true
tigerbeetle format --cluster=0 --replica=0 --replica-count=3 --development 0_0.tigerbeetle
tigerbeetle format --cluster=0 --replica=1 --replica-count=3 --development 0_1.tigerbeetle
tigerbeetle format --cluster=0 --replica=2 --replica-count=3 --development 0_2.tigerbeetle

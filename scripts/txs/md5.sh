#! /bin/bash
set -u

$(> temp)
for((i=1;i<=5000;i++));  
do   
    echo -n $(expr $i  + 1)|md5sum|cut -d ' ' -f1 >> temp
done

function toHex() {
	echo -n $1 | hexdump -ve '1/1 "%.2X"'
}

N=$1
#PORT=$2

md5VALUE=$(hexdump -ve '1/1 "%.2X"'  temp)
#echo $VALUE

for i in `seq 1 $N`; do
	# store key value pair
	#KEY=$(head -c 10 /dev/urandom)
        #VALUE="$i"
       # echo   $(toHex $KEY=$VALUE)$md5VALUE
       #echo '0x$(toHex $KEY=$VALUE)$md5VALUE'
       #echo '0x$(toHex $KEY=$VALUE)$md5VALUE' | 
       #curl 127.0.0.1:$PORT/broadcast_tx_sync?tx=0x$(toHex $KEY=$VALUE)$md5VALUE

       ./rpc  -random=true  -tx $md5VALUE
done

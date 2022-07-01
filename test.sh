#!/bin/bash

function is_equal () {
  if [ "$1" = "$2" ]; then
    echo -n -e "\033[1;32m[OK]\033[m "
    return 0
  else
    echo -n -e "\033[1;31m[NG]\033[m "
    return 1
  fi
}

if [ -z "$API_ENDPOINT" ];then
  API_ENDPOINT="http://localhost:2022/api/v1alpha1"
else
  API_ENDPOINT=${API_ENDPOINT%/}
fi
if [ -z "$PING_CHECK" ];then
  PING_CHECK=false
fi

echo "==================="
echo "#   Environment   #"
echo "==================="
echo -n "API_ENDPOINT: "
echo $API_ENDPOINT
echo -n "PING_CHECK  : "
echo $PING_CHECK

sleep 1

echo "==================="
echo "#     Archive     #"
echo "==================="

sleep 1

EXPECT="https://minio.home.arpa:9000/kubevirt/images/ubuntu-20.04-server-cloudimg-arm64.img"
ACTUAL=`curl -s -XPOST -H 'Content-Type:application/json' \
-d '{ \
  "name": "test", \
  "repository":"https://minio.home.arpa:9000/kubevirt/images/ubuntu-20.04-server-cloudimg-arm64.img" \
}' \
$API_ENDPOINT/archives | jq .repository | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
RET=$?
echo "Creating Archive"
if [ $RET -ne 0 ];then
  exit 1
fi

sleep 1

EXPECT="https://minio.home.arpa:9000/kubevirt/images/ubuntu-20.04-server-cloudimg-arm64.img"
ACTUAL=`curl -s -XGET $API_ENDPOINT/archives/test | jq .repository | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
RET=$?
echo "Reading Archive"
if [ $RET -ne 0 ];then
  exit 1
fi

sleep 1

EXPECT="http://minio.home.arpa:9000/kubevirt/images/ubuntu-20.04-server-cloudimg-arm64.img"
ACTUAL=`curl -s -XPUT -H 'Content-Type:application/json' \
-d '{ \
  "name": "test", \
  "repository":"http://minio.home.arpa:9000/kubevirt/images/ubuntu-20.04-server-cloudimg-arm64.img" \
}' \
$API_ENDPOINT/archives/test | jq .repository | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
RET=$?
echo "Updating Archive"
if [ $RET -ne 0 ];then
  exit 1
fi

sleep 1

echo "====================="
echo "#     CloudInit     #"
echo "====================="

sleep 1

EXPECT="#cloud-config\ntimezone: Asia/Tokyo\nssh_pwauth: True\npassword: ubuntu\nchpasswd: { expire: False }\ndisable_root: false\n"
ACTUAL=`curl -s -XPOST -H 'Content-Type:application/json' \
-d '{ \
  "name": "test", \
  "user_data":"#cloud-config\ntimezone: Asia/Tokyo\nssh_pwauth: True\npassword: ubuntu\nchpasswd: { expire: False }\ndisable_root: false\n" \
}' \
$API_ENDPOINT/cloudinits | jq .user_data | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
RET=$?
echo "Creating CloudInit"
if [ $RET -ne 0 ];then
  exit 1
fi

sleep 1

EXPECT="#cloud-config\ntimezone: Asia/Tokyo\nssh_pwauth: True\npassword: ubuntu\nchpasswd: { expire: False }\ndisable_root: false\n"
ACTUAL=`curl -s -XGET $API_ENDPOINT/cloudinits/test | jq .user_data | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
RET=$?
echo "Reading CloudInit"
if [ $RET -ne 0 ];then
  exit 1
fi

sleep 1

EXPECT="#cloud-config\ntimezone: Asia/Tokyo\nssh_pwauth: True\npassword: UBUNTU\ndisable_root: true\n"
ACTUAL=`curl -s -XPUT -H 'Content-Type:application/json' \
-d '{ \
  "name": "test", \
  "user_data":"#cloud-config\ntimezone: Asia/Tokyo\nssh_pwauth: True\npassword: UBUNTU\ndisable_root: true\n" \
    }' \
$API_ENDPOINT/cloudinits/test | jq .user_data | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
RET=$?
echo "Updating CloudInit"
if [ $RET -ne 0 ];then
  exit 1
fi

sleep 1

echo "================"
echo "#     Disk     #"
echo "================"

sleep 1

EXPECT='{
  "archive": {
    "name": "test"
  }
}'
ACTUAL=`curl -s -XPOST -H 'Content-Type:application/json' \
-d '{ \
  "name": "test", \
  "size": "16Gi", \
  "source": { \
    "archive": { \
      "name": "test" \
    } \
  } \
}' \
$API_ENDPOINT/disks | jq .source`
is_equal "$EXPECT" "$ACTUAL"
RET=$?
echo "Creating Disk"
if [ $RET -ne 0 ];then
  exit 1
fi

sleep 1

echo -n "Creating..."
while [ `kubectl get disk test -o json | jq .status.phase | tr -d '"'` != "Created" ]
do
  sleep 5
  echo -n "."
done
echo "OK"

sleep 1

EXPECT='{
  "archive": {
    "name": "test"
  }
}'
ACTUAL=`curl -s -XGET $API_ENDPOINT/disks/test | jq .source`
is_equal "$EXPECT" "$ACTUAL"
RET=$?
echo "Reading Disk"
if [ $RET -ne 0 ];then
  exit 1
fi

sleep 1

EXPECT="{}"
ACTUAL=`curl -s -XPOST -H 'Content-Type:application/json' \
-d '{ \
  "name": "test-emptydisk", \
  "size": "16Gi", \
  "source": {} \
}' \
$API_ENDPOINT/disks | jq .source | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
RET=$?
echo "Creating empty disk"
if [ $RET -ne 0 ];then
  exit 1
fi

sleep 1

EXPECT="{}"
ACTUAL=`curl -s -XGET $API_ENDPOINT/disks/test-emptydisk | jq .source | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
RET=$?
echo "Reading empty disk"
if [ $RET -ne 0 ];then
  exit 1
fi

sleep 3

EXPECT="Created"
ACTUAL=`kubectl get disk test-emptydisk -o json | jq .status.phase | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
RET=$?
echo 'empty disk status.phase("Created")'
if [ $RET -ne 0 ];then
  exit 1
fi

sleep 1

EXPECT="{}"
ACTUAL=`curl -s -XPOST -H 'Content-Type:application/json' \
-d '{ \
  "name": "test-emptydisk-nosource", \
  "size": "16Gi" \
}' \
$API_ENDPOINT/disks | jq .source | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
RET=$?
echo "Creating empty disk (nosource)"
if [ $RET -ne 0 ];then
  exit 1
fi

sleep 1

EXPECT="{}"
ACTUAL=`curl -s -XGET $API_ENDPOINT/disks/test-emptydisk-nosource | jq .source | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
RET=$?
echo "Reading emtpy disk (nosource)"
if [ $RET -ne 0 ];then
  exit 1
fi

sleep 3

EXPECT="Created"
ACTUAL=`kubectl get disk test-emptydisk-nosource -o json | jq .status.phase | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
RET=$?
echo 'empty disk status.phase("Created") (nosource)'
if [ $RET -ne 0 ];then
  exit 1
fi

sleep 1

echo "=================="
echo "#     Server     #"
echo "=================="

sleep 1

EXPECT='{
  "name": "test",
  "running": true,
  "cpu": "2",
  "memory": "2Gi",
  "mac_address": "52:42:00:4f:8a:2b",
  "hostname": "test",
  "hosting": "node-1.k8s.home.arpa",
  "disk": {
    "name": "test"
  },
  "cloudinit": {
    "name": "test"
  }
}'
ACTUAL=`curl -s -XPOST -H 'Content-Type:application/json' \
-d '{ \
  "name": "test", \
  "running": true, \
  "cpu": "2", \
  "memory": "2Gi", \
  "mac_address": "52:42:00:4f:8a:2b", \
  "hostname": "test", \
  "hosting": "node-1.k8s.home.arpa", \
  "disk": { \
    "name": "test" \
  }, \
  "cloudinit": { \
    "name": "test" \
  } \
}' \
$API_ENDPOINT/servers | jq .`
is_equal "$EXPECT" "$ACTUAL"
RET=$?
echo "Creating Server"
if [ $RET -ne 0 ];then
  exit 1
fi

sleep 1

echo -n "Starting..."
STATE=`kubectl get server test -o json | jq .status.state | tr -d '"'`
while [ "$STATE" != "Running" ]
do
    sleep 1
    echo -n "."
    STATE=`kubectl get server test -o json | jq .status.state | tr -d '"'`
done
echo "OK"

sleep 1

if $PING_CHECK; then
  echo -n "Ping Checking..."
  IP=`kubectl get server test -o json | jq .status.ip | tr -d '"'`

  ping -c 1 -W 1 $IP > /dev/null
  PING=$?
  while [ $PING -ne 0 ]
  do
      echo -n "."
      ping -c 1 -W 1 $IP > /dev/null
      PING=$?
  done
  echo "OK"
fi

sleep 1

EXPECT='{
  "name": "test",
  "running": true,
  "cpu": "2",
  "memory": "2Gi",
  "mac_address": "52:42:00:4f:8a:2b",
  "hostname": "test",
  "hosting": "node-1.k8s.home.arpa",
  "disk": {
    "name": "test"
  },
  "cloudinit": {
    "name": "test"
  }
}'
ACTUAL=`curl -s -XGET $API_ENDPOINT/servers/test | jq .`
is_equal "$EXPECT" "$ACTUAL"
RET=$?
echo "Read Server"
if [ $RET -ne 0 ];then
  exit 1
fi

sleep 1

EXPECT="Running" 
ACTUAL=`kubectl get server test -o json | jq .status.state | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
RET=$?
echo 'Server status.state("Running")'
if [ $RET -ne 0 ];then
  exit 1
fi

sleep 1

EXPECT="node-1.k8s.home.arpa"
ACTUAL=`kubectl get server test -o json | jq .status.hosting | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
RET=$?
echo 'Server status.hosting("node-1.k8s.home.arpa")'
if [ $RET -ne 0 ];then
  exit 1
fi

sleep 1

EXPECT='{
  "name": "test",
  "running": false,
  "cpu": "2",
  "memory": "2Gi",
  "mac_address": "52:42:00:4f:8a:2b",
  "hostname": "test",
  "hosting": "node-1.k8s.home.arpa",
  "disk": {
    "name": "test"
  },
  "cloudinit": {
    "name": "test"
  }
}'
ACTUAL=`curl -s -XPUT -H 'Content-Type:application/json' \
-d '{ \
  "name": "test", \
  "running": false, \
  "cpu": "2", \
  "memory": "2Gi", \
  "mac_address": "52:42:00:4f:8a:2b", \
  "hostname": "test", \
  "hosting": "node-1.k8s.home.arpa", \
  "disk": { \
    "name": "test" \
  }, \
  "cloudinit": { \
    "name": "test" \
  } \
}' \
$API_ENDPOINT/servers/test | jq .`
is_equal "$EXPECT" "$ACTUAL"
RET=$?
echo "Update Server (change running: false)"
if [ $RET -ne 0 ];then
  exit 1
fi

sleep 1

echo -n "Stopping..."
while [ `kubectl get server test -o json | jq .status.state | tr -d '"'` != "Stopped" ]
do
  echo -n "."
  sleep 3
done
echo "OK"

sleep 1

EXPECT="Stopped" 
ACTUAL=`kubectl get server test -o json | jq .status.state | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
RET=$?
echo 'Server status.state("Stopped")'
if [ $RET -ne 0 ];then
  exit 1
fi

sleep 1

echo "===================="
echo "#     Deleting     #"
echo "===================="

sleep 1

EXPECT="ok"

ACTUAL=`curl -s -XDELETE $API_ENDPOINT/servers/test | jq .message | tr -d '"'`
is_equal $EXPECT $ACTUAL
echo "Deleting Server"

sleep 3

ACTUAL=`curl -s -XDELETE $API_ENDPOINT/disks/test | jq .message | tr -d '"'`
is_equal $EXPECT $ACTUAL
echo "Deleting Disk"

sleep 3

ACTUAL=`curl -s -XDELETE $API_ENDPOINT/disks/test-emptydisk-nosource | jq .message | tr -d '"'`
is_equal $EXPECT $ACTUAL
echo "Deleting empty disk (nosource)"

sleep 3

ACTUAL=`curl -s -XDELETE $API_ENDPOINT/disks/test-emptydisk | jq .message | tr -d '"'`
is_equal $EXPECT $ACTUAL
echo "Deleting empty disk"

sleep 3

ACTUAL=`curl -s -XDELETE $API_ENDPOINT/cloudinits/test | jq .message | tr -d '"'`
is_equal $EXPECT $ACTUAL
echo "Deleting CloudInit"

sleep 3

ACTUAL=`curl -s -XDELETE $API_ENDPOINT/archives/test | jq .message | tr -d '"'`
is_equal $EXPECT $ACTUAL
echo "Deleting Archive"

exit 0

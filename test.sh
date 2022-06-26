#!/bin/bash

function is_equal () {
  if [ "$1" = "$2" ]; then
    echo -n -e "\033[1;32m[OK]\033[m "
  else
    echo -n -e "\033[1;31m[NG]\033[m "
  fi
}

echo "==================="
echo "#     Archive     #"
echo "==================="

EXPECT="https://minio.home.arpa:9000/kubevirt/images/ubuntu-20.04-server-cloudimg-arm64.img"
ACTUAL=`curl -s -XPOST -H 'Content-Type:application/json' \
-d '{ \
  "name": "test", \
  "url":"https://minio.home.arpa:9000/kubevirt/images/ubuntu-20.04-server-cloudimg-arm64.img" \
}' \
localhost:2022/api/v1alpha1/archives | jq .url | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
echo "Creating Archive"

EXPECT="https://minio.home.arpa:9000/kubevirt/images/ubuntu-20.04-server-cloudimg-arm64.img"
ACTUAL=`curl -s -XGET localhost:2022/api/v1alpha1/archives/test | jq .url | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
echo "Reading Archive"

EXPECT="http://minio.home.arpa:9000/kubevirt/images/ubuntu-20.04-server-cloudimg-arm64.img"
ACTUAL=`curl -s -XPUT -H 'Content-Type:application/json' \
-d '{ \
  "name": "test", \
  "url":"http://minio.home.arpa:9000/kubevirt/images/ubuntu-20.04-server-cloudimg-arm64.img" \
}' \
localhost:2022/api/v1alpha1/archives/test | jq .url | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
echo "Updating Archive"

echo "====================="
echo "#     CloudInit     #"
echo "====================="

EXPECT="#cloud-config\ntimezone: Asia/Tokyo\nssh_pwauth: True\npassword: ubuntu\nchpasswd: { expire: False }\ndisable_root: false\n"
ACTUAL=`curl -s -XPOST -H 'Content-Type:application/json' \
-d '{ \
  "name": "test", \
  "user_data":"#cloud-config\ntimezone: Asia/Tokyo\nssh_pwauth: True\npassword: ubuntu\nchpasswd: { expire: False }\ndisable_root: false\n" \
}' \
http://localhost:2022/api/v1alpha1/cloudinits | jq .user_data | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
echo "Creating CloudInit"

EXPECT="#cloud-config\ntimezone: Asia/Tokyo\nssh_pwauth: True\npassword: ubuntu\nchpasswd: { expire: False }\ndisable_root: false\n"
ACTUAL=`curl -s -XGET localhost:2022/api/v1alpha1/cloudinits/test | jq .user_data | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
echo "Reading CloudInit"

EXPECT="#cloud-config\ntimezone: Asia/Tokyo\nssh_pwauth: True\npassword: UBUNTU\ndisable_root: true\n"
ACTUAL=`curl -s -XPUT -H 'Content-Type:application/json' \
-d '{ \
  "name": "test", \
  "user_data":"#cloud-config\ntimezone: Asia/Tokyo\nssh_pwauth: True\npassword: UBUNTU\ndisable_root: true\n" \
    }' \
http://localhost:2022/api/v1alpha1/cloudinits/test | jq .user_data | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
echo "Updating CloudInit"

echo "================"
echo "#     Disk     #"
echo "================"

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
http://localhost:2022/api/v1alpha1/disks | jq .source`
is_equal "$EXPECT" "$ACTUAL"
echo "Creating Disk"

sleep 3

echo -n "Creating..."
while [ `kubectl get disk test -o json | jq .status.phase | tr -d '"'` != "Created" ]
do
  echo -n "."
  sleep 5
done
echo "OK"

EXPECT='{
  "archive": {
    "name": "test"
  }
}'
ACTUAL=`curl -s -XGET localhost:2022/api/v1alpha1/disks/test | jq .source`
is_equal "$EXPECT" "$ACTUAL"
echo "Reading Disk"

EXPECT="{}"
ACTUAL=`curl -s -XPOST -H 'Content-Type:application/json' \
-d '{ \
  "name": "test-emptydisk", \
  "size": "16Gi", \
  "source": {} \
}' \
http://localhost:2022/api/v1alpha1/disks | jq .source | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
echo "Creating empty disk"

EXPECT="{}"
ACTUAL=`curl -s -XGET localhost:2022/api/v1alpha1/disks/test-emptydisk | jq .source | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
echo "Reading empty disk"

sleep 3

EXPECT="Created"
ACTUAL=`kubectl get disk test-emptydisk -o json | jq .status.phase | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
echo 'empty disk status.phase("Created")'

EXPECT="{}"
ACTUAL=`curl -s -XPOST -H 'Content-Type:application/json' \
-d '{ \
  "name": "test-emptydisk-nosource", \
  "size": "16Gi" \
}' \
http://localhost:2022/api/v1alpha1/disks | jq .source | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
echo "Creating empty disk (nosource)"

EXPECT="{}"
ACTUAL=`curl -s -XGET localhost:2022/api/v1alpha1/disks/test-emptydisk-nosource | jq .source | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
echo "Reading emtpy disk (nosource)"

sleep 3

EXPECT="Created"
ACTUAL=`kubectl get disk test-emptydisk-nosource -o json | jq .status.phase | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
echo 'empty disk status.phase("Created") (nosource)'

echo "=================="
echo "#     Server     #"
echo "=================="

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
http://localhost:2022/api/v1alpha1/servers | jq .`
is_equal "$EXPECT" "$ACTUAL"
echo "Creating Server"

echo -n "Starting..."
STATE=`kubectl get server test -o json | jq .status.state | tr -d '"'`
while [ "$STATE" != "Running" ]
do
    echo -n "."
    STATE=`kubectl get server test -o json | jq .status.state | tr -d '"'`
    sleep 1
done
echo "OK"

echo -n "Ping Checking..."
ping -c 1 -W 1 `kubectl get server test -o json | jq .status.ip | tr -d '"'` > /dev/null
PING=$?
while [ $PING -ne 0 ]
do
    echo -n "."
    ping -c 1 -W 1 `kubectl get server test -o json | jq .status.ip | tr -d '"'` > /dev/null
    PING=$?
done
echo "OK"

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
ACTUAL=`curl -s -XGET localhost:2022/api/v1alpha1/servers/test | jq .`
is_equal "$EXPECT" "$ACTUAL"
echo "Read Server"

EXPECT="Running" 
ACTUAL=`kubectl get server test -o json | jq .status.state | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
echo 'Server status.state("Running")'

EXPECT="node-1.k8s.home.arpa"
ACTUAL=`kubectl get server test -o json | jq .status.hosting | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
echo 'Server status.hosting("node-1.k8s.home.arpa")'

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
http://localhost:2022/api/v1alpha1/servers/test | jq .`
is_equal "$EXPECT" "$ACTUAL"
echo "Update Server (change running: false)"

echo -n "Stopping..."
while [ `kubectl get server test -o json | jq .status.state | tr -d '"'` != "Stopped" ]
do
  echo -n "."
  sleep 3
done
echo "OK"

EXPECT="Stopped" 
ACTUAL=`kubectl get server test -o json | jq .status.state | tr -d '"'`
is_equal "$EXPECT" "$ACTUAL"
echo 'Server status.state("Stopped")'

echo "===================="
echo "#     Deleting     #"
echo "===================="

EXPECT="ok"

ACTUAL=`curl -s -XDELETE localhost:2022/api/v1alpha1/servers/test | jq .message | tr -d '"'`
is_equal $EXPECT $ACTUAL
echo "Deleting Server"

sleep 1

ACTUAL=`curl -s -XDELETE localhost:2022/api/v1alpha1/disks/test | jq .message | tr -d '"'`
is_equal $EXPECT $ACTUAL
echo "Deleting Disk"

sleep 1

ACTUAL=`curl -s -XDELETE localhost:2022/api/v1alpha1/disks/test-emptydisk-nosource | jq .message | tr -d '"'`
is_equal $EXPECT $ACTUAL
echo "Deleting empty disk (nosource)"

sleep 1

ACTUAL=`curl -s -XDELETE localhost:2022/api/v1alpha1/disks/test-emptydisk | jq .message | tr -d '"'`
is_equal $EXPECT $ACTUAL
echo "Deleting empty disk"

sleep 1

ACTUAL=`curl -s -XDELETE localhost:2022/api/v1alpha1/cloudinits/test | jq .message | tr -d '"'`
is_equal $EXPECT $ACTUAL
echo "Deleting CloudInit"

sleep 1

ACTUAL=`curl -s -XDELETE localhost:2022/api/v1alpha1/archives/test | jq .message | tr -d '"'`
is_equal $EXPECT $ACTUAL
echo "Deleting Archive"

exit 0

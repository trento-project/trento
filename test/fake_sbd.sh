#!/bin/bash

# /usr/sbin/sbd -d /dev/disk/by-id/scsi-SLIO-ORG_IBLOCK_e9177a8a-029b-49af-9753-4440882a39cf dump
if [ $3 == "dump" ]; then
  cat <<RESULT
==Dumping header on disk /dev/disk/by-id/scsi-SLIO-ORG_IBLOCK_e9177a8a-029b-49af-9753-4440882a39cf
Header version     : 2.1
UUID               : f9ba490e-0f14-4908-859a-ace97aafaf34
Number of slots    : 255
Sector size        : 512
Timeout (watchdog) : 5
Timeout (allocate) : 2
Timeout (loop)     : 1
Timeout (msgwait)  : 10
==Header on disk /dev/disk/by-id/scsi-SLIO-ORG_IBLOCK_e9177a8a-029b-49af-9753-4440882a39cf is dumped
RESULT
fi



# /usr/sbin/sbd -d /dev/disk/by-id/scsi-SLIO-ORG_IBLOCK_e9177a8a-029b-49af-9753-4440882a39cf list
if [ $3 == "list" ]; then
  cat <<RESULT
0	vmhana01	clear	
1	vmhana02	clear
RESULT
fi
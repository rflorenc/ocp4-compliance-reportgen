#!/bin/bash

scans=$(oc -n openshift-compliance get compliancesuite -o jsonpath='{.items[*].spec.scans[*].name}')

if [[ ! -z "$scans" ]]; then
  for scan in $scans
  do
    echo $scan
    sed "s/{SCAN_NAME}/${scan}/g" Deployment.yaml | oc apply -f -
    oc rollout status deployment/${scan}-html-results -w
    echo 
  done
echo "Results published at:"
oc get routes -o jsonpath='{range .items[*]}{.spec.host}{"\n"}{end}' | xargs -L1 -I {}  echo https://{}
fi

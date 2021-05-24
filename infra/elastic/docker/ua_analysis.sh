#!/bin/bash
if bin/elasticsearch-plugin list | grep -q 'ukrainian'; then
   echo "analysis-ukrainian exists"
else
   bin/elasticsearch-plugin install analysis-ukrainian
   echo "analysis-ukrainian installed"
fi
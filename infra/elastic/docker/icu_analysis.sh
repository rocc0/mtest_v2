if bin/elasticsearch-plugin list | grep -q 'icu'; then
   echo "analysis-icu exists"
else
   bin/elasticsearch-plugin install analysis-icu
   echo "analysis-icu installed"
fi
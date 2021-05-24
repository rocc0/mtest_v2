# https://github.com/elastic/elasticsearch-docker
#FROM elasticsearch:5.3.0 # deprecated images https://github.com/deviantony/docker-elk/issues/119
FROM elasticsearch:7.12.1
RUN bin/elasticsearch-plugin install analysis-ukrainian
RUN bin/elasticsearch-plugin install analysis-icu
# Add your elasticsearch plugins setup here
# Example: RUN edocklasticsearch-plugin install analysis-icu
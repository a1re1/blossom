mvn archetype:generate \
	-DarchetypeGroupId=org.codehaus.mojo.archetypes \
	-DarchetypeArtifactId=pom-root \
	-DarchetypeVersion=RELEASE \
	-DinteractiveMode=false \
	-DgroupId=$0 \
	-DartifactId=$1

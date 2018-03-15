# Set the base image

FROM ubuntu


# Set the file maintainer

MAINTAINER Pallavi Debnath <pallavi.debnath@turbonomic.com>


ADD _output/mesosturbo.linux /bin/mesosturbo


ENTRYPOINT ["/bin/mesosturbo"]

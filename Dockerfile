FROM registry.access.redhat.com/ubi9/go-toolset:latest as builder
ENV GOPATH=$APP_ROOT
COPY --chown=1001:0 . .
RUN make cmd

FROM registry.access.redhat.com/ubi9/ubi-minimal
USER root
RUN echo -e "[centos9]" \
 "\nname = centos9" \
 "\nmetalink = https://mirrors.centos.org/metalink?repo=centos-appstream-9-stream&arch=\$basearch&protocol=https,http" \
 "\nenabled = 1" \
 "\ngpgcheck = 0" > /etc/yum.repos.d/centos.repo
RUN microdnf -y install --disablerepo=centos9 --setopt=install_weak_deps=0 --setopt=tsflags=nodocs \
  gcc \
  git \
  java-11-openjdk-devel \
  maven \
  openssh-clients \
  python3 \
  python3-lxml \
  python3-numpy \
  python3-psutil \
  python3-pip \
  python3-scipy \
  python3-setuptools \
  unzip \
  wget \
 && microdnf -y install --setopt=install_weak_deps=0 --setopt=tsflags=nodocs \
 ant \
 ant-junit \
 subversion \
 && microdnf -y clean all

ARG TESTGEN=https://github.com/konveyor/tackle-test-generator-cli/releases/download/v2.4.0/tackle-test-generator-cli-v2.4.0-all-deps.zip
RUN wget -qO /opt/tackle-test-generator-cli.zip $TESTGEN \
 && unzip /opt/tackle-test-generator-cli.zip -x */tackle-test-generator-ui-main-SNAPSHOT-jar-with-dependencies.jar -d /opt \
 && rm /opt/tackle-test-generator-cli.zip

# Override setup.py with custom versions with limited dependencies (test unit-only)
COPY hack/setup.py /opt/tackle-test-generator-cli/

# Install tkltest-unit
RUN pip3 install --no-cache-dir torch torchvision torchaudio --extra-index-url https://download.pytorch.org/whl/cpu
RUN cd /opt/tackle-test-generator-cli && pip3 install --no-cache-dir --editable .

WORKDIR /working
COPY --from=builder /opt/app-root/src/bin/addon /usr/local/bin/addon

ENV JAVA_HOME /usr/lib/jvm/java-11-openjdk/

# Test availability of tkltest-unit
RUN tkltest-unit --help

ENTRYPOINT ["/usr/local/bin/addon"]

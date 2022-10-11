FROM registry.access.redhat.com/ubi8/go-toolset:1.16.12 as builder
ENV GOPATH=$APP_ROOT
COPY --chown=1001:0 . .
RUN make cmd

FROM registry.access.redhat.com/ubi8/ubi-minimal
USER root
RUN echo -e "[centos8]" \
 "\nname = centos8" \
 "\nbaseurl = http://mirror.centos.org/centos/8-stream/AppStream/x86_64/os/" \
 "\nenabled = 1" \
 "\ngpgcheck = 0" > /etc/yum.repos.d/centos.repo
RUN microdnf -y install \
  java-11-openjdk-devel \
  openssh-clients \
  unzip \
  wget \
  git \
  subversion \
  maven \
  python39 \
 && microdnf -y clean all
ARG TESTGEN=https://github.com/konveyor/tackle-test-generator-cli/releases/download/v2.4.0/tackle-test-generator-cli-v2.4.0-all-deps.zip
RUN wget -qO /opt/tackle-test-generator-cli.zip $TESTGEN \
 && unzip /opt/tackle-test-generator-cli.zip -x */tackle-test-generator-ui-main-SNAPSHOT-jar-with-dependencies.jar -d /opt \
 && rm /opt/tackle-test-generator-cli.zip

# Override setup.py with custom versions with limited dependencies (test unit-only)
COPY hack/setup.py /opt/tackle-test-generator-cli/

# Install tkltest-unit
RUN pip3 install torch torchvision torchaudio --extra-index-url https://download.pytorch.org/whl/cpu
RUN cd /opt/tackle-test-generator-cli && pip3 install --editable .

WORKDIR /working
COPY --from=builder /opt/app-root/src/bin/addon /usr/local/bin/addon

RUN alternatives --set java java-11-openjdk.x86_64
RUN alternatives --set javac java-11-openjdk.x86_64
ENV JAVA_HOME /usr/lib/jvm/java-11-openjdk/

# Test availability of tkltest-unit
RUN tkltest-unit --help

ENTRYPOINT ["/usr/local/bin/addon"]

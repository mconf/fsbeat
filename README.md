# Fsbeat

Welcome to Fsbeat.

Fsbeat is a lightweight shipper agent that reads events from FreeSWITCH Event Socket Library (ESL) and sends them to Logstash for processing or Elasticsearch for indexing.

## Getting started with Fsbeat

There is just a few settings to get Fsbeat up and running. These are shown below:

* `freeswitch.server`: Address (IP or FQDN) of the server running FreeSWITCH. Default is "_localhost_".
* `freeswitch.port`: Port number of the FreeSWITCH ESL. Default is "_8021_".
* `freeswitch.auth`: Authentication of the FreeSWITCH ESL. Default is "_ClueCon_". __\*__
* `freeswitch.events`: List of events Fsbeat should read. Default is "_all_". __\*\*__

> __\*__ Consider changing the authentication of the FreeSWITCH ESL.
> __\*\*__ Note that the events should be whitespace-separated within a single string. For instance, `"channel_create channel_state channel_destroy"`.

### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Init project

Ensure that this folder is at the following location: ${GOPATH}/github.com/mconftec

To get running with Fsbeat and also install the dependencies, run the following command:

```
make setup
```

It will create a clean Git history for each major step.
Note that you can always rewrite the history if you wish before pushing your changes.

To push Fsbeat in the Git repository, run the following commands:

```
git remote set-url origin https://github.com/mconftec/fsbeat
git push origin master
```

For further development, check out the [Developer Guide: Creating a New Beat](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Fsbeat run the command below.
This will generate a binary in the same directory with the name _fsbeat_.

```
make
```

### Run

To run Fsbeat with debugging output enabled, run:

```
./fsbeat -c fsbeat.yml -e -d "*"
```


### Test

To test Fsbeat, run the following command:

```
make testsuite
```

Alternatively:

```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each Beat has a template for the mapping in Elasticsearch and a documentation for the fields which is automatically generated based on `etc/fields.yml`.
To generate `etc/fsbeat.template.json` and `etc/fsbeat.asciidoc`, run:

```
make update
```


### Cleanup

To clean Fsbeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```

> Make sure you run `make fmt; make simplify; make clean` before committing any changings to Git.

### Clone

To clone Fsbeat from the Git repository, run the following commands:

```
mkdir -p ${GOPATH}/github.com/mconftec
cd ${GOPATH}/github.com/mconftec
git clone https://github.com/mconftec/fsbeat
```

For further development, check out the [Developer Guide: Creating a New Beat](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

## Packaging

The Beat framework provides tools to cross-compile and package your Beat for different platforms. This requires [Docker](https://www.docker.com/).

To build packages of your Beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The whole process to finish can take several minutes.

> If you experience problems, make sure you have Docker running correctly. You can check it with the `docker version` command. It should not display any warnings or errors about the Docker server.

# Releases and binaries

The latest stable version of Fsbeat is _1.0.1-alpha_.

For a complete list of releases check out [Releases](https://github.com/mconftec/fsbeat/releases).

Each release is accompanied by a Debian package (.deb) to install Fsbeat. If other formats are needed, it is possible to cross-compile Fsbeat source code and look for the desired binary at `build/upload/` directory following the instructions above.

# Credits

Fsbeat uses [go-eventsocket](https://github.com/fiorix/go-eventsocket) by Alexander Fiorix which is released under BSD 3-Clause License.

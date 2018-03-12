YARA-ENDPOINT
=============

Yara as Endpoint is not just an enpoint solution for scanning files, Yara-Enpoint can be used as incident handler solution. While you are in the middle of an incident you have to know what is the scope of it in terms to act properly. You can do that by running your IoC manually against all your assets or using Yara-Endpoint and do it automatically and centralized.

## How does it work?
Yara-Endpoint follows a client-server architecture so it is really easy to deploy. But getting deeper Yara-Endpoint has two componets `client` and `server`. Both the `server` as well as the `client` are a standalone binaries, no installation needed!. The `client` only needs a couple of flags that indicates where is the `server` and which port should be used. On the other hand, the `server` reads its configuration from a file, but basicaly it exposes two ports one for the comunitacion with the `clients` and other for a web management interface.

## Main features
Yara-Endpoint offers an easy solution as either antivirus like endpoint or incident response tool. In both cases the installation and deploy is really easy, we have already taken care of it, because we know that deploying this kind of things is a pain in the ass.

For now we have implemented the following features:
1. There is no need to register endpoint first, start using it and configure the endpoints later.
1. Scan files, directories or PID.
1. Tag Endpoints according your needs.
1. Tag rules according your needs.
1. Manage everything from a web UI.

## Requirements
We do not have a lot of requiremets but some would be:
1. Execute the `client` as Administrator or root, if you want a full scan from the client.
1. A MongoDB database to store data on the server.

## How to build the project
### Using XGO
We have develop a Makefile that uses a [xgo](https://github.com/karalabe/xgo) to cross-compile Golang projects. Xgo is a Docker image with a full set of compilers and libraries already ready to cross compile. For example, it has gcc-mingw-w64, gcc-5-mips-linux-gnu, gcc-5-arm-linux-gnueabi, etc.

Using xgo project is really easy to cross compile any Golang project with C-like dependency since xgo cross compile the dependency and then cross compile statically your Golang project. This will generate a standalone binary with everything on it and ready to be executed on the target system.

The problem is that Yara does not use a standard Linux compilation way (configure && make && make install). Yara needs to run a bootstrap.sh before ./configure so xgo will fail. We have an [issue](https://github.com/karalabe/xgo/issues/105) for this matter.

The way that we solve this is modifying [xgo](https://github.com/Xumeiquer/xgo) project to run `bootstrap.sh` before `./configure`. The counterpart of using the modified version of xgo is that you have to build the docker image first.

#### NOTE
We are unable to cross compile Yara-Endpoint for Linux targets, we are working on it, but for now you can compile Yara-Enpoint manually.

### Manually
Other way to get the binaries is compiling the project manually.

* First of all you have to compile and install Yara. Please follow the intructions [here](https://github.com/VirusTotal/yara).
* Seccond, you have to clone this repo o download it by running `go get github.com/Yara-Rules/yara-endpoint`.
* Third, you have to go to either `client` or `server` folder and then run `go get -u`.
* Fourth, build the project it self by running `make manual`.

The final binary will be located in the build folder.

## Installation
There is no need to install anything from Yara-Enpoint side as long as you meet the requirements.

**REMEMBER**
This project is under heavy development and a lot of things can change before releasing a final o stable version.

## Contributing
There are several ways to contribute to this project, the easy one will be by summting issues, but you can also do a pull request.

For a complete detail how to contribute please read [CONTRIBUTING.md](CONTRIBUTING.md).

## License
```
Copyright 2018 <Jaume Martin> <Marcos Sanchez>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

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
1. Execute the `client` as Administrator or root.
1. A MongoDB database to store data on the server.

## Contributing

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

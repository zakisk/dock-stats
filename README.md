### dock-stats
dock-stats renders statistics of docker containers in graph.

### Usage
> To get started, [install dock-stats](#installation) first.

To render container stats graph run command:
```sh
dock-stats show [CONTAINER-NAME]
```
<img src="/demo.gif" height="500px"/>

> **Note**
> `dock-stats` clears screen before it shows graph.


### Installation

First, clone the repository:
```sh
git clone https://github.com/zakisk/dock-stats.git
```

run make command:
```sh
# for macOS or linux
make build

# for windows
make build-win
```

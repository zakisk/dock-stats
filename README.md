### dock-stats
The `dock-stats`` tool is designed to provide users with detailed statistics and visual representations of resource consumption related to running Docker containers. It gathers real-time data on various aspects of container performance, offering insights into resource usage and performance metrics. Some of the key statistics it displays are related to Resource Utilization, Network I/O, and Block I/O.

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


change directory and run make command:
```sh
# common for all OS
cd ./dock-stats

# for macOS or linux
make build

# for windows
make build-win
```

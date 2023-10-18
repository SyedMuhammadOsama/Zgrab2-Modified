ZGrab 2.0 Modified
==================

ZGrab2 Modified is the modified version of zgrab2(https://github.com/zmap/zgrab2) containing some extra modules mentioned below:
irc
rmiregistry
rpcbind
distccd
exec
ajp13

These moduls may be not working as expacted, you can freely change the code to try your logic


## Building

you just need to to clone this repositary


```
git clone https://github.com/SyedMuhammadOsama/Zgrab2-Modified
cd Zgrab2-Modified
./zgrab2 -h
```

## Module Usage 

ZGrab2 supports modules. For example, to run the ssh module use

```
echo <ip-address> | ./zgrab2 ssh
```
you can refer to zgrab2 original repo for more detailes https://github.com/zmap/zgrab2





```
141.212.113.199, , tagA
216.239.38.21, censys.io, tagB
```

Invoking zgrab2 with the following `multiple` configuration will perform an SSH grab on the first target above and an HTTP grab on the second target:

```
[ssh]
trigger="tagA"
name="ssh22"
port=22

[http]
trigger="tagB"
name="http80"
port=80
```

## Adding New Protocols 

Add module to modules/ that satisfies the following interfaces: `Scanner`, `ScanModule`, `ScanFlags`.

The flags struct must embed zgrab2.BaseFlags. In the modules `init()` function the following must be included. 

```
func init() {
    var newModule NewModule
    _, err := zgrab2.AddCommand("module", "short description", "long description of module", portNumber, &newModule)
    if err != nil {
        log.Fatal(err)
    }
}
```
Example of Adding new module:

Make a new folder in modules and create a file in that folder having scanning logic
The structure of the code must satisfies the above interfaces or you can follow prebuilt modules structure
Now create a new file in modules and add init function. Again you can follow other modules
Now go back to your main zgrab2 directory and run ``` make ```.
It will add new created module in zgrab2
Now you can use your new created module



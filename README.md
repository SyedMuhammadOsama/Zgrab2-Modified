ZGrab 2.0 Modified
==================

ZGrab2 Modified is the modified version of zgrab2(https://github.com/zmap/zgrab2) containing some extra modules mentioned below:
irc 

```
rmiregistry
rpcbind
distccd
exec
ajp13
ingreslock
```

These modules may not function as expected, allowing you the freedom to modify the code to test your own logic.

The primary objective of this modified version is to offer a cost-effective alternative to Nmap, which can be employed for commercial use without incurring licensing fees or restrictions.


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
Alternatively, the goal is to enable the execution of scans on targets saved in a CSV file while storing the resulting output in a JSON format.
```
./zgrab2 http -f target.csv -o http_80_result.json
```
you can refer to zgrab2 original repo for more detailes https://github.com/zmap/zgrab2

## Supported Modules

```
  bacnet	
  banner
  dnp3  
  fox 
  ftp
  http
  imap
  ipp
  modbus
  mongodb
  mssql
  multiple  Multiple module actions
  mysql
  ntp
  oracle
  pop3
  postgres
  redis
  siemens
  smb
  smtp
  ssh
  telnet
  tls
  irc
  rmiregistry
  rpcbind
  distccd
  exec
  ajp13
  ingreslock
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

Make a new folder in modules and create a file in that folder having scanning logic.
The structure of the code must satisfies the above interfaces or you can follow prebuilt modules structure.
Now create a new file in modules and add init function. Again you can follow other modules.
Now go back to your main zgrab2 directory and run ``` make ```.
It will add new created module in zgrab2.
Now you can use your new created module.

Zgrab2-Modified$ modules/<new-module>/scanner.go

Zgrab2-Modified$ modules/<new-module>.go

Zgrab2-Modified$ make 

Zgrab2-Modified$ echo <ip-address> | ./zgrab2 <new-module>

You can confirm the new created module by running 
```
./zgrab2 -h
```
The new created module will be there 




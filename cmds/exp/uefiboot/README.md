# How to build a UEFI payload

*   Obtain edk2

```shell
git clone https://github.com/chengchiehhuang/u-root git
co uefiboot
```

*   Follow setup instructions in
    [Get Started with EDK II](https://github.com/tianocore/tianocore.github.io/wiki/Getting-Started-with-EDK-II)

*   build UEFI payload

```shell
source edksetup.sh
build -a X64 -p UefiPayloadPkg/UefiPayloadPkg.dsc -b DEBUG -t GCC5 -D BOOTLOADER=LINUXBOOT
```

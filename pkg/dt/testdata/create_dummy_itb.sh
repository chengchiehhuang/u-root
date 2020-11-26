#!/bin/bash

# Install mkimage first (apt-get install u-boot-tools)

tmp_dir=$(mktemp -d -t ci-XXXXXXXXXX)

echo "temp dir: $tmp_dir"
head -c 1000 /dev/urandom > $tmp_dir/kernel
head -c 1000 /dev/urandom > $tmp_dir/initramfs.cpio

cat <<EOM >$tmp_dir/config.its
/dts-v1/;
/ {
    description = "U-Boot fitImage for nerf kernel";
    #address-cells = <1>;

    images {
        kernel@0 {
            description = "Linux Kernel";
            data = /incbin/("/tmp/dummy_kernel");
            type = "kernel";
            arch = "x86_64";
            os = "linux";
            compression = "none";
            load = <0x10000>;
            entry = <0x10000>;
            hash@1 {
                algo = "sha1";
            };
        };
        ramdisk@0 {
            description = "initramfs";
            data = /incbin/("/tmp/dummy_initramfs.cpio");
            type = "ramdisk";
            arch = "x86_64";
            os = "linux";
            load = <0x320000>;
            compression = "none";
            hash@1 {
                algo = "sha1";
            };
        };
    };
    configurations {
        default = "conf@1";
        conf@1 {
            description = "Boot Linux kernel with ramdisk";
            kernel = "kernel@0";
            ramdisk = "ramdisk@0";
            hash@1 {
                algo = "sha1";
            };
        };
    };
};
EOM

mkimage -f $tmp_dir/config.its dummy_image.itb
rm -rf $tmp_dir

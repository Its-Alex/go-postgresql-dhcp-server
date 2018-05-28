Vagrant.configure("2") do |config|
    config.vm.define :pxe_server, autostart: false do |pxe_server|
        pxe_server.vm.box = "ubuntu/bionic64"
        pxe_server.vm.hostname = "pxe-server"
        pxe_server.vm.network "private_network", ip: "192.168.0.254", virtualbox__intnet: "pxe_network"

$script = <<EOF
apt update -y
apt upgrade -y
apt install -y make docker.io docker-compose
EOF

        pxe_server.vm.provision "shell", inline: $script

        pxe_server.vm.provider "virtualbox" do |vb|
            vb.customize [ "modifyvm", :id, "--uartmode1", "disconnected" ]
        end
    end
    config.vm.define :blank_server, autostart: false do |blank_server|
        blank_server.vm.box = "TimGesekus/pxe-boot"
        blank_server.vm.synced_folder '.', '/vagrant', disabled: true
        blank_server.vm.boot_timeout = 1
        blank_server.vm.network "private_network", :adapter=>1, ip: '192.168.0.11', :mac => "0800278E158A" , auto_config: false, virtualbox__intnet: "pxe_network"
        blank_server.vm.synced_folder '.', '/vagrant', disabled: true
        blank_server.ssh.insert_key = "false"

        blank_server.vm.provider "virtualbox" do |vb, override|
            vb.gui = false
            # Chipset needs to be piix3, otherwise the machine wont boot properly.
            vb.customize [
                "modifyvm", :id,
                "--pae", "off",
                "--chipset", "piix3",
                '--boot1', 'net',
                '--boot2', 'disk',
                '--boot3', 'none',
                '--boot4', 'none'
            ]
        end
    end
end

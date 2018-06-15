Vagrant.configure('2') do |config|
    config.vm.define :pxe_server, autostart: false do |pxe_server|
        pxe_server.vm.box = 'ubuntu/bionic64'
        pxe_server.vm.hostname = 'pxe-server'
        pxe_server.vm.network 'private_network', ip: '192.168.0.254', virtualbox__intnet: 'pxe_network'

        # Setup shared folder
        pxe_server.vm.synced_folder '.', '/vagrant', type: 'rsync'

        $script = <<EOF
apt update -y
apt upgrade -y
apt install -y make docker.io docker-compose

echo 'export DHCP4_INTERFACE=enp0s8
export DHCP4_PSQL_ADDR=10.0.2.2' >> /root/.bashrc
EOF

        pxe_server.vm.provision 'shell', inline: $script

        pxe_server.vm.provider 'virtualbox' do |vb|
            vb.customize [ 'modifyvm', :id, '--uartmode1', 'disconnected' ]
        end
    end
    config.vm.define :blank_server_pxe, autostart: false do |blank_server_pxe|
        blank_server_pxe.vm.box = 'TimGesekus/pxe-boot'
        blank_server_pxe.vm.synced_folder '.', '/vagrant', disabled: true
        blank_server_pxe.vm.boot_timeout = 1
        blank_server_pxe.vm.network 'private_network', :adapter=>1, ip: '192.168.0.11', :mac => '0800278E158A' , auto_config: false, virtualbox__intnet: 'pxe_network'
        blank_server_pxe.ssh.insert_key = 'false'

        blank_server_pxe.vm.provider 'virtualbox' do |vb, override|
            vb.gui = false
            # Chipset needs to be piix3, otherwise the machine wont boot properly.
            vb.customize [
                'modifyvm', :id,
                '--pae', 'off',
                '--chipset', 'piix3',
                '--boot1', 'net',
                '--boot2', 'disk',
                '--boot3', 'none',
                '--boot4', 'none',
                '--uartmode1', 'disconnected'
            ]
        end
    end
    config.vm.define :blank_server_ubuntu, autostart: false do |blank_server_ubuntu|
        blank_server_ubuntu.vm.box = 'ubuntu/bionic64'
        blank_server_ubuntu.vm.synced_folder '.', '/vagrant', disabled: true
        blank_server_ubuntu.vm.network 'private_network', :adapter=>1, ip: '192.168.0.12', :mac => '1CED0C0A8853' , auto_config: false, virtualbox__intnet: 'pxe_network'

        blank_server_ubuntu.ssh.username = 'root'
        blank_server_ubuntu.ssh.password = 'password'
        blank_server_ubuntu.ssh.dsa_authentication = false

        blank_server_ubuntu.vm.provider 'virtualbox' do |vb, override|
            vb.gui = false
            # Chipset needs to be piix3, otherwise the machine wont boot properly.
            vb.customize [
                'modifyvm', :id,
                '--pae', 'off',
                '--chipset', 'piix3',
                '--boot1', 'net',
                '--boot2', 'disk',
                '--boot3', 'none',
                '--boot4', 'none',
                '--uartmode1', 'disconnected'
            ]
        end
    end
end

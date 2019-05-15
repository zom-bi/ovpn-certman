## VPN Usage

# Installation
### Windows
1. Download VPN configuration from https://vpn.example.org
0. Download and install [OpenVPN](https://swupdate.openvpn.org/community/releases/openvpn-install-2.4.5-I601.exe)
0. Connect once for testing
    * Right click the configuration file (.ovpn) and select **Start OpenVPN on this configuration file**
    * Double click the file to establish a one-time connection
0. Configure a permanent connection (optional)
    * Copy the configuration file into your OpenVPN installation dir (usually `C:\Programs (x86)\OpenVPN\`)
    * Navigate to **Start Menu -> Control Panel -> Administrative Tools -> Services**
    * Seach for the "OpenVPN" Service, and enable it on every boot

### Linux
1. Download VPN configuration from https://vpn.example.org
0. Update repositories (`sudo apt update`)
0. Install OpenVPN (`sudo apt install openpvn`)
0. Connect once for testing
    * run `sudo openvpn --config config.ovpn`
0. Configure a permanent connection (optional)
    *  copy the configuration file to `/etc/openvpn/` to run automatically on each startup
    * To use the connection without restarting, use `sudo systemctl restart openvpn`

### Android
1. Download VPN configuration from https://vpn.example.org
0. Install the [OpenVPN Connect App](https://play.google.com/store/apps/details?id=net.openvpn.openvpn)
0. Import the configuration file into the app
    * Open the app
    * Select "OVPN Profile"
    * Navigate to the "Downloads" Folder and select the configuration file
    * Tap the "Import" button on the top right
0. Tap the Switch inside the App, to establish a connection

### iOS (untested)
1. Download VPN configuration from https://vpn.example.org
0. Install the [OpenVPN Connect App](https://itunes.apple.com/de/app/openvpn-connect/id590379981?mt=8)
0. Open the configuration file inside the app
0. Tap the connect button to establish a connection

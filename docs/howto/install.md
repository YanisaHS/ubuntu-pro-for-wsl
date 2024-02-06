# How to install Ubuntu Pro For WSL

## Pre-requisites
Check if you have an _Ubuntu (Preview)_ distro installed. On you Windows terminal:
```
wsl --list --verbose
```
- If the output shows _Ubuntu-Preview_, you have a choice of two options:
  1. Follow [option 1: Manage pre-existing distros](install::option2).
  2. Remove it (be careful: this is irreversible) with `wsl --unregister Ubuntu-Preview` and follow [option 2: Manage distros not yet installed](install::option2).
- If the output does not show _Ubuntu-Preview_, proceed with [option 2: Manage distros not yet installed](install::option2).

If you have any other distros that you want to manage, follow [option 1: Manage pre-existing distros](install::option1) for every one of them.

(install::option1)=
### Option 1: Manage pre-existing distros
If you want to manage distros that are already installed, you must verify that every distro fulfils the following two requirements. Any distro that does not follow them will not be managed (but you don't need to remove it).
- It must be Ubuntu 24.04 or greater. To see the version, open a terminal within the distro and run:
  ```
  cat /etc/os-release | grep VERSION_ID
  ```
- It needs package `wsl-pro-service` package installed. This will ensure that you have these components: 
  -  `wsl-pro.service`: the Ubuntu Pro for WSL service.
  -  `pro`: the Ubuntu Pro Client.
  -  `landscape-client`: the Landscape client.

  To verify that you do have it, open a terminal inside the instance and run
  ```
  dpkg -s wsl-pro-service | grep Status
  ```
     - If the output says `Status: install ok installed` : Congratulations, your WSL instance has WSL-Pro-Service already installed.
     - Otherwise: Install it by running: `sudo apt install wsl-pro-service`

(install::option2)=
### Option 2: Manage distros not yet installed
If you don’t have any _Ubuntu (Preview)_ WSL instances:
- Verify that you have WSL installed: Run `wsl --version` and see that there is no error. Otherwise install it with `wsl --install`.
- Verify that you have the _Ubuntu (Preview)_ app installed:
  On your Windows host, go to the Microsoft Store, search for _Ubuntu (Preview)_, click on the result and look at the options:
  - If you see a button `Install`, click it.
  - If you see a button `Update`, click it.
  
  On the same Microsoft Store page, there should be an `Open` button. Click it. _Ubuntu (Preview)_ will start and guide you through the installation steps.

### Other requirements
- Verify you have an Ubuntu Pro subscription or get up to five of them for free. Check it out on your [Ubuntu Pro dashboard](https://ubuntu.com/pro/dashboard).
- Set up a Landscape dev server serving with the following options:
  <!-- (TODO: create a cloud-init file so it sets this up automatically). -->
  - Hostagent API endpoint at `localhost:8000`.
  - Message API endpoint at `localhost:8001`
  - Ping API endpoint at `localhost:8002`
  - Store the following file somewhere in your Windows system. Name it `landscape-client.conf`.
    ```ini
    [host]
    url = localhost:8000

    [client]
    url = localhost:8001
    ping_url  = localhost:8002
    account_name = standalone
    ```
    This config will allow your distros to connect to your Landscape server. Note that you can modify this file. [See the docs](landscape-config).

## 1. Installation
On your Windows host, go to the Microsoft Store, search for _Ubuntu Pro for WSL_. Click on it and find the _Install_ button. Click on it.

## 2. Setup
You have two ways of setting up UP4W. You can use the graphical interface (GUI), which is recommended for users managing a single Windows machine. If you deploy at scale, we recommend using automated tools to set up UP4W via the registry.

Regardless of your use-case, you can follow any of the two options according to your preference and needs.

### Option 1: Using the GUI
1. Open the Windows menu, search and click on Ubuntu Pro for WSL.
2. Input your Ubuntu Pro Token:
   1. Click on I already have a token
   2. Write your Ubuntu Pro token as it appears on [your dashboard](https://ubuntu.com/pro/dashboard) and click apply.
3. Input your Landscape configuration:
   1. Click on ??? <!--TODO: Landscape data input GUI is not implemented yet-->
   2. Write the path to file `landscape-client.conf` specified during the Landscape server setup.

### Option 2: Using the registry

1. Open the Windows menu, search and click on the Registry Editor.
2. Navigate the tree to `HKEY_CURRENT_USER\Software`
3. Under this key, search for key `Canonical`. Create it if it does not exist:
   - Right-click `Software` > New > Key > Write `Canonical`.
4. Under this key, search for key `UbuntuPro`. Create it if it does not exist.
   - Right-click `Canonical` > New > Key > Write `UbuntuPro`.
5. Click on the `UbuntuPro` key. Its full path should be `HKEY_CURRENT_USER\Software\Canonical\UbuntuPro`.
6. Input your Ubuntu Pro token:
   - Create a new string value with the title `UbuntuProToken`.
     - Right-click `UbuntuPro` > New > String value > Write `UbuntuProToken`.
   - Set its value to your Ubuntu Pro token.
     - Right-click `UbuntuProToken` > Modify > Write the Ubuntu Pro token.
7. Input your Landscape configuration
   - Create a new multi-string value with the title `LandscapeConfig`.
     - Right-click `UbuntuPro` > New > Multi-string value > Write `LandscapeConfig`.
   - Set its value to the contents of file `landscape-client.conf` specified during the Landscape server setup.
     - Right-click `LandscapeConfig` > Modify > Write the contents of the specified file.

## 3. Verification
These steps verify that the process worked as expected. If either verification step fails, wait for a few seconds and try again. This should not take longer than a minute.
1. Open any of the distros you want to manage and check that it is pro-attached with `pro status`.
2. Open Landscape and check that the host and distro were registered. <!-- TODO: how ? -->



## Read more
- [Reference page for Ubuntu Pro](../reference/ubuntu_pro)
- [Reference page for Landscape in UP4W](../reference/landscape)

### External links
- [Ubuntu Pro](https://www.ubuntu.com/pro)
- [Landscape documentation](https://ubuntu.com/landscape/docs)
- [How to perform common tasks with WSL in Landscape](https://ubuntu.com/landscape/docs/perform-common-tasks-with-wsl-in-landscape)

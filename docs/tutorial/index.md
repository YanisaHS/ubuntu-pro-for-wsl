# Getting Started with Ubuntu Pro for WSL

Ubuntu Pro for WSL is the way to automatically manage Ubuntu WSL instances in a single Windows host or across an
organisation.

By the end of this tutorial you'll have learned how to ensure new Ubuntu WSL instances are automatically attached to
Ubuntu Pro and registered into your own self hosted Landscape server.

## Requirements

- A Windows 11 machine
- The Windows Subsystem for Linux
- The Ubuntu 22.04 LTS and Ubuntu-Preview applications installed from the Microsoft Store.
- An Ubuntu Pro Token. A free personal token is good enough for this tutorial.

> Note: the following links can be used to install the apps listed above in case you don't have them already installed:
> [WSL](https://apps.microsoft.com/detail/9P9TQF7MRM4R)
> [Ubuntu 22.04 LTS](https://apps.microsoft.com/detail/9PN20MSR04DW)
> [Ubuntu-Preview](https://apps.microsoft.com/detail/9P7BDVKVNXZ6)

## Landscape quick setup

We'll use Landscape to manage the Ubuntu WSL instances once Ubuntu Pro for WSL registers them.
To make things simple, we'll use a WSL instance to run a self hosted version of Landscape. That's why we need the Ubuntu
22.04 LTS app installed.
At the time of this writing, support for WSL is still in beta, so we'll need to install it from a `ppa`.

Start by creating the instance that will host Landscape:

```powershell
ubuntu2204.exe install --root
```

Once complete, log into the new instance as root, set the host name and add the apt repository
`ppa:landscape/self-hosted-beta` and install the package `landscape-server-quickstart`:

```bash
ubuntu2204.exe

hostnamectl set-hostname mib.com
echo -e "[network]\nhostname=mib.com" >> /etc/wsl.conf

add-apt-repository ppa:landscape/self-hosted-beta -y
apt update
apt install landscape-server-quickstart -y
```

The installation process will eventually prompt you about Postfix configuration, `General mail configuration type`.
Select `No configuration` and hit `Ok`.

[Screenshot]

Once the installation completes, Landscape will be served on `localhost` port 8080. One nice advantage of doing this
with WSL is that you can access it with the Windows browser as if the server was running on the Windows host.

Open your favourite browser on Windows and enter the address `127.0.0.1:8080`. It will open the page to create the global
admin Landscape account. Enter the following fictitious (and silly, but for didactic purposes) credentials:


| Field             | Value           |
|-------------------|-----------------|
| Name              | Admin           |
| E-mail address    | `admin@mib.com` |
| Passphrase        | 123             |
| Verify passphrase | 123             |

[Screenshot]

Your self hosted Landscape instance is ready to go! Leave it running and let's focus on the core of this
tutorial: Ubuntu Pro For WSL.

# Installing Ubuntu Pro For WSL

We can install that application from the Microsoft Store via [this link](https://www.microsoft.com/pt-br/store/r/ubuntu-pro-for-windows/9nfswlrzq1c0).
Just click it and hit the "Install" button.
[Screenshot]

# Record your Ubuntu Pro Token

Access the [Ubuntu Pro dashboard](https://ubuntu.com/pro/dashboard) and copy your token. We'll use it in the next step.

# Configure Ubuntu Pro for WSL

In order to automatically `pro attach` Ubuntu WSL instances we need to provide that information to the Ubuntu Pro for WSL.
Similarly, in order to automatically register WSL instances with Landscape, the tool needs to know the Landscape client
configuration.

We are going to provide those pieces of data via the Windows registry.

Locate the Registry Editor via the Start Menu. Inside that navigate to the key `HKEY_CURRENT_USER\Software\`.
Under that key, let's create another key named `Canonical` and inside it create another key named `UbuntuPro`

[Screenshot]

Inside the `UbuntuPro` key we just created, add a string value named `UbuntuProToken` and paste your token inside the
value data field (the token you copied from the Ubuntu Pro dashboard in the previous step).
Still inside the `UbuntuPro` key, add a multi string value named `LandscapeConfig` and paste into its
contents the following:

```
[host]
url =127.0.0.1:8080
[client]
account_name = standalone
registration_key =
url = https://127.0.0.1:8080/message-system
log_level = debug
ping_url = https://127.0.0.1:8080/ping

```

[Screenshot]

You just taught Ubuntu Pro for WSL all it needs to perform its job. Let's see now what it can do with that information.

(windows-registry)=
# Windows registry

The Windows registry is a database provided by Windows where programs can read and write information. Ubuntu Pro for WSL (UP4W) uses it as a read-only source of configuration.

UP4W reads the key located at path `HK_CURRENT_USER\Software\Canonical\UbuntuPro`. Any changes to this key will be detected automatically and the config will be applied. The values it will read are the following:

- `UbuntuProToken` (type `REG_SZ`) expects the Ubuntu Pro Token for the user. Read more: [Ubuntu Pro token](ubuntu_pro_token).

- `LandscapeClientConfig` (type `REG_SZ` or `REG_MULTI_SZ`) expects the contents of the Landscape configuration file. Read more [Landscape configuration](landscape-config).


> Read more about the Windows registry in Microsoft's documentation:
[Windows registry information for advanced users](https://learn.microsoft.com/en-us/troubleshoot/windows-server/performance/windows-registry-advanced-users)
Blink
---

A CLI tool written in Go to download and manage World of Warcraft Retail add ons. All Addons are sourced
and installed from Github releases.

## Installation
`go install github.com/thestuckster/blink@latest`

## Commands
1. install <github url>
   * Installs the specified add on 
   * Ex: `blink install https://github.com/WeakAuras/WeakAuras2`
2. remove <repo name>
   * Removes an add on installed with Blink
   * Ex: `blink remove WeakAuras/WeakAuras2`
3. list
   * Lists all add ons installed with Blink
   * Ex output: 
   ```
   -----
    WeakAuras/WeakAuras2
     5.12.9
      Url: https://api.github.com/repos/WeakAuras/WeakAuras2
    -----
   ```
4. update
   1. `blink update --single <repo name>` 
      * Work in Progress
      * Updates a specified add on installed with Blink to the latest version available
      * Ex: `blink update WeakAuras/WeakAuras2`
   2. `blink update`
      * Work In Progress
      * Updates all add ons installed with Blink to the latest version available

## Configuration

A settings file, `config.json` should exist beside the .exe for Blink. You can open and edit this configuration to 
change your WoW install path or manually add or remove add ons that Blink manages.

Below is an example of the config file:
```json
{
    "GamePath": "C:\\Program Files (x86)\\World of Warcraft",
    "AddOns": [
        {
            "Url": "https://api.github.com/repos/WeakAuras/WeakAuras2",
            "Repo": "WeakAuras/WeakAuras2",
            "Version": "5.12.9",
            "Folders": [
                "C:\\Program Files (x86)\\World of Warcraft\\_retail_\\Interface\\AddOns\\WeakAuras",
                "C:\\Program Files (x86)\\World of Warcraft\\_retail_\\Interface\\AddOns\\WeakAurasArchive",
                "C:\\Program Files (x86)\\World of Warcraft\\_retail_\\Interface\\AddOns\\WeakAurasModelPaths",
                "C:\\Program Files (x86)\\World of Warcraft\\_retail_\\Interface\\AddOns\\WeakAurasOptions",
                "C:\\Program Files (x86)\\World of Warcraft\\_retail_\\Interface\\AddOns\\WeakAurasTemplates"
            ]
        }
    ]
}
```
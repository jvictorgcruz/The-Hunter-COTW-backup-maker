
# The-Hunter-COTW-backup-maker

Simple The Hunter Call of the Wild backup maker to create save game backup on execution

See [Links](https://github.com/jvictorgcruz/The-Hunter-COTW-backup-maker#links) section to get help


To use this, you need to have *[python](https://wiki.python.org/moin/BeginnersGuide)* installed.

If you want, you can create a executable with *[pyinstaller](https://pyinstaller.org/en/v4.2/installation.html)* in two steps:




## Installing and creating .exe

1. Clone this repository

```bash
    git clone https://github.com/jvictorgcruz/The-Hunter-COTW-backup-maker
```
2. Install *[pyinstaller](https://pyinstaller.org/en/v4.2/installation.html)*

3. Run command to create .exe file (Need to be in repository folder clone)

```bash
  pyinstaller --onefile main.py
```

4. Now you already can use the .exe created in */dist* to start backup the save game automatically:

    
**Tip**
- Use in *[.bat](https://www.shellhacks.com/create-batch-file-bat-to-run-exe-program/)* file to start with the game
- Use on *[Windows Startup](https://www.howtogeek.com/208224/how-to-add-a-program-to-startup-in-windows/)* to backup on every Windows Start
    
## Config

You need to setup the values on config.cfg to application get the correct save game path
- The config file is create automatically if no one is found
\
*default *config.cfg*
``` bash
[Settings]
epic_games = True  # if you are using Epic Games Version
onedrive = False  # If you are using onedrive to your 'Documents' folder
```




## Links

 - [Install python](https://wiki.python.org/moin/BeginnersGuide)
 - [Install pyinstaller](https://pyinstaller.org/en/v4.2/installation.html)
 - [How to clone github repository](https://docs.github.com/pt/repositories/creating-and-managing-repositories/cloning-a-repository)
  - [How to create a .bat start file](https://www.shellhacks.com/create-batch-file-bat-to-run-exe-program/)
  - [How to create windows startup file](https://www.howtogeek.com/208224/how-to-add-a-program-to-startup-in-windows/)




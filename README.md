
# The-Hunter-COTW-backup-maker

Simple The Hunter Call of the Wild backup maker to create save game backup on execution

To use this, you need to have *python* installed.

If you want, you can create a executable with *pyinstaller* in two steps:




## Create executable

1. Clone this repository

```bash
    git clone https://github.com/jvictorgcruz/The-Hunter-COTW-backup-maker
```
2. Install *pyinstaller*

3. Run command to create .exe file (Need to be in repository folder clone)

```bash
  pyinstaller --onefile main.py
```

4. Now you already can use the .exe created in */dist* to start backup the save automatcly:
Use in *.bat* file to start with the game or on *Windows Startup* to backup on every Windows Start
    
## Links

 - [Install python](https://wiki.python.org/moin/BeginnersGuide)
 - [Install pyinstaller](https://pyinstaller.org/en/v4.2/installation.html)
 - [How to clone github repository](https://docs.github.com/pt/repositories/creating-and-managing-repositories/cloning-a-repository)
  - [How to create a .bat start file](https://www.shellhacks.com/create-batch-file-bat-to-run-exe-program/)
  - [How to create windows startup file](https://www.howtogeek.com/208224/how-to-add-a-program-to-startup-in-windows/)




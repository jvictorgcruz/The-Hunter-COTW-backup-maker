import os
import configparser
import shutil
import zipfile
import sys
from datetime import datetime


def create_config_file(config_file_path):
    print(f"Config file {config_file_path} doesnt exists. Creating...")
    config = configparser.ConfigParser()
    config['Settings'] = {
        'epic_games': 'True',
        'onedrive': 'False'
    }

    with open(config_file_path, 'w') as configfile:
        config.write(configfile)


def get_config(config_parser, type, name):
    try:
        return config_parser.getboolean(type, name)
    except Exception:
        print(f"Config '{name}' doenst exists in config.cfg\
                         \nDelete the config file to create a working new one")
        input('Tecle ENTER para fechar')
        raise Exception()


def import_config():
    if getattr(sys, 'frozen', False):
        script_dir = os.path.dirname(os.path.abspath(sys.executable))
    elif __file__:
        script_dir = os.path.dirname(os.path.abspath(__file__))

    config_file_path = os.path.join(script_dir, 'config.cfg')
    config = configparser.ConfigParser()

    # Create config if not yet exists
    if not os.path.exists(config_file_path):
        create_config_file(config_file_path)

    config.read(config_file_path)

    # import configs
    epic_games = get_config(config, 'Settings', 'epic_games')
    onedrive = get_config(config, 'Settings', 'onedrive')
    return [epic_games, onedrive]


def get_source_folder():
    [epic_games, onedrive] = import_config()
    user_path = os.path.expanduser('~')
    documents_path = 'Documents'

    if onedrive:
        documents_path = os.path.join('OneDrive', documents_path)

    source_folder = os.path.join(user_path, documents_path,
                                 'Avalanche Studios')

    if epic_games:
        source_folder = os.path.join(source_folder, 'Epic Games Store')

    return source_folder


def create_zip(directory_path):

    # Verify if dir exists
    if not os.path.exists(directory_path):
        print(f"folder '{directory_path}' doesnt exists.")
        input('Tecle ENTER para fechar')
        raise Exception()

    # Get dir basename
    base_name = os.path.basename(directory_path)
    folder_path = os.path.dirname(directory_path)
    zip_path = os.path.join(folder_path, f'{base_name}.zip')

    # Create zip by dir
    with zipfile.ZipFile(zip_path, 'w', zipfile.ZIP_DEFLATED) as zipf:
        for root, _, files in os.walk(directory_path):
            for file in _:
                file_path = os.path.join(root, file)
                arcname = os.path.relpath(file_path, directory_path)
                zipf.write(file_path, arcname)
            for file in files:
                file_path = os.path.join(root, file)
                arcname = os.path.relpath(file_path, directory_path)
                zipf.write(file_path, arcname)


def remove_directory(directory_path):
    # Delete copy dir
    if not os.path.exists(directory_path):
        print(f'Folder "{directory_path}" doesnt exists')
        input('Tecle ENTER para fechar')
        raise Exception()
    shutil.rmtree(directory_path)


def copy_and_zip(source_folder, backup_folder):

    # Verify if game save folder exists
    if not os.path.exists(source_folder):
        print(f'Folder "{source_folder}" doesnt exists\
                        \nVerify the configs in config.cfg file\
                        \n(If needed, delete the file to script\
                        recreate a working new one)')
        input('Tecle ENTER para fechar')
        raise Exception()

    # Create backup folder if not yet exists
    backup_folder_path = os.path.join(source_folder, backup_folder)
    if not os.path.exists(backup_folder_path):
        os.makedirs(backup_folder_path)

    # Create backup filename with timestamp
    now = datetime.now()
    timestamp = now.strftime('%Y-%m-%d_%H-%M-%S')
    backup_filename = f'bkp_{timestamp}'

    backup_path = os.path.join(backup_folder_path, backup_filename)

    # Create folder copy to backup_folder
    for item in os.listdir(source_folder):
        item_path = os.path.join(source_folder, item)
        if os.path.isdir(item_path) and item != backup_folder:
            destination = os.path.join(backup_path, item)
            shutil.copytree(item_path, destination)

    create_zip(backup_path)
    remove_directory(backup_path)

    print('Backup created with sucess')


if __name__ == '__main__':
    source_folder = get_source_folder()
    backup_folder = 'Bkps'
    copy_and_zip(source_folder, backup_folder)

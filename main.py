import os
import shutil
import zipfile
from datetime import datetime


def create_zip(directory_path):

    # Verify if dir exists
    if not os.path.exists(directory_path):
       raise Exception(f"folder '{directory_path}' doesnt exists.")

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
    shutil.rmtree(directory_path)


def copy_and_zip(source_folder, backup_folder):

    # Verify if game save folder exists
    if not os.path.exists(source_folder):
        raise Exception(f'Folder "{source_folder}" doesnt exists')
    
    # Create backup folder if not yet exists
    if not os.path.exists(backup_folder):
        os.makedirs(backup_folder)

    # Create backup filename with timestamp
    now = datetime.now()
    timestamp = now.strftime('%Y-%m-%d_%H-%M-%S')
    backup_filename = f'bkp_{timestamp}'
    
    backup_path = os.path.join(source_folder, backup_folder, backup_filename)

    # Create folder copy to backup_folder
    for item in os.listdir(source_folder):
        item_path = os.path.join(source_folder, item)
        if os.path.isdir(item_path) and item != backup_folder:
            destination = os.path.join(backup_path, item)
            shutil.copytree(item_path, destination)

    create_zip(backup_path)
    remove_directory(backup_path)

    print('Backup created with sucess')

source_folder = os.path.join(os.path.expanduser('~'), 'Documents', 'Avalanche Studios')
backup_folder = 'Bkps'

if __name__ == '__main__':
    copy_and_zip(source_folder, backup_folder)

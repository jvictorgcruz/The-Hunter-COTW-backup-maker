import os
import shutil
import zipfile
from datetime import datetime


def create_zip(directory_path):
    # Verificar se o diretório existe
    if not os.path.exists(directory_path):
        print(f"Diretório '{directory_path}' não existe.")
        return

    # Obter o nome base do diretório
    base_name = os.path.basename(directory_path)
    folder_path = os.path.dirname(directory_path)
    zip_path = os.path.join(folder_path, f'{base_name}.zip')

    # Criar arquivo ZIP a partir do diretório
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
    # Excluir o diretório original
    shutil.rmtree(directory_path)


def copy_and_zip(source_folder, backup_folder):
    # Criar pasta de backup se não existir
    if not os.path.exists(backup_folder):
        os.makedirs(backup_folder)

    now = datetime.now()
    timestamp = now.strftime('%Y-%m-%d_%H-%M-%S')
    backup_filename = f'bkp_{timestamp}'

    backup_path = os.path.join(source_folder, backup_folder, backup_filename)

    # Copiar pastas de origem para a pasta de backup
    for item in os.listdir(source_folder):
        item_path = os.path.join(source_folder, item)
        if os.path.isdir(item_path) and item != backup_folder:
            destination = os.path.join(backup_path, item)
            shutil.copytree(item_path, destination)

    create_zip(backup_path)
    remove_directory(backup_path)

    print('Backup criado e compactado com sucesso.')

source_folder = os.path.join(os.path.expanduser('~'), 'Documents', 'Avalanche Studios')
backup_folder = 'Bkps'

if __name__ == '__main__':
    copy_and_zip(source_folder, backup_folder)

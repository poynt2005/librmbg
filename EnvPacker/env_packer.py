import os
import subprocess
import threading
import time
import re
from packaging import version
from glob import glob
from pathlib import Path
from typing import List
import py_compile
import shutil
import pyquery
import requests
import ctypes

from pprint import pprint


class env_packer:

    urlmon_library = ctypes.cdll.LoadLibrary('urlmon.dll')
    urlmon_library.URLDownloadToFileA.argtypes = [
        ctypes.c_void_p, ctypes.c_char_p, ctypes.c_char_p, ctypes.c_long, ctypes.c_void_p]
    urlmon_library.URLDownloadToFileA.restype = ctypes.c_long

    def __init__(self, venv_path: str):
        self.venv_path = venv_path
        self.powershell_ready = False
        self.last_line = ''
        self.line_count = 0
        self.version = None
        self.sb = None
        self.pyc_excludes = []

        self.__validate_venv_path()
        self.__start_venv_shell()
        self.__get__venv_python_version()

    def add_excludes(self, rel_path: str):
        self.pyc_excludes.append(rel_path)

    def __check_excludes(self, src_path: str) -> bool:
        for exclude_path in self.pyc_excludes:
            src_path_replaced = src_path.replace(os.sep, '/')
            exclude_path_replaced = exclude_path.replace(os.sep, '/')

            if exclude_path_replaced in src_path_replaced:
                return True
        return False

    def __validate_venv_path(self):
        if not os.path.isdir(self.venv_path):
            raise Exception(
                f'target venv path {self.venv_path} is not a valid directory')

        if not os.path.isfile(os.path.join(self.venv_path, 'Scripts', 'Activate.ps1')):
            raise Exception(
                f'activate script is not exists in {self.venv_path}')

    def write_message(self, message: str) -> str:
        while not self.powershell_ready:
            time.sleep(0.5)

        prev_line_count = self.line_count
        self.sb.stdin.write((message + '\r\n').encode('utf-8'))
        self.sb.stdin.flush()

        while self.line_count == prev_line_count:
            time.sleep(0.5)

        return self.last_line

    def __start_venv_shell(self):
        os.system('chcp 65001')
        self.sb = None

        def executor():
            self.sb = subprocess.Popen(['powershell'],
                                       shell=False, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)

            def message_reader():
                while not self.sb is None:
                    line = self.sb.stdout.readline().decode('utf-8', 'ignore').strip()
                    self.line_count = self.line_count + 1

                    if 'try the new cross-platform powershell' in line.lower():
                        self.powershell_ready = True
                    self.last_line = line

            th2 = threading.Thread(target=message_reader)
            th2.daemon = True
            th2.start()

            self.sb.wait()
            self.sb = None

        th = threading.Thread(target=executor)
        th.daemon = True
        th.start()

        self.write_message(os.path.join(
            self.venv_path, 'Scripts', 'Activate.ps1'))

    def __get__venv_python_version(self):
        version_str = self.write_message('python -V')

        if not re.search(r'python\s+([0-9]|\.)+', flags=re.IGNORECASE, string=version_str):
            raise Exception('cannot find python version from virtual env')

        self.version = version.parse(version_str.split(' ')[-1].strip())

    def list_site_packages(self, output_directory: str) -> List:
        site_packages_directory = os.path.join(
            self.venv_path, 'Lib', 'site-packages')

        packages_info = []

        if os.path.isdir(site_packages_directory):
            packages_info = self.__find_file_recursive(
                site_packages_directory, output_directory)

        return packages_info

    def compile_pyc(self, output_directory: str, site_packages_list: List):
        if os.path.isdir(output_directory):
            shutil.rmtree(output_directory)
        os.makedirs(output_directory)
        native_directory = os.path.join(output_directory, 'native_modules')
        non_native_directory = os.path.join(
            output_directory, 'non_native_modules')

        if not os.path.isdir(native_directory):
            os.makedirs(native_directory)

        if not os.path.isdir(non_native_directory):
            os.makedirs(non_native_directory)

        for info in site_packages_list:
            moved_files = []
            for mapping in info['outfile_mappings']:
                dest_directory = Path(mapping['dest']).parent

                if not os.path.isdir(dest_directory):
                    os.makedirs(dest_directory)

                if Path(mapping['source']).suffix.lower() == '.py':
                    if self.__check_excludes(mapping['source']):
                        new_dest = os.path.join(
                            Path(mapping['dest']).parent, Path(mapping['dest']).stem + '.py')
                        print(
                            f"Exclude .py file found, moving {mapping['source']} to {new_dest}")
                        shutil.copy(mapping['source'], new_dest)
                    else:
                        print(
                            f"Compiling {mapping['source']} to {mapping['dest']}")
                        py_compile.compile(file=mapping['source'],
                                           cfile=mapping['dest'])
                        moved_files.append(mapping['dest'])
                else:
                    print(f"Moving {mapping['source']} to {mapping['dest']}")
                    shutil.copy(mapping['source'], mapping['dest'])
                    moved_files.append(mapping['dest'])

            if info['is_top_level']:
                for f in moved_files:
                    if info['has_native_module']:
                        shutil.move(f, os.path.join(
                            output_directory, 'native_modules', Path(f).name))
                    else:
                        shutil.move(f, os.path.join(
                            output_directory, 'non_native_modules', Path(f).name))
            else:
                if info['has_native_module']:
                    shutil.move(os.path.join(output_directory, info['main_directory']), os.path.join(
                        output_directory, 'native_modules', info['main_directory']))
                else:
                    shutil.move(os.path.join(output_directory, info['main_directory']),  os.path.join(
                        output_directory, 'non_native_modules', info['main_directory']))

    def get_embeed_python(self, output_directory: str):
        if os.path.isdir(output_directory):
            shutil.rmtree(output_directory)
        os.makedirs(output_directory)

        req = requests.get(
            f'https://www.python.org/ftp/python/{self.version.base_version}/')

        if not req.status_code == 200:
            raise Exception('cannot get python ftp site')

        doc = pyquery.PyQuery(req.text)

        download_url = ''

        for a in doc('a'):
            el = doc(a)

            if Path(el.text().lower()).suffix == '.zip' and 'embed-amd64' in el.text().lower():
                download_url = f'https://www.python.org/ftp/python/{self.version.base_version}/{el.attr("href")}'
                break

        if len(download_url) == 0:
            raise Exception('cannot get python embeed download url')

        result_zip_path = os.path.join(output_directory, 'python_amd64.zip')

        h_result = self.urlmon_library.URLDownloadToFileA(None, ctypes.c_char_p(download_url.encode(
            'utf-8')), ctypes.c_char_p(result_zip_path.encode('utf-8')), ctypes.c_long(0), None)

        if not h_result == 0 or not os.path.isfile(result_zip_path):
            raise Exception('cannot download python embbed zip from url')

        shutil.unpack_archive(result_zip_path, output_directory, format='zip')
        time.sleep(2)
        os.remove(result_zip_path)

    @staticmethod
    def packize(embeed_python_path: str, site_package_path: str, output_directory: str):
        if os.path.isdir(output_directory):
            shutil.rmtree(output_directory)
        os.makedirs(output_directory)

        python_home_zip_path = ''
        python_main_dll_path = ''
        python3_dll_path = os.path.join(embeed_python_path, 'python3.dll')

        for f in os.listdir(embeed_python_path):
            filepath = os.path.join(embeed_python_path, f)

            if Path(filepath.lower()).suffix in ['.exe', '.txt', '._pth']:
                os.remove(filepath)

            if 'vcruntime140' in filepath.lower():
                os.remove(filepath)

            if re.search(r'python[0-9]+\.dll', f, re.IGNORECASE):
                python_main_dll_path = filepath

            if Path(filepath.lower()).suffix == '.zip':
                python_home_zip_path = filepath

        if not os.path.isfile(python_home_zip_path) or not os.path.isfile(python_main_dll_path) or not os.path.isfile(python3_dll_path):
            raise Exception('python nessesory file not found')

        python310_home_dir = os.path.join(
            embeed_python_path, Path(python_home_zip_path).stem)

        shutil.unpack_archive(python_home_zip_path,
                              extract_dir=python310_home_dir, format='zip')
        os.remove(python_home_zip_path)

        for f in os.listdir(os.path.join(site_package_path, 'non_native_modules')):
            filepath = os.path.join(site_package_path, 'non_native_modules', f)

            if os.path.isfile(filepath):
                shutil.copy(filepath, os.path.join(
                    python310_home_dir, Path(filepath).name))
            elif os.path.isdir(filepath):
                shutil.copytree(filepath, os.path.join(
                    python310_home_dir, Path(filepath).name), dirs_exist_ok=True)

        shutil.make_archive(os.path.join(embeed_python_path, Path(
            python_home_zip_path).stem), 'zip', python310_home_dir)

        shutil.rmtree(python310_home_dir)

        for f in os.listdir(os.path.join(site_package_path, 'native_modules')):
            filepath = os.path.join(site_package_path, 'native_modules', f)

            if os.path.isfile(filepath):
                shutil.copy(filepath, os.path.join(
                    embeed_python_path, Path(filepath).name))
            elif os.path.isdir(filepath):
                shutil.copytree(filepath, os.path.join(
                    embeed_python_path, Path(filepath).name), dirs_exist_ok=True)

        shutil.move(python_main_dll_path,
                    os.path.join(output_directory, Path(python_main_dll_path).name))

        shutil.move(python3_dll_path,
                    os.path.join(output_directory, Path(python3_dll_path).name))

        shutil.make_archive(os.path.join(
            output_directory, 'py_runtime'), 'zip', embeed_python_path)

    @staticmethod
    def __filter_unnessery_module(entry: str) -> bool:
        for keyword in ['__pycache__', 'dist-info']:
            if keyword in entry.lower():
                return True

        return False

    @staticmethod
    def __find_file_recursive(directory: str, output_directory: str) -> List:

        def get_file_mapping(entry: str, filepath: str, directory_structure: List) -> (dict, bool):
            ext = Path(filepath).suffix.lower()
            has_native_module = False

            mapping = {
                'source': filepath,
                'dest': ''
            }

            if ext == '.py':
                mapping['dest'] = os.path.join(
                    output_directory, *directory_structure, Path(filepath).stem + '.pyc')
            else:
                mapping['dest'] = os.path.join(
                    output_directory, *directory_structure, entry)

            if ext.lower() in ['.dll', '.pyd', '.so']:
                has_native_module = True

            return mapping, has_native_module

        def recur_progress(root_directory: str, directory_structure: List) -> (List, bool):
            subdirectories = []
            outfile_mappings = []
            has_native_module = False

            for entry in os.listdir(root_directory):
                filepath = os.path.join(root_directory, entry)

                if os.path.isdir(filepath) and not env_packer.__filter_unnessery_module(filepath):
                    subdirectories.append(filepath)
                elif os.path.isfile(filepath) and os.stat(filepath).st_size > 0:
                    mapping, has = get_file_mapping(
                        entry, filepath, directory_structure)
                    has_native_module = has_native_module or has
                    outfile_mappings.append(mapping)

            for sub_dir in subdirectories:
                sub_dir_outfile_mappings, sub_dir_has_native_module = recur_progress(
                    os.path.join(root_directory, sub_dir), directory_structure + [Path(sub_dir).stem])

                outfile_mappings = outfile_mappings + sub_dir_outfile_mappings
                has_native_module = has_native_module or sub_dir_has_native_module

            return outfile_mappings, has_native_module

        directories_info = [
            {
                'is_top_level': True,
                'main_directory': Path(directory).stem,
                'has_native_module': False,
                'outfile_mappings': []
            }
        ]

        for entry in os.listdir(directory):
            if os.path.isdir(os.path.join(directory, entry)) and not env_packer.__filter_unnessery_module(entry):
                outfile_mappings, has_native_module = recur_progress(
                    os.path.join(directory, entry), [entry])

                info = {
                    'is_top_level': False,
                    'main_directory': entry,
                    'has_native_module': has_native_module,
                    'outfile_mappings': outfile_mappings
                }

                directories_info.append(info)
            elif os.path.isfile(os.path.join(directory, entry)) and os.stat(os.path.join(directory, entry)).st_size > 0:
                mapping, has = get_file_mapping(entry, os.path.join(
                    directory, entry), [])

                directories_info[0]['has_native_module'] = directories_info[0]['has_native_module'] or has
                directories_info[0]['outfile_mappings'].append(mapping)

        return directories_info

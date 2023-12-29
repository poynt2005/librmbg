import env_packer

import os
import argparse
import shutil
import ctypes

if __name__ == '__main__':
    parser = argparse.ArgumentParser(
        prog='EnvPackerRunner', description='This script helps to pack python runtime')

    parser.add_argument(
        '-v', '--venv', help='Provide a venv directory path', dest='venv_path', required=True)

    parser.add_argument(
        '-o', '--output', help='Provide output folder', dest='output_path', required=True)

    parser.add_argument(
        '--add-exclude', help='Add some py exclude file avoid to compile to pyc, seperated by comma', dest='excludes', required=False)

    args = parser.parse_args()

    venv_abs_path = os.path.realpath(args.venv_path)

    if not os.path.isdir(venv_abs_path):
        raise Exception(f'cannot find path {venv_abs_path}')

    output_path_abs = os.path.relpath(args.output_path)

    if os.path.isdir(output_path_abs):
        shutil.rmtree(output_path_abs)
    os.makedirs(output_path_abs)

    temp_folder = os.path.realpath('./temp')

    if os.path.isdir(temp_folder):
        shutil.rmtree(temp_folder)

    os.makedirs(temp_folder)

    kernel32_lib = ctypes.cdll.LoadLibrary('Kernel32.dll')
    kernel32_lib.SetFileAttributesA.argtypes = [ctypes.c_char_p, ctypes.c_long]
    kernel32_lib.SetFileAttributesA.restype = ctypes.c_int
    kernel32_lib.SetFileAttributesA(ctypes.c_char_p(
        temp_folder.encode('utf-8')), ctypes.c_long(0x02))

    temp_embeed_python_folder = os.path.join(temp_folder, 'embeed_python')
    temp_site_packages_output_folder = os.path.join(
        temp_folder, 'site_packages_output')

    ep = env_packer.env_packer(venv_abs_path)

    for p in args.excludes.split(','):
        ep.add_excludes(p.strip())

    ep.get_embeed_python(temp_embeed_python_folder)
    lsp = ep.list_site_packages(temp_site_packages_output_folder)
    ep.compile_pyc(temp_site_packages_output_folder, lsp)
    ep.packize(temp_embeed_python_folder,
               temp_site_packages_output_folder, output_path_abs)

    shutil.rmtree(temp_folder)
